package role

import (
	"context"
	"seedgo/internal/model"
	"seedgo/internal/modules/perms"
	"seedgo/internal/shared"

	"gorm.io/gorm"
)

type Service struct {
	shared.BaseService[model.Role]
}

func NewService() *Service {
	return &Service{
		BaseService: *shared.NewBaseService[model.Role](),
	}
}

// Create 创建
func (l *Service) Create(ctx context.Context, entity *model.Role) error {
	return l.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 处理关联权限
		if entity.PermissionIds != nil && len(*entity.PermissionIds) > 0 {
			var perms []*model.Permission
			if err := tx.Where("id IN ?", *entity.PermissionIds).Find(&perms).Error; err != nil {
				return err
			}
			entity.Permissions = perms
		}

		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		return nil
	})
}

// Update 更新角色
func (l *Service) Update(ctx context.Context, entity *model.Role) error {
	err := l.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新基本信息
		if err := tx.Omit("created_at").Save(entity).Error; err != nil {
			return err
		}

		// 处理关联权限
		if entity.PermissionIds != nil {
			var perms []*model.Permission
			if len(*entity.PermissionIds) > 0 {
				if err := tx.Where("id IN ?", *entity.PermissionIds).Find(&perms).Error; err != nil {
					return err
				}
			}
			// 替换关联
			if err := tx.Model(entity).Association("Permissions").Replace(perms); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 更新成功后，清除所有用户的权限缓存
	return perms.GetService().ClearPermissionAllCache()
}

func (l *Service) Get(ctx context.Context, id model.ID) (*model.Role, error) {
	var role model.Role
	// 1. 只查角色基础信息
	if err := l.DB.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}

	// 2. 仅查询关联表的 ID
	var ids []model.ID
	err := l.DB.WithContext(ctx).
		Table("role_permission").
		Where("role_id = ?", role.ID).
		Pluck("permission_id", &ids).Error // Pluck 专门用于提取单列

	if err == nil {
		role.PermissionIds = &ids
	}

	return &role, nil
}
