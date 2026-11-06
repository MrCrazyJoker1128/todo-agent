package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"service/internal/service/modelconfig"
	"service/internal/tools"
	"strings"
	"sync"
	"time"

	openaiModel "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	einoAgent "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// StreamEventHandler 流式事件回调。complete 事件携带 taskID（无任务时为 uuid.Nil）。
type StreamEventHandler func(event, message string, chatID, taskID uuid.UUID)

// GetChatList 获取聊天列表纯业务。
func GetChatList(ctx context.Context, userID string, req *api.GetChatListReq) (*api.GetChatListRes, error) {
	chatMessages, total, err := dao.DbChatMessages.GetChatMessagesWithTask(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("get chat messages by user id: %w", err)
	}
	return &api.GetChatListRes{
		List:  chatMessages,
		Total: total,
	}, nil
}

// InsertChatMessage 插入聊天消息纯业务。
func InsertChatMessage(ctx context.Context, userID string, req *api.InsertChatMessageReq) (*api.InsertChatMessageRes, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	var relatedTaskID uuid.UUID
	if req.RelatedTaskId != "" {
		relatedTaskID, err = uuid.Parse(req.RelatedTaskId)
		if err != nil {
			return nil, fmt.Errorf("invalid related task id: %w", err)
		}
	}

	chatMessage := &entity.DbChatMessages{
		UserId:        parsedUserID,
		Content:       req.Content,
		Role:          req.Role,
		RelatedTaskId: relatedTaskID,
		Id:            uuid.New(),
	}
	if err := dao.DbChatMessages.InsertChatMessage(ctx, chatMessage); err != nil {
		return nil, fmt.Errorf("insert chat message: %w", err)
	}
	return &api.InsertChatMessageRes{
		Success: true,
	}, nil
}

// ErrModelAuth 模型侧鉴权失败（Key 无效/过期/欠费等）。上层转 error:MODEL_AUTH，前端弹去设置页换 Key。
var ErrModelAuth = errors.New("模型 Key 无效或已过期，请前往设置页更换")

// isModelAuthError 识别模型侧 401/鉴权错误（ark SDK 透出的 authorized_error/login fail）。
func isModelAuthError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "401") &&
		(strings.Contains(msg, "authorized_error") ||
			strings.Contains(msg, "login fail") ||
			strings.Contains(msg, "authorization"))
}

// taskEditConfirmed 用户是否已确认创建（确认文本由前端 TaskFlowCard 组装）。
func taskEditConfirmed(ctx context.Context) bool {
	content, _ := ctx.Value("user_content").(string)
	return strings.HasPrefix(content, "确认创建任务：")
}

// taskSummaryConfirmed 用户是否已确认生成总结（确认文本由前端 TaskSummaryCard 组装）。
func taskSummaryConfirmed(ctx context.Context) bool {
	content, _ := ctx.Value("user_content").(string)
	return strings.HasPrefix(content, "确认生成总结：")
}

// normalizeSummaryArgs 归一化总结维度建议：非法 period→空；period 为空但起止齐→custom；分块逻辑在 tool 侧算窗。
func normalizeSummaryArgs(args map[string]any) map[string]any {
	str := func(k string) string {
		if v, ok := args[k].(string); ok {
			return strings.TrimSpace(v)
		}
		return ""
	}
	p := strings.ToLower(str("period"))
	if p != "day" && p != "week" && p != "month" && p != "custom" {
		p = ""
	}
	start, end := str("startDate"), str("endDate")
	if p == "" && start != "" && end != "" {
		p = "custom"
	}
	missing := []string{}
	if p == "" {
		missing = append(missing, "period")
	}
	return map[string]any{
		"period":      p,
		"startDate":   start,
		"endDate":     end,
		"taskType":    strings.TrimSpace(str("taskType")),
		"priority":    strings.ToLower(strings.TrimSpace(str("priority"))),
		"missing":     missing,
		"aiSuggested": true,
	}
}

