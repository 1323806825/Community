package utils

import (
	"encoding/json"
	"errors"
	uuid2 "github.com/google/uuid"
	"regexp"
)

func GenerateUUID() string {
	uuid := uuid2.New()
	return uuid.String()
}

type TokenClaim struct {
	UID uint `json:"uint"`
	//登录平台
	LoginPlatform string `json:"login_platform"`
	UA            string `json:"ua"`
}

func ParserToken(tokenStr string) (TokenClaim, error) {
	tokenClaims := TokenClaim{}
	err := json.Unmarshal([]byte(tokenStr), &tokenClaims)
	if err != nil {
		return TokenClaim{}, errors.New("非法token: " + err.Error())
	}
	return tokenClaims, nil

}

func UALoginPlatform(ua string) string {
	isPhone := regexp.MustCompile(`(iphone|Android)`)
	isPad := regexp.MustCompile(`(ipad)`)
	isWatch := regexp.MustCompile(`(Watch)`)
	if isPhone.MatchString(ua) {
		return "Phone"
	}
	if isPad.MatchString(ua) {
		return "Pad"
	}
	if isWatch.MatchString(ua) {
		return "Watch"
	}

	return "Web"
}

func GenerateTokenClaimsFromRaw(uid uint, loginPlatform, ua string) TokenClaim {
	return TokenClaim{
		UID:           uid,
		LoginPlatform: loginPlatform,
		UA:            ua,
	}
}

func (m *TokenClaim) ToJson() string {
	tokenStr, _ := json.Marshal(m)
	return string(tokenStr)
}
