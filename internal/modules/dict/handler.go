package dict

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	logic *Service
	shared.BaseHandler[model.Dict]
}

func NewHandler() *Handler {
	logic := NewService()
	ctrl := &Handler{
		logic: logic,
	}
	ctrl.BaseHandler = *shared.NewBaseHandler(logic, nil, ctrl)
	return ctrl
}

// BeforeList 重写列表查询前的钩子，可以按需添加 Preload 或过滤
func (c *Handler) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		func(d *gorm.DB) *gorm.DB {
			// 列表页通常不需要加载所有字典项，或者根据需要加载
			return d
		},
	}
}

// BatchDelete 批量删除
func (c *Handler) BatchDelete(ctx *gin.Context) {
	// 继承 BaseHandler 的 BatchDelete 实现，或者根据需要重写
	c.BaseHandler.BatchDelete(ctx)
}
