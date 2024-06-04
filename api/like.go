package api

import (
	"blog/service"
	"github.com/gin-gonic/gin"
)

type LikeApi struct {
	BaseApi
	Service *service.LikeService
}

func NewLikeApi() LikeApi {
	return LikeApi{
		NewBaseApi(),
		service.NewLikeApi(),
	}
}

func (m *LikeApi) AddLike(c *gin.Context) {

}

func (m *LikeApi) DeleteLike(c *gin.Context) {

}
