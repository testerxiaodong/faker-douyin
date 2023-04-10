package service

import (
	"faker-douyin/model/dao"
	"faker-douyin/model/entity"
	"faker-douyin/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
}

func (v VideoServiceImpl) Feed(lastTime time.Time, userId uint64) ([]entity.TableVideo, time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceImpl) GetVideo(videoId int64, userId uint64) (entity.TableVideo, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceImpl) Publish(data *multipart.FileHeader, userId uint64, title string) error {
	file, err := data.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		log.Println("open upload file failed", err)
		return err
	}
	// 上传视频
	videoName := uuid.NewString()
	err = utils.VideoFTP(file, videoName)
	if err != nil {
		log.Println("上传视频失败：", err)
		return err
	}
	// 调用ffmpeg生成截图
	imageName := uuid.NewString()
	utils.Ffchan <- utils.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}
	err = dao.InsertTableVideo(title, videoName, imageName, userId)
	if err != nil {
		log.Println("新增视频数据失败：", err)
		return err
	}
	return nil
}

func (v VideoServiceImpl) List(userId uint64, curId uint64) ([]entity.TableVideo, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceImpl) GetVideoIdList(userId uint64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
