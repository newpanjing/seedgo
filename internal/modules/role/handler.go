package role

import (
	"seedgo/internal/model"
	"seedgo/internal/scope"
	"seedgo/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	*shared.BaseHandler[model.Role]
}

func NewHandler() *Handler {

	var ctr = &Handler{}

	ctr.BaseHandler = shared.NewBaseHandler(NewService(), nil, ctr)
	return ctr
}

func (h *Handler) Use(g *gin.RouterGroup) {
	//cruddy
	h.BaseHandler.Use(g)
}

// Get 实现Get方法
func (h *Handler) Get(ctx *gin.Context) {
	id := model.ToID(ctx.Param("id"))
	entity, err := h.Logic.Get(ctx.Request.Context(), id)
	if err != nil {
		scope.Fail(ctx, "Not found")
		return
	}
	scope.OkWithData(ctx, entity)
}