// RunAgentStream 运行 agent 流式对话：保存用户消息、建模、stream 主循环、
// 组装知识库结果、保存助手消息，全程经 onEvent 回吐事件。
// 模型未配置时透出 modelconfig.ErrNotConfigured，由上层转 NEED_MODEL_CONFIG。
func RunAgentStream(ctx context.Context, userID, content string, onEvent StreamEventHandler) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	// user_id 进 ctx供 modelconfig.Resolve 与工具调用取用；用户原文进 ctx 供 create_task 预览/真建判定。
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "user_content", content)

	if done, reply := handleConfirmText(ctx, userID, content, parsedUserID, onEvent); done {
		return reply
	}

	if reply, ok := cancelReply(content); ok {
		if err := saveChatMessage(ctx, uuid.New(), parsedUserID, content, "user", uuid.Nil); err != nil {
			return fmt.Errorf("insert chat message: %w", err)
		}
		onEvent("received", content, uuid.New(), uuid.Nil)
		systemChatID := uuid.New()
		onEvent("send_response", "", systemChatID, uuid.Nil)
		finalContent := reply
		if err := saveChatMessage(ctx, systemChatID, parsedUserID, finalContent, "assistant", uuid.Nil); err != nil {
			return fmt.Errorf("insert chat message failed: %w", err)
		}
		onEvent("complete", finalContent, systemChatID, uuid.Nil)
		return nil
	}

	if err := saveChatMessage(ctx, uuid.New(), parsedUserID, content, "user", uuid.Nil); err != nil {
		return fmt.Errorf("insert chat message: %w", err)
	}

	onEvent("received", content, uuid.New(), uuid.Nil)

	chatModel, cred, err := createChatModel(ctx)
	if err != nil {
		if errors.Is(err, modelconfig.ErrNotConfigured) {
			return err
		}
		return fmt.Errorf("create chat model: %w", err)
	}

	ragCtx := retrieveForChat(ctx, userID, content)

	agent, toolHandler, err := createAgent(ctx, chatModel, getSystemPrompt(cred.SystemPrompt)+libraryCatalog(ctx, userID), ragCtx)
	if err != nil {
		return fmt.Errorf("创建agent失败: %w", err)
	}

	prev := loadRecentHistory(ctx, userID, 30)
	if isCreateIntent(content) {
		// 当前就是创建意图：历史里的一模一样记录是上一次执行，不是本次的重复——
		// 把它们从输入历史里拿掉，模型才敢调 create_task，否则会压成 query_tasks。
		filtered := prev[:0:0]
		for _, m := range prev {
			if m.Role == schema.User && isCreateIntent(m.Content) {
				continue
			}
			if m.Role == schema.Assistant && (strings.Contains(m.Content, "已创建任务") || strings.Contains(m.Content, "task-edit")) {
				continue
			}
			filtered = append(filtered, m)
		}
		prev = filtered
	}
	chatMsg := append(prev, &schema.Message{Role: schema.User, Content: content})

	callback := react.BuildAgentCallback(nil, nil)
	sr, err := agent.Stream(ctx, chatMsg, einoAgent.WithComposeOptions(compose.WithCallbacks(callback), compose.WithCallbacks(toolHandler)))
	if err != nil {
		if isModelAuthError(err) {
			return ErrModelAuth
		}
		return fmt.Errorf("stream failed: %w", err)
	}

	systemChatID := uuid.New()
	onEvent("send_response", "", systemChatID, uuid.Nil)

	toolCallBuffer := make(map[int]*schema.ToolCall)
	finalContent := ""

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if isModelAuthError(err) {
				return ErrModelAuth
			}
			return fmt.Errorf("recv failed: %w", err)
		}

		handleToolCalls(ctx, msg, toolCallBuffer, userID, toolHandler)
		onEvent("send_response", msg.Content, systemChatID, uuid.Nil)
		finalContent += msg.Content
	}

	knowledgeResult := toolHandler.GetKnowledgeResult()
	knowledgeError := toolHandler.GetKnowledgeError()

	if knowledgeResult != "" {
		var kr map[string]interface{}
		json.Unmarshal([]byte(knowledgeResult), &kr)

		// 处理 add_knowledge 结果
		if title, ok := kr["title"].(string); ok {
			finalContent += "\n\n✅ 知识库添加成功：\n- 标题：" + title
			if author, ok := kr["author"].(string); ok && author != "" {
				finalContent += "\n- 作者：" + author
			}
		}

		// 处理 search_knowledge 结果
		if context, ok := kr["context"].(string); ok && context != "" {
			finalContent += "\n\n" + context
		}
	}

	if knowledgeError != "" {
		finalContent += "\n\n❌ 添加知识库失败：" + knowledgeError
	}

	// 任务人机验证 envelope：前端解析渲染候选/结果卡，原文同时入库供历史回放。
	if candidates := toolHandler.GetTaskCandidates(); candidates != "" {
		finalContent += "\n\n```task-candidates\n" + candidates + "\n```\n"
	}
	if listText := toolHandler.GetTaskListText(); listText != "" {
		finalContent += "\n\n" + listText
	}
	if updateResult := toolHandler.GetTaskUpdateResult(); updateResult != "" {
		finalContent += "\n\n```task-updated\n" + updateResult + "\n```\n"
	}
	if preview := toolHandler.GetTaskEditPreview(); preview != "" {
		finalContent += "\n\n```task-edit\n" + preview + "\n```\n"
	}
	if summaryPreview := toolHandler.GetTaskSummaryPreview(); summaryPreview != "" {
		finalContent += "\n\n```task-summary-confirm\n" + summaryPreview + "\n```\n"
	}
	if summaryResult := toolHandler.GetTaskSummaryResult(); summaryResult != "" {
		finalContent += "\n\n```task-summary\n" + summaryResult + "\n```\n"
	}

	if strings.TrimSpace(finalContent) != "" {
		// <think> 推理块剥掉再入库/回吐：那是过程不是结论，存了污染历史、发了污染展示。
		finalContent = strings.TrimSpace(thinkRe.ReplaceAllString(finalContent, ""))
		// 写入前过一遍断裂修复：被丢弃的空壳不入库，防空 content 行污染下轮历史。
		checked := FixToolPairs([]*schema.Message{
			{Role: schema.User, Content: content},
			{Role: schema.Assistant, Content: finalContent},
		})
		for _, m := range checked {
			if m.Role == schema.Assistant && strings.TrimSpace(m.Content) != "" {
				if err := saveChatMessage(ctx, systemChatID, parsedUserID, m.Content, "assistant", toolHandler.GetTaskID()); err != nil {
					return fmt.Errorf("insert chat message failed: %w", err)
				}
			}
		}
	}

	onEvent("complete", finalContent, systemChatID, toolHandler.GetTaskID())
	return nil
}

