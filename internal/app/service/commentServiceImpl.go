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
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	DataRepo             *dao.DataRepo
	VideoCommentRabbitMQ *rabbitmq.VideoCommentRabbitMQ
	CommentVideoRabbitMQ *rabbitmq.CommentVideoRabbitMQ
	UserService          UserService
}

func (c *CommentServiceImpl) GetCommentById(commentId int64) (*entity.Comment, error) {
	comment, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.ID.Eq(commentId)).First()
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (c *CommentServiceImpl) GetCommentCountByVideoId(videoId int64) (int64, error) {
	// 1. 先在缓存中查找
	count, err := c.DataRepo.Rdb.SCard(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, videoId)).Result()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	// 2. 缓存中有数据，直接返回
	if count > 0 {
		return count, nil
	}
	// 3. 缓存中不存在，在数据库中找
	cntDao, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.VideoID.Eq(videoId)).Count()
	if err != nil {
		// 查询出错，返回err
		log.AppLogger.Error(err.Error())
		return 0, err
	}
	// 4. 异步更新缓存
	go func() {
		if cntDao > 0 {
			//查询评论id list
			cList, _ := c.DataRepo.Db.Comment.Select(c.DataRepo.Db.Comment.ID).Where(c.DataRepo.Db.Comment.VideoID.Eq(videoId)).Find()
			//设置key值过期时间
			_, err = c.DataRepo.Rdb.Expire(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, videoId),
				time.Duration(consts.ExpireTime+rand.Intn(360))*time.Second).Result()
			if err != nil {
				log.AppLogger.Error(err.Error())
			}
			//评论id循环存入redis
			for _, commentId := range cList {
				c.insertRedisVideoCommentId(videoId, commentId.ID)
			}

			log.AppLogger.Debug("count comment save ids in redis")
		}
	}()
	//5. 返回结果
	return cntDao, nil
}

func (c *CommentServiceImpl) InsertComment(userId int64, videoId int64, commentContent string) (*entity.Comment, error) {
	// 1.判断视频是否存在
	_, err := c.DataRepo.Db.Video.Where(c.DataRepo.Db.Video.ID.Eq(videoId)).First()
	if err != nil {
		return nil, err
	}
	// 2. 更新数据库
	var comment entity.Comment
	comment.VideoID = videoId
	comment.UserID = userId
	comment.CommentContent = commentContent
	err = c.DataRepo.Db.Comment.Create(&comment)
	if err != nil {
		return &comment, err
	}
	// 3. 同步删除缓存
	c.insertRedisVideoCommentId(videoId, comment.ID)
	return &comment, nil
}

func (c *CommentServiceImpl) DeleteComment(userId int64, commentId int64) error {
	// 1. 先查询评论是否存在，存在的话就获取评论对应的videoId，不存在返回error信息
	comment, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.ID.Eq(commentId)).First()
	// 查询失败，或者评论不存在，直接返回error信息
	if err != nil {
		return err
	}
	// 2. 判断该评论是否是该用户发起的
	if comment.UserID != userId {
		return errors.New("current comment is not created by this user")
	}
	// 3. 删除数据库评论信息
	_, err = c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.ID.Eq(commentId)).Delete()
	if err != nil {
		return err
	}
	// 4.1 同步删除redis缓存数据：视频id的评论id
	result, err := c.DataRepo.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, comment.VideoID)).Result()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	if result > 0 {
		_, err := c.DataRepo.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, comment.VideoID), commentId).Result()
		if err != nil {
			// 失败异步补偿机制
			go c.VideoCommentRabbitMQ.Publish(rabbitmq.CommentMessage{
				CommentDealType: consts.CommentDelMode,
				VideoId:         comment.VideoID,
				CommentId:       commentId,
			})
		}
	}
	// 4.2 同步删除redis缓存数据：评论id对应视频id
	result, err = c.DataRepo.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId)).Result()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	if result > 0 {
		_, err := c.DataRepo.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId), comment.VideoID).Result()
		if err != nil {
			// 失败异步补偿机制
			go c.CommentVideoRabbitMQ.Publish(rabbitmq.CommentMessage{
				CommentDealType: consts.CommentDelMode,
				VideoId:         comment.VideoID,
				CommentId:       commentId,
			})
		}
	}
	return nil
}

func (c *CommentServiceImpl) CommentList(videoId int64) ([]*response.CommentInfoRes, error) {
	commentTableList, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.VideoID.Eq(videoId)).Find()
	fmt.Println(commentTableList)
	// 查询失败，返回
	if err != nil {
		return nil, err
	}
	// 评论数为零直接返回，同时防止对空指针进行操作
	if len(commentTableList) == 0 {
		return nil, nil
	}
	// 并发调用UserService，提升性能
	commentList := make([]*response.CommentInfoRes, 0, len(commentTableList))
	var wg sync.WaitGroup
	wg.Add(len(commentTableList))
	for _, commentTable := range commentTableList {
		var oneCommentInfo response.CommentInfoRes
		// 传入循环变量作为临时变量，防止bug
		go func(commentTable *entity.Comment) {
			c.oneComment(commentTable, &oneCommentInfo)
			commentList = append(commentList, &oneCommentInfo)
			wg.Done()
		}(commentTable)
		fmt.Println("one comment info", oneCommentInfo)
	}
	wg.Wait()
	// 根据id倒序，也就是根据创建时间倒序
	sort.Sort(response.CommentList(commentList))
	return commentList, nil
}

func (c *CommentServiceImpl) oneComment(comment *entity.Comment, commentInfo *response.CommentInfoRes) {
	userInfo, err := c.UserService.GetByID(comment.UserID)
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("UserService.GetByID failed, user_id：%d", comment.UserID))
	}
	commentInfo.Id = comment.ID
	commentInfo.UserInfo.ID = userInfo.ID
	commentInfo.UserInfo.Username = userInfo.Username
	commentInfo.Content = comment.CommentContent
	log.AppLogger.Debug(fmt.Sprintf("get oneComment info success, CommentInfo: %v", commentInfo))
}

func (c *CommentServiceImpl) insertRedisVideoCommentId(videoId int64, commentId int64) {
	// redis中新增videoId对应的commentId
	_, err := c.DataRepo.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, videoId), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		// 新增失败，把要操作的数据放入rabbitmq进行重试
		c.VideoCommentRabbitMQ.Publish(rabbitmq.CommentMessage{
			CommentDealType: consts.CommentAddMode,
			VideoId:         videoId,
			CommentId:       commentId,
		})
		log.AppLogger.Error(fmt.Sprintf("add video_id-comment_id failed, videoId: %d, commentId: %d", videoId, commentId))
	}
	// redis中新增commentId对应的videoId
	_, err = c.DataRepo.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId), strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		// 新增失败，把要操作的数据放入rabbitmq进行重试
		c.CommentVideoRabbitMQ.Publish(rabbitmq.CommentMessage{
			CommentDealType: consts.CommentAddMode,
			VideoId:         videoId,
			CommentId:       commentId,
		})
		log.AppLogger.Error(fmt.Sprintf("add video_id-comment_id failed, videoId: %d, commentId: %d", videoId, commentId))
	}
}
