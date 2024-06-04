package router

import (
	"blog/api"
	"github.com/gin-gonic/gin"
)

func InitQuestionRouters() {
	//注册路由
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		questionApi := api.NewQuestionApi()

		rgPublicQuestion := rgPublic.Group("question")
		//公共组question组绑定
		{
			rgPublicQuestion.GET("/list", questionApi.GetQuestionList)
			rgPublicQuestion.GET("/detail/:id", questionApi.GetQuestionById)
		}
		rgAuthQuestion := rgAuth.Group("question")
		{
			rgAuthQuestion.GET("/my", questionApi.GetMyQuestion)
			rgAuthQuestion.POST("/add", questionApi.AddQuestion)
			rgAuthQuestion.POST("/upload/picture", questionApi.UploadPicture)
			rgAuthQuestion.POST("/update/title", questionApi.UpdateTitle)
			rgAuthQuestion.POST("/update/content", questionApi.UpdateContent)
			rgAuthQuestion.POST("/delete/picture", questionApi.DeleteQuestionPicture)
			rgAuthQuestion.POST("/delete", questionApi.DeleteQuestion)
			//rgAuthQuestion.POST("/add/subscribe", questionApi.AddQuestionSubscribe)
			rgAuthQuestion.POST("/delete/subscribe", questionApi.DeleteQuestionSubscribe)
			rgAuthQuestion.GET("/subscribe/my", questionApi.GetMyQuestionSubscribe)
		}
	})
}
