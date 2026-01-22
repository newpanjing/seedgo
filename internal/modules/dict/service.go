package dict

import (
	"context"
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"gorm.io/gorm"
)

type Service struct {
	*shared.BaseService[model.Dict]
}

func NewService() *Service {
	return &Service{
		BaseService: shared.NewBaseService[model.Dict](),
	}
}

// Create 创建字典及其项
func (s *Service) Create(ctx context.Context, entity *model.Dict) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		return nil
	})
}

// Update 更新字典及其项
func (s *Service) Update(ctx context.Context, entity *model.Dict) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新基本信息，忽略 createdAt
		if err := tx.Model(entity).Omit("CreatedAt").Updates(entity).Error; err != nil {
			return err
		}

		// 如果提供了 Items，则全量替换
		if entity.Items != nil {
			// 由于 DictItem 的 DictID 是 not null，Replace 默认行为是将移除的项 DictID 设为 null，这会报错
			// 所以我们需要手动处理：
			// 1. 删除不在新列表中的旧项
			// 2. 更新或插入新项

			// 收集新列表中的 ID
			var newItemIDs []model.ID
			for _, item := range entity.Items {
				if item.ID != 0 {
					newItemIDs = append(newItemIDs, item.ID)
				}
				// 确保关联 ID 正确
				item.DictID = entity.ID
			}

			// 删除该字典下，不在新 ID 列表中的项
			deleteQuery := tx.Where("dict_id = ?", entity.ID)
			if len(newItemIDs) > 0 {
				deleteQuery = deleteQuery.Where("id NOT IN ?", newItemIDs)
			}
			if err := deleteQuery.Delete(&model.DictItem{}).Error; err != nil {
				return err
			}

			// 保存新项（更新或插入）
			for _, item := range entity.Items {
				if err := tx.Save(item).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Get 获取详情并包含关联项
func (s *Service) Get(ctx context.Context, id model.ID) (*model.Dict, error) {
	var entity model.Dict
	err := s.DB.WithContext(ctx).Preload("Items").First(&entity, id).Error
	return &entity, err
}
