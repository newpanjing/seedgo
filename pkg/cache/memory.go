package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type item struct {
	Value     []byte
	ExpiresAt int64 // Unix nanoseconds
}

// MemoryCache 基于内存的缓存实现 (简单的 map + 读写锁)
type MemoryCache struct {
	items   map[string]*item
	mu      sync.RWMutex
	lastErr error
	stop    chan struct{}
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]*item),
		stop:  make(chan struct{}),
	}
	// 启动定时清理任务
	go c.cleanupLoop()
	return c
}

// cleanupLoop 定时清理过期 key
func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stop:
			return
		}
	}
}

func (c *MemoryCache) cleanup() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.ExpiresAt > 0 && now > v.ExpiresAt {
			delete(c.items, k)
		}
	}
}

func (c *MemoryCache) Set(key string, value any, ttl ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.lastErr = err
		return err
	}

	var duration time.Duration
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	expiresAt := int64(0)
	if duration > 0 {
		expiresAt = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		Value:     data,
		ExpiresAt: expiresAt,
	}
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) Get(key string, dest any) error {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		c.lastErr = ErrCacheMiss
		return ErrCacheMiss
	}

	// 检查过期
	if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
		// 懒惰删除
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()

		c.lastErr = ErrCacheMiss
		return ErrCacheMiss
	}

	err := json.Unmarshal(item.Value, dest)
	c.lastErr = err
	return err
}

func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) DeleteMulti(keys []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
	}
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) DeletePrefix(prefix string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		// 简单的前缀匹配
		if len(key) >= len(prefix) && key[0:len(prefix)] == prefix {
			delete(c.items, key)
		}
	}
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) Has(key string) bool {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return false
	}
	if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
		c.mu.Lock()
		delete(c.items, key) // 懒惰删除
		c.mu.Unlock()
		return false
	}
	return true
}

func (c *MemoryCache) Expire(key string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		c.lastErr = ErrCacheMiss
		return ErrCacheMiss
	}

	if ttl > 0 {
		item.ExpiresAt = time.Now().Add(ttl).UnixNano()
	} else {
		item.ExpiresAt = 0
	}
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) Ttl(key string) time.Duration {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return -1
	}

	if item.ExpiresAt == 0 {
		return 0
	}

	now := time.Now().UnixNano()
	if now > item.ExpiresAt {
		return -1
	}

	return time.Duration(item.ExpiresAt - now)
}

func (c *MemoryCache) Call(key string, dest any, f func() (any, error), ttl ...time.Duration) error {
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

func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*item)
	c.lastErr = nil
	return nil
}

func (c *MemoryCache) Error() error {
	return c.lastErr
}

// Close 停止清理任务 (非接口方法，仅用于资源释放)
func (c *MemoryCache) Close() {
	close(c.stop)
}