// createChatModel 创建 ChatModel（OpenAI 兼容协议对标 eino-learn：裸 BaseURL/APIKey/Model 三字段）。
// Key 来源：当前用户的模型配置（AES-256-GCM 密文入库此处解密到内存）。
// 未配置返回哨兵错上层转 400 引导去设置页。
func createChatModel(ctx context.Context) (model.ToolCallingChatModel, *modelconfig.ModelCredential, error) {
	cred, err := modelconfig.Resolve(ctx)
	if err != nil {
		return nil, nil, err
	}
	keyPrefix := cred.APIKey
	if len(keyPrefix) > 12 {
		keyPrefix = keyPrefix[:12]
	}
	g.Log().Info(ctx, "[createChatModel] baseURL=%s model=%s apiKey=%s", cred.BaseURL, cred.Model, keyPrefix)
	chatModel, err := openaiModel.NewChatModel(ctx, &openaiModel.ChatModelConfig{
		BaseURL: cred.BaseURL,
		APIKey:  cred.APIKey,
		Model:   cred.Model,
	})
	if err != nil {
		return nil, nil, err
	}
	return chatModel, cred, nil
}

func getSystemPrompt(custom string) string {
	ragRules := "知识型问题先看系统提示末尾的 <retrieved> 块：有相关内容就引用作答并标注来源（标题/链接），不要再调 search_knowledge；<retrieved> 为空或不相关则按需调 search_knowledge（libraryId 按工具描述中的库列表填，ID 照抄不编造，不确定就留空）；闲聊/技术讨论直接回答，禁止为凑检索而调工具；多轮对话中用户说\"刚才那个/改成X/还要\"时，指代上文最近一次提到的任务/知识库/时间，没有才问；除最后一条外的 user 消息全是历史，只执行最后一条的指令，历史动作不得重复执行；问名字/身份/刚才说过的话只看历史，禁止调 search_knowledge。"
	if trimmed := strings.TrimSpace(custom); trimmed != "" {
		return trimmed + "\n\n" + ragRules + "\n\n" + dateContextBlock()
	}
	return defaultSystemPrompt() + dateContextBlock()
}

