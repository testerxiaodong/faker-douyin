package request

import (
	"faker-douyin/model/common"
)

type PublishVideoReq struct {
	Title string `json:"title" binding:"required"`
}

type VideoFeedReq struct {
	LastTime common.LocalTime `json:"last_time" binding:"required"`
}

type VideoListReq struct {
	UserId int64 `json:"user_id" binding:"required"`
}
