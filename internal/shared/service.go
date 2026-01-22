package shared

import (
	"context"
	"fmt"
	"seedgo/internal/global"
	"seedgo/internal/model"
	"seedgo/pkg"

	"gorm.io/gorm"
)

type IBaseService[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id model.ID) error
	Get(ctx context.Context, id model.ID) (*T, error)
	List(ctx context.Context) ([]T, error)
	Page(ctx context.Context, query pkg.QueryPage, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error)
}

// BaseService 约束 T 必须实现 model.Entity 接口
type BaseService[T any] struct {
	DB *gorm.DB
}

func NewBaseService[T any]() *BaseService[T] {
	return &BaseService[T]{
		DB: global.DB,
	}
}

func (s *BaseService[T]) Create(ctx context.Context, entity *T) error {
	return s.DB.WithContext(ctx).Create(entity).Error
}

func (s *BaseService[T]) Update(ctx context.Context, entity *T) error {
	return s.DB.WithContext(ctx).Omit("created_at").Save(entity).Error
}

func (s *BaseService[T]) Delete(ctx context.Context, id model.ID) error {
	var entity T
	return s.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (s *BaseService[T]) Get(ctx context.Context, id model.ID) (*T, error) {
	var entity T
	err := s.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (s *BaseService[T]) List(ctx context.Context) ([]T, error) {
	var entities []T
	err := s.DB.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (s *BaseService[T]) Page(ctx context.Context, query pkg.QueryPage, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	var entities []T
	var total int64

	var entity T
	// Count 只需要 Where 条件，Preload 会被忽略，但为了保持一致性我们还是带上 scopes
	s.DB.WithContext(ctx).Model(&entity).Scopes(scopes...).Count(&total)

	offset := (*query.Page - 1) * *query.PageSize
	// Find 需要 Preload，显式开启 Session 确保无状态残留
	order := fmt.Sprintf("%s %s", *query.SortBy, *query.SortDesc)
	err := s.DB.WithContext(ctx).Model(&entity).Scopes(scopes...).Order(order).Offset(offset).Limit(*query.PageSize).Session(&gorm.Session{}).Find(&entities).Error

	return entities, total, err
}
