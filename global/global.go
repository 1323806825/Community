package global

import (
	"blog/config"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *config.RedisClient
)
