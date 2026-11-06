// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_task_summaries。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbTaskSummaries is the golang structure of table db_task_summaries for DAO operations like Where/Data.
type DbTaskSummaries struct {
	g.Meta    `orm:"table:db_task_summaries, do:true"`
	Id        any
	UserId    any
	Period    any
	StartDate *gtime.Time
	EndDate   *gtime.Time
	TaskType  any
	Priority  any
	Title     any
	Content   any
	Stats     any
	CreatedAt *gtime.Time
	UpdatedAt *gtime.Time
}