// dateContextBlock 日期时间上下文块：拼在 custom 之后、内置模板之内共用。
func dateContextBlock() string {
	now := time.Now()
	currentDate := now.Format("2006-01-02")
	year := now.Year()

	// 下个月15号 (15th of next month)
	nextMonth := now.AddDate(0, 1, 0)
	nextMonthDay15 := fmt.Sprintf("%d-%02d-15", nextMonth.Year(), nextMonth.Month())

	return `**时间转换（北京时间 UTC+8）：**
- 日期格式：YYYY-MM-DDTHH:MM:SS+08:00
- 今天是 ` + currentDate + `

**日期理解规则：**
- "今天" = 今天是 ` + currentDate + `
- "明天" = 今天 + 1天
- "本周X" = 本周的星期X（周一~周日），如果今天已过本周X，则指下周X
- "下周X" = 下周的星期X
- "X月X日" 或 "X号" = 今年的该日期（如已过则为明年），如"5月20日" = ` + fmt.Sprintf("%d-05-20", year) + `
- "下个月X号" = 下个月的第X天，如"下个月15号" = ` + nextMonthDay15 + `
- 用户说"25号"、"20日"等，需结合上下文判断月份（本月/下月/特定月份）

**重要：**
- 提取日期时，必须根据上述规则准确计算
- 只说"X号"时，优先理解为当前月，如果当天已过则理解为下个月
- 日期时间默认设置为当天 09:00:00+08:00

`
}

// defaultSystemPrompt 内置工具说明模板（custom 为空时使用）。
func defaultSystemPrompt() string {
	return `你是一个任务管理助手。可用工具及参数由系统自动提供（每个工具的名称、用途、参数说明以工具定义为准），按需调用。

**判断规则（关键）：**
- 当用户明确表达任务或知识库相关意图时，调用相应工具
- 当用户询问技术问题、讨论概念、聊天交流时，**直接回答问题，不要调用任何工具**
- 例如："今天天气怎么样" → 直接回答；"添加任务明天开会" → 调用创建任务的工具；"删除刚才创建的任务" → 调用取消任务的工具；"把这个网址加入知识库" → 调用添加知识库的工具
- 知识型问题（"是什么/怎么用/查一下X"且用户库可能相关）→ 先看系统提示末尾的 <retrieved> 块：有相关内容就引用作答并标注来源（标题/链接），不要再调 search_knowledge；<retrieved> 为空或明确不相关 → 按需调 search_knowledge（libraryId 按工具描述中的库列表填，ID 照抄不编造；不确定查哪个库就留空查全库）
- 闲聊/技术讨论/观点问题 → 直接回答，禁止为凑检索而调工具
- 多轮对话中用户说"刚才那个/改成X/还要"时，指代上文最近一次提到的任务/知识库/时间，没有才问
- 输入中除最后一条外的 user 消息全是历史，只供指代查阅；只执行最后一条 user 消息的指令，历史中的任何动作（创建/修改/删除/总结）一律不得重复执行
- 问名字/身份/刚才说过的话：答案只在对话历史里找，禁止调 search_knowledge（向量库里没有个人信息，调了就是污染）；历史里没有才说不知道
- 所有输出（含思考过程、引导语、回复）一律使用中文；用户用英文提问也用中文回答除非用户明确要求英文。

**任务新增/修改/删除人机验证流程（关键）：**
- 用户表达改任务意图（改标题/时间/优先级/类型/状态等）时，先调 query_tasks：只用标题/内容关键词填 keyword 定位任务；startDate/endDate 仅当用户明确说任务的当前时间范围（如"下周的任务""今天的任务"）才填，用户说的目标新时间（如"改到11点"）是改动值不是查询条件，禁止填入；收到结果后只输出一句引导语（候选卡由前端渲染，不要复述列表）。
- 用户要总结/日报/周报/月报时调 summarize_tasks，所有能推断的维度都填上作 AI 建议（period 从话术取 今天/昨日→day 本周→week 本月→month 明确日期→起止；taskType/priority 有则填无则留空不准编造）；系统一定先弹确认卡（用户可改），AI 禁止假设维度已确认；用户发 确认生成总结：{...} 时按 JSON 原样重调 summarize_tasks（本轮真执行）；取消生成总结 则不调工具、正常回话；执行成功后按返回统计写一段中文自然语言总结（数字以工具返回为准不准编），envelope 由系统自动拼接、不要手写。
- 用户只想查看任务（如"查找明天的任务""下周有什么任务""查询任务"）时同样调 query_tasks但只展示结果：收到结果后用自然语言列出任务（标题+时间），不要输出任何 envelope、不调其他工具、不引导确认。
- 用户表达新增意图时直接调 create_task（title/description/dueDate/priority/taskType 按理解填全）；系统会自动拦截转预览卡由用户确认，不会直接入库；用户发 确认创建任务：{...} 时再次调 create_task，系统此次才会真正入库。
- 删除流程：用户表达删任务意图时同样先调 query_tasks 展示候选项；用户回复要删的任务后，输出 task-delete-confirm envelope（含该任务 id/title/dueDate）加一句二次确认引导语，不调工具；用户发 确认删除任务：{"id":"..."} 时调 cancel_task（只传 taskId，禁止传 chatId 或空参走兜底）；用户发 取消删除任务/直接聊别的则不调工具。
- envelope 格式（fence 名精确）：候选 task-candidates 由工具结果自动拼接；新增预览 task-edit 由系统拦截自动拼接；修改结果 task-updated 自动拼接；二次确认 task-delete-confirm 含 id/title/dueDate；删除结果 task-deleted 含 id/title（调完 cancel_task 后输出）。

`
}

