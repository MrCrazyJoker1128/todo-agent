package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"service/internal/dao"
	"service/internal/middleware"
	"service/internal/model/entity"
	"time"

	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

type AddKnowledgeArgs struct {
	URL string `json:"url" jsonschema:"知识库链接，必填"`
}

type AddKnowledgeResult struct {
	Id          string `json:"id" jsonschema:"知识库ID"`
	Title       string `json:"title" jsonschema:"文章标题"`
	Author      string `json:"author" jsonschema:"作者"`
	Description string `json:"description" jsonschema:"内容描述"`
	Success     bool   `json:"success" jsonschema:"是否成功"`
	Message     string `json:"message" jsonschema:"结果消息"`
}

func AddKnowledgeTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "add_knowledge",
			Desc: "添加知识库。当用户说添加知识库、保存到知识库、记一个知识等意图时调用。需要用户提供要保存的网页URL。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"url": {
						Desc:     "要添加到知识库的网页URL",
						Type:     "string",
						Required: true,
					},
				}),
		},
		addKnowledge,
	)
}

func addKnowledge(ctx context.Context, input *AddKnowledgeArgs) (string, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return "", errors.New("user_id not found in context")
	}

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return "", errors.New("无效的用户ID")
	}

	if input.URL == "" {
		return "", errors.New("URL不能为空")
	}

	if !isValidURL(input.URL) {
		return "", errors.New("无效的URL格式")
	}

	if !isURLReachable(ctx, input.URL) {
		return "", errors.New("该URL无法访问，请检查链接是否正确")
	}

	existingUrl, err := dao.DbKnowledgeUrls.GetUrlByUserAndUrl(ctx, userId, input.URL)
	if err != nil {
		g.Log().Debug(ctx, "check existing url: %v", err)
	}
	if existingUrl != nil {
		return "", errors.New("该知识库已存在，请勿重复添加")
	}

	contentData, err := extractWebContent(ctx, input.URL)
	if err != nil {
		return "", fmt.Errorf("提取内容失败: %w", err)
	}

	urlId := uuid.New()
	now := gtime.Now()
	hashStr := fmt.Sprintf("%s-%d", input.URL, time.Now().UnixNano())

	urlRecord := &entity.DbKnowledgeUrls{
		Id:           urlId,
		UserId:       userIdUUID,
		Url:          input.URL,
		Title:        contentData.Title,
		Description:  contentData.Description,
		ContentType:  "article",
		Domain:       extractDomain(input.URL),
		FaviconUrl:   "",
		Author:       contentData.Author,
		PublishDate:  nil,
		WordCount:    len(contentData.Content),
		ChunkCount:  0,
		ContentHash: hashStr,
		Status:      "pending",
		ErrorMessage: "",
		IsDeleted:  false,
		DeletedAt:  nil,
		CrawledAt:  now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = dao.DbKnowledgeUrls.InsertUrl(ctx, urlRecord)
	if err != nil {
		return "", fmt.Errorf("保存知识库失败: %w", err)
	}

	go processKnowledgeChunksAsync(urlId, userIdUUID, contentData, input.URL, hashStr, now)

	result := AddKnowledgeResult{
		Id:          urlId.String(),
		Title:       contentData.Title,
		Author:      contentData.Author,
		Description: contentData.Description,
		Success:     true,
		Message:     "知识库添加成功",
	}

	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

func isValidURL(urlStr string) bool {
	parsed, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return false
	}
	return parsed.URL.Scheme != "" && parsed.URL.Host != ""
}

func isURLReachable(ctx context.Context, urlStr string) bool {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("HEAD", urlStr, nil)
	if err != nil {
		return false
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		getReq, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return false
		}
		getReq = getReq.WithContext(ctx)
		resp, err = client.Do(getReq)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode < 400
	}

	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func extractDomain(rawURL string) string {
	parsed, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return ""
	}
	return parsed.URL.Host
}

type WebContentData struct {
	Title       string
	Author      string
	Description string
	Content     string
}

func extractWebContent(ctx context.Context, urlStr string) (*WebContentData, error) {
	loader, err := url.NewLoader(ctx, &url.LoaderConfig{})
	if err != nil {
		return nil, fmt.Errorf("创建加载器失败: %w", err)
	}

	docs, err := loader.Load(ctx, document.Source{URI: urlStr})
	if err != nil {
		return nil, fmt.Errorf("加载网页失败: %w", err)
	}

	if len(docs) == 0 {
		return nil, errors.New("未获取到内容")
	}

	rawContent := docs[0].Content

	contentData, err := callExtractLLM(ctx, urlStr)
	if err != nil {
		return nil, err
	}

	if contentData.Title == "" {
		contentData.Title = "未提取到标题"
	}
	if contentData.Content == "" {
		contentData.Content = rawContent
		if len(contentData.Content) > 20000 {
			contentData.Content = contentData.Content[:20000]
		}
	}
	if contentData.Description == "" {
		contentData.Description = "内容提取完成"
	}

	return contentData, nil
}

func callExtractLLM(ctx context.Context, urlStr string) (*WebContentData, error) {
	extractTool := ExtractContentTool()
	resultStr, err := extractTool.InvokableRun(ctx, fmt.Sprintf(`{"url":"%s"}`, urlStr))
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	var extractResult ExtractContentResult
	if err := json.Unmarshal([]byte(resultStr), &extractResult); err != nil {
		return nil, fmt.Errorf("解析LLM结果失败: %w", err)
	}

	return &WebContentData{
		Title:       extractResult.Title,
		Author:      extractResult.Author,
		Description: extractResult.Description,
		Content:     extractResult.Content,
	}, nil
}

func processKnowledgeChunksAsync(urlId uuid.UUID, userId uuid.UUID, contentData *WebContentData, pageUrl, hashStr string, now *gtime.Time) {
	ctx := context.Background()
	dao.DbKnowledgeUrls.UpdateStatus(ctx, urlId.String(), "processing", "")
	dao.DbKnowledgeUrls.UpdateStatus(ctx, urlId.String(), "completed", "")
}