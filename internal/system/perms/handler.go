package perms

import (
	"fmt"
	"github.com/newpanjing/seedgo/internal/model"
	"github.com/newpanjing/seedgo/internal/shared"
	"github.com/newpanjing/seedgo/pkg/request"
	"github.com/newpanjing/seedgo/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logic *Service
	shared.BaseCtrl[model.Permission]
}

// 实现这个方法可以重新
func (c *Handler) Use(g *gin.RouterGroup) {
	g.GET("/tree", c.GetTree)
	c.BaseCtrl.Use(g)
	log.Println("Registering perms routes")
}

// Update 更新权限信息
// 参数:
//
//	ctx: gin上下文对象，包含HTTP请求信息和路由参数
//
// 返回值:
//
//	无返回值，通过HTTP响应返回结果
func (c *Handler) Update(ctx *gin.Context) {
	//permissionIds 处理权限相关

	// 获取路径参数中的ID并转换为模型ID类型
	idStr := ctx.Param("id")
	id := model.ToID(idStr)

	// 绑定JSON请求体到Permission数据传输对象
	var dto model.Permission
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, fmt.Sprintf("Invalid parameters:%v", err.Error()))
		return
	}
	dto.ID = id

	// 调用业务逻辑层更新权限信息
	if err := c.logic.Update(ctx.Request.Context(), &dto); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *Handler) List(ctx *gin.Context) {
	//获取当前用户，封装一个方法
	user := request.GetCurrentUser(ctx)
	if user == nil {
		response.FailWithCode(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tree, err := c.logic.GetTree(user)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.OkWithData(ctx, response.PageResult{
		Total: 0,
		Items: tree,
	})
}

// GetTree /**
func (c *Handler) GetTree(ctx *gin.Context) {
	//获取当前用户，封装一个方法
	user := request.GetCurrentUser(ctx)
	if user == nil {
		response.FailWithCode(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tree, err := c.logic.GetTree(user)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.OkWithData(ctx, tree)
}

func NewHandler() *Handler {
	logic := NewService()
	ctrl := &Handler{
		logic: logic,
	}
	ctrl.BaseCtrl = *shared.NewBaseCtrl(logic, nil, ctrl)

	return ctrl
}
