package tools

import (
	"context"
	"errors"
	"fmt"
	"service/internal/service/modelconfig"
	"strings"

	"github.com/cloudwego/eino-ext/components/document/loader/url"
	openaiModel "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

// ExtractContentArgs 提取网页内容参数
type ExtractContentArgs struct {
	URL string `json:"url" jsonschema:"网页URL，必填"`
}

// ExtractContentResult 提取内容结果
type ExtractContentResult struct {
	Title       string `json:"title" jsonschema:"网页标题"`
	Author      string `json:"author" jsonschema:"作者"`
	Description string `json:"description" jsonschema:"内容精简描述"`
	Content     string `json:"content" jsonschema:"主要内容，Markdown格式"`
	URL         string `json:"url" jsonschema:"来源URL"`
}

// ExtractContentTool 创建提取网页内容的工具
func ExtractContentTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "extract_content",
			Desc: "从URL中提取网页的标题、作者、主要内容，并生成精简描述。返回的内容为Markdown格式。当用户需要获取网页内容摘要、分析文章、提取关键信息时调用。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"url": {
						Desc:     "要提取内容的网页URL",
						Type:     "string",
						Required: true,
					},
				}),
		},
		extractContent,
	)
}

// createMiniMaxChatModel 创建 MiniMax ChatModel (与 chat.go 一致对标 eino-learn 裸三字段)
// Key 来源：当前用户的模型配置（AES-256-GCM 密文入库此处解密到内存）。
func createMiniMaxChatModel(ctx context.Context) (model.ChatModel, error) {
	cred, err := modelconfig.Resolve(ctx)
	if err != nil {
		return nil, err
	}
	return openaiModel.NewChatModel(ctx, &openaiModel.ChatModelConfig{
		BaseURL: cred.BaseURL,
		APIKey:  cred.APIKey,
		Model:   cred.Model,
	})
}

// extractContent 执行提取内容
func extractContent(ctx context.Context, input *ExtractContentArgs) (string, error) {
	if input.URL == "" {
		return "", errors.New("URL不能为空")
	}

	// 1. 加载网页内容
	loader, err := url.NewLoader(ctx, &url.LoaderConfig{})
	if err != nil {
		return "", fmt.Errorf("创建URL加载器失败: %w", err)
	}

	docs, err := loader.Load(ctx, document.Source{URI: input.URL})
	if err != nil {
		return "", fmt.Errorf("加载网页失败: %w", err)
	}

	if len(docs) == 0 {
		return "", errors.New("未获取到内容")
	}

	rawContent := docs[0].Content

	if !isArticleContent(rawContent) {
		return "", errors.New("该地址不是文章内容，无法提取。请输入文章详情页的URL，而不是网站首页或列表页")
	}

	// 2. 使用 MiniMax LLM 提取关键信息并转换为 Markdown
	chatModel, err := createMiniMaxChatModel(ctx)
	if err != nil {
		return "", fmt.Errorf("创建LLM失败: %w", err)
	}

	// 3. 构建提取prompt - 要求输出Markdown格式
	systemPrompt := `你是一个内容提取专家。请从网页内容中提取以下信息，并以JSON格式返回：
1. title: 网页标题
2. author: 作者（如果找不到则写"未知"）
3. content: 主要内容，转换为Markdown格式，保留标题、段落、列表、代码块等格式，删除广告、导航等无关内容
4. description: 用50字以内概括文章主题

返回格式（必须是有效JSON）：
{"title":"标题","author":"作者","description":"简短描述","content":"Markdown格式的主要内容"}

注意：content字段必须是Markdown格式，包含标题(##)、段落、列表等 Markdown 语法。`

	// 截取前20000字符避免超出token限制
	maxContent := rawContent
	if len(maxContent) > 20000 {
		maxContent = maxContent[:20000]
	}

	messages := []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
		{Role: schema.User, Content: "请提取以下网页内容的关键信息，转换为Markdown格式：\n\n" + maxContent},
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM生成失败: %w", err)
	}

	// 4. 解析LLM返回结果
	resultStr := resp.Content
	resultStr = strings.TrimSpace(resultStr)

	// 去掉可能的markdown code block
	if strings.HasPrefix(resultStr, "```") {
		lines := strings.SplitN(resultStr, "\n", 2)
		if len(lines) == 2 {
			resultStr = lines[1]
		}
	}
	if strings.HasSuffix(strings.TrimSpace(resultStr), "```") {
		resultStr = strings.TrimSuffix(resultStr, "```")
	}
	resultStr = strings.TrimSpace(resultStr)

	g.Log().Info(ctx, "extract_content result length: %d", len(resultStr))

	// 构建完整结果
	fullResult := &ExtractContentResult{
		URL:         input.URL,
		Description: "内容提取完成",
	}

	// 解析JSON
	fullResult = parseExtractResult(resultStr, input.URL)

	// 确保有默认值
	if fullResult.Title == "" {
		fullResult.Title = "未提取到标题"
	}
	if fullResult.Content == "" {
		fullResult.Content = maxContent
	}
	if fullResult.Description == "" {
		fullResult.Description = "内容提取完成"
	}

	return fmt.Sprintf(`{"title":"%s","author":"%s","description":"%s","content":"%s","url":"%s"}`,
		escapeJSON(fullResult.Title), escapeJSON(fullResult.Author),
		escapeJSON(fullResult.Description), escapeJSON(fullResult.Content),
		fullResult.URL), nil
}

