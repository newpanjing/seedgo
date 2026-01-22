package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 基于 Redis 的缓存实现
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	// Redis SET 支持 []byte, string, int 等，但为了统一对象行为，我们统一序列化为 JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(val, dest)
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *RedisCache) Clear(ctx context.Context) error {
	// FlushDB 清空当前 DB
	return c.client.FlushDB(ctx).Err()
}