func createAgent(ctx context.Context, chatModel model.ToolCallingChatModel, systemPrompt, ragCtx string) (*react.Agent, *toolResultHandler, error) {
	createTaskTool := tools.CreateTask()
	cancelTaskTool := tools.CancelTask()
	addKnowledgeTool := tools.AddKnowledgeTool()
	searchKnowledgeTool := tools.SearchKnowledgeTool()
	queryTasksTool := tools.QueryTasks()
	updateTaskTool := tools.UpdateTaskTool()
	summarizeTasksTool := tools.SummarizeTasks()
	toolHandler := &toolResultHandler{}

	fullPrompt := systemPrompt + "\n\n<retrieved>\n" + ragCtx + "\n</retrieved>"
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               []tool.BaseTool{createTaskTool, cancelTaskTool, addKnowledgeTool, searchKnowledgeTool, queryTasksTool, updateTaskTool, summarizeTasksTool},
			ExecuteSequentially: false,
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage(fullPrompt))
			res = append(res, input...)
			return res
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return agent, toolHandler, nil
}

// handleToolCalls 处理流式 tool_calls
func handleToolCalls(ctx context.Context, msg *schema.Message, toolCallBuffer map[int]*schema.ToolCall, userId string, toolHandler *toolResultHandler) {
	g.Log().Info(ctx, "[handleToolCalls] 收到 tool_calls, 数量: %d, userId: %s", len(msg.ToolCalls), userId)
	for _, tc := range msg.ToolCalls {
		if tc.Index != nil {
			idx := int(*tc.Index)
			if existing, ok := toolCallBuffer[idx]; ok {
				existing.Function.Arguments += tc.Function.Arguments
			} else {
				tcCopy := tc
				toolCallBuffer[idx] = &tcCopy
			}

			if strings.HasSuffix(tc.Function.Arguments, "}") {
				completeTC := toolCallBuffer[idx]
				executeToolCall(ctx, completeTC, userId, toolHandler)
				delete(toolCallBuffer, idx)
			}
		}
	}
}

