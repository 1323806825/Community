package model

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitDB() (*gorm.DB, error) {
	logMode := logger.Info
	if !viper.GetBool("mode.develop") {
		logMode = logger.Error
	}

	db, err := gorm.Open(mysql.Open(viper.GetString("db.dsn")), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if db != nil {
		db.InstanceSet("gorm:table_options", "ENGINE=MyISAM")
	}

	sqlDb, _ := db.DB()

	err = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&UploadFile{})
	_ = db.AutoMigrate(&Question{})
	_ = db.AutoMigrate(&Subscribe{})
	_ = db.AutoMigrate(&Comment{})
	_ = db.AutoMigrate(&Like{})
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(viper.GetInt("db.SetMaxIdleConns"))
	sqlDb.SetMaxOpenConns(viper.GetInt("db.SetMaxOpenConns"))
	sqlDb.SetConnMaxLifetime(10 * time.Hour)
	return db, nil
}
