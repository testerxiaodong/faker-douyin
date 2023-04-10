package router

import (
	v1 "faker-douyin/api/v1"
	"faker-douyin/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin/v1")
	apiRouter.POST("/user/register/", v1.Register)
	apiRouter.POST("/user/login/", v1.Login)
	apiRouter.POST("/video/publish/", jwt.Auth(), v1.Publish)
}
