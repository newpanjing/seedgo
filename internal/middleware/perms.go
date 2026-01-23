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
			scope.FailWithCode(c, http.StatusUnauthorized, "请重新登录")
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
		// 匹配权限树中的路径
		if !hasPermission(tree, c.Request.Method, currentPath) {
			scope.FailWithCode(c, http.StatusForbidden, "没有权限访问")
			c.Abort()
			return
		}
		c.Next()
	}
}

// hasPermission 检查用户是否拥有访问权限
func hasPermission(perms []*model.Permission, method, path string) bool {
	if len(perms) == 0 {
		return false
	}

	// 遍历权限树，检查是否有匹配的路径
	for _, p := range perms {
		matched := false
		// 1.如果path 包含p.Path 或者 p.Path 包含 urls的时候，就匹配到了
		if p.Path != "" && strings.Contains(path, p.Path) {
			matched = true
		} else if p.PermissionUrls != "" {
			urls := strings.Split(p.PermissionUrls, ",")
			for _, u := range urls {
				if strings.Contains(path, strings.TrimSpace(u)) {
					matched = true
					break
				}
			}
		}

		if matched {
			// 2. 如果请求的是GET 方法，就直接返回true
			if method == "GET" {
				return true
			}

			// 如果是其他类型的方法POST、PUT、DELETE，在判断子项的code是否有:create,:update,:delete
			var suffix string
			switch method {
			case "POST":
				// 3. 例外处理：请求路径：/system/roles/batch-delete，默认为删除权限
				if strings.HasSuffix(path, "/batch-delete") {
					suffix = ":delete"
				} else {
					suffix = ":create"
				}
			case "PUT":
				suffix = ":update"
			case "DELETE":
				suffix = ":delete"
			}

			if suffix != "" {
				for _, child := range p.Children {
					if strings.HasSuffix(child.PermissionCode, suffix) {
						return true
					}
				}
			}
		}

		if len(p.Children) > 0 {
			if hasPermission(p.Children, method, path) {
				return true
			}
		}
	}

	return false
}
