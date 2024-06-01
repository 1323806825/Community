package dto

import (
	"blog/global"
	"blog/model"
	"blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var TokenHeader = "Authorization"
var TokenPrefix = "TokenPrefix"

type UserLoginDTO struct {
	//默认错误信息message， 针对某一个tag的错误单独给信息 tag_err
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (m *UserLoginDTO) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("alphanumunderscore", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(fl.Field().String())
	})
	return validate.Struct(m)
}

type UserPublicAuthDTO struct {
	IP string
	UA string
}

func (m *UserPublicAuthDTO) PutAuth(c *gin.Context) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	m.IP = clientIP
	m.UA = userAgent
}

type UserAddDTO struct {
	ID       uint
	Avatar   string
	Email    string `json:"email" form:"email" validate:"required,email"`
	Username string `json:"username" form:"username" validate:"required,alphanumunderscore,min=3,max=30"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=32"`
}

func (m *UserAddDTO) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("alphanumunderscore", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(fl.Field().String())
	})
	return validate.Struct(m)
}

func (m *UserAddDTO) ConvertToModel(iUser *model.User) {
	iUser.Password = m.Password
	iUser.Email = m.Email
	iUser.Username = m.Username
	iUser.Nickname = m.Username
}

type UserAuthDTO struct {
	Uid uint
	IP  string
}

func (m *UserAuthDTO) PutAuth(c *gin.Context) {
	IP := c.ClientIP()
	token := c.GetHeader(TokenHeader)
	token = token[len(TokenPrefix):]

	keys, _ := global.RedisClient.GetKeysAndValue("*" + token + "*")

	var authTokenValue string
	for k := range keys {
		authTokenValue = keys[k]
	}

	//解析token
	TokenClaim, err := utils.ParserToken(authTokenValue)
	if err != nil {
		return
	}

	m.IP = IP
	m.Uid = TokenClaim.UID
}

type UserUpdateUsernameDTO struct {
	Username string `json:"username" form:"username" validate:"required , alphanumunderscore, min=3,max=30"`
}

func (m *UserUpdateUsernameDTO) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("alphanumunderscore", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(fl.Field().String())
	})
	return validate.Struct(m)
}

type UserUpdateSignatureDTO struct {
	Signature string `json:"signature" form:"signature" validate:"required,max=30"`
}

func (m *UserUpdateSignatureDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *UserUpdateSignatureDTO) ConvertToModel(user *model.User) {
	user.Signature = m.Signature
}

type UserUpdateEmailDTO struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

func (m *UserUpdateEmailDTO) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("alphanumunderscore", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(fl.Field().String())
	})
	return validate.Struct(m)
}

type InnerFileInfo struct {
	FileName string
	FilePath string
	FileSize int64
	FileType string
}

type UpdateNicknameDTO struct {
	Nickname string `json:"nickname" form:"nickname" validate:"required,min=1,max=30"`
}

func (m *UpdateNicknameDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required,min=6,max=32"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required,min=6,max=32"`
}

func (m *UpdatePasswordDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type UserSearchDTO struct {
	Keyword string `json:"keyword" form:"keyword" `
	Paginate
}

func (m *UserSearchDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}
