package cache

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type item struct {
	Value     []byte
	ExpiresAt int64 // Unix nanoseconds
}

// MemoryCache 基于内存的缓存实现 (简单的 map + 读写锁)
// 注意: 这是一个简单的实现，没有定期清理过期 key 的协程 (Get 时懒惰删除)。
// 生产环境建议使用 go-cache 等成熟库，或者自行添加 cleanup goroutine。
type MemoryCache struct {
	items map[string]*item
	mu    sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]*item),
	}
}

func (c *MemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiresAt := int64(0)
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		Value:     data,
		ExpiresAt: expiresAt,
	}
	return nil
}

func (c *MemoryCache) Get(ctx context.Context, key string, dest any) error {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return ErrCacheMiss
	}

	// 检查过期
	if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
		c.Del(ctx, key)
		return ErrCacheMiss
	}

	return json.Unmarshal(item.Value, dest)
}

func (c *MemoryCache) Del(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Has(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return false, nil
	}
	if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
		c.Del(ctx, key) // 懒惰删除
		return false, nil
	}
	return true, nil
}

func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*item)
	return nil
}

// PrefixMemoryCache 支持前缀的装饰器（简单模拟）
// 实际可以直接在 key 上加前缀，这里不做复杂封装
func (c *MemoryCache) Keys(pattern string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var keys []string
	for k := range c.items {
		if strings.Contains(k, pattern) {
			keys = append(keys, k)
		}
	}
	return keys
}
