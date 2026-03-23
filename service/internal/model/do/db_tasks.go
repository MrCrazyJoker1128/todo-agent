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
	g.Meta      `orm:"table:db_tasks, do:true"`
	Id          any         // 主键：任务唯一标识
	UserId      any         // 外键：任务归属的用户 ID，级联删除
	Title       any         // 任务标题：简短概括任务内容
	Description any         // 任务详细描述：可包含多行文本
	Status      any         // 任务状态：todo(待办), doing(进行中), done(完成), cancelled(取消)
	Priority    any         // 任务优先级：low(低), medium(中), high(高)
	DueDate     *gtime.Time // 任务截止时间：用于排序和提醒
	ParentId    any         // 父任务 ID：实现任务拆解/子任务功能，NULL 表示根任务
	Meta        any         // 扩展元数据：JSONB 格式，存储来源、标签、原始输入等非结构化数据
	Embedding   any         // 向量数据：存储任务内容的 Embedding，用于语义搜索 (如"找关于钱的任务")
	IsDeleted   any         // 软删除标志：逻辑删除，防止误删数据
	CreatedAt   *gtime.Time // 任务创建时间
	UpdatedAt   *gtime.Time // 任务最后更新时间
}
