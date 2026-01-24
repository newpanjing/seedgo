package middleware

import (
	"bytes"
	"context"
	"io"
	"seedgo/internal/model"
	"seedgo/internal/modules/log"
	"seedgo/internal/scope"
	"time"

	"github.com/gin-gonic/gin"
)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// Filter out irrelevant methods
		if c.Request.Method == "OPTIONS" || c.Request.Method == "HEAD" {
			return
		}

		// Ignore operation log module itself
		if len(c.Request.URL.Path) >= 21 && c.Request.URL.Path[:21] == "/api/system/operation" {
			return
		}

		// Capture data
		latency := time.Since(startTime).Milliseconds()
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		errMsg := c.Errors.String()

		var userID model.ID
		var username string
		var tenantID model.ID

		userCtx := scope.GetCurrentUser(c)
		if userCtx != nil {
			userID = userCtx.ID
			username = userCtx.Username
			tenantID = userCtx.TenantID
		}

		// Truncate body
		bodyStr := string(bodyBytes)
		if len(bodyStr) > 2000 {
			bodyStr = bodyStr[:2000] + "..."
		}

		// Async save
		go func() {
			opLog := &model.OperationLog{
				BaseTenantModel: model.BaseTenantModel{
					TenantModel: model.TenantModel{TenantID: tenantID},
				},
				UserID:       userID,
				Username:     username,
				Method:       method,
				Path:         path,
				Query:        query,
				Body:         bodyStr,
				IP:           clientIP,
				UserAgent:    userAgent,
				Status:       status,
				Latency:      latency,
				ErrorMessage: errMsg,
			}
			opLog.SetOperationTime()

			// Use background context
			_ = log.GetService().Create(context.Background(), opLog)
		}()
	}
}
