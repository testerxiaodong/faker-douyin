package v1

import (
	"faker-douyin/internal/app/model/common"
	"faker-douyin/internal/app/model/dto/request"
	"faker-douyin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/importcjj/sensitive"
	"strconv"
)

type CommentController struct {
	Filter         *sensitive.Filter
	CommentService service.CommentService
}

// CommentAction POST /douyin/v1/comment/action/ 发表评论和删除评论
func (cc *CommentController) CommentAction(c *gin.Context) {
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
	// 1. 删除逻辑
	if commentActionReq.ActionType == 1 {
		// 删除评论
		err = cc.CommentService.DeleteComment(int64(userId), commentActionReq.CommentId)
		if err != nil {
			common.FailWithMessage(err.Error(), c)
			return
		}
		common.OkWithMessage("删除评论成功", c)
		return
	}
	// 2. 新增逻辑
	if commentActionReq.ActionType == 2 {
		// 敏感词判断
		result, _ := cc.Filter.FindIn(commentActionReq.CommentContent)
		if result {
			common.FailWithMessage("评论包含敏感词，操作失败", c)
			return
		}
		// 插入评论
		comment, err := cc.CommentService.InsertComment(int64(userId), commentActionReq.VideoId, commentActionReq.CommentContent)
		if err != nil {
			common.FailWithMessage(err.Error(), c)
			return
		}
		common.OkWithDetailed(comment, "新增评论成功", c)
		return
	}
	common.FailWithMessage("action type error", c)
}

// CommentList GET /douyin/v1/comment/list/ 获取评论列表
func (cc *CommentController) CommentList(c *gin.Context) {
	var commentListReq request.CommentListReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&commentListReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 获取评论列表
	commentList, err := cc.CommentService.CommentList(commentListReq.VideoId)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(commentList, c)
}
