package log

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	*shared.BaseHandler[model.OperationLog]
}

func NewHandler() *Handler {
	h := &Handler{}
	h.BaseHandler = shared.NewBaseHandler[model.OperationLog](GetService(), nil, h)
	return h
}

// Use 注册路由，支持查询和删除
func (h *Handler) Use(g *gin.RouterGroup) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.DELETE("/:id", h.Delete)
	g.POST("/batch-delete", h.BatchDelete)
}

func (h *Handler) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	var scopes []func(*gorm.DB) *gorm.DB
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	if startDate != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("operation_time >= ?", startDate)
		})
	}
	if endDate != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("operation_time <= ?", endDate)
		})
	}
	return scopes
}
