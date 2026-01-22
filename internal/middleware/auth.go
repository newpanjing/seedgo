package middleware

import (
	"context"
	"github.com/newpanjing/seedgo/pkg/global"
	"github.com/newpanjing/seedgo/pkg/request"
	"github.com/newpanjing/seedgo/pkg/response"
	"github.com/newpanjing/seedgo/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简单的权限验证中间件示例
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在公开路径中
		path := c.Request.URL.Path
		for _, p := range global.Config.Auth.PublicPaths {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		// 这里可以实现Token验证逻辑
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailWithCode(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// 去除 Bearer 前缀
		if after, ok := strings.CutPrefix(token, "Bearer "); ok {
			token = after
		}

		// 验证 Token
		claims, err := utils.ParseToken(token)
		if err != nil {
			// global.Logger.Error("Token validation failed", "error", err)
			response.FailWithCode(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		userCtx := &request.UserContext{
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
