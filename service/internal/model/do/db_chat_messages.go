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
	Id            any         // 主键：消息唯一标识
	UserId        any         // 外键：消息归属的用户 ID
	Role          any         // 消息角色：区分是用户发送还是 AI 回复
	Content       any         // 消息内容：具体的文本对话内容
	RelatedTaskId any         // 关联任务 ID：当对话涉及创建/修改任务时，关联到此任务
	Embedding     any         // 向量数据：存储消息内容的 Embedding，用于长期记忆检索
	CreatedAt     *gtime.Time // 消息发送时间
}
