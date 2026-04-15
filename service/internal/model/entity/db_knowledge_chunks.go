// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbKnowledgeChunks is the golang structure for table db_knowledge_chunks.
type DbKnowledgeChunks struct {
	Id               uuid.UUID   `json:"id"               orm:"id"                 description:"主键：分块唯一标识"`                    // 主键：分块唯一标识
	UrlId            uuid.UUID   `json:"urlId"            orm:"url_id"             description:"外键：关联的 URL 记录"`                // 外键：关联的 URL 记录
	UserId           uuid.UUID   `json:"userId"           orm:"user_id"            description:"外键：归属用户 ID（冗余字段，加速查询）"`        // 外键：归属用户 ID（冗余字段，加速查询）
	ChunkIndex       int         `json:"chunkIndex"       orm:"chunk_index"        description:"分块序号：用于排序"`                    // 分块序号：用于排序
	Content          string      `json:"content"          orm:"content"            description:"分块原文内容"`                       // 分块原文内容
	ContentHash      string      `json:"contentHash"      orm:"content_hash"       description:"分块哈希：用于去重"`                    // 分块哈希：用于去重
	WordCount        int         `json:"wordCount"        orm:"word_count"         description:"分块字数"`                         // 分块字数
	MilvusVectorId   uuid.UUID   `json:"milvusVectorId"   orm:"milvus_vector_id"   description:"Milvus 向量 ID"`                 // Milvus 向量 ID
	VectorSyncStatus string      `json:"vectorSyncStatus" orm:"vector_sync_status" description:"向量同步状态：pending|synced|failed"` // 向量同步状态：pending|synced|failed
	Metadata         string      `json:"metadata"         orm:"metadata"           description:"元数据：JSONB 格式，存储标题、段落等"`        // 元数据：JSONB 格式，存储标题、段落等
	IsDeleted        bool        `json:"isDeleted"        orm:"is_deleted"         description:"软删除标记"`                        // 软删除标记
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:"创建时间"`                         // 创建时间
}
