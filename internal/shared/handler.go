package shared

import (
	"log"
	"seedgo/internal/model"
	"seedgo/internal/scope"
	"seedgo/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IHandler interface {
	Use(g *gin.RouterGroup)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	List(ctx *gin.Context)
	// BeforeList 给List查询之前添加条件
	BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB
	BatchDelete(ctx *gin.Context)
}

type BaseHandler[T any] struct {
	Logic IBaseService[T]
	Hook  HandlerHook[T]
	Impl  IHandler
}

// NewBaseHandler 创建一个NewBase的方法
func NewBaseHandler[T any](logic IBaseService[T], hook HandlerHook[T], impl IHandler) *BaseHandler[T] {
	var finalHook = hook
	if hook == nil {
		finalHook = &DefaultHandlerHook[T]{}
	}
	return &BaseHandler[T]{
		Logic: logic,
		Hook:  finalHook,
		Impl:  impl,
	}
}

func (c *BaseHandler[T]) Use(g *gin.RouterGroup) {

	g.POST("/batch-delete", c.Impl.BatchDelete)
	g.POST("", c.Impl.Create)
	g.PUT("/:id", c.Impl.Update)
	g.DELETE("/:id", c.Impl.Delete)
	g.GET("/:id", c.Impl.Get)
	g.GET("", c.Impl.List)
}

func (c *BaseHandler[T]) Create(ctx *gin.Context) {

	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		scope.Fail(ctx, "Invalid parameters")
		return
	}

	if err := c.Logic.Create(ctx.Request.Context(), &entity); err != nil {
		log.Printf("errors:%s", err.Error())
		scope.Fail(ctx, err.Error())
		return
	}
	scope.Ok(ctx)
}

func (c *BaseHandler[T]) Update(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		scope.Fail(ctx, "Invalid parameters")
		return
	}
	err := pkg.SetFieldValue(&entity, "ID", id)
	if err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	if err := c.Logic.Update(ctx.Request.Context(), &entity); err != nil {
		scope.Fail(ctx, err.Error())
		return
	}
	scope.Ok(ctx)
}

func (c *BaseHandler[T]) Delete(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	if err := c.Logic.Delete(ctx.Request.Context(), id); err != nil {
		scope.Fail(ctx, err.Error())
		return
	}
	scope.Ok(ctx)
}

func (c *BaseHandler[T]) Get(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	entity, err := c.Logic.Get(ctx.Request.Context(), id)
	if err != nil {
		scope.Fail(ctx, "Not found")
		return
	}
	scope.OkWithData(ctx, entity)
}
func (c *BaseHandler[T]) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{}
}
func (c *BaseHandler[T]) List(ctx *gin.Context) {

	scopes := c.Impl.BeforeList(ctx)
	queryPage := pkg.BindQuery(ctx)
	if items, total, err := c.Logic.Page(ctx.Request.Context(), *queryPage, scopes...); err != nil {
		scope.Fail(ctx, err.Error())
		return
	} else {
		scope.OkWithData(ctx, scope.PageResult{
			Items: items,
			Total: total,
		})
	}
}

func (c *BaseHandler[T]) BatchDelete(ctx *gin.Context) {
	//post ids
	type Post struct {
		ids []string
	}
	post := &Post{}
	if err := ctx.ShouldBindJSON(&post); err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	//组装id 用id in 删除
	scope.Ok(ctx)
}
