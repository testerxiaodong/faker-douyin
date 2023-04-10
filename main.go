package main

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/router"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	global.LoadConfig()
	dao.Init()
	utils.InitFtp()
	utils.InitSSH()
	gin.SetMode(global.Config.Server.AppMode)
	r := gin.Default()
	router.InitRouter(r)
	_ = r.Run(global.Config.Server.HttpPort)
}
