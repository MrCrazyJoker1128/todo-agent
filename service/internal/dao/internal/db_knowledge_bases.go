// ==========================================================================
// 手写（非 gf gen dao 产物）：表 db_knowledge_bases。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbKnowledgeBasesDao is the data access object for the table db_knowledge_bases.
type DbKnowledgeBasesDao struct {
	table    string
	group    string
	columns  DbKnowledgeBasesColumns
	handlers []gdb.ModelHandler
}

// DbKnowledgeBasesColumns defines and stores column names for the table db_knowledge_bases.
type DbKnowledgeBasesColumns struct {
	Id          string
	UserId      string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

// dbKnowledgeBasesColumns holds the columns for the table db_knowledge_bases.
var dbKnowledgeBasesColumns = DbKnowledgeBasesColumns{
	Id:          "id",
	UserId:      "user_id",
	Name:        "name",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewDbKnowledgeBasesDao creates and returns a new DAO object for table data access.
func NewDbKnowledgeBasesDao(handlers ...gdb.ModelHandler) *DbKnowledgeBasesDao {
	return &DbKnowledgeBasesDao{
		group:    "default",
		table:    "db_knowledge_bases",
		columns:  dbKnowledgeBasesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbKnowledgeBasesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbKnowledgeBasesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbKnowledgeBasesDao) Columns() DbKnowledgeBasesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbKnowledgeBasesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbKnowledgeBasesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *DbKnowledgeBasesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
