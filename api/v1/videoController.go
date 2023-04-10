package v1

import (
	"faker-douyin/model/common"
	"faker-douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Publish(c *gin.Context) {
	// userId在jwt中间件中已被存入Context
	userId, err := strconv.ParseInt(c.GetString("userId"), 10, 64)
	if err != nil {
		log.Println("parse userId failed", err)
	}
	title := c.PostForm("title")
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("服务端接收视频文件失败", err)
		common.FailWithMessage(err.Error(), c)
		return
	}
	vsi := service.VideoServiceImpl{}
	err = vsi.Publish(file, uint64(userId), title)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("upload success", c)
}
