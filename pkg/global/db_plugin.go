// 全局多租户插件和delete_at过滤
package global

import (
	"seedgo/internal/model"
	"seedgo/pkg/request"
	"reflect"

	"gorm.io/gorm"
)

type TenantPlugin struct{}

func (p *TenantPlugin) Name() string { return "TenantPlugin" }

func (p *TenantPlugin) Initialize(db *gorm.DB) error {
	// 覆盖查询、更新、删除
	db.Callback().Query().Before("gorm:query").Register("tenant:filter", p.filter)
	db.Callback().Delete().Before("gorm:delete").Register("tenant:filter", p.filter)

	db.Callback().Update().Before("gorm:update").Register("tenant:filter", p.filter)
	db.Callback().Create().Before("gorm:create").Register("tenant:filter", p.create)
	return nil
}

func (p *TenantPlugin) create(db *gorm.DB) {
	if db.Statement.Schema != nil {
		// 检查是否跳过租户过滤
		if skip, ok := db.Get("skip_tenant_filter"); ok && skip.(bool) {
			return
		}

		value := db.Statement.Context.Value("user")
		// 查找模型中是否存在 TenantID 字段
		if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
			// 从 Context 获取租户 ID
			if value != nil {
				user := value.(*request.UserContext)
				//超级用户可以看所有数据
				if !user.IsSuper {
					//给实体设置
					// 3. 处理数据注入 (支持单条和批量)
					switch db.Statement.ReflectValue.Kind() {
					case reflect.Slice, reflect.Array:
						// 批量插入场景
						for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
							item := db.Statement.ReflectValue.Index(i)
							//设置租户id
							p.setTenantFieldValue(item, user.TenantID)
						}
					case reflect.Struct:
						// 单条插入场景
						p.setTenantFieldValue(db.Statement.ReflectValue, user.TenantID)
					}
				}
			}
		}
	}
}

func (p *TenantPlugin) setTenantFieldValue(val reflect.Value, tenantID model.ID) {
	// 如果是指针，获取指向的值
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return
	}

	// 查找并设置 TenantID 字段
	if field := val.FieldByName("TenantID"); field.IsValid() && field.CanSet() {
		// 设置字段值
		if field.Kind() == reflect.Int64 || field.Type().ConvertibleTo(reflect.TypeOf(tenantID)) {
			field.Set(reflect.ValueOf(tenantID).Convert(field.Type()))
		}
	}
}

func (p *TenantPlugin) filter(db *gorm.DB) {

	if db.Statement.Schema != nil {
		// 检查是否跳过租户过滤
		if skip, ok := db.Get("skip_tenant_filter"); ok && skip.(bool) {
			return
		}

		value := db.Statement.Context.Value("user")

		// 查找模型中是否存在 TenantID 字段
		if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
			// 从 Context 获取租户 ID
			if value != nil {
				user := value.(*request.UserContext)
				//超级用户可以看所有数据
				if !user.IsSuper {
					// 只有当查询条件中没有 tenant_id 时才添加
					// 注意：这只是一个简单的检查，复杂的 SQL 可能无法覆盖
					// 更好的方式是在 Context 中设置一个标志位来跳过租户过滤，或者检查 db.Statement.Vars
					db.Where("tenant_id = ?", user.TenantID)
				}
			}
		}
	}
}
