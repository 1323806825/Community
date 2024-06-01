package service

import (
	"blog/dao"
	"blog/dto"
	"blog/global"
	"blog/model"
	"blog/utils"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var LoginRedisTokenKeyID = "community:tk-{id}-{token}"
var userService *UserService

type UserService struct {
	BaseService
	Dao *dao.UserDao
}

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			Dao: dao.NewUserDao(),
		}
	}
	return userService
}

var LoginRedisTokenID = "community:tk-{id}-{token}"

func GenerateAndCacheToken(uid uint, loginPlatform string, ua string) (uuid string, expired int64, err error) {
	uuid = utils.GenerateUUID()

	//生成redis的tokenID
	rdKey := utils.GeneralRedisKey(LoginRedisTokenID, "{id}", strconv.Itoa(int(uid)), viper.GetString("custom.prefix"))
	rdKey = utils.GeneralRedisKey(rdKey, "{token}", uuid, nil)

	//把值赋给Claims的结构内传给value
	value := utils.GenerateTokenClaimsFromRaw(uid, loginPlatform, ua)
	//设置过期时间
	expired = time.Now().Add(10 * time.Hour).Unix()
	//
	err = global.RedisClient.Set(rdKey, value.ToJson(), 10*time.Hour)
	if err != nil {
		return "", 0, err
	}
	return uuid, expired, nil
}

func (m *UserService) Login(iUserDTO *dto.UserLoginDTO, authDTO *dto.UserPublicAuthDTO) (model.AuthUser, string, int64, error) {
	var errResult error
	var token string
	var expired int64
	var iUser model.User

	err := iUserDTO.Validate()
	if err != nil {
		return model.AuthUser{}, token, expired, utils.ParseValidateError(err)
	}

	//进入dao层获取username
	iUser, err = userService.Dao.GetUserByName(iUserDTO.Username)
	//匹配密码
	if err != nil && !utils.CompareHashAndPassword(iUser.Password, iUserDTO.Password) {
		errResult = errors.New("incorrect username or password")
	} else {
		//先检查以前登录的情况
		tokens, _ := global.RedisClient.GetKeysAndValue("*:tk-" + strconv.Itoa(int(iUser.ID)) + "-*")
		loginPa := utils.UALoginPlatform(authDTO.UA)

		//删除同一平台的token，实现单点登录
		for k, v := range tokens {
			//先解析token
			tokenClaims, err := ParseToken(v)
			if err != nil {
				continue
			}
			//检查是否同一平台
			if tokenClaims.LoginPlatform == loginPa {
				//删除
				_ = global.RedisClient.Delete(k)
			}
		}
		token, expired, err = GenerateAndCacheToken(iUser.ID, utils.UALoginPlatform(authDTO.UA), authDTO.UA)
		if err != nil {
			errResult = errors.New("unable to create a token for this account")
		}
	}
	return iUser.InAuth(), token, expired, errResult
}

func (m *UserService) AddUser(iUserAddDTO *dto.UserAddDTO) error {
	//dto层检验validate
	if err := iUserAddDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	//dao层名字是否存在
	if _, err := m.Dao.CheckUserNameExist(iUserAddDTO.Username); err != nil {
		return errors.New("username already exists")
	}
	//dao层邮箱是否存在
	if m.Dao.CheckUserEmailExist(iUserAddDTO.Email) {
		return errors.New("email address is tied to another user")
	}
	//dao层添加user信息
	return m.Dao.AddUser(iUserAddDTO)
}

func ParseToken(tokenStr string) (utils.TokenClaim, error) {
	tokenClaim := utils.TokenClaim{}
	err := json.Unmarshal([]byte(tokenStr), &tokenClaim)
	if err != nil {
		return utils.TokenClaim{}, errors.New("非法token： " + err.Error())
	}
	return tokenClaim, nil
}

func (m *UserService) GetMyUser(iUserAuthDTO *dto.UserAuthDTO) (model.AuthUser, error) {
	user, err := m.Dao.GetMyUser(iUserAuthDTO.Uid)
	return user, err
}

