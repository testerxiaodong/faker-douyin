package main

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/router"
	"faker-douyin/utils"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func main() {
	global.LoadConfig()
	dao.Init()
	utils.InitFtp()
	utils.InitSSH()
	gin.SetMode(global.Config.Server.AppMode)
	r := gin.New()
	logger, _ := zap.NewProduction()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	router.InitRouter(r)
	_ = r.Run(global.Config.Server.HttpPort)
}
