package router

import (
	"blog/api"
	"github.com/gin-gonic/gin"
)

func InitSubscribeRouters() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		SubscribeApi := api.NewSubscribeApi()
		LikeApi := api.NewLikeApi()

		rgAuthSubscribe := rgAuth.Group("subscribe")
		{
			rgAuthSubscribe.POST("/add", SubscribeApi.AddSubscribe)
			rgAuthSubscribe.POST("/delete", SubscribeApi.DeleteSubscribe)
		}

		rgAuthLike := rgAuth.Group("like")
		{
			rgAuthLike.POST("/add", LikeApi.AddLike)
			rgAuthLike.POST("/add", LikeApi.DeleteLike)
		}
	})
}