// executeToolCall 执行工具调用
func executeToolCall(ctx context.Context, tc *schema.ToolCall, userId string, toolHandler *toolResultHandler) {
	argsJSON := tc.Function.Arguments
	ctxWithUserId := context.WithValue(ctx, "user_id", userId)

	g.Log().Info(ctx, "[executeToolCall] 收到工具调用: name=%s, args=%s", tc.Function.Name, argsJSON)

	if tc.Function.Name == "create_task" {
		if !taskEditConfirmed(ctx) {
			// 人机验证：首轮只预览不入库。把模型参数包成 task-edit envelope 给前端渲染确认卡；
			// 用户发 确认创建任务 后下一轮 agent 再调本工具真建（见 taskEditConfirmed 判定）。
			g.Log().Info(ctx, "[create_task] 预览拦截, args=%s", argsJSON)
			toolHandler.mu.Lock()
			toolHandler.taskEditPreview = argsJSON
			toolHandler.mu.Unlock()
			return
		}
		g.Log().Info(ctx, "[create_task] 确认后真建")
		createTaskTool := tools.CreateTask()
		result, toolErr := createTaskTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "create_task failed: %v", toolErr)
			return
		}
		g.Log().Info(ctx, "[create_task] 执行成功, result=%s", result)

		var taskIdStr string
		json.Unmarshal([]byte(result), &taskIdStr)
		if id, err := uuid.Parse(taskIdStr); err == nil {
			toolHandler.mu.Lock()
			toolHandler.createdTaskID = id
			toolHandler.mu.Unlock()
			g.Log().Info(ctx, "[create_task] taskID=%s", id)
		}
	} else if tc.Function.Name == "cancel_task" {
		cancelTaskTool := tools.CancelTask()
		_, toolErr := cancelTaskTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "cancel_task failed: %v", toolErr)
		}
	} else if tc.Function.Name == "add_knowledge" {
		addKnowledgeTool := tools.AddKnowledgeTool()
		result, toolErr := addKnowledgeTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "add_knowledge failed: %v", toolErr)
			toolHandler.mu.Lock()
			toolHandler.knowledgeError = toolErr.Error()
			toolHandler.mu.Unlock()
		} else {
			toolHandler.mu.Lock()
			toolHandler.knowledgeResult = result
			toolHandler.mu.Unlock()
		}
	} else if tc.Function.Name == "search_knowledge" {
		g.Log().Info(ctx, "[search_knowledge] 调用工具, 参数: %s, userId: %s", argsJSON, userId)
		searchTool := tools.SearchKnowledgeTool()
		result, toolErr := searchTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "[search_knowledge] 执行失败: %v", toolErr)
		} else {
			g.Log().Info(ctx, "[search_knowledge] 执行成功, 结果: %s", result)
			toolHandler.mu.Lock()
			toolHandler.knowledgeResult = result
			toolHandler.mu.Unlock()
		}
	} else if tc.Function.Name == "query_tasks" {
		g.Log().Info(ctx, "[query_tasks] 调用工具, 参数: %s, userId: %s", argsJSON, userId)
		queryTool := tools.QueryTasks()
		result, toolErr := queryTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "[query_tasks] 执行失败: %v", toolErr)
		} else {
			g.Log().Info(ctx, "[query_tasks] 执行成功, 结果: %s", result)
			// 纯查看（无改/删意图）不拼候选卡：列表转自然语言文本，由模型直接答。
			// 改前定位（原文含改/编辑/修改/删除/删）才拼 task-candidates 卡走人机确认。
			if isTaskLocateIntent(ctx) {
				toolHandler.mu.Lock()
				toolHandler.taskCandidates = result
				toolHandler.mu.Unlock()
			} else {
				toolHandler.mu.Lock()
				toolHandler.taskListText = formatTaskListText(result)
				toolHandler.mu.Unlock()
			}
		}
	} else if tc.Function.Name == "update_task" {
		g.Log().Info(ctx, "[update_task] 调用工具, 参数: %s, userId: %s", argsJSON, userId)
		updateTool := tools.UpdateTaskTool()
		result, toolErr := updateTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "[update_task] 执行失败: %v", toolErr)
		} else {
			g.Log().Info(ctx, "[update_task] 执行成功, 结果: %s", result)
			toolHandler.mu.Lock()
			toolHandler.taskUpdateResult = result
			toolHandler.mu.Unlock()
		}
	} else if tc.Function.Name == "summarize_tasks" {
		g.Log().Info(ctx, "[summarize_tasks] 调用工具, 参数: %s, userId: %s", argsJSON, userId)
		// 确认门：首轮一律只出预览卡（无论 AI 建议齐不齐）；确认文本回环才真执行。
		if !taskSummaryConfirmed(ctx) {
			var preview map[string]any
			if err := json.Unmarshal([]byte(argsJSON), &preview); err != nil {
				g.Log().Error(ctx, "[summarize_tasks] 预览参数解析失败: %v", err)
				return
			}
			normalized := normalizeSummaryArgs(preview)
			previewJSON, _ := json.Marshal(normalized)
			toolHandler.mu.Lock()
			toolHandler.taskSummaryPreview = string(previewJSON)
			toolHandler.mu.Unlock()
			return
		}
		summaryTool := tools.SummarizeTasks()
		result, toolErr := summaryTool.InvokableRun(ctxWithUserId, argsJSON)
		if toolErr != nil {
			g.Log().Error(ctx, "[summarize_tasks] 执行失败: %v", toolErr)
		} else {
			g.Log().Info(ctx, "[summarize_tasks] 执行成功")
			toolHandler.mu.Lock()
			toolHandler.taskSummaryResult = result
			toolHandler.mu.Unlock()
		}
	}
}

