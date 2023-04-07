package router

import (
	v1 "faker-douyin/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin/v1")
	apiRouter.POST("/user/register/", v1.Register)
}
