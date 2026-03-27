package utility

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

const DefaultExpire = 24 * time.Hour

func Redis() *gredis.Redis {
	return g.Redis()
}

func RedisSet(ctx context.Context, key string, value any, expire ...time.Duration) error {
	if len(expire) > 0 && expire[0] > 0 {
		expireSeconds := int64(expire[0].Seconds())
		_, err := g.Redis().Set(ctx, key, value, gredis.SetOption{
			TTLOption: gredis.TTLOption{
				EX: &expireSeconds,
			},
		})
		return err
	}
	_, err := g.Redis().Set(ctx, key, value)
	return err
}

func RedisGet(ctx context.Context, key string) (string, error) {
	result, err := g.Redis().Get(ctx, key)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func RedisDel(ctx context.Context, key string) error {
	_, err := g.Redis().Del(ctx, key)
	return err
}
