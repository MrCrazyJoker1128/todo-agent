package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"service/internal/middleware"
	"service/internal/service/modelconfig"
	"service/internal/utility"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/os/glog"
)

type SearchKnowledgeArgs struct {
	Query     string `json:"query" jsonschema:"用户问题/查询内容，必填"`
	LibraryId string `json:"libraryId" jsonschema:"限定某知识库ID，可空"`
}

type SearchKnowledgeResult struct {
	Results []ChunkResult `json:"results" jsonschema:"检索到的知识库片段"`
	Count   int           `json:"count" jsonschema:"结果数量"`
	Context string        `json:"context" jsonschema:"构建的上下文，用于LLM回答"`
	Success bool          `json:"success" jsonschema:"是否成功"`
	Message string        `json:"message" jsonschema:"结果消息"`
}

type ChunkResult struct {
	Content string  `json:"content" jsonschema:"片段内容"`
	Title   string  `json:"title" jsonschema:"来源标题"`
	Author  string  `json:"author" jsonschema:"来源作者"`
	URL     string  `json:"url" jsonschema:"来源URL"`
	Score   float64 `json:"score" jsonschema:"相似度分数"`
}

func SearchKnowledgeTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "search_knowledge",
			Desc: "搜索知识库。当用户询问问题，需要从已添加的知识库中查找相关信息时调用。返回与问题相关的知识库片段。libraryId 候选见系统提示中的知识库列表（ID 照抄不编造；不确定查哪个库就留空查全库）。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"query": {
						Desc:     "用户的问题或查询内容",
						Type:     "string",
						Required: true,
					},
					"libraryId": {
						Desc:     "限定某知识库ID，可空；候选见系统提示",
						Type:     "string",
						Required: false,
					},
				}),
		},
		searchKnowledge,
	)
}

func searchKnowledge(ctx context.Context, input *SearchKnowledgeArgs) (string, error) {
	glog.Info(ctx, "[searchKnowledge] 开始搜索, query: %s", input.Query)

	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		glog.Error(ctx, "[searchKnowledge] user_id not found in context")
		return "", errors.New("user_id not found in context")
	}
	// 历史回放类（名字/刚才说过的话）答案在对话历史里：tool 被模型误调时直接短路，
	// 不 embed 不查 ES。返回正常结果形，content 写明原因，上层自然转述。
	if isHistoryRecallQuery(input.Query) {
		glog.Info(ctx, "[searchKnowledge] 历史回放类查询，直接短路")
		result := SearchKnowledgeResult{
			Success: true,
			Message: "答案在对话历史中，不在知识库",
			Count:   0,
			Results: []ChunkResult{},
			Context: "该问题问的是对话历史（名字/身份/刚才说过的话），答案在上文对话里，不在知识库中。请只根据对话历史回答。",
		}
		resJSON, _ := json.Marshal(result)
		return string(resJSON), nil
	}
	glog.Info(ctx, "[searchKnowledge] userId: %s", userId)

	if input.Query == "" {
		return "", errors.New("查询内容不能为空")
	}

	query := rewriteQuery(input.Query)

	embedding, err := getQueryEmbedding(ctx, query)
	if err != nil {
		glog.Error(ctx, "[searchKnowledge] 获取嵌入失败: %v", err)
		return "", fmt.Errorf("获取嵌入失败: %w", err)
	}
	glog.Info(ctx, "[searchKnowledge] 嵌入生成成功, 向量维度: %d", len(embedding))

	chunks, err := searchChunksByVector(ctx, userId, input.LibraryId, embedding)
	if err != nil {
		glog.Error(ctx, "[searchKnowledge] ES搜索失败: %v", err)
		return "", fmt.Errorf("搜索失败: %w", err)
	}
	glog.Info(ctx, "[searchKnowledge] 搜索完成, 找到 %d 条结果", len(chunks))

	if len(chunks) == 0 {
		result := SearchKnowledgeResult{
			Success: true,
			Message: "未找到相关内容",
			Count:   0,
			Results: []ChunkResult{},
			Context: "没有找到相关的知识库内容。",
		}
		resJSON, _ := json.Marshal(result)
		return string(resJSON), nil
	}

	contextStr := buildContext(chunks)
	results := make([]ChunkResult, 0, len(chunks))
	for _, chunk := range chunks {
		results = append(results, ChunkResult{
			Content: chunk.Content,
			Title:   chunk.Title,
			Author:  chunk.Author,
			URL:     chunk.URL,
			Score:   chunk.Score,
		})
	}

	result := SearchKnowledgeResult{
		Success: true,
		Message: fmt.Sprintf("找到%d条相关内容", len(chunks)),
		Count:   len(chunks),
		Results: results,
		Context: contextStr,
	}

	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

func getQueryEmbedding(ctx context.Context, query string) ([]float64, error) {
	embedder, err := modelconfig.NewEmbedder(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("创建嵌入器失败: %w", err)
	}

	embeddings, err := embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("嵌入失败: %w", err)
	}

	if len(embeddings) == 0 {
		return nil, errors.New("未生成嵌入向量")
	}

	return embeddings[0], nil
}

