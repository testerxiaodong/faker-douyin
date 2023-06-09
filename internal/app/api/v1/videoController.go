package v1

import (
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/model/common"
	"faker-douyin/internal/app/model/dto/request"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/service"
	"faker-douyin/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type VideoController struct {
	V service.VideoService
}

func (v *VideoController) Publish(c *gin.Context) {
	// userId在jwt中间件中已被存入Context
	userId, err := strconv.ParseInt(c.GetString("userId"), 10, 64)
	if err != nil {
		log.AppLogger.Error(err.Error())
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
		log.AppLogger.Error(err.Error())
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 插入视频
	var video response.PublishVideoRes
	video, err = v.V.Publish(file, userId, title)
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
		log.AppLogger.Debug(err.Error())
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 获取视频流
	videoInfo, lastTime, err := v.V.Feed(time.Time(videoFeedReq.LastTime))
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
		log.AppLogger.Error(err.Error())
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 获取视频列表
	videoList, err := v.V.List(VideoListReq.UserId)
	if err != nil {
		log.AppLogger.Error(err.Error())
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(response.VideoListRes{VideosInfo: videoList}, c)
	log.AppLogger.Debug("获取用户视频列表成功")
}
