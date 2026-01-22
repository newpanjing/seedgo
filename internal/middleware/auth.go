package middleware

import (
	"context"
	"net/http"
	global2 "seedgo/internal/global"
	"seedgo/internal/scope"
	"seedgo/internal/shared"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简单的权限验证中间件示例
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在公开路径中
		path := c.Request.URL.Path
		for _, p := range global2.Config.Auth.PublicPaths {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		// 获取JWT Token
		token := c.GetHeader("Authorization")
		if token == "" {
			scope.FailWithCode(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// 去除 Bearer 前缀
		if after, ok := strings.CutPrefix(token, "Bearer "); ok {
			token = after
		}

		// 验证 Token
		claims, err := shared.ParseToken(token)
		if err != nil {
			// global.Logger.Error("Token validation failed", "error", err)
			scope.FailWithCode(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		userCtx := &scope.UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			TenantID: claims.TenantID,
			IsSuper:  claims.Super,
		}
		ctx := context.WithValue(c.Request.Context(), "tenant_id", claims.TenantID)
		ctx = context.WithValue(ctx, "userId", claims.UserID)
		ctx = context.WithValue(ctx, "user", userCtx)

		c.Request = c.Request.WithContext(ctx)
		c.Set("user", userCtx)

		c.Next()

	}
}
