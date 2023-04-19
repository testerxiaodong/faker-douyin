package service

import (
	"context"
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	DataRepo    *dao.DataRepo
	UserService UserService
}

func (c *CommentServiceImpl) CommentInfo(commentId int64) (*entity.Comment, error) {
	comment, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.ID.Eq(commentId)).First()
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (c *CommentServiceImpl) Count(videoId int64) (int64, error) {
	// 先在缓存中查找
	count, err := c.DataRepo.Rdb.SCard(context.Background(), strconv.Itoa(int(videoId))).Result()
	if err != nil {
		//return 0, err
		fmt.Println(err)
	}
	// 缓存中有数据，直接返回
	if count > 0 {
		return count, nil
	}
	// 在数据库中找
	//cntDao, err := dao.Count(videoId)
	cntDao, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.VideoID.Eq(videoId)).Count()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	if cntDao > 0 {
		//查询评论id list
		cList, _ := c.DataRepo.Db.Comment.Select(c.DataRepo.Db.Comment.ID).Where(c.DataRepo.Db.Comment.VideoID.Eq(videoId)).Find()
		//设置key值过期时间
		_, err = c.DataRepo.Rdb.Expire(context.Background(), strconv.Itoa(int(videoId)),
			time.Duration(consts.OneMonth)*time.Second).Result()
		if err != nil {
			log.AppLogger.Error(err.Error())
		}
		//评论id循环存入redis
		for _, commentId := range cList {
			c.insertRedisVideoCommentId(videoId, commentId.ID)
		}
		log.AppLogger.Debug("count comment save ids in redis")
	}
	//返回结果
	return cntDao, nil
}

func (c *CommentServiceImpl) InsertComment(userId int64, videoId int64, commentContent string) (*entity.Comment, error) {
	// 先插入数据库
	var comment entity.Comment
	comment.VideoID = videoId
	comment.UserID = userId
	comment.CommentContent = commentContent
	err := c.DataRepo.Db.Comment.Create(&comment)
	if err != nil {
		return &comment, err
	}
	// 再更新缓存
	c.insertRedisVideoCommentId(videoId, comment.ID)
	return &comment, nil
}

func (c *CommentServiceImpl) DeleteComment(commentId int64) error {
	// 先删除数据库数据
	resultInfo, err := c.DataRepo.Db.Comment.Where(c.DataRepo.Db.Comment.ID.Eq(commentId)).Delete()
	if err != nil {
		return err
	}
	// 处理resultInfo.Error
	if resultInfo.Error != nil {
		return resultInfo.Error
	}
	log.AppLogger.Info(fmt.Sprintf("dao.DeleteCommentById成功, comment_id: %d", commentId))
	// 先看redis中是否有数据
	_, err = c.DataRepo.Rdb.Exists(context.Background(), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		log.AppLogger.Info(fmt.Sprintf("key not exist in comment_iv-video_id  %d", commentId))
		return err
	}
	log.AppLogger.Info(fmt.Sprintf("redis中存在key：comment_id: %d", commentId))
	// 有数据，直接删redis数据
	videoId, err := c.DataRepo.Rdb.Get(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId)).Result()
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("get videoId from comment:video failed, comment_id: %d", commentId))
		return err
	}
	_, err = c.DataRepo.Rdb.Del(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId)).Result()
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("delete comment:video failed, comment_id: %d", commentId))
		return err
	}
	videoIdInt64, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		log.AppLogger.Error("strconv ParseInt failed")
		return err
	}
	_, err = c.DataRepo.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, videoIdInt64), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("remove video:comment failed, video_id: %d comment_id: %d", videoIdInt64, commentId))
		return err
	}
	log.AppLogger.Info(fmt.Sprintf("delete comment_id: %d video_id: %d success", commentId, videoIdInt64))
	fmt.Println("delete ", commentId, videoId)
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
	commentInfo.Id = uint64(comment.ID)
	commentInfo.UserInfo.ID = userInfo.ID
	commentInfo.UserInfo.Username = userInfo.Username
	commentInfo.Content = comment.CommentContent
	log.AppLogger.Debug(fmt.Sprintf("get oneComment info success, CommentInfo: %v", commentInfo))
}

func (c *CommentServiceImpl) insertRedisVideoCommentId(videoId int64, commentId int64) {
	_, err := c.DataRepo.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, videoId), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		// 新增失败，暂时先上报日志，之后引入重试机制
		log.AppLogger.Error(fmt.Sprintf("add video_id-comment_id failed, videoId: %d, commentId: %d", videoId, commentId))
	}
	_, err = c.DataRepo.Rdb.Set(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, commentId), strconv.FormatInt(videoId, 10), 0).Result()
	if err != nil {
		// 新增失败，暂时先上报日志，之后引入重试机制
		log.AppLogger.Error(fmt.Sprintf("add comment_id-video_id failed, videoId: %d, commentId: %d", videoId, commentId))
	}
}
