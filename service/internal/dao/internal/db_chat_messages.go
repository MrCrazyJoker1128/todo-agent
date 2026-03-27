// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbChatMessagesDao is the data access object for the table db_chat_messages.
type DbChatMessagesDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  DbChatMessagesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// DbChatMessagesColumns defines and stores column names for the table db_chat_messages.
type DbChatMessagesColumns struct {
	Id            string // 主键：消息唯一标识 (UUID)
	UserId        string // 外键：归属用户 ID
	Role          string // 消息角色：user|assistant
	Content       string // 消息内容：文本
	RelatedTaskId string // 关联任务 ID：如果消息涉及任务
	CreatedAt     string // 消息创建时间
}

// dbChatMessagesColumns holds the columns for the table db_chat_messages.
var dbChatMessagesColumns = DbChatMessagesColumns{
	Id:            "id",
	UserId:        "user_id",
	Role:          "role",
	Content:       "content",
	RelatedTaskId: "related_task_id",
	CreatedAt:     "created_at",
}

// NewDbChatMessagesDao creates and returns a new DAO object for table data access.
func NewDbChatMessagesDao(handlers ...gdb.ModelHandler) *DbChatMessagesDao {
	return &DbChatMessagesDao{
		group:    "default",
		table:    "db_chat_messages",
		columns:  dbChatMessagesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbChatMessagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbChatMessagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbChatMessagesDao) Columns() DbChatMessagesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbChatMessagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbChatMessagesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DbChatMessagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
