package chat

import (
	"context"
	"encoding/json"
	"errors"
	"service/api"
	chatservice "service/internal/service/chat"
	"service/internal/service/modelconfig"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"
)

type ChatController struct {
}

func NewChatController() *ChatController {
	return &ChatController{}
}

// ==================== SSE 薄层 ====================

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

// friendlyModelError 把模型调用错误翻译成用户可读文案进对话流。日志保留原文此处只做展示。
func friendlyModelError(err error) string {
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "certificate has expired") || strings.Contains(msg, "certificate is not yet valid") || strings.Contains(msg, "tls:"):
		return "模型服务连接失败（对方证书过期或不可信）请检查 BaseURL 是否正确或稍后重试"
	case strings.Contains(msg, "connection refused"):
		return "模型服务连接被拒绝请检查 BaseURL 与网络"
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded"):
		return "模型服务响应超时请稍后重试"
	case strings.Contains(msg, "429") || strings.Contains(msg, "rate limit"):
		return "模型服务限流请稍后重试"
	default:
		return "模型调用失败请稍后重试（详情见服务端日志）"
	}
}

func (c *ChatController) GetChatList(ctx context.Context, req *api.GetChatListReq) (res *api.GetChatListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	if userId == "" {
		g.Log().Error(ctx, "没有找到用户，请先登录")
		return nil, errors.New("没有找到用户，请先登录")
	}
	return chatservice.GetChatList(ctx, userId, req)
}

func (c *ChatController) InsertChatMessage(ctx context.Context, req *api.InsertChatMessageReq) (res *api.InsertChatMessageRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	if userId == "" {
		g.Log().Error(ctx, "没有找到用户，请先登录")
		return nil, errors.New("没有找到用户，请先登录")
	}
	return chatservice.InsertChatMessage(ctx, userId, req)
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

	err = chatservice.RunAgentStream(ctx, userId, req.Content, func(event, message string, chatID, taskID uuid.UUID) {
		_ = sendSSEEvent(r, event, chatID, message, taskID)
	})
	if err != nil {
		g.Log().Error(ctx, "send chat message error: %v", err)
		if errors.Is(err, modelconfig.ErrNotConfigured) {
			// 前端据此弹“去设置页配置 Key”引导。
			_ = sendSSEEvent(r, "error:NEED_MODEL_CONFIG", uuid.New(), modelconfig.ErrNotConfigured.Error(), uuid.Nil)
			sendErrorResponse(r, "NEED_MODEL_CONFIG")
			return nil, err
		}
		if errors.Is(err, chatservice.ErrModelAuth) {
			// 前端据此弹“Key 无效去设置页更换”引导。
			_ = sendSSEEvent(r, "error:MODEL_AUTH", uuid.New(), chatservice.ErrModelAuth.Error(), uuid.Nil)
			sendErrorResponse(r, "MODEL_AUTH")
			return nil, err
		}
		// 通用错误同样先发 SSE 事件进对话流，前端渲染为助手错误气泡再按 HTTP 语义收尾。
		_ = sendSSEEvent(r, "error:MODEL_CALL", uuid.New(), friendlyModelError(err), uuid.Nil)
		sendErrorResponse(r, err.Error())
		return nil, err
	}
	return nil, nil
}

func (c *ChatController) EchoWebSocket(r *ghttp.Request) {
	ctx := r.Context()
	g.Log().Info(ctx, "[EchoWebSocket] WebSocket 连接建立")

	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error(ctx, "[EchoWebSocket] WebSocket 升级失败: %v", err)
		r.Response.WriteStatus(500)
		return
	}
	defer ws.Close()

	for {
		msgType, msgBytes, err := ws.ReadMessage()
		if err != nil {
			if err.Error() == "EOF" || err.Error() == "aborted" {
				g.Log().Info(ctx, "[EchoWebSocket] 客户端关闭连接")
			} else {
				g.Log().Error(ctx, "[EchoWebSocket] 读取消息失败: %v", err)
			}
			break
		}

		content := string(msgBytes)
		g.Log().Info(ctx, "[EchoWebSocket] 收到消息: %s", content)

		response := "我收到了" + content
		if err := ws.WriteMessage(msgType, []byte(response)); err != nil {
			g.Log().Error(ctx, "[EchoWebSocket] 发送消息失败: %v", err)
			break
		}
		g.Log().Info(ctx, "[EchoWebSocket] 发送响应: %s", response)
	}
}
