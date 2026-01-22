package shared

import (
	"context"
	"fmt"
	"github.com/newpanjing/seedgo/pkg/global"
	"github.com/newpanjing/seedgo/internal/model"
	"github.com/newpanjing/seedgo/pkg/request"

	"gorm.io/gorm"
)

type IBaseLogic[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id model.ID) error
	Get(ctx context.Context, id model.ID) (*T, error)
	List(ctx context.Context) ([]T, error)
	Page(ctx context.Context, query request.QueryPage, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error)
}

// 约束 T 必须实现 model.Entity 接口
type BaseLogic[T any] struct {
	DB *gorm.DB
}

func NewBaseLogic[T any]() *BaseLogic[T] {
	return &BaseLogic[T]{
		DB: global.DB,
	}
}

func (s *BaseLogic[T]) Create(ctx context.Context, entity *T) error {
	return s.DB.WithContext(ctx).Create(entity).Error
}

func (s *BaseLogic[T]) Update(ctx context.Context, entity *T) error {
	return s.DB.WithContext(ctx).Omit("created_at").Save(entity).Error
}

func (s *BaseLogic[T]) Delete(ctx context.Context, id model.ID) error {
	var entity T
	return s.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (s *BaseLogic[T]) Get(ctx context.Context, id model.ID) (*T, error) {
	var entity T
	err := s.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (s *BaseLogic[T]) List(ctx context.Context) ([]T, error) {
	var entities []T
	err := s.DB.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (s *BaseLogic[T]) Page(ctx context.Context, query request.QueryPage, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
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