func (m *UserService) UpdateUsername(iUserUpdateUsernameDTO *dto.UserUpdateUsernameDTO, iAuthDTO *dto.UserAuthDTO) error {
	if err := iUserUpdateUsernameDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	return m.Dao.UpdateUsername(iAuthDTO.Uid, iUserUpdateUsernameDTO.Username)
}

func (m *UserService) UpdateSignature(iUserUpdateSignatureDTO *dto.UserUpdateSignatureDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := iUserUpdateSignatureDTO.Validate(); err != nil {
		return err
	}
	return m.Dao.UpdateSignature(iUserUpdateSignatureDTO.Signature, iUserAuth.Uid)
}

func (m *UserService) UpdateEmail(iUserUpdateEmailDTO *dto.UserUpdateEmailDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := iUserUpdateEmailDTO.Validate(); err != nil {
		return err
	}
	return m.Dao.UpdateEmail(iUserUpdateEmailDTO.Email, iUserAuth.Uid)
}

func (m *UserService) UpdateAvatar(innerFile dto.InnerFileInfo, iUserAuthDTO *dto.UserAuthDTO) (int, error) {
	file, err := m.Dao.SaveFile(innerFile, iUserAuthDTO.Uid)
	if err != nil {
		return 0, err
	} else {
		return int(file.ID), m.Dao.UpdateAvatar(file.ID, iUserAuthDTO.Uid)
	}

}

func (m *UserService) UpdateNickname(UpdateNicknameDTO *dto.UpdateNicknameDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := UpdateNicknameDTO.Validate(); err != nil {
		return err
	}
	return m.Dao.UpdateNickname(UpdateNicknameDTO.Nickname, iUserAuth.Uid)
}

func (m *UserService) UpdatePassword(iUserUpdatePasswordDTO *dto.UpdatePasswordDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iUserUpdatePasswordDTO.Validate(); err != nil {
		return err
	}

	if err := m.Dao.CheckPasswordRepeat(iUserUpdatePasswordDTO.OldPassword, iUserUpdatePasswordDTO.NewPassword, iUserAuthDTO.Uid); err != nil {
		return err
	}

	return m.Dao.UpdatePassword(iUserUpdatePasswordDTO.NewPassword, iUserAuthDTO.Uid)
}

func (m *UserService) RenewToken(token string, iUserAuthDTO *dto.UserAuthDTO) (string, int64, error) {
	//获取token实体
	oldTokenKey := utils.GeneralRedisKey(LoginRedisTokenKeyID, "{id}", strconv.Itoa(int(iUserAuthDTO.Uid)), viper.GetString("custom.prefix"))
	oldTokenKey = utils.GeneralRedisKey(oldTokenKey, "{token}", token, nil)
	odv, err := global.RedisClient.Get(oldTokenKey)
	if err != nil {
		return "", 0, errors.New("irregular token")
	}
	//断言odv
	if _, ok := odv.(string); !ok {
		return "", 0, errors.New("irregular token")
	}
	//解析token
	tokenClaims, err := utils.ParserToken(odv.(string))
	if err != nil {
		return "", 0, errors.New("irregular token")
	}
	//更新token
	newToken, newExpired, err := GenerateAndCacheToken(iUserAuthDTO.Uid, tokenClaims.LoginPlatform, tokenClaims.UA)
	if err != nil {
		return "", 0, errors.New("token generation error")
	}
	//删除旧token
	_ = global.RedisClient.Delete(oldTokenKey)
	return newToken, newExpired, nil
}

func (m *UserService) FindUserByKeyword(iUserSearchDTO *dto.UserSearchDTO) ([]model.SafeUser, int64, error) {
	if err := iUserSearchDTO.Validate(); err != nil {
		return nil, 0, err
	}
	if iUserSearchDTO.PageNum <= 0 {
		iUserSearchDTO.PageNum = 1
	}
	if iUserSearchDTO.PageSize <= 0 {
		iUserSearchDTO.PageSize = 10
	}
	users, nTotal, err := m.Dao.FindUserById(iUserSearchDTO.Keyword, iUserSearchDTO.Paginate)
	return users, nTotal, err
}
