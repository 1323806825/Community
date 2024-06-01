package main

import (
	"blog/config"
	"blog/global"
	"blog/model"
	"blog/router"
)

func main() {
	config.InitConfig()

	db, err := model.InitDB()
	global.DB = db
	if err != nil {
		panic(err.Error())
	}
	rdRedisClient, err := config.InitRedis()
	global.RedisClient = rdRedisClient
	if err != nil {
		panic(err.Error())
	}

	router.InitRouter()
}
