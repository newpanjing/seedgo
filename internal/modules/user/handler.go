package user

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	logic *Service
	shared.BaseHandler[model.User]
}

func NewHandler() *Handler {
	logic := NewService()
	ctrl := &Handler{
		logic: logic,
	}
	ctrl.BaseHandler = *shared.NewBaseHandler(logic, nil, ctrl)
	return ctrl
}

// BeforeList 重写BeforeList
func (c *Handler) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		func(d *gorm.DB) *gorm.DB {
			return d.Preload("Roles", func(db *gorm.DB) *gorm.DB {
				return db.Set("skip_tenant_filter", true)
			}).Preload("Tenant")
		},
	}
}
