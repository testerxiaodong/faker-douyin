package router

import (
	v1 "faker-douyin/api/v1"
	"faker-douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin/v1")
	apiRouter.POST("/user/register/", v1.Register)
	apiRouter.POST("/user/login/", v1.Login)
	apiRouter.GET("/user/", middleware.Auth(), v1.UserInfo)
	apiRouter.POST("/video/publish/", middleware.Auth(), v1.Publish)
	apiRouter.GET("/video/feed/", v1.Feed)
	apiRouter.GET("/video/publish/list/", v1.List)
	apiRouter.POST("/comment/action/", middleware.Auth(), v1.CommentAction)
	apiRouter.GET("/comment/list/", v1.CommentList)
}
