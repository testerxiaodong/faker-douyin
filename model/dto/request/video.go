package request

import (
	"faker-douyin/model/common"
)

type PublishVideoReq struct {
	Title string `json:"title"`
}

type VideoFeedReq struct {
	LastTime common.LocalTime `json:"last_time"`
}

type VideoListReq struct {
	UserId uint64 `json:"user_id"`
}
