package app

import (
	"github.com/newpanjing/seedgo/internal/middleware"
	"github.com/newpanjing/seedgo/internal/system/auth"
	"github.com/newpanjing/seedgo/internal/system/perms"
	"github.com/newpanjing/seedgo/internal/system/role"
	"github.com/newpanjing/seedgo/internal/system/tenant"
	"github.com/newpanjing/seedgo/internal/system/user"
	"github.com/newpanjing/seedgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.Cors())

	g := r.Group("/app")
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
