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
	"service/internal/tools"
	"strings"
	"sync"
	"time"

	arkModel "github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	einoAgent "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"
)

type ChatController struct {
}

func NewChatController() *ChatController {
	return &ChatController{}
}

// ==================== 辅助方法 ====================

// setupSSEHeaders 设置 SSE 响应头
func setupSSEHeaders(r *ghttp.Request) {
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
}

// sendErrorResponse 发送错误响应
func sendErrorResponse(r *ghttp.Request, errMsg string) {
	r.Response.WriteStatus(500)
	r.Response.Write(errMsg)
	r.Response.Flush()
}

// createChatModel 创建 ChatModel
func createChatModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	timeout := 60 * time.Second
	return arkModel.NewChatModel(ctx, &arkModel.ChatModelConfig{
		BaseURL: "https://api.minimaxi.com/v1",
		APIKey:  "sk-cp-tSDxyBrROyLiq2RwCat12rqHbrRnQUxeK7kLdY4g3jXstA8TG-L8vXkGv5aroe9O7pTf01OS5O_L2aHpmRRHg9ZRgWHGCbUavS7l9bN2mWDsruhxpejvyR0",
		Model:   "MiniMax-M2.7",
		Timeout: &timeout,
	})
}

// getSystemPrompt 获取系统提示词
func getSystemPrompt() string {
	now := time.Now()
	currentDate := now.Format("2006-01-02")
	year := now.Year()

	// 下个月15号 (15th of next month)
	nextMonth := now.AddDate(0, 1, 0)
	nextMonthDay15 := fmt.Sprintf("%d-%02d-15", nextMonth.Year(), nextMonth.Month())

	return `你是一个任务管理助手。你有以下工具可以使用：

1. create_task - 创建新任务
   - 当用户说"添加任务"、"新增待办"、"创建任务"、"安排会议"、"提醒我"、"帮我记一下"时调用
   - 参数：title(必填), description, due_date(ISO 8601), priority(HIGH/MEDIUM/LOW), task_type(work/personal/project/adhoc)

2. cancel_task - 取消/删除任务
   - 当用户说"删除任务"、"取消任务"、"移除任务"时调用
   - 参数（按优先级）：taskId > chatId > 不带参数（取消最近任务）
   - taskId：直接指定任务ID
   - chatId：通过会话ID查找该会话关联的任务并取消
   - 不带参数：自动取消用户最近创建的任务

3. add_knowledge - 添加知识库
   - 当用户说"添加知识库"、"保存到知识库"、"记一个知识"、"把这个网页添加到知识库"时调用
   - 参数：url(必填，要添加的网页URL)
   -会自动提取标题、作者、内容，并保存到知识库

4. search_knowledge - 搜索知识库
   - 当用户询问问题，需要从已添加的知识库中查找相关信息时调用
   - 参数：query(必填，用户的问题)
   - 返回与问题相关的知识库片段（最多5条）

**判断规则（关键）：**
- 当用户明确表达任务或知识库相关意图时，调用相应工具
- 当用户询问技术问题、讨论概念、聊天交流时，**直接回答问题，不要调用任何工具**
- 例如："今天天气怎么样" → 直接回答；"添加任务明天开会" → 调用 create_task；"删除刚才创建的任务" → 调用 cancel_task；"把这个网址加入知识库" → 调用 add_knowledge

**时间转换（北京时间 UTC+8）：**
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

`}
func createAgent(ctx context.Context, chatModel model.ToolCallingChatModel, systemPrompt string) (*react.Agent, *toolResultHandler, error) {
	createTaskTool := tools.CreateTask()
	cancelTaskTool := tools.CancelTask()
	addKnowledgeTool := tools.AddKnowledgeTool()
	searchKnowledgeTool := tools.SearchKnowledgeTool()
	toolHandler := &toolResultHandler{}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               []tool.BaseTool{createTaskTool, cancelTaskTool, addKnowledgeTool, searchKnowledgeTool},
			ExecuteSequentially: false,
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage(systemPrompt))
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
		g.Log().Info(ctx, "[create_task] 开始执行")
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
	}
}

