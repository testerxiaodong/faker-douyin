package service

import (
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/model/dto/response"
	"fmt"
)

type LikeServiceImpl struct {
	DateRepo     *dao.DataRepo
	UserService  UserService
	VideoService VideoService
}

func (l *LikeServiceImpl) GetVideoLikeCount(videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LikeServiceImpl) UserIsLikeVideo(userId int64, videoId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LikeServiceImpl) GetUserTotalLikedCount(userId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LikeServiceImpl) UserLikeAction(userId int64, videoId int64, actionType int64) error {
	// 判断当前用户是否存在
	_, err := l.UserService.GetByID(userId)
	if err != nil {
		log.AppLogger.Error(err.Error())
		return err
	}
	// 判断当前视频是否存在
	_, err = l.VideoService.GetVideoById(videoId)
	if err != nil {
		log.AppLogger.Error(err.Error())
		return err
	}
	// 取消点赞
	if actionType == 1 {
		// 操作数据库
		resultInfo, err := l.DateRepo.Db.Like.Where(l.DateRepo.Db.Like.UserID.Eq(userId), l.DateRepo.Db.Like.VideoID.Eq(videoId)).Delete()
		if err != nil {
			log.AppLogger.Error(err.Error())
			return err
		}
		log.AppLogger.Info(fmt.Sprintf("resultInfo: %v", resultInfo))
		// 操作redis
	}
	return nil
}

func (l *LikeServiceImpl) GetUserLikeList(userId int64) ([]*response.VideoInfoRes, error) {
	//TODO implement me
	panic("implement me")
}
