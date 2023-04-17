package _const

import "time"

// OneDayOfHours 时间
const (
	OneDayOfHours = 60 * 60 * 24
	OneMinute     = 60 * 1
	OneMonth      = 60 * 60 * 24 * 30
	OneYear       = 365 * 60 * 60 * 24
	ExpireTime    = time.Hour * 48 // 设置Redis数据热度消散时间。
)

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// PlayUrlPrefix 存储的图片和视频的链接
const (
	PlayUrlPrefix  = "http://192.168.18.3/"
	CoverUrlPrefix = "http://192.168.18.3/images/"
)

const (
	MaxMsgCount      = 100
	SSHHeartbeatTime = 10 * 60
)

const (
	HeartbeatTime = 120
)

const Secret = "douyin"