// saveChatMessage 保存聊天消息到数据库
func saveChatMessage(ctx context.Context, id, userId uuid.UUID, content, role string, relatedTaskId uuid.UUID) error {
	chatMessage := &entity.DbChatMessages{
		Id:            id,
		UserId:        userId,
		Content:       content,
		Role:          role,
		RelatedTaskId: relatedTaskId,
	}
	return dao.DbChatMessages.InsertChatMessage(ctx, chatMessage)
}

type toolResultHandler struct {
	callbacks.HandlerBuilder
	createdTaskID      uuid.UUID
	knowledgeResult    string
	knowledgeError     string
	taskCandidates     string
	taskListText       string
	taskUpdateResult   string
	taskEditPreview    string
	taskSummaryPreview string
	taskSummaryResult  string
	mu                 sync.RWMutex
}

func (h *toolResultHandler) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	return ctx
}

func (h *toolResultHandler) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	if info.Name == "create_task" {
		if toolOutput, ok := output.(*tool.CallbackOutput); ok {
			var taskIDStr string
			if err := json.Unmarshal([]byte(toolOutput.Response), &taskIDStr); err == nil {
				if id, err := uuid.Parse(taskIDStr); err == nil {
					h.mu.Lock()
					h.createdTaskID = id
					h.mu.Unlock()
				}
			}
		}
	}
	return ctx
}

func (h *toolResultHandler) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	g.Log().Error(ctx, "[AgentError] Name=%s, Component=%s, Error=%v", info.Name, info.Component, err)
	return ctx
}

func (h *toolResultHandler) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}

func (h *toolResultHandler) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	defer output.Close()
	return ctx
}

func (h *toolResultHandler) GetTaskID() uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.createdTaskID
}

func (h *toolResultHandler) IsTaskCreated() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.createdTaskID != uuid.Nil
}

func (h *toolResultHandler) GetKnowledgeResult() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.knowledgeResult
}

func (h *toolResultHandler) GetKnowledgeError() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.knowledgeError
}

func (h *toolResultHandler) GetTaskCandidates() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskCandidates
}

func (h *toolResultHandler) GetTaskListText() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskListText
}

// isCreateIntent 当前这条消息本身就是创建意图（首轮直建或确认回环）：
// 此时历史里即使有一模一样的创建记录，也不算"重复执行"——用户明确要求再建一次。
// 供 system 提示词引用，压住"历史动作不得重复执行"对当前指令的误伤。
// isTaskLocateIntent 原文是否含改/删意图：有则 query_tasks 结果走候选卡确认，
// 无则走纯文本列表直接答。确认回环文本（确认修改任务/我选删除任务/确认删除任务…）
// 自带改删字样，必须走定位分支，否则确认后反而列成文本不执行。
func isTaskLocateIntent(ctx context.Context) bool {
	content, _ := ctx.Value("user_content").(string)
	for _, kw := range []string{"改", "编辑", "修改", "删除", "删", "取消"} {
		if strings.Contains(content, kw) {
			return true
		}
	}
	return false
}

