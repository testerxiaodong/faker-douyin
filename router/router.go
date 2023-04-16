package router

import (
	v1 "faker-douyin/api/v1"
	"faker-douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin/v1")

	userController := v1.UserController{}
	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", middleware.Auth(), userController.UserInfo)

	videoController := v1.VideoController{}
	apiRouter.POST("/video/publish/", middleware.Auth(), videoController.Publish)
	apiRouter.GET("/video/feed/", videoController.Feed)
	apiRouter.GET("/video/publish/list/", videoController.List)

	commentController := v1.CommentController{}
	apiRouter.POST("/comment/action/", middleware.Auth(), commentController.CommentAction)
	apiRouter.GET("/comment/list/", commentController.CommentList)
}
