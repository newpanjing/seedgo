package request

import (
	"github.com/gin-gonic/gin"
)

// GetCurrentUser 从 Context 中获取当前登录用户
// 由于 AuthMiddleware 已经保证了用户的存在和合法性，这里不再返回 error
func GetCurrentUser(c *gin.Context) *UserContext {
	// 1. 从 Context 获取 user 实体
	obj, exists := c.Get("user")
	if !exists {
		// 理论上不会走到这里，如果走到说明中间件配置有误
		return nil
	}

	user, ok := obj.(*UserContext)
	if !ok {
		return nil
	}

	return user
}
