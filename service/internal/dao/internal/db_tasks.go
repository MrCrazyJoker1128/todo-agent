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
	Id          string // 主键：任务唯一标识
	UserId      string // 外键：任务归属的用户 ID，级联删除
	Title       string // 任务标题：简短概括任务内容
	Description string // 任务详细描述：可包含多行文本
	Status      string // 任务状态：todo(待办), doing(进行中), done(完成), cancelled(取消)
	Priority    string // 任务优先级：low(低), medium(中), high(高)
	DueDate     string // 任务截止时间：用于排序和提醒
	ParentId    string // 父任务 ID：实现任务拆解/子任务功能，NULL 表示根任务
	Meta        string // 扩展元数据：JSONB 格式，存储来源、标签、原始输入等非结构化数据
	Embedding   string // 向量数据：存储任务内容的 Embedding，用于语义搜索 (如"找关于钱的任务")
	IsDeleted   string // 软删除标志：逻辑删除，防止误删数据
	CreatedAt   string // 任务创建时间
	UpdatedAt   string // 任务最后更新时间
}

// dbTasksColumns holds the columns for the table db_tasks.
var dbTasksColumns = DbTasksColumns{
	Id:          "id",
	UserId:      "user_id",
	Title:       "title",
	Description: "description",
	Status:      "status",
	Priority:    "priority",
	DueDate:     "due_date",
	ParentId:    "parent_id",
	Meta:        "meta",
	Embedding:   "embedding",
	IsDeleted:   "is_deleted",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
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
