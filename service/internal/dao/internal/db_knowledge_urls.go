// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbKnowledgeUrlsDao is the data access object for the table db_knowledge_urls.
type DbKnowledgeUrlsDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  DbKnowledgeUrlsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// DbKnowledgeUrlsColumns defines and stores column names for the table db_knowledge_urls.
type DbKnowledgeUrlsColumns struct {
	Id           string // 主键：URL 记录唯一标识
	UserId       string // 外键：归属用户 ID
	Url          string // 网页 URL 完整地址
	Title        string // 文章标题：用于列表展示
	Description  string // 内容摘要：用于列表预览
	ContentType  string // 内容类型：article|pdf|doc|etc
	Domain       string // 域名：用于按网站分组展示
	FaviconUrl   string // 网站图标 URL：用于列表展示
	Author       string // 文章作者
	PublishDate  string // 文章原始发布时间
	WordCount    string // 字数统计：用于展示内容规模
	ChunkCount   string // 分块数量：向量化后的片段数
	ContentHash  string // 内容哈希：用于检测内容变化和去重
	Status       string // 处理状态：pending|processing|completed|failed
	ErrorMessage string // 错误信息：处理失败时记录
	IsDeleted    string // 软删除标记：true 表示已删除
	DeletedAt    string // 删除时间：软删除时记录
	CrawledAt    string // 抓取完成时间
	CreatedAt    string // 用户收藏时间：用于排序
	UpdatedAt    string // 最后更新时间
}

// dbKnowledgeUrlsColumns holds the columns for the table db_knowledge_urls.
var dbKnowledgeUrlsColumns = DbKnowledgeUrlsColumns{
	Id:           "id",
	UserId:       "user_id",
	Url:          "url",
	Title:        "title",
	Description:  "description",
	ContentType:  "content_type",
	Domain:       "domain",
	FaviconUrl:   "favicon_url",
	Author:       "author",
	PublishDate:  "publish_date",
	WordCount:    "word_count",
	ChunkCount:   "chunk_count",
	ContentHash:  "content_hash",
	Status:       "status",
	ErrorMessage: "error_message",
	IsDeleted:    "is_deleted",
	DeletedAt:    "deleted_at",
	CrawledAt:    "crawled_at",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewDbKnowledgeUrlsDao creates and returns a new DAO object for table data access.
func NewDbKnowledgeUrlsDao(handlers ...gdb.ModelHandler) *DbKnowledgeUrlsDao {
	return &DbKnowledgeUrlsDao{
		group:    "default",
		table:    "db_knowledge_urls",
		columns:  dbKnowledgeUrlsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbKnowledgeUrlsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbKnowledgeUrlsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbKnowledgeUrlsDao) Columns() DbKnowledgeUrlsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbKnowledgeUrlsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbKnowledgeUrlsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DbKnowledgeUrlsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
