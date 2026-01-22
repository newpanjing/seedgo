package tenant

import (
	"context"
	"errors"
	"seedgo/internal/model"
	"seedgo/internal/shared"
	"seedgo/pkg/utils"

	"gorm.io/gorm"
)

type TenantLogic struct {
	*shared.BaseService[model.Tenant]
}

func NewTenantLogic() *TenantLogic {
	return &TenantLogic{
		shared.NewBaseService[model.Tenant](),
	}
}
func (s *TenantLogic) Create(ctx context.Context, entity *model.Tenant) error {
	//入参：{"status":1,"username":"user_x7t46eus","password":"ydeux3agAa1!","phone":"15688979878","realName":"656","name":"123213"}
	//判断用户名和手机号在用户表中是否存在，不存在就创建用户关联，存在了就抛出异常。
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 检查用户名是否存在
		var count int64
		if err := tx.Model(&model.User{}).Where("username = ?", entity.Username).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("username already exists")
		}

		// 2. 检查手机号是否存在
		if entity.Phone != "" {
			if err := tx.Model(&model.User{}).Where("phone = ?", entity.Phone).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("phone number already exists")
			}
		}

		// 3. 创建租户
		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		// 4. 创建管理员用户
		hash, err := utils.HashPassword(entity.Password)
		if err != nil {
			return err
		}

		// 处理可选字段
		var realName *string
		if entity.RealName != "" {
			realName = &entity.RealName
		}

		var phone *string
		if entity.Phone != "" {
			phone = &entity.Phone
		}

		isMain := int8(1) // 主账号
		user := &model.User{
			Username:     entity.Username,
			PasswordHash: hash,
			RealName:     realName,
			Phone:        phone,
			BaseTenantModel: model.BaseTenantModel{
				TenantModel: model.TenantModel{
					TenantID: entity.ID, // 关联刚创建的租户ID
				},
			},
			IsMain:  &isMain,
			IsSuper: nil, // 默认为 false
		}
		// 显式跳过租户过滤，因为此时可能还没有上下文或是在创建新租户
		// 但由于 User 创建依赖 TenantID，我们手动设置了 TenantID
		// 这里的关键是确保 TenantPlugin 不会干扰或者正确处理
		// 通常 User 创建时 TenantPlugin 会尝试填充 TenantID，但这里我们手动指定了，插件应该会尊重（取决于插件逻辑，通常是如果不为空则不覆盖）
		// 不过看之前的插件逻辑：setTenantFieldValue 会覆盖。
		// 所以这里可能需要 skip_tenant_filter
		if err := tx.Set("skip_tenant_filter", true).Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}
