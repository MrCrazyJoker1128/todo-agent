// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbTasks is the golang structure of table db_tasks for DAO operations like Where/Data.
type DbTasks struct {
	g.Meta         `orm:"table:db_tasks, do:true"`
	Id             any         // 主键：任务唯一标识 (UUID)
	UserId         any         // 外键：归属用户 ID
	Title          any         // 任务标题：1-255 字符
	Description    any         // 任务详细描述：可选
	Status         any         // 任务状态：todo|doing|done|cancelled
	Priority       any         // 优先级：low|medium|high
	DueDate        *gtime.Time // 截止时间：带时区
	MilvusVectorId any         // Milvus 向量 ID (UUID 类型，与 Milvus 主键对应)
	TaskType       any         // 任务类型：work|personal|project|adhoc
	PriorityLevel  any         // 优先级级别：HIGH|MEDIUM|LOW
	IsDeleted      any         // 软删除标记：true 表示已删除
	CreatedAt      *gtime.Time // 任务创建时间
	UpdatedAt      *gtime.Time // 任务最后更新时间
}
