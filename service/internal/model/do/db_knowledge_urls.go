// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbKnowledgeUrls is the golang structure of table db_knowledge_urls for DAO operations like Where/Data.
type DbKnowledgeUrls struct {
	g.Meta       `orm:"table:db_knowledge_urls, do:true"`
	Id           any         // 主键：URL 记录唯一标识
	UserId       any         // 外键：归属用户 ID
	Url          any         // 网页 URL 完整地址
	Title        any         // 文章标题：用于列表展示
	Description  any         // 内容摘要：用于列表预览
	ContentType  any         // 内容类型：article|pdf|doc|etc
	Domain       any         // 域名：用于按网站分组展示
	FaviconUrl   any         // 网站图标 URL：用于列表展示
	Author       any         // 文章作者
	PublishDate  *gtime.Time // 文章原始发布时间
	WordCount    any         // 字数统计：用于展示内容规模
	ChunkCount   any         // 分块数量：向量化后的片段数
	ContentHash  any         // 内容哈希：用于检测内容变化和去重
	Status       any         // 处理状态：pending|processing|completed|failed
	ErrorMessage any         // 错误信息：处理失败时记录
	IsDeleted    any         // 软删除标记：true 表示已删除
	DeletedAt    *gtime.Time // 删除时间：软删除时记录
	CrawledAt    *gtime.Time // 抓取完成时间
	CreatedAt    *gtime.Time // 用户收藏时间：用于排序
	UpdatedAt    *gtime.Time // 最后更新时间
}
