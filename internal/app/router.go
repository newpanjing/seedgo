package app

import (
	"seedgo/internal/middleware"
	"seedgo/internal/system/auth"
	"seedgo/internal/system/perms"
	"seedgo/internal/system/role"
	"seedgo/internal/system/tenant"
	"seedgo/internal/system/user"
	"seedgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.Cors())

	g := r.Group("/api")
	g.GET("/ping", func(c *gin.Context) {
		response.OkWithData(c, gin.H{
			"message": "pong",
			"time":    time.Now().UnixMilli(),
		})
	})

	//需要权限的
	g.Use(middleware.AuthMiddleware())
	{

		//认证
		auth.NewAuthCtrl().Use(g.Group("/auth"))
		//权限资源
		perms.NewHandler().Use(g.Group("system/permissions"))
		//用户
		user.NewUserCtrl().Use(g.Group("system/users"))
		//角色
		role.NewRoleCtrl().Use(g.Group("system/roles"))
		//租户
		tenant.NewTenantCtrl().Use(g.Group("tenant/tenants"))
	}

	return r
}
