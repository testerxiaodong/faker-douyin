package consts

// OneDayOfHours 时间
const (
	OneDayOfHours = 60 * 60 * 24
	ExpireTime    = 60 * 60 * 48 // 设置Redis数据热度消散时间。
)

const (
	RedisCommentVideoPrefix     = "comment_id:video_id:"
	RedisVideoCommentPrefix     = "video_id:comment_id:"
	RedisUserLikeVideoPrefix    = "user_id:video_id:"
	RedisVideoLikedByUserPrefix = "video_id:user_id"
)

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// PlayUrlPrefix 存储的图片和视频的链接
const (
	PlayUrlPrefix  = "http://192.168.18.3/videos/"
	CoverUrlPrefix = "http://192.168.18.3/images/"
)

const (
	MaxMsgCount      = 100
	MaxFailCount     = 3
	SSHHeartbeatTime = 10 * 60
)

const (
	CommentDelMode        = 0
	CommentAddMode        = 1
	VideoCommentQueue     = "video_comment"
	CommentVideoQueue     = "comment_video"
	LikeCancelMode        = 0
	LikeConfirmMode       = 1
	VideoLikedByUserQueue = "video_liked_by_user"
	UserLikeVideoQueue    = "user_like_video"
)

const (
	HeartbeatTime = 120
)

const Secret = "douyin"
