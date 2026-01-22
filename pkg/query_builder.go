package pkg

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Searchable 搜索接口，内部通过断言判断模型是否支持搜索
type Searchable interface {
	SearchFields() []string
}

// BuildQueryScope 根据 gin.Context 中的参数自动构建 GORM 查询 Scope
// 支持:
// 1. keyword: 关键字搜索，模型需实现 Searchable 接口
// 2. 字段过滤: ?name=admin -> WHERE name = 'admin'
// 3. 范围查询: ?age_gt=18 -> WHERE age > 18
// 4. 模糊查询: ?name_like=ad -> WHERE name LIKE '%ad%'
// 5. IN查询: ?id_in=1,2,3 -> WHERE id IN (1,2,3)
func BuildQueryScope(c *gin.Context, model any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if model == nil {
			return db
		}

		// 1. 处理关键字搜索
		keyword := c.Query("keyword")
		if keyword != "" {
			// 断言是否实现了 SearchFields 方法
			if s, ok := model.(interface{ SearchFields() []string }); ok {
				fields := s.SearchFields()
				if len(fields) > 0 {
					var conditions []string
					var args []any
					for _, field := range fields {
						conditions = append(conditions, fmt.Sprintf("`%s` LIKE ?", CamelToSnake(field)))
						args = append(args, "%"+keyword+"%")
					}
					db = db.Where(strings.Join(conditions, " OR "), args...)
				}
			}
		}

		// 2. 排除保留字段后的自动过滤
		excluded := map[string]bool{
			"keyword":   true,
			"page":      true,
			"pageSize":  true,
			"sortBy":    true,
			"sortOrder": true,
		}

		// 获取所有查询参数
		queries := c.Request.URL.Query()
		for key, values := range queries {
			if excluded[key] || len(values) == 0 || values[0] == "" {
				continue
			}

			value := values[0]
			column := ""

			// 处理后缀扩展查询
			if strings.HasSuffix(key, "_like") {
				column = CamelToSnake(strings.TrimSuffix(key, "_like"))
				db = db.Where(fmt.Sprintf("`%s` LIKE ?", column), "%"+value+"%")
			} else if strings.HasSuffix(key, "_gt") {
				column = CamelToSnake(strings.TrimSuffix(key, "_gt"))
				db = db.Where(fmt.Sprintf("`%s` > ?", column), value)
			} else if strings.HasSuffix(key, "_ge") {
				column = CamelToSnake(strings.TrimSuffix(key, "_ge"))
				db = db.Where(fmt.Sprintf("`%s` >= ?", column), value)
			} else if strings.HasSuffix(key, "_lt") {
				column = CamelToSnake(strings.TrimSuffix(key, "_lt"))
				db = db.Where(fmt.Sprintf("`%s` < ?", column), value)
			} else if strings.HasSuffix(key, "_le") {
				column = CamelToSnake(strings.TrimSuffix(key, "_le"))
				db = db.Where(fmt.Sprintf("`%s` <= ?", column), value)
			} else if strings.HasSuffix(key, "_in") {
				column = CamelToSnake(strings.TrimSuffix(key, "_in"))
				db = db.Where(fmt.Sprintf("`%s` IN ?", column), strings.Split(value, ","))
			} else {
				// 默认等于查询
				column = CamelToSnake(key)
				db = db.Where(fmt.Sprintf("`%s` = ?", column), value)
			}
		}

		return db
	}
}
