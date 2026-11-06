// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_task_summaries，见 sql/migration_task_summaries.sql。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbTaskSummaries is the golang structure for table db_task_summaries.
type DbTaskSummaries struct {
	Id        uuid.UUID   `json:"id"        orm:"id"         description:"主键 UUID，应用侧赋值"`
	UserId    uuid.UUID   `json:"userId"    orm:"user_id"    description:"归属用户 ID，级联删除"`
	Period    string      `json:"period"    orm:"period"     description:"周期：day|week|month|custom"`
	StartDate *gtime.Time `json:"startDate" orm:"start_date" description:"统计窗口起点"`
	EndDate   *gtime.Time `json:"endDate"   orm:"end_date"   description:"统计窗口终点"`
	TaskType  string      `json:"taskType"  orm:"task_type"  description:"任务类型过滤，空=全部"`
	Priority  string      `json:"priority"  orm:"priority"   description:"紧急程度过滤（priority 小写版），空=全部"`
	Title     string      `json:"title"     orm:"title"      description:"总结标题"`
	Content   string      `json:"content"   orm:"content"    description:"服务端拼好的总结正文 markdown"`
	Stats     string      `json:"stats"     orm:"stats"      description:"统计 JSON 字符串"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
}
