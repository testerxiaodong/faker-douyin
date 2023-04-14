package v1

import (
	"faker-douyin/model/common"
	"faker-douyin/model/dto/request"
	"faker-douyin/service"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CommentAction(c *gin.Context) {
	var commentActionReq request.CommentActionReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&commentActionReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	id, _ := c.Get("userId")
	// 从上下文中获取用户id
	userId, err := strconv.Atoi(id.(string))
	// 删除逻辑
	if commentActionReq.ActionType == 1 {
		csi := service.CommentServiceImpl{}
		// 获取评论信息
		comment, err := csi.CommentInfo(commentActionReq.CommentId)
		if err != nil {
			common.FailWithMessage(err.Error(), c)
			return
		}
		// 当前评论非当前用户发表，无权删除
		if comment.UserId != uint64(userId) {
			common.FailWithMessage("current comment is not created by this user", c)
			return
		}
		// 删除评论
		err = csi.DeleteComment(commentActionReq.CommentId)
		if err != nil {
			common.FailWithMessage(err.Error(), c)
			return
		}
		common.OkWithMessage("删除评论成功", c)
		return
	}
	// 新增逻辑
	if commentActionReq.ActionType == 2 {
		// 敏感词判断
		result, _ := utils.Filter.FindIn(commentActionReq.CommentContent)
		if result {
			common.FailWithMessage("评论包含敏感词，操作失败", c)
			return
		}
		csi := service.CommentServiceImpl{}
		// 插入评论
		comment, err := csi.InsertComment(uint64(userId), commentActionReq.VideoId, commentActionReq.CommentContent)
		if err != nil {
			common.FailWithMessage(err.Error(), c)
			return
		}
		common.OkWithDetailed(comment, "新增评论成功", c)
		return
	}
	common.FailWithMessage("action type error", c)
}

func CommentList(c *gin.Context) {
	var commentListReq request.CommentListReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&commentListReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	csi := service.CommentServiceImpl{}
	// 获取评论列表
	commentList, err := csi.CommentList(commentListReq.VideoId)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(commentList, c)
}
