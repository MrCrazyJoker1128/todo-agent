package chat

import (
	"context"
	"regexp"
	"strings"

	"service/internal/dao"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

// taskFenceRe 兜掉落库文本里的任务 envelope fence（```task-* … ```），
// 那是给前端渲染卡片的，拼进模型输入只会让它从 JSON 里瞎捡意图。
var taskFenceRe = regexp.MustCompile("(?s)```task-[a-z-]+\\n.*?```")

// thinkRe 兜掉 <think> 推理块：那是上轮的思考过程，不是结论，重读会把猜测当事实。
var thinkRe = regexp.MustCompile("(?s)<think>.*?(?:</think>|$)")

// confirmRe 确认指令原文：用户点卡片发出的"确认创建任务：{...}"等，
// 是一次性动作指令，重放给模型会被当成"用户又确认了一次"直接执行。
var confirmRe = regexp.MustCompile(`^(确认创建任务|确认修改任务|我选删除任务|确认删除任务|取消任务修改|取消删除任务|确认生成总结|取消生成总结)[:：].*`)

// loadRecentHistory 读用户最近 N 条落库消息，转 agent 输入历史。
// dao 按 created_at asc，取尾部 N 条即最近；空 content 行跳过；失败返回空 slice 不中断对话。
func loadRecentHistory(ctx context.Context, userID string, n int) []*schema.Message {
	msgs, err := dao.DbChatMessages.GetLatestMessages(ctx, userID, n)
	if err != nil {
		g.Log().Warning(ctx, "[loadRecentHistory] 读历史失败，降级为无历史: %v", err)
		return nil
	}
	if len(msgs) == 0 {
		return nil
	}
	prev := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		text := strings.TrimSpace(taskFenceRe.ReplaceAllString(m.Content, ""))
		text = strings.TrimSpace(thinkRe.ReplaceAllString(text, ""))
		// 确认指令是动作，回放只留"做过什么"的事实：转成陈述句，模型不再当新指令执行。
		// 取消类是"没做"的意思，直接整条丢掉。
		if m.Role != "assistant" && confirmRe.MatchString(strings.TrimSpace(m.Content)) {
			raw := strings.TrimSpace(m.Content)
			if strings.HasPrefix(raw, "取消任务修改") || strings.HasPrefix(raw, "取消删除任务") || strings.HasPrefix(raw, "取消生成总结") {
				continue
			}
			text = "（历史动作：" + beautifyConfirmLine(raw) + "，已执行完毕，不要重复执行）"
		}
		if text == "" {
			continue
		}
		if m.Role == "assistant" {
			// assistant 历史只保留首句结论：复述/任务列表全砍，模型要用数据自己调工具查。
			// 尾部留任务细节的截断方向是错的——复述恰好长在尾部。
			text = firstSentence(text)
			if text == "" {
				continue
			}
			prev = append(prev, schema.AssistantMessage(text, nil))
		} else {
			prev = append(prev, schema.UserMessage(text))
		}
	}
	return prev
}

// firstSentence 取首句（按 。！？\n 切）。首句通常是结论，后面全是复述。
func firstSentence(s string) string {
	for i, r := range []rune(s) {
		if r == '。' || r == '！' || r == '？' || r == '\n' {
			head := strings.TrimSpace(string([]rune(s)[:i+1]))
			if head != "" {
				return head
			}
		}
	}
	return strings.TrimSpace(s)
}

// beautifyConfirmLine 确认指令转人话（后端版，与前端 beautifyConfirm 同语义，历史喂模型用）。
func beautifyConfirmLine(raw string) string {
	if i := strings.IndexAny(raw, ":："); i >= 0 && i+3 <= len(raw) {
		return strings.TrimSpace(raw[:i]) + "：" + strings.TrimSpace(raw[i+3:])
	}
	return raw
}

// isHistoryRecall 问的是对话历史本身（刚才/上文/聊天记录），答案在输入历史里，不该走向量检索。
func isHistoryRecall(content string) bool {
	for _, kw := range []string{"刚才", "刚刚", "上文", "之前问", "之前说", "历史", "聊天记录", "回顾一下", "名字", "我叫", "你叫"} {
		if strings.Contains(content, kw) {
			return true
		}
	}
	return false
}

// isChitchat 纯闲聊（自我介绍/打招呼/短句寒暄）：无任务词、无知识疑问词，不值得走向量检索。
// 阈值保守：宁可多检一次，不漏检知识问题。命中任务/疑问词的一律放行。
func isChitchat(content string) bool {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return true
	}
	// 有明确信号的放行：任务词、疑问词、知识动词、网址。
	for _, kw := range []string{
		"任务", "待办", "提醒", "日程", "会议", "开会",
		"什么", "怎么", "如何", "为何", "为什么", "?", "？",
		"查", "搜", "找", "总结", "日报", "周报", "月报",
		"http", "www", ".com", ".cn",
		"改", "编辑", "修改", "删除", "删", "取消",
	} {
		if strings.Contains(trimmed, kw) {
			return false
		}
	}
	// 无信号且短句（≤10 rune）：自我介绍/打招呼类，直接跳过。
	return len([]rune(trimmed)) <= 10
}

// cancelReply 取消类终止语短路：不进 agent，直接回一句。
// 这类话是"停下上一件事"，不是新意图；进模型只会误召回全量任务。
func cancelReply(content string) (string, bool) {
	trimmed := strings.TrimSpace(content)
	switch {
	case strings.HasPrefix(trimmed, "取消任务修改"):
		return "已取消修改，刚才的内容没有入库。", true
	case strings.HasPrefix(trimmed, "取消删除任务"):
		return "已取消删除，任务保留。", true
	case strings.HasPrefix(trimmed, "取消生成总结"):
		return "已取消生成总结。", true
	}
	return "", false
}

// FixToolPairs 修复断裂的工具调用对（照抄官方 eino-examples/eino-learn/internal/service/history/history.go）。
// 有 ToolCalls 的 assistant 消息必须有对应 ID 的 Tool 结果消息，否则 MiniMax 报 2013。
// 修法：收集所有 ToolMessage 的 tool_call_id，无结果的 ToolCalls 删掉；
// 删光且无 Content 的整条丢弃，有 Content 的留文本。
// 当前落库只有文本无 ToolCalls，本函数是 no-op；envelope 回环跨轮调工具时保命用。
func FixToolPairs(msgs []*schema.Message) []*schema.Message {
	hasResult := map[string]bool{}
	for _, m := range msgs {
		if m.Role == schema.Tool && m.ToolCallID != "" {
			hasResult[m.ToolCallID] = true
		}
	}
	var out []*schema.Message
	for _, m := range msgs {
		if len(m.ToolCalls) == 0 {
			out = append(out, m)
			continue
		}
		var kept []schema.ToolCall
		for _, tc := range m.ToolCalls {
			if hasResult[tc.ID] {
				kept = append(kept, tc)
			}
		}
		if len(kept) == 0 && m.Content == "" {
			continue
		}
		m.ToolCalls = kept
		out = append(out, m)
	}
	return out
}
