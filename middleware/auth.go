package middleware

import (
	"blog/api"
	"blog/global"
	"blog/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	ErrorTokenInvalid = 499
	TOKENHEADER       = "Authorization"
	TOKENPREFIX       = "Bearer:"
)

func tokenErr(c *gin.Context, code ...int) {
	nCode := ErrorTokenInvalid
	if len(code) > 0 {
		nCode = code[0]
	}
	api.Unauthorized(c, api.ResponseJson{
		Status:  http.StatusUnauthorized,
		Code:    nCode,
		Message: "authorization invalid",
	})
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader(TOKENHEADER)

		if token == "" || !strings.HasPrefix(token, TOKENPREFIX) {
			//无token或者格式有误
			tokenErr(c)
			return
		}

		token = token[len(TOKENPREFIX):]

		//从redis获取相似token
		keys, err := global.RedisClient.GetKeysAndValue("*" + token + "*")
		if err != nil || len(keys) == 0 {
			tokenErr(c)
			return
		}

		//拿匹配的第一个
		var authTokenValue string
		for key := range keys {
			authTokenValue = keys[key]
			break
		}

		//解析token
		TokenClaims, err := ParseToken(authTokenValue)
		if TokenClaims.UID == 0 {
			tokenErr(c)
			return
		}

		c.Next()
	}
}

func ParseToken(tokenStr string) (utils.TokenClaim, error) {
	tokenClaim := utils.TokenClaim{}
	err := json.Unmarshal([]byte(tokenStr), &tokenClaim)
	if err != nil {
		return utils.TokenClaim{}, errors.New("非法token： " + err.Error())
	}
	return tokenClaim, nil
}
