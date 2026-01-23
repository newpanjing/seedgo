package api

import (
	"seedgo/internal/middleware"
	"seedgo/internal/modules/auth"
	"seedgo/internal/modules/dict"
	"seedgo/internal/modules/perms"
	"seedgo/internal/modules/role"
	"seedgo/internal/modules/tenant"
	"seedgo/internal/modules/user"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	r.Use(middleware.Cors())

	g := r.Group("/api")

	//认证
	auth.NewHandler().Use(g.Group("/auth"))

	// 权限校验
	g.Use(middleware.AuthMiddleware(), middleware.PermissionsMiddleware())
	{
		//权限资源
		perms.NewHandler().Use(g.Group("system/permissions"))
		//用户
		user.NewHandler().Use(g.Group("system/users"))
		//角色
		role.NewHandler().Use(g.Group("system/roles"))
		//租户
		tenant.NewHandler().Use(g.Group("tenant/tenants"))
		//字典
		dict.NewHandler().Use(g.Group("system/dicts"))
	}

	return r
}
