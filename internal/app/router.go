package app

import (
	"seedgo/internal/middleware"
	"seedgo/internal/system/auth"
	"seedgo/internal/system/perms"
	"seedgo/internal/system/role"
	"seedgo/internal/system/tenant"
	"seedgo/internal/system/user"

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

	//需要权限的
	g.Use(middleware.AuthMiddleware())
	{

		//认证
		auth.NewHandler().Use(g.Group("/auth"))
		//权限资源
		perms.NewHandler().Use(g.Group("system/permissions"))
		//用户
		user.NewHandler().Use(g.Group("system/users"))
		//角色
		role.NewHandler().Use(g.Group("system/roles"))
		//租户
		tenant.NewHandler().Use(g.Group("tenant/tenants"))
	}

	return r
}
