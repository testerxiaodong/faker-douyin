package response

import (
	"faker-douyin/model/entity"
	"time"
)

type PublishVideoRes struct {
	Video entity.TableVideo `json:"video"`
}

type VideoInfoRes struct {
	Video  entity.TableVideo
	Author UserInfoRes `json:"author"`
}

type VideoFeedRes struct {
	VideosInfo []VideoInfoRes `json:"videos_info"`
	LastTime   time.Time      `json:"last_time"`
}

type VideoListRes struct {
	VideosInfo []VideoInfoRes `json:"videos_info"`
}
