package service

import (
	"faker-douyin/model/entity"
	"mime/multipart"
	"time"
)

type VideoService interface {
	// Feed
	// 通过传入时间戳，当前用户的id，返回对应的视频切片数组，以及视频数组中最早的发布时间
	Feed(lastTime time.Time, userId uint64) ([]entity.TableVideo, time.Time, error)

	// GetVideo
	// 传入视频id获得对应的视频对象
	GetVideo(videoId int64, userId uint64) (entity.TableVideo, error)

	// Publish
	// 将传入的视频流保存在文件服务器中，并存储在mysql表中
	// 5.23 加入title
	Publish(data *multipart.FileHeader, userId uint64, title string) error

	// List
	// 通过userId来查询对应用户发布的视频，并返回对应的视频切片数组
	List(userId uint64, curId uint64) ([]entity.TableVideo, error)

	// GetVideoIdList
	// 通过一个作者id，返回该用户发布的视频id切片数组
	GetVideoIdList(userId uint64) ([]int64, error)
}