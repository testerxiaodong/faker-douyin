package service

import (
	"context"
	"errors"
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/middleware/rabbitmq"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
	"fmt"
)

type LikeServiceImpl struct {
	DateRepo                 *dao.DataRepo
	UserLikeVideoRabbitMQ    *rabbitmq.UserLikeVideoRabbitMQ
	VideoLikedByUserRabbitMQ *rabbitmq.VideoLikedByUserRabbitMQ
	UserService              UserService
	VideoService             VideoService
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
		// 1. 操作数据库
		resultInfo, err := l.DateRepo.Db.Like.Where(l.DateRepo.Db.Like.UserID.Eq(userId), l.DateRepo.Db.Like.VideoID.Eq(videoId)).Delete()
		if err != nil {
			log.AppLogger.Error(err.Error())
			return err
		}
		if resultInfo.RowsAffected == 0 {
			return errors.New("用户未点赞该视频，取消点赞失败")
		}
		log.AppLogger.Info(fmt.Sprintf("resultInfo: %v", resultInfo))
		// 2.1 操作redis：删除用户点赞的视频
		result, err := l.DateRepo.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, userId)).Result()
		if err != nil {
			log.AppLogger.Error(err.Error())
			// 删除整个key
			l.DateRepo.Rdb.Del(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, userId))
		}
		// 存在的时候，才删除，不存在，可能是key过期了
		if result > 0 {
			_, err = l.DateRepo.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, userId), videoId).Result()
			if err != nil {
				go l.UserLikeVideoRabbitMQ.Publish(rabbitmq.LikeMessage{
					LikeDealType: consts.CommentDelMode,
					UserId:       userId,
					VideoId:      videoId,
				})
			}
		}
		// 2.2 操作redis：删除视频被点赞的用户
		result, err = l.DateRepo.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoLikedByUserPrefix, videoId)).Result()
		if err != nil {
			log.AppLogger.Error(err.Error())
			// 删除整个key
			l.DateRepo.Rdb.Del(context.Background(), dao.GetRedisKeyByPrefix(consts.VideoLikedByUserQueue, videoId))
		}
		if result > 0 {
			_, err = l.DateRepo.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoLikedByUserPrefix, videoId), userId).Result()
			if err != nil {
				go l.VideoLikedByUserRabbitMQ.Publish(rabbitmq.LikeMessage{
					LikeDealType: consts.LikeCancelMode,
					UserId:       userId,
					VideoId:      videoId,
				})
			}
		}
		// 点赞逻辑
	} else if actionType == 1 {
		// 1. 先操作数据库
		var like entity.Like
		like.UserID = userId
		like.VideoID = videoId
		err := l.DateRepo.Db.Like.Create(&like)
		if err != nil {
			return err
		}
		// 2.1 新增redis数据：用户点赞视频
		result, err := l.DateRepo.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, userId)).Result()
		if err != nil || result == 0 {
			// 	重新构建缓存
		}
		_, err = l.DateRepo.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, userId), videoId).Result()
		// 3.1 新增失败，异步消息补偿
		if err != nil {
			go l.UserLikeVideoRabbitMQ.Publish(rabbitmq.LikeMessage{
				LikeDealType: consts.LikeConfirmMode,
				UserId:       userId,
				VideoId:      videoId,
			})
		}
		// 2.2 新增redis数据：视频被点赞用户
	}
	return nil
}

func (l *LikeServiceImpl) GetUserLikeList(userId int64) ([]*response.VideoInfoRes, error) {
	//TODO implement me
	panic("implement me")
}
