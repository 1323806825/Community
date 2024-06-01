package router

import (
	"blog/api"
	"github.com/gin-gonic/gin"
)

func InitUserRouters() {
	//注册路由
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		UserApi := api.NewUserApi()

		rgPublicUser := rgPublic.Group("user")
		{
			rgPublicUser.POST("/login", UserApi.Login)
			rgPublicUser.POST("/register", UserApi.AddUser)
		}
		rgAuthUser := rgAuth.Group("user")
		{
			rgAuthUser.GET("/find", UserApi.FindUserByKeyword)
			rgAuthUser.GET("/renew_token", UserApi.RefreshToken)
			rgAuthUser.GET("/my", UserApi.GetMyUser)

			rgAuthUser.PUT("/update/username", UserApi.UpdateUsername)
			rgAuthUser.PUT("/update/password", UserApi.UpdatePassword)
			rgAuthUser.PUT("/update/signature", UserApi.UpdateSignature)
			rgAuthUser.PUT("/update/nickname", UserApi.UpdateNickname)
			rgAuthUser.PUT("/update/email", UserApi.UpdateEmail)
			rgAuthUser.PUT("/update/avatar", UserApi.UpdateAvatar)
		}
	})
}
