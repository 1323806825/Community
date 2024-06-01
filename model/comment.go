package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	OwnerID         *uint      `gorm:"foreignKey:OwnerID;autoForeignKey"`
	Owner           User       `gorm:"references:ID"`
	Content         string     `gorm:"type:text;not null; comment:内容"`
	ParentCommentID uint       `gorm:"not null;comment:父评论ID"` //表示所属楼层
	QuestionID      uint       `gorm:"not null;comment:所属问题ID"`
	PictureID       *uint      `gorm:"foreignKey:PictureID;autoForeignKey"`
	Picture         UploadFile `gorm:"references:ID"`
	//LikeCount  uint   `gorm:"default:0;comment:点赞量"`
}
