package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache: key not found")
)

// Cache 定义通用缓存接口
type Cache interface {
	// Set 设置缓存
	// value: 支持任意类型，内部会自动序列化
	// ttl: 过期时间，0 表示永不过期
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Get 获取缓存
	// dest: 接收值的变量指针 (e.g. &user)
	Get(ctx context.Context, key string, dest any) error

	// Del 删除缓存
	Del(ctx context.Context, key string) error

	// Has 判断是否存在
	Has(ctx context.Context, key string) (bool, error)

	// Clear 清空缓存 (慎用)
	Clear(ctx context.Context) error
}
