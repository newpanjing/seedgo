package auth

import (
	"seedgo/internal/form"
	"seedgo/internal/modules/user"
	"seedgo/internal/scope"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logic *user.Service
}

func NewHandler() *Handler {
	return &Handler{
		logic: user.NewService(),
	}
}

func (h *Handler) Use(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
}

func (h *Handler) GetMe(ctx *gin.Context) {
	user := scope.GetCurrentUser(ctx)
	//通过用户ID查询详情，包含角色，排除密码
	userModel, err := h.logic.GetProfile(user.ID)
	if err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	scope.OkWithData(ctx, userModel)
}

func (h *Handler) Login(ctx *gin.Context) {
	var dto form.LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		scope.Fail(ctx, "Invalid parameters")
		return
	}

	vo, err := h.logic.Login(ctx.Request.Context(), dto)
	if err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	scope.OkWithData(ctx, vo)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	var dto form.UpdateProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		scope.Fail(ctx, "Invalid parameters")
		return
	}

	user := scope.GetCurrentUser(ctx)
	if err := h.logic.UpdateProfile(ctx.Request.Context(), user.ID, dto); err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	scope.Ok(ctx)
}

func (h *Handler) Logout(ctx *gin.Context) {
	user := scope.GetCurrentUser(ctx)
	if user != nil {
		_ = h.logic.Logout(user.ID)
	}
	scope.Ok(ctx)
}

func (h *Handler) ChangePassword(ctx *gin.Context) {
	var dto form.ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		scope.Fail(ctx, "Invalid parameters")
		return
	}

	user := scope.GetCurrentUser(ctx)
	if err := h.logic.ChangePassword(ctx.Request.Context(), user.ID, dto); err != nil {
		scope.Fail(ctx, err.Error())
		return
	}

	scope.Ok(ctx)
}
