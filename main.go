package main

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/router"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitFilter()
	global.LoadConfig()
	dao.Init()
	utils.InitFtp()
	utils.InitSSH()
	gin.SetMode(global.Config.Server.AppMode)
	r := gin.New()
	//logger, _ := zap.NewProduction()
	//r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	//r.Use(ginzap.RecoveryWithZap(logger, true))
	router.InitRouter(r)
	_ = r.Run(global.Config.Server.HttpPort)
}
