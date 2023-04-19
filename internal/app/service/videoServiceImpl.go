package service

import (
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
	"faker-douyin/internal/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
	DataRepo       *dao.DataRepo
	FtpClient      *utils.FtpClient
	FfmpegClient   *utils.FfmpegClient
	UserService    UserService
	CommentService CommentService
}

func (v *VideoServiceImpl) Feed(lastTime time.Time) ([]response.VideoInfoRes, time.Time, error) {
	if v.UserService == nil {
		log.AppLogger.Fatal("UserService is nil")
		return nil, time.Time{}, nil
	}
	videos := make([]response.VideoInfoRes, 0, consts.VideoCount)
	tableVideos, err := v.DataRepo.Db.Video.Where(v.DataRepo.Db.Video.CreatedAt.Lt(lastTime)).Limit(consts.VideoCount).Order(v.DataRepo.Db.Video.CreatedAt.Desc()).Find()
	fmt.Println("feed videos: ", tableVideos)
	if err != nil {
		fmt.Println("dao.GetVideosByLastTime 失败", err)
		return videos, time.Time{}, err
	}
	if len(tableVideos) == 0 {
		return videos, time.Time{}, nil
	}
	for _, video := range tableVideos {
		videoAuthor, err := v.UserService.GetByID(video.AuthorID)
		if err != nil {
			log.AppLogger.Error("UserService.GetByID 失败")
			return videos, time.Time{}, err
		}
		log.AppLogger.Info("UserService.GetByID success")
		commentCount, err := v.CommentService.Count(video.ID)
		if err != nil {
			fmt.Println(err)
		}
		var singleVideo response.VideoInfoRes
		singleVideo.Author = *videoAuthor
		singleVideo.Video = *video
		singleVideo.CommentCount = commentCount
		videos = append(videos, singleVideo)
	}
	log.AppLogger.Info("feed success")
	return videos, tableVideos[len(videos)-1].CreatedAt, nil
}

func (v *VideoServiceImpl) GetVideo(videoId int64, userId int64) (entity.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string) (response.PublishVideoRes, error) {
	var video response.PublishVideoRes
	file, err := data.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		log.AppLogger.Error(err.Error())
		return video, err
	}
	// 上传视频
	videoName := uuid.NewString()
	err = v.FtpClient.VideoFTP(file, videoName)
	if err != nil {
		return video, err
	}
	// 调用ffmpeg生成截图
	imageName := uuid.NewString()
	v.FfmpegClient.Ffchan <- utils.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}
	var newVideo entity.Video
	newVideo.AuthorID = userId
	newVideo.Title = title
	newVideo.PlayURL = consts.PlayUrlPrefix + videoName + ".mp4"
	newVideo.CoverURL = consts.CoverUrlPrefix + imageName + ".jpg"
	err = v.DataRepo.Db.Video.Create(&newVideo)
	if err != nil {
		log.AppLogger.Error(err.Error())
		return response.PublishVideoRes{}, err
	}
	video.Video = newVideo
	return video, nil
}

func (v *VideoServiceImpl) List(userId int64) ([]response.VideoInfoRes, error) {
	var userVideoList []response.VideoInfoRes
	if v.UserService == nil {
		log.AppLogger.Error("UserService is nil")
		return userVideoList, nil
	}
	videos, err := v.DataRepo.Db.Video.Where(v.DataRepo.Db.Video.AuthorID.Eq(userId)).Find()
	if err != nil {
		log.AppLogger.Error(err.Error())
		return userVideoList, err
	}
	if len(videos) == 0 {
		return userVideoList, nil
	}
	for _, video := range videos {
		var userVideoInfo response.VideoInfoRes
		userVideoInfo.Video = *video
		user, err := v.UserService.GetByID(video.AuthorID)
		if err != nil {
			log.AppLogger.Error(err.Error())
			return userVideoList, err
		}
		userVideoInfo.Author = *user
		userVideoList = append(userVideoList, userVideoInfo)
	}
	log.AppLogger.Info("get user video list success")
	return userVideoList, nil
}

func (v *VideoServiceImpl) GetVideoIdList(userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
