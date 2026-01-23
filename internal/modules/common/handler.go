package common

import (
	"seedgo/internal/form"
	"seedgo/internal/model"
	"seedgo/internal/modules/perms"
	"seedgo/internal/modules/role"
	"seedgo/internal/modules/user"
	"seedgo/internal/scope"
	"seedgo/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Use(g *gin.RouterGroup) {
	//公共接口，只要登录了都可以获取，不受权限控制

	//用户相关
	g.GET("user/profile", h.GetProfile)
	g.POST("user/profile", h.UpdateProfile)
	g.POST("user/change-password", h.ChangePassword)

	//权限树获取
	g.GET("user/permissions", h.GetPermissions)

	//角色获取
	options := g.Group("options")
	{
		//角色获取
		options.GET("roles", h.GetRoles)
	}

}

// GetRoles 获取角色下拉框
func (h Handler) GetRoles(ctx *gin.Context) {
	// 自动从查询参数构建 Where 条件
	var scopes []func(*gorm.DB) *gorm.DB
	scopes = append(scopes, pkg.BuildQueryScope(ctx, &model.Role{}))
	service := role.Instance()
	queryPage := pkg.BindQuery(ctx)
	items, total, err := service.Options(ctx.Request.Context(), *queryPage, scopes...)
	if err != nil {
		scope.Fail(ctx, err.Error())
		return
	}
	//items 转换为 SelectOption 类型
	var options []SelectOption
	for _, item := range items {
		options = append(options, SelectOption{
			Label: item.Name,
			Value: item.ID,
		})
	}
	scope.OkWithData(ctx, scope.PageResult{
		Items: options,
		Total: total,
	})
}

// GetProfile 获取用户信息
func (h Handler) GetProfile(c *gin.Context) {
	u := scope.GetCurrentUser(c)
	out, err := user.GetService().GetProfile(u.ID)
	if err != nil {
		scope.Fail(c, err.Error())
		return
	}
	scope.OkWithData(c, out)
}

func (h Handler) UpdateProfile(c *gin.Context) {
	var dto form.UpdateProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		scope.Fail(c, err.Error())
		return
	}
	err := user.GetService().UpdateProfile(c.Request.Context(), scope.GetCurrentUser(c).ID, dto)
	if err != nil {
		scope.Fail(c, err.Error())
		return
	}
	scope.Ok(c)
}

func (h Handler) ChangePassword(c *gin.Context) {
	var dto form.ChangePasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		scope.Fail(c, err.Error())
		return
	}
	err := user.GetService().ChangePassword(c.Request.Context(), scope.GetCurrentUser(c).ID, dto)
	if err != nil {
		scope.Fail(c, err.Error())
		return
	}
	scope.Ok(c)
}

func (h Handler) GetPermissions(c *gin.Context) {
	permissions, err := perms.GetService().GetCacheTree(scope.GetCurrentUser(c))
	if err != nil {
		scope.Fail(c, err.Error())
		return
	}
	scope.OkWithData(c, permissions)
}
