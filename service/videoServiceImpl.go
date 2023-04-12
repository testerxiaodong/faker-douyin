package service

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
	"faker-douyin/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
	UserService
}

func (v *VideoServiceImpl) Feed(lastTime time.Time) ([]response.VideoInfoRes, time.Time, error) {
	if v.UserService == nil {
		fmt.Println("VideoService is nil")
		return nil, time.Time{}, nil
	}
	videos := make([]response.VideoInfoRes, 0, global.VideoCount)
	tableVideos, err := dao.GetVideosByLastTime(lastTime)
	fmt.Println("feed videos: ", tableVideos)
	if err != nil {
		fmt.Println("dao.GetVideosByLastTime 失败", err)
		return videos, time.Time{}, err
	}
	if len(tableVideos) == 0 {
		return videos, time.Time{}, nil
	}
	for _, video := range tableVideos {
		videoAuthor, err := v.GetTableUserById(video.AuthorId)
		fmt.Println("UserService.GetTableUserById 成功", videoAuthor)
		if err != nil {
			fmt.Printf("UserServie.GetTableUserById 失败，user_id: %d", video.AuthorId)
			return videos, time.Time{}, err
		}
		var singleVideo response.VideoInfoRes
		singleVideo.Author = videoAuthor
		singleVideo.Video = video
		videos = append(videos, singleVideo)
	}
	return videos, tableVideos[len(videos)-1].CreatedAt, nil
}

func (v *VideoServiceImpl) GetVideo(videoId int64, userId uint64) (entity.TableVideo, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VideoServiceImpl) Publish(data *multipart.FileHeader, userId uint64, title string) (response.PublishVideoRes, error) {
	var video response.PublishVideoRes
	file, err := data.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		log.Println("open upload file failed", err)
		return video, err
	}
	// 上传视频
	videoName := uuid.NewString()
	err = utils.VideoFTP(file, videoName)
	if err != nil {
		log.Println("上传视频失败：", err)
		return video, err
	}
	// 调用ffmpeg生成截图
	imageName := uuid.NewString()
	utils.Ffchan <- utils.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}
	tableVideo, err := dao.InsertTableVideo(title, videoName, imageName, userId)
	video.Video = tableVideo
	if err != nil {
		log.Println("新增视频数据失败：", err)
		return video, err
	}
	return video, nil
}

func (v *VideoServiceImpl) List(userId uint64) ([]response.VideoInfoRes, error) {
	var userVideoList []response.VideoInfoRes
	if v.UserService == nil {
		return userVideoList, nil
	}
	videos, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		fmt.Println("dao.GetVideosByAuthorId 失败，user_id：", userId)
		return userVideoList, err
	}
	if len(videos) == 0 {
		return userVideoList, nil
	}
	for _, video := range videos {
		var userVideoInfo response.VideoInfoRes
		userVideoInfo.Video = video
		user, err := v.UserService.GetTableUserById(video.AuthorId)
		if err != nil {
			fmt.Println("UserService.GetTableUserById 失败，user_id：", video.AuthorId)
			return userVideoList, err
		}
		userVideoInfo.Author = user
		userVideoList = append(userVideoList, userVideoInfo)
	}
	return userVideoList, nil
}

func (v *VideoServiceImpl) GetVideoIdList(userId uint64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
