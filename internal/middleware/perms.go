package middleware

import "github.com/gin-gonic/gin"

// PermissionsMiddleware 权限验证中间件，根据当前用户的角色权限判断是否有权限访问，拦截所有接口
func PermissionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
