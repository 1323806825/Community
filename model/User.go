package model

import (
	"blog/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string     `gorm:"type:varchar(32);unique;default:null;comment:用户名"`
	Password  string     `gorm:"type:varchar(512);default:null;comment:密码"`
	AvatarID  *uint      `gorm:"foreignKey:AvatarID;autoForeignKey;comment:头像文件id"`
	Avatar    UploadFile `gorm:"references:ID"`
	Email     string     `gorm:"type:varchar(128);unique;comment:用户邮箱"`
	Nickname  string     `gorm:"type:varchar(32);default:null;comment:用户昵称"`
	Signature string     `gorm:"type:varchar(1024);default:null;comment:用户签名"`
}

type AuthUser struct {
	CreatedAt int64  `json:"create_time"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarID  *uint  `json:"avatar_id"`
	Signature string `json:"signature"`
}

type SafeUser struct {
	Nickname  string `json:"nickname"`
	Username  string `json:"username"`
	AvatarID  *uint  `json:"avatar_id"`
	Signature string `json:"signature"`
}

func (m *AuthUser) InPublic() SafeUser {
	return SafeUser{
		Nickname:  m.Nickname,
		Username:  m.Username,
		AvatarID:  m.AvatarID,
		Signature: m.Signature,
	}
}

func (m *User) InAuth() AuthUser {
	return AuthUser{
		CreatedAt: m.CreatedAt.Unix(),
		Username:  m.Username,
		Email:     m.Email,
		Nickname:  m.Nickname,
		AvatarID:  m.AvatarID,
		Signature: m.Signature,
	}
}

func (m *User) InPublic(avatarUrl string) SafeUser {
	return SafeUser{
		Nickname:  m.Nickname,
		Username:  m.Username,
		AvatarID:  m.AvatarID,
		Signature: m.Signature,
	}
}

func (m *User) Encrypt() error {
	hash, err := utils.Encrypt(m.Password)
	if err == nil {
		m.Password = hash
	}
	return err
}

func (m *User) BeforeCreate(_ *gorm.DB) error {
	return m.Encrypt()
}
