package api

import (
	"blog/dto"
	"blog/model"
	"blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserApi struct {
	BaseApi
	Service *service.UserService
}

func NewUserApi() UserApi {
	return UserApi{
		NewBaseApi(),
		service.NewUserService(),
	}

}

func (m *UserApi) Login(c *gin.Context) {
	//检验请求信息valida是否合理
	var iUserLoginDTO dto.UserLoginDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserLoginDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}
	//读取IP,user-agent
	var iPublicAuthDTO dto.UserPublicAuthDTO
	iPublicAuthDTO.PutAuth(c)

	//服务层实现登录
	iAuthUser, token, expired, err := m.Service.Login(&iUserLoginDTO, &iPublicAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	type reData struct {
		IAuthUserInfo model.AuthUser `json:"user"`
		Token         string         `json:"token"`
		Expired       int64          `json:"expired"`
	}

	m.OK(ResponseJson{
		Data: reData{
			Token:         token,
			IAuthUserInfo: iAuthUser,
			Expired:       expired,
		},
	})
}

func (m *UserApi) AddUser(c *gin.Context) {
	var iUserAddDTO dto.UserAddDTO
	err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserAddDTO,
		BindAll: true,
	}).GetError()
	if err != nil {
		return
	}

	var iUserPublicAuth dto.UserPublicAuthDTO
	iUserPublicAuth.PutAuth(c)

	err = m.Service.AddUser(&iUserAddDTO)

	if err != nil {
		m.Fail(ResponseJson{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	type reData struct {
		Username string `json:"username"`
		ID       uint   `json:"id"`
		Email    string `json:"email"`
	}

	m.OK(ResponseJson{
		Data: reData{
			Username: iUserAddDTO.Username,
			ID:       iUserAddDTO.ID,
			Email:    iUserAddDTO.Email,
		},
	})
}

func (m *UserApi) GetMyUser(c *gin.Context) {
	var iUserAuthDTO dto.UserAuthDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserAuthDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	iUserAuthDTO.PutAuth(c)

	iAuthUser, err := m.Service.GetMyUser(&iUserAuthDTO)
	//server层实现GetMyUser

	if err != nil {
		m.Fail(ResponseJson{
			Code:    501,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Message: "get user success",
		Data:    iAuthUser,
	})
}

func (m *UserApi) UpdateUsername(c *gin.Context) {
	var iUserUpdateUsernameDTO dto.UserUpdateUsernameDTO

	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserUpdateUsernameDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdateUsername(&iUserUpdateUsernameDTO, &iUserAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    501,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Message: "update username success",
		Data: struct {
			NewUsername string `json:"new_username"`
		}{
			NewUsername: iUserUpdateUsernameDTO.Username,
		},
	})
}

func (m *UserApi) UpdatePassword(c *gin.Context) {
	var iUserUpdatePasswordDTO dto.UpdatePasswordDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserUpdatePasswordDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdatePassword(&iUserUpdatePasswordDTO, &iUserAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "update password success",
	})
}

func (m *UserApi) RefreshToken(c *gin.Context) {
	var iUserAuthDTO dto.UserAuthDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserAuthDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	iUserAuthDTO.PutAuth(c)

	t := c.GetHeader("Authorization")
	token := strings.ReplaceAll(t, "Bearer:", "")
	newToken, newExpired, err := m.Service.RenewToken(token, &iUserAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	type reData struct {
		Token   string `json:"token"`
		Expired int64  `json:"expired"`
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "refresh token success",
		Data: reData{
			Token:   newToken,
			Expired: newExpired,
		},
	})
}

func (m *UserApi) UpdateSignature(c *gin.Context) {
	var iUserUpdateSignature dto.UserUpdateSignatureDTO

	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserUpdateSignature,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.UpdateSignature(&iUserUpdateSignature, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    501,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "update signature success",
		Data:    iUserUpdateSignature.Signature,
	})
}

func (m *UserApi) UpdateNickname(c *gin.Context) {
	var iUserUpdateNicknameDTO dto.UpdateNicknameDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserUpdateNicknameDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdateNickname(&iUserUpdateNicknameDTO, &iUserAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "update nickname success",
		Data:    iUserUpdateNicknameDTO.Nickname,
	})

}

func (m *UserApi) UpdateEmail(c *gin.Context) {
	var iUserUpdateEmailDTO dto.UserUpdateEmailDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserUpdateEmailDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	//var iUserPublicAuthDTO dto.UserPublicAuthDTO
	//iUserPublicAuthDTO.PutAuth(c)

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdateEmail(&iUserUpdateEmailDTO, &iUserAuthDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "update email success",
		Data:    iUserUpdateEmailDTO.Email,
	})
}

func (m *UserApi) UpdateAvatar(c *gin.Context) {
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	//运行文件目录下uploads文件夹
	uploadsPath := "./upload."

	file, err := m.SaveFile("avatar", uploadsPath, FileUploadLimit{
		MaxSize:   1024 * 1024 * 5,
		MaxCount:  1,
		AllowType: []string{"jpg", "png", "jpeg", "svg"},
	})
	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	fileID, err := m.Service.UpdateAvatar(file[0], &iUserAuth)
	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "update avatar success",
		Data:    fileID,
	})
}

func (m *UserApi) FindUserByKeyword(c *gin.Context) {
	var iUserSearchDTO dto.UserSearchDTO

	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUserSearchDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	iSafeUsers, nTotal, err := m.Service.FindUserByKeyword(&iUserSearchDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "success",
		Data:    iSafeUsers,
		Total:   nTotal,
	})
}
