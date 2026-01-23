package global

import (
	"seedgo/pkg/cache"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Config *Configuration
	Cache  cache.Cache
)
