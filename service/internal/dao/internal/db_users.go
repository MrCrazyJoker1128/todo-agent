// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbUsersDao is the data access object for the table db_users.
type DbUsersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DbUsersColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DbUsersColumns defines and stores column names for the table db_users.
type DbUsersColumns struct {
	Id           string // 主键：用户唯一标识 (UUID)
	Email        string // 登录邮箱：唯一索引
	PasswordHash string // 密码哈希：bcrypt 加密
	Nickname     string // 用户昵称：可选
	CreatedAt    string // 账户创建时间
	UpdatedAt    string // 账户最后更新时间
}

// dbUsersColumns holds the columns for the table db_users.
var dbUsersColumns = DbUsersColumns{
	Id:           "id",
	Email:        "email",
	PasswordHash: "password_hash",
	Nickname:     "nickname",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewDbUsersDao creates and returns a new DAO object for table data access.
func NewDbUsersDao(handlers ...gdb.ModelHandler) *DbUsersDao {
	return &DbUsersDao{
		group:    "default",
		table:    "db_users",
		columns:  dbUsersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbUsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbUsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbUsersDao) Columns() DbUsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbUsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbUsersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DbUsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
