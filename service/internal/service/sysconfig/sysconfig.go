// Package sysconfig 提供系统级 key-value 配置读写（表 db_system_configs）。
// 首个用途：model_key_enc_key —— 模型 Key 的 KEK 入库，替代环境变量/配置文件散落。
package sysconfig

import (
	"context"
	"errors"
	"strings"

	"service/internal/dao"
	"service/internal/utility"
)

func init() {
	// 启动即注册：utility 解密默认走库里 KEK，env/config 只做兜底。
	// DB 不可用（迁移前/单测）时返回 ""，调用方回退原有路径。
	utility.KEKProvider = func(ctx context.Context) string {
		v, err := dao.DbSystemConfigs.GetValue(ctx, KEKKey)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(v)
	}
}

const KEKKey = "model_key_enc_key"

// GetKEK 取入库的 KEK；未配置返回 ("", nil)，由调用方回退到 env/config。
func GetKEK(ctx context.Context) (string, error) {
	v, err := dao.DbSystemConfigs.GetValue(ctx, KEKKey)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(v), nil
}

// SetKEK 写入/轮换 KEK。空串直接拒绝——清 KEK 会锁死所有用户 Key，必须走运维 SQL 手动执行。
func SetKEK(ctx context.Context, raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("kek 不能为空")
	}
	return dao.DbSystemConfigs.SetValue(ctx, KEKKey, strings.TrimSpace(raw))
}

// HasKEK 是否已配置入库 KEK（只返回布尔，不回显值）。
func HasKEK(ctx context.Context) (bool, error) {
	v, err := GetKEK(ctx)
	if err != nil {
		return false, err
	}
	return v != "", nil
}
