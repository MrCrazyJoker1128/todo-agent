// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbTasks is the golang structure for table db_tasks.
type DbTasks struct {
	Id          uuid.UUID   `json:"id"          orm:"id"          description:"主键：任务唯一标识"`                                          // 主键：任务唯一标识
	UserId      uuid.UUID   `json:"userId"      orm:"user_id"     description:"外键：任务归属的用户 ID，级联删除"`                                 // 外键：任务归属的用户 ID，级联删除
	Title       string      `json:"title"       orm:"title"       description:"任务标题：简短概括任务内容"`                                      // 任务标题：简短概括任务内容
	Description string      `json:"description" orm:"description" description:"任务详细描述：可包含多行文本"`                                     // 任务详细描述：可包含多行文本
	Status      string      `json:"status"      orm:"status"      description:"任务状态：todo(待办), doing(进行中), done(完成), cancelled(取消)"` // 任务状态：todo(待办), doing(进行中), done(完成), cancelled(取消)
	Priority    string      `json:"priority"    orm:"priority"    description:"任务优先级：low(低), medium(中), high(高)"`                   // 任务优先级：low(低), medium(中), high(高)
	DueDate     *gtime.Time `json:"dueDate"     orm:"due_date"    description:"任务截止时间：用于排序和提醒"`                                     // 任务截止时间：用于排序和提醒
	ParentId    uuid.UUID   `json:"parentId"    orm:"parent_id"   description:"父任务 ID：实现任务拆解/子任务功能，NULL 表示根任务"`                     // 父任务 ID：实现任务拆解/子任务功能，NULL 表示根任务
	Meta        string      `json:"meta"        orm:"meta"        description:"扩展元数据：JSONB 格式，存储来源、标签、原始输入等非结构化数据"`                 // 扩展元数据：JSONB 格式，存储来源、标签、原始输入等非结构化数据
	Embedding   string      `json:"embedding"   orm:"embedding"   description:"向量数据：存储任务内容的 Embedding，用于语义搜索 (如\"找关于钱的任务\")"`       // 向量数据：存储任务内容的 Embedding，用于语义搜索 (如"找关于钱的任务")
	IsDeleted   bool        `json:"isDeleted"   orm:"is_deleted"  description:"软删除标志：逻辑删除，防止误删数据"`                                  // 软删除标志：逻辑删除，防止误删数据
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"任务创建时间"`                                             // 任务创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"任务最后更新时间"`                                           // 任务最后更新时间
}
