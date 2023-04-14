package service

import (
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
)

type CommentService interface {
	/*
		供videoService调用获取视频点赞数供前端展示
	*/
	Count(videoId uint64) (uint64, error)
	/*
		CommentService内部调用
	*/
	CommentInfo(commentId uint64) (entity.TableComment, error)
	InsertComment(userId uint64, videoId uint64, commentContent string) (entity.TableComment, error)
	DeleteComment(commentId uint64) error
	CommentList(videoId uint64) ([]response.CommentInfoRes, error)
}
