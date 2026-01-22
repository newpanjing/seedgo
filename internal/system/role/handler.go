package role

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"
	"seedgo/pkg/response"

	"github.com/gin-gonic/gin"
)

type RoleCtrl struct {
	*shared.BaseCtrl[model.Role]
}

func NewRoleCtrl() *RoleCtrl {

	var ctr = &RoleCtrl{}

	ctr.BaseCtrl = shared.NewBaseCtrl(NewRoleLogic(), nil, ctr)
	return ctr
}

func (r *RoleCtrl) Use(g *gin.RouterGroup) {
	//注册基本crud
	r.BaseCtrl.Use(g)
	g.GET("/c", func(ctx *gin.Context) {})
	g.GET("/d", func(ctx *gin.Context) {})
	g.GET("/e", func(ctx *gin.Context) {})
}

// 实现Get方法
func (c *RoleCtrl) Get(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	entity, err := c.Logic.Get(ctx.Request.Context(), id)
	if err != nil {
		response.Fail(ctx, "Not found")
		return
	}
	response.OkWithData(ctx, entity)
}
