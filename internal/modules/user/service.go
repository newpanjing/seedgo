package user

import (
	"context"
	"errors"
	"seedgo/internal/form"
	"seedgo/internal/model"
	"seedgo/internal/modules/perms"
	"seedgo/internal/shared"
	"seedgo/pkg"
	"sync"
	"time"

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

// 单例模式
var (
	instance *Service
	once     sync.Once
)

// GetService 获取单例实例
func GetService() *Service {
	once.Do(func() {
		instance = NewService()
	})
	return instance
}

func (s *Service) Login(ctx context.Context, dto form.LoginDTO) (*form.LoginVO, error) {
	user, err := s.FindByUsername(dto.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !pkg.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != nil && *user.Status == 0 {
		return nil, errors.New("user is disabled")
	}

	// Update login info
	now := time.Now()
	user.LastLoginAt = &now
	// user.LastLoginIP = ... // context is needed to get IP, or passed in DTO
	s.Update(ctx, user)

	isSuper := false
	if user.IsSuper != nil {
		isSuper = *user.IsSuper
	}

	token, err := shared.GenerateToken(user.ID, user.Username, user.TenantID, isSuper)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &form.LoginVO{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetProfile(userID model.ID) (*model.User, error) {
	return s.FindByIdWithRoles(userID)
}

func (s *Service) UpdateProfile(ctx context.Context, uid model.ID, dto form.UpdateProfileDTO) error {
	user, err := s.Get(ctx, uid)
	if err != nil {
		return err
	}
	user.RealName = &dto.Name
	user.Phone = &dto.Phone
	user.Email = &dto.Email
	return s.Update(ctx, user)
}

func (s *Service) ChangePassword(ctx context.Context, uid model.ID, dto form.ChangePasswordDTO) error {
	user, err := s.Get(ctx, uid)
	if err != nil {
		return err
	}
	if !pkg.CheckPasswordHash(dto.OldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}
	hash, err := pkg.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.Update(ctx, user)
}

func (s *Service) Logout(uid model.ID) error {
	return perms.GetService().ClearPermissionCache(uid)
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
