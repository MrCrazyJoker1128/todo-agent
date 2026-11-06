// ==========================================================================
// 手写（非 gf gen dao 产物）：表 db_system_configs。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbSystemConfigsDao is the data access object for the table db_system_configs.
type DbSystemConfigsDao struct {
	table    string
	group    string
	columns  DbSystemConfigsColumns
	handlers []gdb.ModelHandler
}

// DbSystemConfigsColumns defines and stores column names for the table db_system_configs.
type DbSystemConfigsColumns struct {
	Key       string
	Value     string
	UpdatedAt string
}

// dbSystemConfigsColumns holds the columns for the table db_system_configs.
var dbSystemConfigsColumns = DbSystemConfigsColumns{
	Key:       "key",
	Value:     "value",
	UpdatedAt: "updated_at",
}

// NewDbSystemConfigsDao creates and returns a new DAO object for table data access.
func NewDbSystemConfigsDao(handlers ...gdb.ModelHandler) *DbSystemConfigsDao {
	return &DbSystemConfigsDao{
		group:    "default",
		table:    "db_system_configs",
		columns:  dbSystemConfigsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbSystemConfigsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbSystemConfigsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbSystemConfigsDao) Columns() DbSystemConfigsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbSystemConfigsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbSystemConfigsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *DbSystemConfigsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
