package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	Question   Question `gorm:"references:ID"`
	QuestionID *uint    `gorm:"foreignKey:QuestionID;autoForeignKey"`
	Member     User     `gorm:"references:ID"`
	MemberID   *uint    `gorm:"foreignKey:MemberID;autoForeignKey"`
}
