// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbChatMessages is the golang structure of table db_chat_messages for DAO operations like Where/Data.
type DbChatMessages struct {
	g.Meta        `orm:"table:db_chat_messages, do:true"`
	Id            any         // 主键：消息唯一标识 (UUID)
	UserId        any         // 外键：归属用户 ID
	Role          any         // 消息角色：user|assistant
	Content       any         // 消息内容：文本
	RelatedTaskId any         // 关联任务 ID：如果消息涉及任务
	CreatedAt     *gtime.Time // 消息创建时间
}