// sendSSEEvent 发送 SSE 事件
func sendSSEEvent(r *ghttp.Request, event string, chatId uuid.UUID, message string, taskId uuid.UUID) error {
	responseData := &api.SSEMessage{
		ChatId:  chatId,
		Event:   event,
		Message: message,
		TaskId:  taskId,
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		g.Log().Error(r.Context(), "marshal response data error: %v", err)
		sendErrorResponse(r, "marshal response data error")
		return err
	}
	r.Response.Write([]byte("data: " + string(jsonData) + "\n\n"))
	r.Response.Flush()
	return nil
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
	knowledgeResult  string
	knowledgeError  string
	mu            sync.RWMutex
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
func (c *ChatController) GetChatList(ctx context.Context, req *api.GetChatListReq) (res *api.GetChatListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	if userId == "" {
		g.Log().Error(ctx, "没有找到用户，请先登录")
		return nil, errors.New("没有找到用户，请先登录")
	}
	chatMessages, total, err := dao.DbChatMessages.GetChatMessagesWithTask(ctx, userId, req)
	if err != nil {
		g.Log().Error(ctx, "get chat messages by user id error: %v", err)
		return nil, err
	}
	return &api.GetChatListRes{
		List:  chatMessages,
		Total: total,
	}, nil
}

func (c *ChatController) InsertChatMessage(ctx context.Context, req *api.InsertChatMessageReq) (res *api.InsertChatMessageRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	if userId == "" {
		g.Log().Error(ctx, "没有找到用户，请先登录")
		return nil, errors.New("没有找到用户，请先登录")
	}
	parsedUserId := uuid.MustParse(userId)

	var relatedTaskId uuid.UUID
	if req.RelatedTaskId != "" {
		relatedTaskId = uuid.MustParse(req.RelatedTaskId)
	}

	chatMessage := &entity.DbChatMessages{
		UserId:        parsedUserId,
		Content:       req.Content,
		Role:          req.Role,
		RelatedTaskId: relatedTaskId,
		Id:            uuid.New(),
	}
	err = dao.DbChatMessages.InsertChatMessage(ctx, chatMessage)
	if err != nil {
		g.Log().Error(ctx, "insert chat message error: %v", err)
		return nil, err
	}
	return &api.InsertChatMessageRes{
		Success: true,
	}, nil
}

func (c *ChatController) SendChatMessage(ctx context.Context, req *api.SendChatMessageReq) (res *api.SendChatMessageRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()

	setupSSEHeaders(r)

	if userId == "" {
		g.Log().Error(ctx, "没有找到用户，请先登录")
		sendErrorResponse(r, "没有找到用户，请先登录")
		return nil, errors.New("没有找到用户，请先登录")
	}

	parsedUserId := uuid.MustParse(userId)

	if err := saveChatMessage(ctx, uuid.New(), parsedUserId, req.Content, req.Role, uuid.Nil); err != nil {
		g.Log().Error(ctx, "insert chat message error: %v", err)
		sendErrorResponse(r, "insert chat message error")
		return nil, err
	}

	sendSSEEvent(r, "received", uuid.New(), req.Content, uuid.Nil)

	chatModel, err := createChatModel(ctx)
	if err != nil {
		g.Log().Error(ctx, "create chat model error: %v", err)
		sendErrorResponse(r, "create chat model error")
		return nil, err
	}

	agent, toolHandler, err := createAgent(ctx, chatModel, getSystemPrompt())
	if err != nil {
		g.Log().Error(ctx, "创建agent失败: %v", err)
		sendErrorResponse(r, "创建agent失败")
		return nil, err
	}

	ctx = context.WithValue(ctx, "user_id", userId)
	chatMsg := []*schema.Message{{Role: schema.User, Content: req.Content}}

	callback := react.BuildAgentCallback(nil, nil)
	sr, err := agent.Stream(ctx, chatMsg, einoAgent.WithComposeOptions(compose.WithCallbacks(callback), compose.WithCallbacks(toolHandler)))
	if err != nil {
		g.Log().Error(ctx, "Stream failed, err=%v", err)
		sendErrorResponse(r, "Stream failed")
		return nil, err
	}

	systemChatId := uuid.New()
	sendSSEEvent(r, "send_response", systemChatId, "", uuid.Nil)

	toolCallBuffer := make(map[int]*schema.ToolCall)
	finalContent := ""

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			g.Log().Error(ctx, "Recv failed, err=%v", err)
			sendErrorResponse(r, "Recv failed")
			return nil, err
		}

		handleToolCalls(ctx, msg, toolCallBuffer, userId, toolHandler)
		sendSSEEvent(r, "send_response", systemChatId, msg.Content, uuid.Nil)
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

	if err := saveChatMessage(ctx, systemChatId, parsedUserId, finalContent, "assistant", toolHandler.GetTaskID()); err != nil {
		g.Log().Error(ctx, "insert chat message failed: %v", err)
		sendErrorResponse(r, "insert chat message failed")
		return nil, err
	}

	sendSSEEvent(r, "complete", systemChatId, finalContent, toolHandler.GetTaskID())
	return nil, nil
}

type loggerCallback struct {
	callbacks.HandlerBuilder
}

func (cb *loggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	return ctx
}

func (cb *loggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	return ctx
}

func (cb *loggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	g.Log().Error(ctx, "[AgentError] %v", err)
	return ctx
}

func (cb *loggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	defer output.Close()
	return ctx
}

func (cb *loggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
