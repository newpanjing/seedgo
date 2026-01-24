package log

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"
	"sync"
)

type Service struct {
	*shared.BaseService[model.OperationLog]
}

func NewService() *Service {
	return &Service{
		BaseService: shared.NewBaseService[model.OperationLog](),
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
