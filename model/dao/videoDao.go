package dao

import (
	"faker-douyin/global"
	"faker-douyin/model/entity"
	"time"
)

// GetVideosByAuthorId 根据作者id获取作者所有视频信息
func GetVideosByAuthorId(authorId uint64) ([]entity.TableVideo, error) {
	var videos []entity.TableVideo
	if err := Db.Where("author_id = ?", authorId).Find(&videos).Error; err != nil {
		return videos, err
	}
	return videos, nil
}

// GetVideoById 根据视频id获取视频信息
func GetVideoById(videoId uint64) (entity.TableVideo, error) {
	var video entity.TableVideo
	if err := Db.Find(&video, videoId).Error; err != nil {
		return video, err
	}
	return video, nil
}

// GetVideosByLastTime 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]entity.TableVideo, error) {
	var videos []entity.TableVideo
	if err := Db.Where("created_at < ?", lastTime).Order("created_at desc").Limit(global.VideoCount).Find(&videos).Error; err != nil {
		return videos, err
	}
	return videos, nil
}

// InsertTableVideo 插入视频数据
func InsertTableVideo(title string, videoName string, imageName string, authorId uint64) (entity.TableVideo, error) {
	var video entity.TableVideo
	video.Title = title
	video.PlayUrl = global.PlayUrlPrefix + videoName + ".mp4"
	video.CoverUrl = global.CoverUrlPrefix + imageName + ".jpg"
	video.AuthorId = authorId
	if err := Db.Save(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}
