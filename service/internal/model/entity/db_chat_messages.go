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
	Id            uuid.UUID   `json:"id"            orm:"id"              description:"主键：消息唯一标识 (UUID)"`    // 主键：消息唯一标识 (UUID)
	UserId        uuid.UUID   `json:"userId"        orm:"user_id"         description:"外键：归属用户 ID"`          // 外键：归属用户 ID
	Role          string      `json:"role"          orm:"role"            description:"消息角色：user|assistant"` // 消息角色：user|assistant
	Content       string      `json:"content"       orm:"content"         description:"消息内容：文本"`             // 消息内容：文本
	RelatedTaskId uuid.UUID   `json:"relatedTaskId" orm:"related_task_id" description:"关联任务 ID：如果消息涉及任务"`    // 关联任务 ID：如果消息涉及任务
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"消息创建时间"`              // 消息创建时间
}
