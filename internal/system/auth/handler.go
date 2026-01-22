package auth

import (
	"seedgo/pkg/request"
	"seedgo/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authLogic *Service
}

func NewHandler() *Handler {
	return &Handler{
		authLogic: NewService(),
	}
}

func (a *Handler) Use(g *gin.RouterGroup) {
	g.POST("/login", a.Login)
	g.GET("/me", a.GetMe)
	g.POST("/profile", a.UpdateProfile)
	g.POST("/logout", a.Logout)
	g.POST("/change-password", a.ChangePassword)
}

func (a *Handler) GetMe(ctx *gin.Context) {
	user := request.GetCurrentUser(ctx)
	//通过用户ID查询详情，包含角色，排除密码
	userModel, err := a.authLogic.GetProfile(user.ID)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.OkWithData(ctx, userModel)
}

func (c *Handler) Login(ctx *gin.Context) {
	var dto LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, "Invalid parameters")
		return
	}

	vo, err := c.authLogic.Login(ctx.Request.Context(), dto)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.OkWithData(ctx, vo)
}

func (a *Handler) UpdateProfile(ctx *gin.Context) {
	var dto UpdateProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, "Invalid parameters")
		return
	}

	user := request.GetCurrentUser(ctx)
	if err := a.authLogic.UpdateProfile(ctx.Request.Context(), user.ID, dto); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Ok(ctx)
}

func (a *Handler) Logout(ctx *gin.Context) {
	user := request.GetCurrentUser(ctx)
	if user != nil {
		_ = a.authLogic.Logout(user.ID)
	}
	response.Ok(ctx)
}

func (a *Handler) ChangePassword(ctx *gin.Context) {
	var dto ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, "Invalid parameters")
		return
	}

	user := request.GetCurrentUser(ctx)
	if err := a.authLogic.ChangePassword(ctx.Request.Context(), user.ID, dto); err != nil {
		response.Fail(ctx, err.Error())
		return
	}

	response.Ok(ctx)
}
