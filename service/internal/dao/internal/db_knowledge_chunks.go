// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbKnowledgeChunksDao is the data access object for the table db_knowledge_chunks.
type DbKnowledgeChunksDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  DbKnowledgeChunksColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// DbKnowledgeChunksColumns defines and stores column names for the table db_knowledge_chunks.
type DbKnowledgeChunksColumns struct {
	Id               string // 主键：分块唯一标识
	UrlId            string // 外键：关联的 URL 记录
	UserId           string // 外键：归属用户 ID（冗余字段，加速查询）
	ChunkIndex       string // 分块序号：用于排序
	Content          string // 分块原文内容
	ContentHash      string // 分块哈希：用于去重
	WordCount        string // 分块字数
	MilvusVectorId   string // Milvus 向量 ID
	VectorSyncStatus string // 向量同步状态：pending|synced|failed
	Metadata         string // 元数据：JSONB 格式，存储标题、段落等
	IsDeleted        string // 软删除标记
	CreatedAt        string // 创建时间
}

// dbKnowledgeChunksColumns holds the columns for the table db_knowledge_chunks.
var dbKnowledgeChunksColumns = DbKnowledgeChunksColumns{
	Id:               "id",
	UrlId:            "url_id",
	UserId:           "user_id",
	ChunkIndex:       "chunk_index",
	Content:          "content",
	ContentHash:      "content_hash",
	WordCount:        "word_count",
	MilvusVectorId:   "milvus_vector_id",
	VectorSyncStatus: "vector_sync_status",
	Metadata:         "metadata",
	IsDeleted:        "is_deleted",
	CreatedAt:        "created_at",
}

// NewDbKnowledgeChunksDao creates and returns a new DAO object for table data access.
func NewDbKnowledgeChunksDao(handlers ...gdb.ModelHandler) *DbKnowledgeChunksDao {
	return &DbKnowledgeChunksDao{
		group:    "default",
		table:    "db_knowledge_chunks",
		columns:  dbKnowledgeChunksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbKnowledgeChunksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbKnowledgeChunksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbKnowledgeChunksDao) Columns() DbKnowledgeChunksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbKnowledgeChunksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbKnowledgeChunksDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DbKnowledgeChunksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
