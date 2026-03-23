// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbChatMessages is the golang structure for table db_chat_messages.
type DbChatMessages struct {
	Id            uuid.UUID   `json:"id"            orm:"id"              description:"主键：消息唯一标识"`                       // 主键：消息唯一标识
	UserId        uuid.UUID   `json:"userId"        orm:"user_id"         description:"外键：消息归属的用户 ID"`                   // 外键：消息归属的用户 ID
	Role          string      `json:"role"          orm:"role"            description:"消息角色：区分是用户发送还是 AI 回复"`            // 消息角色：区分是用户发送还是 AI 回复
	Content       string      `json:"content"       orm:"content"         description:"消息内容：具体的文本对话内容"`                  // 消息内容：具体的文本对话内容
	RelatedTaskId uuid.UUID   `json:"relatedTaskId" orm:"related_task_id" description:"关联任务 ID：当对话涉及创建/修改任务时，关联到此任务"`    // 关联任务 ID：当对话涉及创建/修改任务时，关联到此任务
	Embedding     string      `json:"embedding"     orm:"embedding"       description:"向量数据：存储消息内容的 Embedding，用于长期记忆检索"` // 向量数据：存储消息内容的 Embedding，用于长期记忆检索
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"消息发送时间"`                          // 消息发送时间
}
