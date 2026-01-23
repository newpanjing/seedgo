package cache

import (
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
	Set(key string, value any, ttl ...time.Duration) error

	// Get 获取缓存
	// dest: 接收值的变量指针 (e.g. &user)
	Get(key string, dest any) error

	// Delete 删除缓存
	Delete(key string) error

	// DeleteMulti 批量删除缓存
	DeleteMulti(keys []string) error

	// DeletePrefix 删除指定前缀的缓存
	DeletePrefix(prefix string) error

	// Has 判断缓存是否存在
	Has(key string) bool

	// Expire 设置过期时间
	Expire(key string, ttl time.Duration) error

	// Ttl 获取剩余过期时间
	Ttl(key string) time.Duration

	// Call 获取缓存，如果不存在则执行回调函数并设置缓存
	// dest: 接收返回值的指针
	Call(key string, dest any, f func() (any, error), ttl ...time.Duration) error

	// Clear 清空缓存 (慎用)
	Clear() error

	// Error 获取最后一次错误
	Error() error
}

// Use 接收一个缓存实现并返回 Cache 接口
// 这是一个统一的入口点，允许用户传入任意实现了 Cache 接口的实例
// 未来可以在这里添加统一的中间件逻辑（如日志、监控、错误处理等）
func Use(impl Cache) Cache {
	return impl
}
