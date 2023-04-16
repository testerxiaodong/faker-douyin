package response

import (
	"faker-douyin/model/entity"
	"time"
)

type PublishVideoRes struct {
	Video entity.Video `json:"video"`
}

type VideoInfoRes struct {
	Video        entity.Video
	Author       UserInfoRes `json:"author"`
	CommentCount int64       `json:"comment_count"`
}

type VideoFeedRes struct {
	VideosInfo []VideoInfoRes `json:"videos_info"`
	LastTime   time.Time      `json:"last_time"`
}

type VideoListRes struct {
	VideosInfo []VideoInfoRes `json:"videos_info"`
}
