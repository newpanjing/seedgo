package main

import (
	"seedgo/internal/model"
	"seedgo/pkg/global"
	"log"
)

func main() {
	// 初始化配置和数据库
	global.InitConfig("config/local.yaml")
	global.InitDB()

	// 自动迁移
	// 这里列出所有需要迁移的 Model
	// 例如: &model.User{}
	// 请在添加业务 Model 后，导入 internal/model 包并将 Model 实例放入 AutoMigrate
	// 1. 先迁移租户表
	global.DB.AutoMigrate(&model.Tenant{})

	// 2. 检查并创建默认租户
	var count int64
	global.DB.Model(&model.Tenant{}).Count(&count)
	var defaultTenantID model.ID
	if count == 0 {
		t := model.Tenant{
			Name:   "Default Tenant",
			Status: 1,
		}
		global.DB.Create(&t)
		log.Printf("Created default tenant with ID: %d", t.ID)
		defaultTenantID = t.ID
	} else {
		var t model.Tenant
		global.DB.First(&t)
		defaultTenantID = t.ID
	}

	// 修复存量数据：将所有用户的 tenant_id 更新为默认租户 ID (确保外键约束通过)
	if defaultTenantID > 0 {
		// 使用原生 SQL 确保更新成功，不依赖 GORM 模型状态
		result := global.DB.Exec("UPDATE `user` SET tenant_id = ? WHERE tenant_id = 0 OR tenant_id IS NULL", defaultTenantID)
		if result.Error != nil {
			log.Printf("Failed to update legacy users: %v", result.Error)
		} else {
			log.Printf("Updated %d legacy users to tenant ID %d", result.RowsAffected, defaultTenantID)
		}
	}

	// 3. 迁移其他表
	err := global.DB.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
