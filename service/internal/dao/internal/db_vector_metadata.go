// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbVectorMetadataDao is the data access object for the table db_vector_metadata.
type DbVectorMetadataDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  DbVectorMetadataColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// DbVectorMetadataColumns defines and stores column names for the table db_vector_metadata.
type DbVectorMetadataColumns struct {
	Id             string // 主键：映射记录唯一标识
	UserId         string // 用户 ID：数据隔离
	SourceTable    string // 源表名：db_tasks 或 db_chat_messages
	SourceId       string // 源表主键 ID
	MilvusVectorId string // Milvus 向量 ID（由 Milvus 自动生成）
	SyncStatus     string // 同步状态：synced|pending|failed
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// dbVectorMetadataColumns holds the columns for the table db_vector_metadata.
var dbVectorMetadataColumns = DbVectorMetadataColumns{
	Id:             "id",
	UserId:         "user_id",
	SourceTable:    "source_table",
	SourceId:       "source_id",
	MilvusVectorId: "milvus_vector_id",
	SyncStatus:     "sync_status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewDbVectorMetadataDao creates and returns a new DAO object for table data access.
func NewDbVectorMetadataDao(handlers ...gdb.ModelHandler) *DbVectorMetadataDao {
	return &DbVectorMetadataDao{
		group:    "default",
		table:    "db_vector_metadata",
		columns:  dbVectorMetadataColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbVectorMetadataDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbVectorMetadataDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbVectorMetadataDao) Columns() DbVectorMetadataColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbVectorMetadataDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbVectorMetadataDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DbVectorMetadataDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
