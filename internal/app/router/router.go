package router

import (
	"faker-douyin/internal/app/api/v1"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(wire.Struct(new(Router), "*"))

type Router struct {
	UserController    v1.UserController
	VideoController   v1.VideoController
	CommentController v1.CommentController
}

func (r *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("api")
	api1 := g.Group("v1")
	api1.Use(middleware.GinLogger(log.AppLogger))
	api1.Use(middleware.GinRecovery(log.AppLogger, true))
	{
		gUser := api1.Group("user")
		{
			gUser.POST("register", r.UserController.Register)
			gUser.POST("login", r.UserController.Login)
			gUser.GET("", middleware.Auth(), r.UserController.UserInfo)
		}
		gVideo := api1.Group("video")
		{
			gVideo.POST("publish", middleware.Auth(), r.VideoController.Publish)
			gVideo.GET("feed", r.VideoController.Feed)
			gVideo.GET("publish/list", middleware.Auth(), r.VideoController.List)
		}
		gComment := api1.Group("comment")
		{
			gComment.POST("action", middleware.Auth(), r.CommentController.CommentAction)
			gComment.GET("list", r.CommentController.CommentList)
		}
	}
}
