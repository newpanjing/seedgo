package pkg

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Searchable 搜索接口，内部通过断言判断模型是否支持搜索
type Searchable interface {
	SearchFields() []string
}

// modelFieldInfo 存储模型字段的数据库列名和类型
type modelFieldInfo struct {
	DBName string
	Type   reflect.Kind
}

// modelCache 缓存模型结构信息，避免重复反射
var modelCache = sync.Map{}

// getModelFieldMap 解析模型结构体，返回 json/字段名 -> {DB列名, 类型} 的映射
func getModelFieldMap(model any) map[string]modelFieldInfo {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	// 检查缓存
	if cache, ok := modelCache.Load(t); ok {
		return cache.(map[string]modelFieldInfo)
	}

	fieldMap := make(map[string]modelFieldInfo)
	parseStruct(t, fieldMap)

	// 存入缓存
	modelCache.Store(t, fieldMap)
	return fieldMap
}

// parseStruct 递归解析结构体字段（处理内嵌结构体）
func parseStruct(t reflect.Type, fieldMap map[string]modelFieldInfo) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 处理匿名嵌入字段 (如 gorm.Model, BaseModel)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			parseStruct(field.Type, fieldMap)
			continue
		}

		// 获取 JSON 标签作为 Key
		jsonTag := field.Tag.Get("json")
		key := field.Name // 默认 Key 为字段名
		if jsonTag != "" && jsonTag != "-" {
			parts := strings.Split(jsonTag, ",")
			key = parts[0]
		}

		// 获取 GORM 标签中的 column
		dbName := ""
		gormTag := field.Tag.Get("gorm")
		if gormTag != "" {
			tags := strings.Split(gormTag, ";")
			for _, tag := range tags {
				if strings.HasPrefix(tag, "column:") {
					dbName = strings.TrimPrefix(tag, "column:")
					break
				}
			}
		}

		// 如果没有 column 标签，使用下划线命名
		if dbName == "" {
			dbName = CamelToSnake(field.Name)
		}

		// 存入映射: jsonName -> info, FieldName -> info (兼容两种传入方式)
		info := modelFieldInfo{DBName: dbName, Type: field.Type.Kind()}
		fieldMap[key] = info
		if key != field.Name {
			fieldMap[field.Name] = info
		}
	}
}

// BuildQueryScope 根据 gin.Context 中的参数自动构建 GORM 查询 Scope
// 支持:
//  1. keyword: 关键字搜索，模型需实现 Searchable 接口
//  2. 字段过滤: ?name=admin -> WHERE name = 'admin' (字符串默认LIKE，数值默认=)
//  3. 高阶用法:
//     name__in=a,b,c -> WHERE name IN ('a','b','c')
//     name__contain=a -> WHERE name LIKE '%a%'
//     name__gt=1 -> WHERE name > 1
//     name__gte=1 -> WHERE name >= 1
//     name__lt=1 -> WHERE name < 1
//     name__lte=1 -> WHERE name <= 1
func BuildQueryScope(c *gin.Context, model any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if model == nil {
			return db
		}

		// 解析模型字段映射
		fieldMap := getModelFieldMap(model)

		// 1. 处理关键字搜索
		keyword := c.Query("keyword")
		if keyword != "" {
			if s, ok := model.(interface{ SearchFields() []string }); ok {
				fields := s.SearchFields()
				if len(fields) > 0 {
					var conditions []string
					var args []any
					for _, field := range fields {
						// 尝试从映射中获取 DB 列名，如果没有则默认转换
						colName := CamelToSnake(field)
						if info, ok := fieldMap[field]; ok {
							colName = info.DBName
						} else if info, ok := fieldMap[CamelToSnake(field)]; ok {
							colName = info.DBName
						}

						conditions = append(conditions, fmt.Sprintf("`%s` LIKE ?", colName))
						args = append(args, "%"+keyword+"%")
					}
					db = db.Where(strings.Join(conditions, " OR "), args...)
				}
			} else {
				//输出警告
				log.Printf("[QueryBuilder] Warning: model '%T' does not implement Searchable interface, keyword search is disabled.", model)
			}
		}

		// 2. 排除保留字段
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

			// 解析 Key 和 Suffix
			fieldName := key
			suffix := ""

			// 支持双下划线后缀
			if idx := strings.LastIndex(key, "__"); idx != -1 {
				fieldName = key[:idx]
				suffix = key[idx+2:]
			}

			// 检查字段是否存在于模型中
			fieldInfo, exists := fieldMap[fieldName]
			if !exists {
				// 尝试匹配下划线格式（如果传入的是 snake_case 但 json 是 camelCase）
				// 这里为了严格匹配，只信任 fieldMap 中的 key (json tag or struct field name)
				log.Printf("[QueryBuilder] Warning: field '%s' (key: %s) not found in model, ignoring.", fieldName, key)
				continue
			}

			column := fieldInfo.DBName

			// 构建查询条件
			switch suffix {
			case "in":
				db = db.Where(fmt.Sprintf("`%s` IN ?", column), strings.Split(value, ","))
			case "contain":
				db = db.Where(fmt.Sprintf("`%s` LIKE ?", column), "%"+value+"%")
			case "gt":
				db = db.Where(fmt.Sprintf("`%s` > ?", column), value)
			case "gte":
				db = db.Where(fmt.Sprintf("`%s` >= ?", column), value)
			case "lt":
				db = db.Where(fmt.Sprintf("`%s` < ?", column), value)
			case "lte":
				db = db.Where(fmt.Sprintf("`%s` <= ?", column), value)
			case "":
				// 默认情况
				// 按照文档：字符串=like，数值=等于
				// 但实际上 SQL 示例 `where name='abc'` 暗示了精准匹配。
				// 鉴于 `__contain` 显式存在，默认精准匹配更为合理且常用。
				// 如果确实需要默认模糊匹配，可以取消下面注释。
				/*
					if fieldInfo.Type == reflect.String {
						db = db.Where(fmt.Sprintf("`%s` LIKE ?", column), "%"+value+"%")
					} else {
						db = db.Where(fmt.Sprintf("`%s` = ?", column), value)
					}
				*/
				db = db.Where(fmt.Sprintf("`%s` = ?", column), value)
			default:
				// 未知后缀，作为普通字段处理（或者忽略？）
				// 目前假设如果没有匹配到后缀，就把它当做字段名的一部分？
				// 但上面 split 逻辑已经切分了。如果 suffix 是未知的（比如 name__foo），
				// 最好是忽略或者报错，这里选择忽略并打日志
				log.Printf("[QueryBuilder] Warning: unknown suffix '%s' for key '%s', using equality match.", suffix, key)
				db = db.Where(fmt.Sprintf("`%s` = ?", column), value)
			}
		}

		return db
	}
}
