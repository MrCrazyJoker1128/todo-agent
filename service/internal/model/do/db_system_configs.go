// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_system_configs。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbSystemConfigs is the golang structure of table db_system_configs for DAO operations like Where/Data.
type DbSystemConfigs struct {
	g.Meta    `orm:"table:db_system_configs, do:true"`
	Key       any
	Value     any
	UpdatedAt *gtime.Time
}
