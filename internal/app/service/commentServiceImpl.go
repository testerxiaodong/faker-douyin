package service

import (
	"context"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
	"faker-douyin/internal/pkg/const"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	dataRepo    *dao.DataRepo
	UserService UserService
}

func NewCommentService(dataRepo *dao.DataRepo, userService UserService) CommentService {
	return &CommentServiceImpl{
		dataRepo:    dataRepo,
		UserService: userService,
	}
}

func (c *CommentServiceImpl) CommentInfo(commentId int64) (*entity.Comment, error) {
	comment, err := dao.Comment.Where(dao.Comment.ID.Eq(commentId)).First()
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (c *CommentServiceImpl) Count(videoId int64) (int64, error) {
	// 先在缓存中查找
	count, err := c.dataRepo.Rdb.SCard(context.Background(), strconv.Itoa(int(videoId))).Result()
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
	cntDao, err := dao.Video.Where(dao.Video.ID.Eq(videoId)).Count()
	if err != nil {
		fmt.Println(err)
	}
	if cntDao > 0 {
		//查询评论id list
		cList, _ := dao.Comment.Select(dao.Video.ID).Where(dao.Comment.VideoID.Eq(videoId)).Find()
		//设置key值过期时间
		_, err = c.dataRepo.Rdb.Expire(context.Background(), strconv.Itoa(int(videoId)),
			time.Duration(_const.OneMonth)*time.Second).Result()
		if err != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		//评论id循环存入redis
		for _, commentId := range cList {
			c.insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.FormatInt(commentId.ID, 10))
		}
		log.Println("count comment save ids in redis")
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
	err := dao.Comment.Create(&comment)
	if err != nil {
		return &comment, err
	}
	// 再更新缓存
	c.insertRedisVideoCommentId(strconv.FormatInt(videoId, 10), strconv.FormatUint(uint64(comment.ID), 10))
	return &comment, nil
}

func (c *CommentServiceImpl) DeleteComment(commentId int64) error {
	// 先删除数据库数据
	resultInfo, err := dao.Comment.Where(dao.Comment.ID.Eq(commentId)).Delete()
	if err != nil {
		return err
	}
	// 处理resultInfo.Error
	if resultInfo.Error != nil {
		return resultInfo.Error
	}
	fmt.Println("dao.DeleteCommentById成功，comment_id: ", commentId)
	// 先看redis中是否有数据
	_, err = c.dataRepo.Rdb.Exists(context.Background(), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		fmt.Println("key not exist in comment_iv-video_id ", commentId)
	}
	fmt.Println("redis中存在key：comment_id ", commentId)
	// 有数据，直接删redis数据
	videoId, err := c.dataRepo.Rdb.Get(context.Background(), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		fmt.Println("get value from comment_id-video_id failed, ", commentId)
	}
	_, err = c.dataRepo.Rdb.Del(context.Background(), strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		fmt.Println(err)
	}
	_, err = c.dataRepo.Rdb.SRem(context.Background(), videoId, strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("delete ", commentId, videoId)
	return nil
}

func (c *CommentServiceImpl) CommentList(videoId int64) ([]*response.CommentInfoRes, error) {
	commentTableList, err := dao.Comment.Where(dao.Comment.VideoID.Eq(videoId)).Find()
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
		fmt.Println("UserService.GetTableUserById failed, user_id：", comment.UserID)
	}
	commentInfo.Id = uint64(comment.ID)
	commentInfo.UserInfo.ID = userInfo.ID
	commentInfo.UserInfo.Username = userInfo.Username
	commentInfo.Content = comment.CommentContent
	fmt.Println(commentInfo)
}

func (c *CommentServiceImpl) insertRedisVideoCommentId(videoId string, commentId string) {
	_, err := c.dataRepo.Rdb.SAdd(context.Background(), videoId, commentId).Result()
	if err != nil {
		// 新增失败，暂时先上报日志，之后引入重试机制
		fmt.Println("add video_id-comment_id failed", videoId, commentId)
	}
	_, err = c.dataRepo.Rdb.Set(context.Background(), commentId, videoId, 0).Result()
	if err != nil {
		// 新增失败，暂时先上报日志，之后引入重试机制
		fmt.Println("save comment_id-video_id failed, ", commentId, videoId)
	}
}
