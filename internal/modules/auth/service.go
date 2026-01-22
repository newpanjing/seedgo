package auth

import (
	"context"
	"errors"
	"seedgo/internal/model"
	"seedgo/internal/shared"
	"seedgo/pkg"
	"time"
)

type Service struct {
	shared.BaseService[model.User]
}

func NewService() *Service {
	return &Service{
		BaseService: *shared.NewBaseService[model.User](),
	}
}

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginVO struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type UpdateProfileDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (*LoginVO, error) {
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

	return &LoginVO{
		Token: token,
		User:  user,
	}, nil
}

// FindByUsername find user by username
func (s *Service) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.DB.First(&user, "username = ?", username).Error
	return &user, err
}

func (s *Service) GetProfile(userID model.ID) (*model.User, error) {
	return s.FindByIdWithRoles(userID)
}

func (s *Service) FindByIdWithRoles(userID model.ID) (*model.User, error) {
	var user model.User
	err := s.DB.Preload("Roles").Omit("passwordHash").First(&user, userID).Error
	return &user, err
}
func (s *Service) UpdateProfile(ctx context.Context, uid model.ID, dto UpdateProfileDTO) error {
	user, err := s.Get(ctx, uid)
	if err != nil {
		return err
	}
	user.RealName = &dto.Name
	user.Phone = &dto.Phone
	user.Email = &dto.Email
	return s.Update(ctx, user)
}

func (s *Service) ChangePassword(ctx context.Context, uid model.ID, dto ChangePasswordDTO) error {
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
	// Implement any cleanup if necessary (e.g. invalidate token in redis)

	return nil
}
