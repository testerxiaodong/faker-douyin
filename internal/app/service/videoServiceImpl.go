package service

import (
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
	"faker-douyin/internal/pkg/const"
	utils2 "faker-douyin/internal/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
	dataRepo       *dao.DataRepo
	UserService    UserService
	CommentService CommentService
}

func NewVideoService(userService UserService, commentService CommentService) VideoService {
	return &VideoServiceImpl{
		UserService:    userService,
		CommentService: commentService,
	}
}

func (v *VideoServiceImpl) Feed(lastTime time.Time) ([]response.VideoInfoRes, time.Time, error) {
	if v.UserService == nil {
		fmt.Println("VideoService is nil")
		return nil, time.Time{}, nil
	}
	videos := make([]response.VideoInfoRes, 0, _const.VideoCount)
	tableVideos, err := v.dataRepo.Db.Video.Where(v.dataRepo.Db.Video.CreatedAt.Lt(lastTime)).Limit(_const.VideoCount).Find()
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
		fmt.Println("UserService.GetTableUserById 成功", videoAuthor)
		if err != nil {
			fmt.Printf("UserServie.GetTableUserById 失败，user_id: %d", video.AuthorID)
			return videos, time.Time{}, err
		}
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
		log.Println("open upload file failed", err)
		return video, err
	}
	// 上传视频
	videoName := uuid.NewString()
	err = utils2.VideoFTP(file, videoName)
	if err != nil {
		log.Println("上传视频失败：", err)
		return video, err
	}
	// 调用ffmpeg生成截图
	imageName := uuid.NewString()
	utils2.Ffchan <- utils2.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}
	var newVideo entity.Video
	newVideo.AuthorID = userId
	newVideo.Title = title
	newVideo.PlayURL = _const.PlayUrlPrefix + videoName + ".mp4"
	newVideo.CoverURL = _const.CoverUrlPrefix + imageName + ".jpg"
	err = v.dataRepo.Db.Video.Create(&newVideo)
	if err != nil {
		fmt.Println("dao.Video.Create failed, video_info: ", newVideo)
		return response.PublishVideoRes{}, err
	}
	video.Video = newVideo
	if err != nil {
		log.Println("新增视频数据失败：", err)
		return video, err
	}
	return video, nil
}

func (v *VideoServiceImpl) List(userId int64) ([]response.VideoInfoRes, error) {
	var userVideoList []response.VideoInfoRes
	if v.UserService == nil {
		return userVideoList, nil
	}
	videos, err := v.dataRepo.Db.Video.Where(v.dataRepo.Db.Video.AuthorID.Eq(userId)).Find()
	if err != nil {
		fmt.Println("dao.Video.Where(dao.Video.AuthorID.Eq(userId)) 失败，user_id：", userId)
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
			fmt.Println("UserService.GetTableUserById 失败，user_id：", video.AuthorID)
			return userVideoList, err
		}
		userVideoInfo.Author = *user
		userVideoList = append(userVideoList, userVideoInfo)
	}
	return userVideoList, nil
}

func (v *VideoServiceImpl) GetVideoIdList(userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
