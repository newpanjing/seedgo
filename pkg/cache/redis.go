package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 基于 Redis 的缓存实现
type RedisCache struct {
	client  *redis.Client
	lastErr error
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Set(key string, value any, ttl ...time.Duration) error {
	var duration time.Duration
	if len(ttl) > 0 {
		duration = ttl[0]
	}
	// Redis SET 支持 []byte, string, int 等，但为了统一对象行为，我们统一序列化为 JSON
	data, err := json.Marshal(value)
	if err != nil {
		c.lastErr = err
		return err
	}
	err = c.client.Set(context.Background(), key, data, duration).Err()
	c.lastErr = err
	return err
}

func (c *RedisCache) Get(key string, dest any) error {
	val, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.lastErr = ErrCacheMiss
			return ErrCacheMiss
		}
		c.lastErr = err
		return err
	}
	err = json.Unmarshal(val, dest)
	c.lastErr = err
	return err
}

func (c *RedisCache) Delete(key string) error {
	err := c.client.Del(context.Background(), key).Err()
	c.lastErr = err
	return err
}

func (c *RedisCache) DeleteMulti(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	err := c.client.Del(context.Background(), keys...).Err()
	c.lastErr = err
	return err
}

func (c *RedisCache) DeletePrefix(prefix string) error {
	ctx := context.Background()
	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		// 批量删除，防止一次积压太多
		if len(keys) >= 100 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				c.lastErr = err
				return err
			}
			keys = keys[:0]
		}
	}
	if err := iter.Err(); err != nil {
		c.lastErr = err
		return err
	}
	if len(keys) > 0 {
		err := c.client.Del(ctx, keys...).Err()
		c.lastErr = err
		return err
	}
	return nil
}

func (c *RedisCache) Has(key string) bool {
	n, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		c.lastErr = err
		return false
	}
	c.lastErr = nil
	return n > 0
}

func (c *RedisCache) Expire(key string, ttl time.Duration) error {
	err := c.client.Expire(context.Background(), key, ttl).Err()
	c.lastErr = err
	return err
}

func (c *RedisCache) Ttl(key string) time.Duration {
	d, err := c.client.TTL(context.Background(), key).Result()
	if err != nil {
		c.lastErr = err
		return -1
	}
	c.lastErr = nil
	return d
}

func (c *RedisCache) Call(key string, dest any, f func() (any, error), ttl ...time.Duration) error {
	// 1. Check if key exists
	if c.Has(key) {
		return c.Get(key, dest)
	}

	// 2. Execute function
	val, err := f()
	if err != nil {
		return err
	}

	// 3. Set cache
	if err := c.Set(key, val, ttl...); err != nil {
		return err
	}

	// 4. Populate dest
	return c.Get(key, dest)
}

func (c *RedisCache) Clear() error {
	// FlushDB 清空当前 DB
	err := c.client.FlushDB(context.Background()).Err()
	c.lastErr = err
	return err
}

func (c *RedisCache) Error() error {
	return c.lastErr
}
