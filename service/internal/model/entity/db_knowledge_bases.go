// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_knowledge_bases，见 sql/migration_knowledge_bases.sql。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbKnowledgeBases is the golang structure for table db_knowledge_bases.
type DbKnowledgeBases struct {
	Id          uuid.UUID   `json:"id"          orm:"id"          description:"主键 UUID，应用侧赋值"`
	UserId      uuid.UUID   `json:"userId"      orm:"user_id"     description:"归属用户 ID，级联删除"`
	Name        string      `json:"name"        orm:"name"        description:"库名"`
	Description string      `json:"description" orm:"description" description:"库描述"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`
}
