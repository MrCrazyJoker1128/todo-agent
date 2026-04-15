// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbKnowledgeChunks is the golang structure of table db_knowledge_chunks for DAO operations like Where/Data.
type DbKnowledgeChunks struct {
	g.Meta           `orm:"table:db_knowledge_chunks, do:true"`
	Id               any         // 主键：分块唯一标识
	UrlId            any         // 外键：关联的 URL 记录
	UserId           any         // 外键：归属用户 ID（冗余字段，加速查询）
	ChunkIndex       any         // 分块序号：用于排序
	Content          any         // 分块原文内容
	ContentHash      any         // 分块哈希：用于去重
	WordCount        any         // 分块字数
	MilvusVectorId   any         // Milvus 向量 ID
	VectorSyncStatus any         // 向量同步状态：pending|synced|failed
	Metadata         any         // 元数据：JSONB 格式，存储标题、段落等
	IsDeleted        any         // 软删除标记
	CreatedAt        *gtime.Time // 创建时间
}
