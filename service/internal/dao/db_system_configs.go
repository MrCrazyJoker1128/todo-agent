// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_system_configs 的读写方法。
// =================================================================================

package dao

import (
	"context"

	"service/internal/dao/internal"
	"service/internal/model/do"
	"service/internal/model/entity"
)

// dbSystemConfigsDao is the data access object for the table db_system_configs.
type dbSystemConfigsDao struct {
	*internal.DbSystemConfigsDao
}

var (
	// DbSystemConfigs is a globally accessible object for table db_system_configs operations.
	DbSystemConfigs = dbSystemConfigsDao{internal.NewDbSystemConfigsDao()}
)

// GetValue 按 key 取配置值，无行返回 ("", nil)。
func (dao *dbSystemConfigsDao) GetValue(ctx context.Context, key string) (string, error) {
	var cfg *entity.DbSystemConfigs
	err := dao.Ctx(ctx).Where(do.DbSystemConfigs{Key: key}).Scan(&cfg)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", nil
	}
	return cfg.Value, nil
}

// SetValue 新增或更新配置值（key 主键，Upsert 语义）。
func (dao *dbSystemConfigsDao) SetValue(ctx context.Context, key, value string) error {
	count, err := dao.Ctx(ctx).Where(do.DbSystemConfigs{Key: key}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = dao.Ctx(ctx).Data(do.DbSystemConfigs{Key: key, Value: value}).Insert()
		return err
	}
	_, err = dao.Ctx(ctx).Where(do.DbSystemConfigs{Key: key}).Data(do.DbSystemConfigs{Value: value}).Update()
	return err
}
