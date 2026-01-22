package global

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	if Config == nil {
		log.Fatal("Config not initialized")
	}

	dsn := Config.Database.Dsn
	if dsn == "" {
		log.Fatal("Database DSN is empty")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//使用插件
	err = DB.Use(&TenantPlugin{})
	if err != nil {
		log.Fatalf("❌ Failed to register TenantPlugin: %v", err)
	} else {
		log.Println("👏 Database connected successfully with TenantPlugin")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	//插件

	log.Println("Database connected successfully")
}
