package user

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserCtrl struct {
	logic *UserLogic
	shared.BaseCtrl[model.User]
}

func NewUserCtrl() *UserCtrl {
	logic := NewUserLogic()
	ctrl := &UserCtrl{
		logic: logic,
	}
	ctrl.BaseCtrl = *shared.NewBaseCtrl(logic, nil, ctrl)
	return ctrl
}

// 重写BeforeList
func (c *UserCtrl) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		func(d *gorm.DB) *gorm.DB {
			return d.Preload("Roles", func(db *gorm.DB) *gorm.DB {
				return db.Set("skip_tenant_filter", true)
			}).Preload("Tenant")
		},
	}
}
