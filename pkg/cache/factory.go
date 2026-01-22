package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Type string

const (
	TypeMemory Type = "memory"
	TypeRedis  Type = "redis"
)

type Config struct {
	Type        Type
	RedisOption *redis.Options // 仅当 Type == "redis" 时需要
}

// NewCache 创建缓存实例
func NewCache(cfg Config) (Cache, error) {
	switch cfg.Type {
	case TypeRedis:
		if cfg.RedisOption == nil {
			return nil, fmt.Errorf("redis option is required for redis cache")
		}
		client := redis.NewClient(cfg.RedisOption)
		return NewRedisCache(client), nil
	case TypeMemory:
		return NewMemoryCache(), nil
	default:
		// 默认返回内存缓存
		return NewMemoryCache(), nil
	}
}
