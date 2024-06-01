package model

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	OwnerID        *uint      `gorm:"foreignKey:OwnerID;autoForeignKey"`
	Owner          User       `gorm:"references:ID"`
	Title          string     `gorm:"type:varchar(255);not null ;comment:标题"`
	Content        string     `gorm:"type:text;not null; comment:内容"`
	PictureID      *uint      `gorm:"foreignKey:PictureID;autoForeignKey"`
	Picture        UploadFile `gorm:"references:ID"`
	SubscribeCount int64      `gorm:"type:bigint;default:0"`
}
