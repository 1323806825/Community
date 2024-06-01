package router

import (
	"blog/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type IfnRegisterRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var gfnRoutes []IfnRegisterRoute

func RegisterRoute(fn IfnRegisterRoute) {
	if fn == nil {
		return
	}
	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	r := gin.Default()

	r.Use(middleware.Cors())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "404 :: not found",
			"data": map[string]interface{}{
				"request": c.Request.URL.Path,
				"method":  c.Request.Method,
				"time":    time.Now().Format("2006-01-02 15:04:05"),
			},
		})

	})

	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")
	rgAuth.Use(middleware.Auth())
	InitBaseRouter()

	for _, fnRegisterRoute := range gfnRoutes {
		//key为  func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)
		//依次注册路由
		fnRegisterRoute(rgPublic, rgAuth)
	}

	//启动监听
	port := "9091"

	//创建服务
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	//r.Run(":9091")

	go func() {
		fmt.Printf("开始监听服务端口:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			//错误非空且不是关闭状态
			fmt.Printf("出错:%s", err.Error())
			return
		}
	}()

	//等着信号
	<-ctx.Done()
	//开始停止的相关操作(5秒超时)
	ctx, ctxShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxShutDown()
	if err := server.Shutdown(ctx); err != nil {
		//处理关闭服务时出错
		return
	}
}

func InitBaseRouter() {
	InitUserRouters()
	InitQuestionRouters()
	InitSubscribeRouters()
	InitCommentRouters()
}
