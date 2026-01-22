package pkg

import (
	"github.com/gin-gonic/gin"
)

type QueryPage struct {
	Page     *int `json:"page"`
	PageSize *int `json:"pageSize"`

	SortBy   *string `json:"sortBy"`
	SortDesc *string `json:"sortOrder"`
}

// BindQuery 从 gin.Context 中解析 query 参数并填充到 QueryPage
func BindQuery(c *gin.Context) *QueryPage {

	page, pageSize := QueryDefaultPage(c)
	query := &QueryPage{
		Page:     &page,
		PageSize: &pageSize,
	}
	if orderBy := c.Query("sortBy"); orderBy != "" {
		// 驼峰转下划线
		orderBy = CamelToSnake(orderBy)
		query.SortBy = &orderBy
	} else {
		id := "id"
		query.SortBy = &id
	}
	if orderDesc := c.Query("sortOrder"); orderDesc != "" {
		query.SortDesc = &orderDesc
	} else {
		desc := "desc"
		query.SortDesc = &desc
	}

	return query
}
