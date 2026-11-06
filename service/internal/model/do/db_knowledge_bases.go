// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_knowledge_bases。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbKnowledgeBases is the golang structure of table db_knowledge_bases for DAO operations like Where/Data.
type DbKnowledgeBases struct {
	g.Meta      `orm:"table:db_knowledge_bases, do:true"`
	Id          any
	UserId      any
	Name        any
	Description any
	CreatedAt   *gtime.Time
	UpdatedAt   *gtime.Time
}
