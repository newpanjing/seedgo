package main

import (
	"fmt"
	"log"
	"seedgo/internal/api"
	"seedgo/internal/db"
	"seedgo/internal/global"
	"seedgo/pkg/cache"
)

func main() {
	// 1. 初始化配置
	global.InitConfig("config/local.yaml")

	// 2. 初始化数据库
	db.InitDB()

	// 3. 初始化缓存
	global.Cache = cache.Use(cache.NewMemoryCache())

	// 4. 初始化路由
	r := api.InitRouter()

	// 4. 启动服务
	port := global.Config.Server.Port
	if port == 0 {
		port = 3000
	}
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
	log.Println("Server started on port", port)
}
