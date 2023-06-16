package service

import "faker-douyin/internal/app/model/dto/response"

type (
	LikeService interface {
		/*
		   1.点赞模块自己调用。
		*/

		// GetUserLikeList 获取用户点赞视频列表
		GetUserLikeList(userId int64) ([]*response.VideoInfoRes, error)
		// UserLikeAction 用户点赞/取消点赞
		UserLikeAction(userId int64, videoId int64, actionType int64) error
		/*
		   2.视频模块调用。
		*/
		// GetVideoLikeCount 获取视频点赞数
		GetVideoLikeCount(videoId int64) (int64, error)
		// UserIsLikeVideo 当前用户是否点赞该视频
		UserIsLikeVideo(userId int64, videoId int64) (bool, error)
		// GetUserTotalLikedCount 获取用户被点赞次数
		GetUserTotalLikedCount(userId int64) (int64, error)
	}
)
