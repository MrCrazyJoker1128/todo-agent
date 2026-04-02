package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

type GetChatListReq struct {
	g.Meta    `path:"/list" method:"get" tags:"Chat" summary:"获取聊天列表"`
	PageIndex int `json:"pageIndex" description:"页码"`
	PageSize  int `json:"pageSize" description:"每页数量"`
}

type ChatMessageWithTask struct {
	ChatId        uuid.UUID   `json:"chatId"`
	Role          string      `json:"role"`
	ChatContent   string      `json:"chatContent"`
	RelatedTaskId *uuid.UUID  `json:"relatedTaskId"`
	ChatCreatedAt *gtime.Time `json:"chatCreatedAt"`
	TaskId        *uuid.UUID  `json:"taskId"`
	Title         *string     `json:"title"`
	Description   *string     `json:"description"`
	Status        *string     `json:"status"`
	Priority      *string     `json:"priority"`
	DueDate       *gtime.Time `json:"dueDate"`
	TaskType      *string     `json:"taskType"`
	PriorityLevel *string     `json:"priorityLevel"`
	TaskCreatedAt *gtime.Time `json:"taskCreatedAt"`
	TaskUpdatedAt *gtime.Time `json:"taskUpdatedAt"`
}

type GetChatListRes struct {
	List  []*ChatMessageWithTask `json:"list"`
	Total int                    `json:"total"`
}

type InsertChatMessageReq struct {
	g.Meta        `path:"/insert" method:"post" tags:"Chat" summary:"插入聊天消息"`
	Content       string `json:"content" description:"消息内容"`
	Role          string `json:"role" description:"消息角色"`
	RelatedTaskId string `json:"relatedTaskId" description:"关联任务ID"`
}

type InsertChatMessageRes struct {
	Success bool `json:"success"`
}

type SendChatMessageReq struct {
	g.Meta  `path:"/send" method:"post" tags:"Chat" summary:"发送聊天消息"`
	Content string   `json:"content" description:"消息内容"`
	Role    string   `json:"role" description:"消息角色"`
	Actions []string `json:"actions" description:"消息操作"`
}

type SSEMessage struct {
	Event   string    `json:"event" description:"事件类型 例如：收到请求 received, 等待响应 wait_response,发送响应 send_response,完成响应 complete"`
	Message string    `json:"message" description:"sse 消息内容"`
	ChatId  uuid.UUID `json:"chatId" description:"聊天消息ID"`
	TaskId  uuid.UUID `json:"taskId" description:"创建的任务ID，仅在 complete 事件时有值"`
}

type SendChatMessageRes struct {
	Success bool `json:"success"`
}
