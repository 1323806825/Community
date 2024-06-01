package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	c := cors.Config{
		//AllowAllOrigins: true,
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "token"},
		AllowMethods:     []string{"PUT", "DELETE", "POST", "HEAD", "OPTION", "PATCH"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}
	return cors.New(c)
}
