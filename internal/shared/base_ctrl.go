package shared

import (
	"github.com/newpanjing/seedgo/internal/model"
	"github.com/newpanjing/seedgo/pkg/request"
	"github.com/newpanjing/seedgo/pkg/response"
	"github.com/newpanjing/seedgo/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IController interface {
	Use(g *gin.RouterGroup)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
	List(ctx *gin.Context)
	//给List查询之前添加条件
	BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB
	BatchDelete(ctx *gin.Context)
}

type BaseCtrl[T any] struct {
	Logic IBaseLogic[T]
	Hook  HandlerHook[T]
	Impl  IController
}

// 创建一个NewBase的方法
func NewBaseCtrl[T any](logic IBaseLogic[T], hook HandlerHook[T], impl IController) *BaseCtrl[T] {
	var finalHook = hook
	if hook == nil {
		finalHook = &DefaultHandlerHook[T]{}
	}
	return &BaseCtrl[T]{
		Logic: logic,
		Hook:  finalHook,
		Impl:  impl,
	}
}

func (c *BaseCtrl[T]) Use(g *gin.RouterGroup) {

	g.POST("/batch-delete", c.Impl.BatchDelete)
	g.POST("", c.Impl.Create)
	g.PUT("/:id", c.Impl.Update)
	g.DELETE("/:id", c.Impl.Delete)
	g.GET("/:id", c.Impl.Get)
	g.GET("", c.Impl.List)
}

func (c *BaseCtrl[T]) Create(ctx *gin.Context) {

	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.Fail(ctx, "Invalid parameters")
		return
	}

	if err := c.Logic.Create(ctx.Request.Context(), &entity); err != nil {
		log.Printf("errors:%s", err.Error())
		response.Fail(ctx, err.Error())
		return
	}
	response.Ok(ctx)
}

func (c *BaseCtrl[T]) Update(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.Fail(ctx, "Invalid parameters")
		return
	}
	utils.SetFieldValue(&entity, "ID", id)

	if err := c.Logic.Update(ctx.Request.Context(), &entity); err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Ok(ctx)
}

func (c *BaseCtrl[T]) Delete(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	if err := c.Logic.Delete(ctx.Request.Context(), id); err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Ok(ctx)
}

func (c *BaseCtrl[T]) Get(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	entity, err := c.Logic.Get(ctx.Request.Context(), id)
	if err != nil {
		response.Fail(ctx, "Not found")
		return
	}
	response.OkWithData(ctx, entity)
}
func (c *BaseCtrl[T]) BeforeList(ctx *gin.Context) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{}
}
func (c *BaseCtrl[T]) List(ctx *gin.Context) {

	scopes := c.Impl.BeforeList(ctx)
	queryPage := request.BindQuery(ctx)
	if items, total, err := c.Logic.Page(ctx.Request.Context(), *queryPage, scopes...); err != nil {
		response.Fail(ctx, err.Error())
		return
	} else {
		response.OkWithData(ctx, response.PageResult{
			Items: items,
			Total: total,
		})
	}
}

func (c *BaseCtrl[T]) BatchDelete(ctx *gin.Context) {
	//post ids
	type Post struct {
		ids []string
	}
	post := &Post{}
	if err := ctx.ShouldBindJSON(&post); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	//组装id 用id in 删除
	response.Ok(ctx)
}
