package main

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/router"
	"github.com/gin-gonic/gin"
)

func main() {
	global.LoadConfig()
	dao.Init()
	gin.SetMode(global.Config.Server.AppMode)
	r := gin.Default()
	router.InitRouter(r)
	r.Run(global.Config.Server.HttpPort)
}
