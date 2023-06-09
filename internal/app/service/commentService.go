package service

import (
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
)

type CommentService interface {
	/*
		供videoService调用获取视频点赞数供前端展示
	*/
	Count(videoId int64) (int64, error)
	/*
		CommentService内部调用
	*/
	CommentInfo(commentId int64) (*entity.Comment, error)
	InsertComment(userId int64, videoId int64, commentContent string) (*entity.Comment, error)
	DeleteComment(commentId int64) error
	CommentList(videoId int64) ([]*response.CommentInfoRes, error)
}
