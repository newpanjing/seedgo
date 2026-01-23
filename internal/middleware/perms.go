package middleware

import (
	"net/http"
	"seedgo/internal/global"
	"seedgo/internal/model"
	"seedgo/internal/modules/perms"
	"seedgo/internal/scope"
	"strings"

	"github.com/gin-gonic/gin"
)

// PermissionsMiddleware 权限验证中间件，根据当前用户的角色权限判断是否有权限访问，拦截所有接口
func PermissionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path

		// 1. 检查是否是公开路径 (PublicPaths)
		// 如果是公开路径，跳过权限验证 (AuthMiddleware 已经处理了是否需要登录，这里只关心是否需要权限)
		for _, p := range global.Config.Auth.PublicPaths {
			if strings.HasPrefix(currentPath, p) {
				c.Next()
				return
			}
		}

		// 2. 检查排除的 URL 列表 (ExcludePaths)
		// 这些路径虽然需要登录，但不需要特定的权限配置 (如 logout, profile)
		for _, excludePath := range global.Config.Permission.ExcludePaths {
			if strings.HasPrefix(currentPath, excludePath) {
				c.Next()
				return
			}
		}

		// 3. 获取当前用户
		user := scope.GetCurrentUser(c)
		if user == nil {
			// 用户未登录或上下文丢失，返回 401
			scope.FailWithCode(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// 4. 超级用户跳过验证
		if user.IsSuper {
			c.Next()
			return
		}

		service := perms.GetService()
		tree, err := service.GetCacheTree(user)
		if err != nil {
			scope.Fail(c, err.Error())
			c.Abort()
			return
		}
		if !hasPermission(tree, c.Request.Method, currentPath) {
			scope.FailWithCode(c, http.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}

// hasPermission 检查用户是否拥有访问权限
func hasPermission(perms []*model.Permission, method, path string) bool {

	//TODO: 用权限树匹配用户的角色权限，菜单、权限树获取，都从一个地方获取，重新登录后，刷新缓存

	return false
}
