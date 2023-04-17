package v1

import (
	"faker-douyin/internal/app/model/common"
	"faker-douyin/internal/app/model/dto/request"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/service"
	"faker-douyin/internal/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type VideoController struct {
	v service.VideoService
}

func (v *VideoController) Publish(c *gin.Context) {
	// userId在jwt中间件中已被存入Context
	userId, err := strconv.ParseInt(c.GetString("userId"), 10, 64)
	if err != nil {
		log.Println("parse userId failed", err)
	}
	// 从表单中获取视频标题
	title := c.PostForm("title")
	// 视频标题不能为空
	if len(title) == 0 {
		common.FailWithMessage("视频标题不能为空", c)
		return
	}
	// 视频标题不能含有敏感词
	result, _ := utils.Filter.FindIn(title)
	if result {
		common.FailWithMessage("标题含有敏感词", c)
		return
	}
	// 获取视频文件头
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("服务端接收视频文件失败", err)
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 插入视频
	var video response.PublishVideoRes
	video, err = v.v.Publish(file, userId, title)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithDetailed(video, "upload success", c)
}

func (v *VideoController) Feed(c *gin.Context) {
	var videoFeedReq request.VideoFeedReq
	// 请求参数绑定与校验
	err := c.ShouldBindJSON(&videoFeedReq)
	if err != nil {
		log.Println("绑定参数失败")
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 获取视频流
	videoInfo, lastTime, err := v.v.Feed(time.Time(videoFeedReq.LastTime))
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(response.VideoFeedRes{VideosInfo: videoInfo, LastTime: lastTime}, c)
}

func (v *VideoController) List(c *gin.Context) {
	var VideoListReq request.VideoListReq
	// 请求参数绑定与校验
	err := c.ShouldBindJSON(&VideoListReq)
	if err != nil {
		fmt.Println("参数绑定失败")
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 获取视频列表
	videoList, err := v.v.List(VideoListReq.UserId)
	if err != nil {
		fmt.Println("VideoService.List 失败，user_id：", VideoListReq.UserId)
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(response.VideoListRes{VideosInfo: videoList}, c)
}
