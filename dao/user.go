package dao

import (
	"blog/dto"
	"blog/model"
	"blog/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

var userDao *UserDao

type UserDao struct {
	BaseDao
}

func NewUserDao() *UserDao {
	if userDao == nil {
		userDao = &UserDao{
			NewBaseDao(),
		}
	}
	return userDao
}

func (m *UserDao) GetUserByName(stUsername string) (model.User, error) {
	var iUser model.User
	err := m.Orm.Model(&model.User{}).Where("username = ?", stUsername).Find(&iUser).Error
	return iUser, err
}

func (m *UserDao) CheckUserNameExist(stUserName string) (bool, error) {
	//定义一个int64的nTotal
	//查询同名的count数量
	//判断count>0
	var user model.User
	fmt.Println(m)
	err := m.Orm.Model(&model.User{}).Where("username = ?", stUserName).Find(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("not data")
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (m *UserDao) CheckUserEmailExist(stEmail string) bool {
	//同上
	var nTotal int64
	//m.Orm.Where("email = ?", stEmail).Count(&nTotal)
	m.Orm.Where("email = ?", stEmail).Count(&nTotal)
	return nTotal > 0
}

func (m *UserDao) AddUser(iUserDTO *dto.UserAddDTO) error {
	var iUser model.User
	iUserDTO.ConvertToModel(&iUser)

	err := m.Orm.Save(&iUser).Error

	if err != nil {
		iUserDTO.ID = iUser.ID
		iUserDTO.Password = ""
	}

	return err
}

func (m *UserDao) GetMyUser(uid uint) (model.AuthUser, error) {
	var iUser model.User
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).First(&iUser).Error
	if err != nil {
		return model.AuthUser{}, err
	}
	var iAuthUser model.AuthUser
	iAuthUser = iUser.InAuth()
	return iAuthUser, nil
}

func (m *UserDao) UpdateUsername(uid uint, UpdateUsername string) error {
	err := m.Orm.
		Model(&model.User{}).
		Where("id = ?", uid).
		Update("username", UpdateUsername).
		Error
	return err
}

func (m *UserDao) UpdateSignature(signature string, uid uint) error {
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).Update("signature", signature).Error
	return err
}

func (m *UserDao) UpdateEmail(email string, uid uint) error {
	err := m.Orm.Model(&model.User{}).
		Where("id = ?", uid).
		Update("email", email).
		Error
	return err
}

func (m *UserDao) UpdateAvatar(fid uint, uid uint) error {
	return m.Orm.Model(&model.User{}).Where("id = ? ", uid).Update("avatar_id", fid).Error
}

func (m *UserDao) UpdateNickname(nickname string, uid uint) error {
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).Update("nickname", nickname).Error
	return err
}

func (m *UserDao) CheckPasswordRepeat(oldPassword string, newPassword string, uid uint) error {
	var iUser model.User
	//判断该用户是否存在
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).Find(&iUser).Error
	if err != nil {
		return errors.New("user dont exist")
	} else {
		//判断该原密码是否匹配
		if !utils.CompareHashAndPassword(iUser.Password, oldPassword) {
			return errors.New("old_password is incorrect")
		} else {
			//判断原密码与新密码是否相等
			if strings.Compare(iUser.Password, newPassword) == 0 {
				return errors.New("new password cannot be the same as the old password")
			} else {
				//检验完成后进行密码加密
				newPassword, err = utils.Encrypt(newPassword)
				if err != nil {
					return errors.New("password encryption failed")
				}
				err = m.Orm.Model(&model.User{}).Where("id = ?", uid).Update("password", newPassword).Error
				return err
			}
		}
	}
}

func (m *UserDao) UpdatePassword(newPassword string, uid uint) error {
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).Update("password", newPassword).Error
	return err
}

func (m *UserDao) FindUserById(keyword string, paginate dto.Paginate) ([]model.SafeUser, int64, error) {
	var iSafeUsers []model.SafeUser
	var iUsers []model.User
	var nTotal int64
	err := m.Orm.Model(&model.User{}).Where("username like ? or nickname like ?", "%"+keyword+"%", "%"+keyword+"%").
		Scopes(m.Paginate(paginate)).
		Find(&iUsers).
		Offset(-1).Limit(-1).
		Count(&nTotal).
		Error
	if err != nil {
		return nil, 0, err
	}
	for _, user := range iUsers {
		authUser := user.InAuth()
		safeUser := authUser.InPublic()
		iSafeUsers = append(iSafeUsers, safeUser)
	}
	return iSafeUsers, nTotal, nil
}
