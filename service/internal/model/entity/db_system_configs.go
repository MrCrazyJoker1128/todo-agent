// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_system_configs，见 sql/migration_system_configs.sql。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DbSystemConfigs is the golang structure for table db_system_configs.
type DbSystemConfigs struct {
	Key       string      `json:"key"       orm:"key"        description:"配置键，如 model_key_enc_key"`
	Value     string      `json:"value"     orm:"value"      description:"配置值，明文存储"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"最后更新时间"`
}
