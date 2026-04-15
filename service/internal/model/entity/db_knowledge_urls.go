// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbKnowledgeUrls is the golang structure for table db_knowledge_urls.
type DbKnowledgeUrls struct {
	Id           uuid.UUID   `json:"id"           orm:"id"            description:"主键：URL 记录唯一标识"`                            // 主键：URL 记录唯一标识
	UserId       uuid.UUID   `json:"userId"       orm:"user_id"       description:"外键：归属用户 ID"`                               // 外键：归属用户 ID
	Url          string      `json:"url"          orm:"url"           description:"网页 URL 完整地址"`                              // 网页 URL 完整地址
	Title        string      `json:"title"        orm:"title"         description:"文章标题：用于列表展示"`                              // 文章标题：用于列表展示
	Description  string      `json:"description"  orm:"description"   description:"内容摘要：用于列表预览"`                              // 内容摘要：用于列表预览
	ContentType  string      `json:"contentType"  orm:"content_type"  description:"内容类型：article|pdf|doc|etc"`                 // 内容类型：article|pdf|doc|etc
	Domain       string      `json:"domain"       orm:"domain"        description:"域名：用于按网站分组展示"`                             // 域名：用于按网站分组展示
	FaviconUrl   string      `json:"faviconUrl"   orm:"favicon_url"   description:"网站图标 URL：用于列表展示"`                          // 网站图标 URL：用于列表展示
	Author       string      `json:"author"       orm:"author"        description:"文章作者"`                                     // 文章作者
	PublishDate  *gtime.Time `json:"publishDate"  orm:"publish_date"  description:"文章原始发布时间"`                                 // 文章原始发布时间
	WordCount    int         `json:"wordCount"    orm:"word_count"    description:"字数统计：用于展示内容规模"`                            // 字数统计：用于展示内容规模
	ChunkCount   int         `json:"chunkCount"   orm:"chunk_count"   description:"分块数量：向量化后的片段数"`                            // 分块数量：向量化后的片段数
	ContentHash  string      `json:"contentHash"  orm:"content_hash"  description:"内容哈希：用于检测内容变化和去重"`                         // 内容哈希：用于检测内容变化和去重
	Status       string      `json:"status"       orm:"status"        description:"处理状态：pending|processing|completed|failed"` // 处理状态：pending|processing|completed|failed
	ErrorMessage string      `json:"errorMessage" orm:"error_message" description:"错误信息：处理失败时记录"`                             // 错误信息：处理失败时记录
	IsDeleted    bool        `json:"isDeleted"    orm:"is_deleted"    description:"软删除标记：true 表示已删除"`                         // 软删除标记：true 表示已删除
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间：软删除时记录"`                              // 删除时间：软删除时记录
	CrawledAt    *gtime.Time `json:"crawledAt"    orm:"crawled_at"    description:"抓取完成时间"`                                   // 抓取完成时间
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"用户收藏时间：用于排序"`                              // 用户收藏时间：用于排序
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"最后更新时间"`                                   // 最后更新时间
}
