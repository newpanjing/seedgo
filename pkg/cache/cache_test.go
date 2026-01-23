package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

func TestMemoryCache(t *testing.T) {
	c := Use(NewMemoryCache())

	// Test Set & Get
	err := c.Set("user", User{Name: "test", Age: 18})
	assert.NoError(t, err)

	var u User
	err = c.Get("user", &u)
	assert.NoError(t, err)
	assert.Equal(t, "test", u.Name)
	assert.Equal(t, 18, u.Age)

	// Test Ttl
	c.Set("ttl_key", "val", 1*time.Second)
	ttl := c.Ttl("ttl_key")
	assert.True(t, ttl > 0)
	assert.True(t, ttl <= 1*time.Second)

	// Test Expire
	c.Expire("user", 1*time.Second)
	ttl = c.Ttl("user")
	assert.True(t, ttl > 0)

	// Test Delete
	c.Delete("user")
	assert.False(t, c.Has("user"))

	// Test Call
	var uCall User
	err = c.Call("call_key", &uCall, func() (any, error) {
		return User{Name: "call", Age: 20}, nil
	}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "call", uCall.Name)
	assert.Equal(t, 20, uCall.Age)

	// 验证 Call 是否缓存
	assert.True(t, c.Has("call_key"))

	// 验证值
	var u2 User
	c.Get("call_key", &u2)
	assert.Equal(t, "call", u2.Name)

	// Test Call with simple type
	var simpleVal string
	err = c.Call("simple_key", &simpleVal, func() (any, error) {
		return "simple_val", nil
	}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "simple_val", simpleVal)
}

func TestMemoryCache_Cleanup(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("expire_key", "val", 100*time.Millisecond)
	assert.True(t, c.Has("expire_key"))

	time.Sleep(200 * time.Millisecond)

	// 懒惰删除测试
	assert.False(t, c.Has("expire_key"))

	// 定时清理测试 (需要等待 ticker，这里简单手动触发一下或者信赖逻辑)
	// c.cleanup() 是私有的，无法直接测，只能测效果
	// 但 Has 已经覆盖了过期逻辑。
}

func TestMemoryCache_JSON(t *testing.T) {
	c := NewMemoryCache()

	// 测试 json 序列化 map
	m := map[string]string{"a": "b"}
	c.Set("map", m)

	var out map[string]string
	c.Get("map", &out)
	assert.Equal(t, "b", out["a"])

	// 测试 json 反序列化到 any
	// c.Get("map", &anyVar) -> anyVar will be map[string]interface{}
	var anyVar any
	c.Get("map", &anyVar)

	// anyVar should be map[string]interface{}
	// assert.IsType(t, map[string]interface{}{}, anyVar)
	// 注意：json.Unmarshal 数字默认是 float64
}
