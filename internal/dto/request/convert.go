package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueryInt 获取查询参数并转为 int，转换失败则返回默认值
func QueryInt(ctx *gin.Context, key string, defaultValue int) int {
	valStr := ctx.Query(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// QueryDefaultPage 获取分页参数的快捷方式
func QueryDefaultPage(ctx *gin.Context) (page int, pageSize int) {
	page = QueryInt(ctx, "page", 1)
	pageSize = QueryInt(ctx, "pageSize", 10)

	// 基础限流保护（可选）
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return
}
