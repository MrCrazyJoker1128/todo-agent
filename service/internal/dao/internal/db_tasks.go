// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbTasksDao is the data access object for the table db_tasks.
type DbTasksDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DbTasksColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DbTasksColumns defines and stores column names for the table db_tasks.
type DbTasksColumns struct {
	Id             string // 主键：任务唯一标识 (UUID)
	UserId         string // 外键：归属用户 ID
	Title          string // 任务标题：1-255 字符
	Description    string // 任务详细描述：可选
	Status         string // 任务状态：todo|doing|done|cancelled
	Priority       string // 优先级：low|medium|high
	DueDate        string // 截止时间：带时区
	MilvusVectorId string // Milvus 向量 ID (UUID 类型，与 Milvus 主键对应)
	TaskType       string // 任务类型：work|personal|project|adhoc
	PriorityLevel  string // 优先级级别：HIGH|MEDIUM|LOW
	IsDeleted      string // 软删除标记：true 表示已删除
	CreatedAt      string // 任务创建时间
	UpdatedAt      string // 任务最后更新时间
}

// dbTasksColumns holds the columns for the table db_tasks.
var dbTasksColumns = DbTasksColumns{
	Id:             "id",
	UserId:         "user_id",
	Title:          "title",
	Description:    "description",
	Status:         "status",
	Priority:       "priority",
	DueDate:        "due_date",
	MilvusVectorId: "milvus_vector_id",
	TaskType:       "task_type",
	PriorityLevel:  "priority_level",
	IsDeleted:      "is_deleted",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewDbTasksDao creates and returns a new DAO object for table data access.
func NewDbTasksDao(handlers ...gdb.ModelHandler) *DbTasksDao {
	return &DbTasksDao{
		group:    "default",
		table:    "db_tasks",
		columns:  dbTasksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbTasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbTasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbTasksDao) Columns() DbTasksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbTasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbTasksDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DbTasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
