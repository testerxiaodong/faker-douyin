package response

// CommentInfoRes 获取评论列表中单个评论信息
type CommentInfoRes struct {
	Id       uint64      `json:"id"`
	UserInfo UserInfoRes `json:"user_info"` // 评论的用户信息，方便前端展示
	Content  string      `json:"content"`
}

// CommentList 实现sort.Interface接口
type CommentList []CommentInfoRes

func (c CommentList) Len() int {
	return len(c)
}

func (c CommentList) Less(i, j int) bool {
	return c[i].Id > c[j].Id
}

func (c CommentList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
