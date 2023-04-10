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
const PlayUrlPrefix = "http://192.168.18.3/"
const CoverUrlPrefix = "http://192.168.18.3/images/"

const MaxMsgCount = 100
const SSHHeartbeatTime = 10 * 60

const ValidComment = 0   //评论状态：有效
const InvalidComment = 1 //评论状态：取消
const DateTime = "2006-01-02 15:04:05"

const IsLike = 0     //点赞的状态
const Unlike = 1     //取消赞的状态
const LikeAction = 1 //点赞的行为
const Attempts = 3   //操作数据库的最大尝试次数

const DefaultRedisValue = -1 //redis中key对应的预设值，防脏读