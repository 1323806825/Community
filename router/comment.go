package router

import (
	"blog/api"
	"github.com/gin-gonic/gin"
)

func InitCommentRouters() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		commentApi := api.NewCommentApi()

		rgPublicComment := rgPublic.Group("/comment")
		{
			rgPublicComment.GET("/list", commentApi.GetCommentList)
		}

		rgAuthComment := rgAuth.Group("/comment")
		{
			rgAuthComment.POST("/post/question", commentApi.CommentPostQuestion)
			rgAuthComment.POST("/post/comment", commentApi.CommentPostComment)
			rgAuthComment.GET("/my", commentApi.GetMyComment)
			rgAuthComment.POST("/upload/picture", commentApi.UploadPicture)
			rgAuthComment.POST("/delete/picture", commentApi.DeleteCommentPicture)
			rgAuthComment.POST("/delete", commentApi.DeleteComment)
			rgAuthComment.GET("/like-count", commentApi.GetMyCommentLikeCount)
		}
	})
}
