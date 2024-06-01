package model

type Subscribe struct {
	Question   Question `gorm:"references:ID"`
	QuestionID *uint    `gorm:"foreignKey:QuestionID;autoForeignKey"`
	//gorm.db.model(&model.Subscribe).where("").count()
	Member   User  `gorm:"references:ID"`
	MemberID *uint `gorm:"foreignKey:MemberID;autoForeignKey"`
	//Subscribe bool //true//false
	//Like      bool
}
