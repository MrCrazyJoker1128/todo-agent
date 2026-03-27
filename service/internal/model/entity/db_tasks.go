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
	Id             uuid.UUID   `json:"id"             orm:"id"               description:"主键：任务唯一标识 (UUID)"`                     // 主键：任务唯一标识 (UUID)
	UserId         uuid.UUID   `json:"userId"         orm:"user_id"          description:"外键：归属用户 ID"`                           // 外键：归属用户 ID
	Title          string      `json:"title"          orm:"title"            description:"任务标题：1-255 字符"`                        // 任务标题：1-255 字符
	Description    string      `json:"description"    orm:"description"      description:"任务详细描述：可选"`                            // 任务详细描述：可选
	Status         string      `json:"status"         orm:"status"           description:"任务状态：todo|doing|done|cancelled"`       // 任务状态：todo|doing|done|cancelled
	Priority       string      `json:"priority"       orm:"priority"         description:"优先级：low|medium|high"`                  // 优先级：low|medium|high
	DueDate        *gtime.Time `json:"dueDate"        orm:"due_date"         description:"截止时间：带时区"`                             // 截止时间：带时区
	MilvusVectorId uuid.UUID   `json:"milvusVectorId" orm:"milvus_vector_id" description:"Milvus 向量 ID (UUID 类型，与 Milvus 主键对应)"` // Milvus 向量 ID (UUID 类型，与 Milvus 主键对应)
	TaskType       string      `json:"taskType"       orm:"task_type"        description:"任务类型：work|personal|project|adhoc"`     // 任务类型：work|personal|project|adhoc
	PriorityLevel  string      `json:"priorityLevel"  orm:"priority_level"   description:"优先级级别：HIGH|MEDIUM|LOW"`                // 优先级级别：HIGH|MEDIUM|LOW
	IsDeleted      bool        `json:"isDeleted"      orm:"is_deleted"       description:"软删除标记：true 表示已删除"`                     // 软删除标记：true 表示已删除
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:"任务创建时间"`                               // 任务创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:"任务最后更新时间"`                             // 任务最后更新时间
}
