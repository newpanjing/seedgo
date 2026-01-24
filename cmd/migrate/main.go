package main

import (
	"log"
	"seedgo/internal/db"
	"seedgo/internal/global"
	"seedgo/internal/model"
	"seedgo/pkg"
)

func main() {
	// 初始化配置和数据库
	global.InitConfig("config/local.yaml")
	db.InitDB()

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

	// 确保User表存在
	global.DB.AutoMigrate(&model.User{})

	// 修复存量数据：将所有用户的 tenant_id 更新为默认租户 ID (确保外键约束通过)
	if defaultTenantID > 0 {
		// 使用原生 SQL 确保更新成功，不依赖 GORM 模型状态
		result := global.DB.Exec("UPDATE `user` SET tenant_id = ? WHERE tenant_id = 0 OR tenant_id IS NULL", defaultTenantID)
		if result.Error != nil {
			log.Printf("Failed to update legacy users: %v", result.Error)
		} else {
			log.Printf("Updated %d legacy users to tenant ID %d", result.RowsAffected, defaultTenantID)
		}

		// 检查并创建超级用户
		createDefaultSuperUser(defaultTenantID)
	}

	// 3. 迁移其他表
	err := global.DB.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.OperationLog{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func createDefaultSuperUser(tenantID model.ID) {
	var count int64
	// 检查是否存在超级用户
	if err := global.DB.Model(&model.User{}).Where("is_super = ?", true).Count(&count).Error; err != nil {
		log.Printf("Failed to check super user: %v", err)
		return
	}

	if count == 0 {
		password, err := pkg.HashPassword("123456")
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return
		}

		isSuper := true
		user := model.User{
			Username:     "admin",
			PasswordHash: password,
			IsSuper:      &isSuper,
		}
		user.TenantID = tenantID
		user.Status = new(int8)
		*user.Status = 1 // 正常状态

		if err := global.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create super user: %v", err)
		} else {
			log.Printf("Created default super user 'admin' with ID: %d", user.ID)
		}
	}
}
