// ==========================================================================
// 手写（非 gf gen dao 产物）：表 db_task_summaries。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbTaskSummariesDao is the data access object for the table db_task_summaries.
type DbTaskSummariesDao struct {
	table    string
	group    string
	columns  DbTaskSummariesColumns
	handlers []gdb.ModelHandler
}

// DbTaskSummariesColumns defines and stores column names for the table db_task_summaries.
type DbTaskSummariesColumns struct {
	Id        string
	UserId    string
	Period    string
	StartDate string
	EndDate   string
	TaskType  string
	Priority  string
	Title     string
	Content   string
	Stats     string
	CreatedAt string
	UpdatedAt string
}

// dbTaskSummariesColumns holds the columns for the table db_task_summaries.
var dbTaskSummariesColumns = DbTaskSummariesColumns{
	Id:        "id",
	UserId:    "user_id",
	Period:    "period",
	StartDate: "start_date",
	EndDate:   "end_date",
	TaskType:  "task_type",
	Priority:  "priority",
	Title:     "title",
	Content:   "content",
	Stats:     "stats",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewDbTaskSummariesDao creates and returns a new DAO object for table data access.
func NewDbTaskSummariesDao(handlers ...gdb.ModelHandler) *DbTaskSummariesDao {
	return &DbTaskSummariesDao{
		group:    "default",
		table:    "db_task_summaries",
		columns:  dbTaskSummariesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbTaskSummariesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbTaskSummariesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbTaskSummariesDao) Columns() DbTaskSummariesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbTaskSummariesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbTaskSummariesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *DbTaskSummariesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
