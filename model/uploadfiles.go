package model

import "gorm.io/gorm"

type UploadFile struct {
	gorm.Model
	FileName string `gorm:"type:varchar(128);not null ;comment:文件名"`
	FileSize int64  `gorm:"type:bigint; not null ; comment:文件大小"`
	FileType string `gorm:"type:varchar(64); not null; comment : 文件类型"`
	FilePath string `gorm:"type:varchar(128);not null; comment : 文件路径"`
	Uid      uint   `gorm:"type:int; not null;comment:上传者ID"`
	IsPublic bool   `gorm:"type:tinyint(1);noe null; default:0 ; comment:是否公开"`
}
