// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbVectorMetadata is the golang structure for table db_vector_metadata.
type DbVectorMetadata struct {
	Id             uuid.UUID   `json:"id"             orm:"id"               description:"主键：映射记录唯一标识"`                     // 主键：映射记录唯一标识
	UserId         uuid.UUID   `json:"userId"         orm:"user_id"          description:"用户 ID：数据隔离"`                      // 用户 ID：数据隔离
	SourceTable    string      `json:"sourceTable"    orm:"source_table"     description:"源表名：db_tasks 或 db_chat_messages"` // 源表名：db_tasks 或 db_chat_messages
	SourceId       uuid.UUID   `json:"sourceId"       orm:"source_id"        description:"源表主键 ID"`                         // 源表主键 ID
	MilvusVectorId int64       `json:"milvusVectorId" orm:"milvus_vector_id" description:"Milvus 向量 ID（由 Milvus 自动生成）"`     // Milvus 向量 ID（由 Milvus 自动生成）
	SyncStatus     string      `json:"syncStatus"     orm:"sync_status"      description:"同步状态：synced|pending|failed"`      // 同步状态：synced|pending|failed
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:"创建时间"`                            // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:"更新时间"`                            // 更新时间
}