type chunkData struct {
	Content string
	Title   string
	Author  string
	URL     string
	Source  string
	Score   float64
}

func searchChunksByVector(ctx context.Context, userId, libraryId string, queryVector []float64) ([]chunkData, error) {
	esClient, err := utility.ES()
	if err != nil {
		return nil, fmt.Errorf("获取ES客户端失败: %w", err)
	}

	filters := []map[string]interface{}{
		{
			"term": map[string]interface{}{
				"user_id.keyword": userId,
			},
		},
	}
	if strings.TrimSpace(libraryId) != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"library_id.keyword": libraryId,
			},
		})
	}
	query := map[string]interface{}{
		"size": 5,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filters,
			},
		},
		"knn": map[string]interface{}{
			"field":          "content_vector",
			"query_vector":   queryVector,
			"k":              5,
			"num_candidates": 20,
		},
		"_source": []string{"content", "title", "author", "url_id", "source"},
	}

	res, err := esClient.Search(
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex("knowledge_chunks"),
		esClient.Search.WithBody(formatJSON(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("ES搜索失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES搜索错误: %s", res.String())
	}

	var response struct {
		Hits struct {
			Hits []struct {
				Score  float64 `json:"_score"`
				Source struct {
					Content string `json:"content"`
					Title   string `json:"title"`
					Author  string `json:"author"`
					URLID   string `json:"url_id"`
					Source  string `json:"source"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析ES响应失败: %w", err)
	}

	chunks := make([]chunkData, 0, len(response.Hits.Hits))
	for _, hit := range response.Hits.Hits {
		urlInfo, _ := getURLInfo(ctx, hit.Source.URLID)
		chunks = append(chunks, chunkData{
			Content: hit.Source.Content,
			Title:   hit.Source.Title,
			Author:  hit.Source.Author,
			URL:     urlInfo,
			Source:  hit.Source.Source,
			Score:   hit.Score,
		})
	}

	return chunks, nil
}

func getURLInfo(ctx context.Context, urlId string) (string, error) {
	type urlRecord struct {
		Url string `json:"url"`
	}
	res, err := utility.ESGet(ctx, "knowledge_urls", urlId)
	if err != nil {
		return "", nil
	}

	var urlData struct {
		Source urlRecord `json:"_source"`
	}
	if err := json.Unmarshal(res, &urlData); err != nil {
		return "", nil
	}
	return urlData.Source.Url, nil
}

func buildContext(chunks []chunkData) string {
	context := "以下是知识库中相关的内容（按相似度从高到低排列）：\n\n"
	for i, chunk := range chunks {
		if i == 0 {
			context += "【来源1·最相关】\n"
		} else {
			context += fmt.Sprintf("【来源%d】\n", i+1)
		}
		if chunk.Title != "" {
			context += fmt.Sprintf("标题：%s\n", chunk.Title)
		}
		if chunk.Author != "" {
			context += fmt.Sprintf("作者：%s\n", chunk.Author)
		}
		if chunk.URL != "" {
			context += fmt.Sprintf("链接：%s\n", chunk.URL)
		} else if chunk.Source != "" {
			context += fmt.Sprintf("来源：%s\n", chunk.Source)
		}
		context += fmt.Sprintf("内容：%s\n\n", chunk.Content)
	}
	return context
}

// historyRecallKw 历史回放信号词：命中则答案在对话历史里，不走向量检索。
var historyRecallKw = []string{"刚才", "刚刚", "上文", "之前问", "之前说", "历史", "聊天记录", "回顾一下", "名字", "我叫", "你叫"}

// isHistoryRecallQuery tool 侧历史回放判定（与 chat 侧 isHistoryRecall 同语义，防 import cycle 各自一份）。
func isHistoryRecallQuery(q string) bool {
	for _, kw := range historyRecallKw {
		if strings.Contains(q, kw) {
			return true
		}
	}
	return false
}

// rewriteQuery 去口语前缀再_embed。正则只处理行首，不调模型。
func rewriteQuery(q string) string {
	q = strings.TrimSpace(q)
	// 前端快捷标签（/添加任务 /查询知识库 …）：是 UI 动作名不是问题内容，去掉。
	// 标签与正文可能无空格分隔，按已知标签表逐个试。
	for _, tag := range []string{
		"/添加任务", "/编辑任务", "/删除任务", "/查询任务",
		"/添加知识库", "/查询知识库",
	} {
		if strings.HasPrefix(q, tag) {
			q = strings.TrimSpace(strings.TrimPrefix(q, tag))
			break
		}
	}
	if strings.HasPrefix(q, "/") {
		q = strings.TrimSpace(strings.TrimPrefix(q, "/"))
	}
	prefixes := []string{
		"帮我查一下", "帮我找一下", "请问", "请帮我",
		"查一下", "找一下", "查询", "搜索",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(q, p) {
			q = strings.TrimSpace(strings.TrimPrefix(q, p))
			break
		}
	}
	if q == "" {
		return q
	}
	return q
}

func formatJSON(v interface{}) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
