package request

// CommentActionReq 新增评论或者删除评论
type CommentActionReq struct {
	ActionType     uint64 `json:"action_type" binding:"required"` // 1表示删除 2表示新增
	CommentId      int64  `json:"comment_id"`
	VideoId        int64  `json:"video_id"`
	CommentContent string `json:"comment_content"`
}

// CommentListReq 获取视频评论评论列表
type CommentListReq struct {
	VideoId int64 `json:"video_id" binding:"required"`
}
