// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbVectorMetadata is the golang structure of table db_vector_metadata for DAO operations like Where/Data.
type DbVectorMetadata struct {
	g.Meta         `orm:"table:db_vector_metadata, do:true"`
	Id             any         // 主键：映射记录唯一标识
	UserId         any         // 用户 ID：数据隔离
	SourceTable    any         // 源表名：db_tasks 或 db_chat_messages
	SourceId       any         // 源表主键 ID
	MilvusVectorId any         // Milvus 向量 ID（由 Milvus 自动生成）
	SyncStatus     any         // 同步状态：synced|pending|failed
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