// parseExtractResult 解析LLM返回的JSON结果
func parseExtractResult(jsonStr, url string) *ExtractContentResult {
	result := &ExtractContentResult{URL: url}

	startIdx := strings.Index(jsonStr, "{")
	endIdx := strings.LastIndex(jsonStr, "}")
	if startIdx == -1 || endIdx == -1 {
		return result
	}

	jsonContent := jsonStr[startIdx : endIdx+1]

	result.Title = extractJSONField(jsonContent, "title")
	result.Author = extractJSONField(jsonContent, "author")
	result.Description = extractJSONField(jsonContent, "description")
	result.Content = extractJSONField(jsonContent, "content")

	return result
}

// extractJSONField 从JSON中提取字段值
func extractJSONField(json, field string) string {
	fieldStart := strings.Index(json, field)
	if fieldStart == -1 {
		return ""
	}

	// 找到冒号的位置
	colonIdx := strings.Index(json[fieldStart:], ":")
	if colonIdx == -1 {
		return ""
	}
	colonIdx += fieldStart

	// 找到值的开始位置
	valueStart := colonIdx + 1
	for valueStart < len(json) && (json[valueStart] == ' ' || json[valueStart] == '"') {
		if json[valueStart] == '"' {
			valueStart++
			break
		}
		valueStart++
	}

	// 找到值的结束位置
	valueEnd := valueStart
	isEscape := false
	for valueEnd < len(json) {
		if isEscape {
			isEscape = false
		} else if json[valueEnd] == '\\' {
			isEscape = true
		} else if json[valueEnd] == '"' {
			break
		}
		valueEnd++
	}

	return json[valueStart:valueEnd]
}

// escapeJSON 转义JSON字符串
func isArticleContent(content string) bool {
	navCount := strings.Count(content, "<nav") + strings.Count(content, "<header")
	articleCount := strings.Count(content, "<article") + strings.Count(content, "<main")
	listCount := strings.Count(content, "<ul") + strings.Count(content, "<ol") + strings.Count(content, "<li")

	hasArticleStructure := articleCount > 0 || strings.Contains(content, "article") || strings.Contains(content, "post")
	hasMultipleLinks := strings.Count(content, "<a ") > 20
	hasNavStructure := navCount > 2 || strings.Contains(content, "navigation")

	if hasNavStructure && hasMultipleLinks && articleCount == 0 {
		return false
	}

	if listCount > 30 && articleCount == 0 {
		return false
	}

	return len(content) > 500 && (hasArticleStructure || !hasNavStructure || articleCount > 0)
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
