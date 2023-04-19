package app

import (
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/router"
	"github.com/gin-gonic/gin"
)

func InitGinEngine(router *router.Router, config *config.Config) *gin.Engine {
	gin.SetMode(config.Server.Mode)
	app := gin.New()
	router.RegisterAPI(app)
	return app
}
