package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseJson struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total"`
}

func BuildStatus(resp ResponseJson, code int) int {
	if resp.Code == 0 {
		return code
	}
	return resp.Code
}

func Unauthorized(ctx *gin.Context, resp ResponseJson) {
	ctx.AbortWithStatusJSON(BuildStatus(resp, http.StatusUnauthorized), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	ctx.AbortWithStatusJSON(BuildStatus(resp, http.StatusOK), resp)
}

func OK(ctx *gin.Context, resp ResponseJson) {
	ctx.AbortWithStatusJSON(BuildStatus(resp, http.StatusOK), resp)
}

func ServerFail(ctx *gin.Context, resp ResponseJson) {
	ctx.AbortWithStatusJSON(BuildStatus(resp, http.StatusInternalServerError), resp)
}