// isCreateIntent 当前消息是否创建意图：以 /添加任务 开头，或含"添加任务/新建任务/创建任务"
// 且非"确认创建任务"回环（回环走 handleConfirmText 硬执行，不进模型）。
func isCreateIntent(content string) bool {
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "确认创建任务") {
		return false
	}
	if strings.HasPrefix(trimmed, "/添加任务") {
		return true
	}
	for _, kw := range []string{"添加任务", "新建任务", "创建任务", "新增任务"} {
		if strings.Contains(trimmed, kw) {
			return true
		}
	}
	return false
}

// formatTaskListText query_tasks 结果转自然语言列表（纯查看场景，不拼卡）。
func formatTaskListText(result string) string {
	var parsed struct {
		Tasks []struct {
			Title    string `json:"title"`
			DueDate  string `json:"dueDate"`
			Status   string `json:"status"`
			Priority string `json:"priority"`
		} `json:"tasks"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil || len(parsed.Tasks) == 0 {
		return "当前没有符合条件的任务。"
	}
	var sb strings.Builder
	sb.WriteString("任务列表：\n")
	for i, t := range parsed.Tasks {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, t.Title))
		if t.DueDate != "" {
			sb.WriteString("（" + t.DueDate + "）")
		}
		if t.Status != "" {
			sb.WriteString("[" + t.Status + "]")
		}
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (h *toolResultHandler) GetTaskUpdateResult() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskUpdateResult
}

func (h *toolResultHandler) GetTaskEditPreview() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskEditPreview
}

func (h *toolResultHandler) GetTaskSummaryPreview() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskSummaryPreview
}

func (h *toolResultHandler) GetTaskSummaryResult() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.taskSummaryResult
}

// retrieveForChat 前置检索：用户全库兜底查一轮，取 context 拼进 system。
// 无库/失败返回空串，不中断对话（检索是增强不是门禁）。
// 历史回放类问题（问刚才/上文/聊天记录）不检索——答案在对话历史里，向量库查它浪费 18s 还污染 <retrieved>。
func retrieveForChat(ctx context.Context, userID, content string) string {
	if isHistoryRecall(content) {
		g.Log().Info(ctx, "[retrieveForChat] 跳过：历史回放类问题")
		return ""
	}
	if isChitchat(content) {
		g.Log().Info(ctx, "[retrieveForChat] 跳过：闲聊无检索信号")
		return ""
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return ""
	}
	bases, err := dao.DbKnowledgeBases.ListByUserID(ctx, uid)
	if err != nil || len(bases) == 0 {
		return ""
	}
	args, _ := json.Marshal(map[string]string{"query": content})
	result, err := tools.SearchKnowledgeTool().InvokableRun(ctx, string(args))
	if err != nil {
		g.Log().Warning(ctx, "[retrieveForChat] 检索失败，降级为无资料: %v", err)
		return ""
	}
	var parsed struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		g.Log().Warning(ctx, "[retrieveForChat] 检索结果解析失败: %v", err)
		return ""
	}
	if len(parsed.Context) > 4000 {
		parsed.Context = parsed.Context[:4000]
	}
	return parsed.Context
}

// libraryCatalog 拼用户知识库目录进 system：模型按名选库时 ID 照抄。
// 无库返回空串；失败记日志不阻断对话。dao 直调（service 层互引会循环）。
func libraryCatalog(ctx context.Context, userID string) string {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return ""
	}
	bases, err := dao.DbKnowledgeBases.ListByUserID(ctx, uid)
	if err != nil {
		g.Log().Warning(ctx, "[libraryCatalog] 查库列表失败: %v", err)
		return ""
	}
	if len(bases) == 0 {
		return "\n\n知识库：暂无（search_knowledge 无资料可查，知识型问题直接回答）"
	}
	var sb strings.Builder
	sb.WriteString("\n\n知识库：")
	for i, b := range bases {
		if i > 0 {
			sb.WriteString("；")
		}
		sb.WriteString(b.Name + "(" + b.Id.String() + ")")
	}
	sb.WriteString("；调 search_knowledge 时 libraryId 从上列 ID 照抄，不确定查哪个库就留空查全库")
	return sb.String()
}
