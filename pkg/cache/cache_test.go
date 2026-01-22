package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type TestUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache()

	t.Run("SetAndGet_String", func(t *testing.T) {
		key := "test_str"
		val := "hello world"
		err := c.Set(ctx, key, val, 10*time.Second)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		var getVal string
		err = c.Get(ctx, key, &getVal)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if getVal != val {
			t.Errorf("Expected %v, got %v", val, getVal)
		}
	})

	t.Run("SetAndGet_Struct", func(t *testing.T) {
		key := "test_user"
		user := TestUser{ID: 1, Name: "Alice"}
		err := c.Set(ctx, key, user, 10*time.Second)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		var getUser TestUser
		err = c.Get(ctx, key, &getUser)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if getUser.ID != user.ID || getUser.Name != user.Name {
			t.Errorf("Expected %+v, got %+v", user, getUser)
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		key := "expire_key"
		val := "gone"
		// 设置 1 毫秒过期
		err := c.Set(ctx, key, val, 1*time.Millisecond)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 等待过期
		time.Sleep(10 * time.Millisecond)

		var getVal string
		err = c.Get(ctx, key, &getVal)
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "del_key"
		c.Set(ctx, key, "value", 0)
		err := c.Del(ctx, key)
		if err != nil {
			t.Fatalf("Del failed: %v", err)
		}

		var getVal string
		err = c.Get(ctx, key, &getVal)
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss after delete, got %v", err)
		}
	})

	t.Run("Has", func(t *testing.T) {
		key := "has_key"
		c.Set(ctx, key, "value", 0)

		exists, _ := c.Has(ctx, key)
		if !exists {
			t.Error("Expected key to exist")
		}

		exists, _ = c.Has(ctx, "non_existent")
		if exists {
			t.Error("Expected key to not exist")
		}
	})
}

func TestRedisCache(t *testing.T) {
	// 尝试连接本地 Redis，如果连不上则跳过测试
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping TestRedisCache")
		return
	}
	defer client.Close()

	c := NewRedisCache(client)

	t.Run("SetAndGet_Struct", func(t *testing.T) {
		key := "redis_user"
		user := TestUser{ID: 99, Name: "RedisUser"}
		err := c.Set(ctx, key, user, 10*time.Second)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		var getUser TestUser
		err = c.Get(ctx, key, &getUser)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if getUser.Name != user.Name {
			t.Errorf("Expected %s, got %s", user.Name, getUser.Name)
		}

		c.Del(ctx, key)
	})
}
