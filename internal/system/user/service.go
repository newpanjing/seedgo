package user

import (
	"context"
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"gorm.io/gorm"
)

type Service struct {
	*shared.BaseService[model.User]
}

func NewService() *Service {
	return &Service{
		BaseService: shared.NewBaseService[model.User](),
	}
}

// Create 创建
func (s *Service) Create(ctx context.Context, entity *model.User) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 处理关联角色
		if entity.RoleIds != nil && len(*entity.RoleIds) > 0 {
			var roles []*model.Role
			if err := tx.Where("id IN ?", *entity.RoleIds).Find(&roles).Error; err != nil {
				return err
			}
			entity.Roles = roles
		}

		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		return nil
	})
}

// Update 更新
func (s *Service) Update(ctx context.Context, entity *model.User) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新基本信息
		txChain := tx.Omit("created_at")
		// 如果密码为空，不更新密码
		if entity.PasswordHash == "" {
			txChain = txChain.Omit("password_hash")
		}

		if err := txChain.Model(entity).Updates(entity).Error; err != nil {
			return err
		}

		// 处理关联角色
		if entity.RoleIds != nil {
			var roles []*model.Role
			if len(*entity.RoleIds) > 0 {
				if err := tx.Where("id IN ?", *entity.RoleIds).Find(&roles).Error; err != nil {
					return err
				}
			}
			// 替换关联
			if err := tx.Model(entity).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.DB.Where("username = ?", username).Preload("Roles").First(&user).Error
	return &user, err
}

func (s *Service) FindByIdWithRoles(id model.ID) (*model.User, error) {
	var user model.User
	err := s.DB.Preload("Roles").Omit("passwordHash").First(&user, id).Error
	return &user, err
}
