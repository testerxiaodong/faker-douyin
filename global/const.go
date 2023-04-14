package global

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
	IsLike            = 0  //点赞的状态
	Unlike            = 1  //取消赞的状态
	LikeAction        = 1  //点赞的行为
	Attempts          = 3  //操作数据库的最大尝试次数
	DefaultRedisValue = -1 //redis中key对应的预设值，防脏读
)
