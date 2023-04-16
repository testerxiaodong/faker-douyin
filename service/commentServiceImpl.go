package service

import (
	"faker-douyin/model/dao"
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
	"fmt"
	"sort"
	"sync"
)

type CommentServiceImpl struct {
	UserService
}

func (c CommentServiceImpl) CommentInfo(commentId uint64) (entity.TableComment, error) {
	comment, err := dao.GetCommentById(commentId)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (c CommentServiceImpl) Count(videoId uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommentServiceImpl) InsertComment(userId uint64, videoId uint64, commentContent string) (entity.TableComment, error) {
	comment, err := dao.InsertComment(userId, videoId, commentContent)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (c CommentServiceImpl) DeleteComment(commentId uint64) error {
	err := dao.DeleteCommentById(commentId)
	if err != nil {
		return err
	}
	return nil
}

func (c CommentServiceImpl) CommentList(videoId uint64) ([]response.CommentInfoRes, error) {
	commentTableList, err := dao.GetCommentList(videoId)
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
	commentList := make([]response.CommentInfoRes, 0, len(commentTableList))
	var wg sync.WaitGroup
	wg.Add(len(commentTableList))
	for _, commentTable := range commentTableList {
		var oneCommentInfo response.CommentInfoRes
		// 传入循环变量作为临时变量，防止bug
		go func(commentTable entity.TableComment) {
			oneComment(&commentTable, &oneCommentInfo)
			commentList = append(commentList, oneCommentInfo)
			wg.Done()
		}(commentTable)
		fmt.Println("one comment info", oneCommentInfo)
	}
	wg.Wait()
	// 根据id倒序，也就是根据创建时间倒序
	sort.Sort(response.CommentList(commentList))
	return commentList, nil
}

func oneComment(comment *entity.TableComment, commentInfo *response.CommentInfoRes) {
	usi := UserServiceImpl{}
	userInfo, err := usi.GetTableUserById(comment.UserId)
	if err != nil {
		fmt.Println("UserService.GetTableUserById failed, user_id：", comment.UserId)
	}
	commentInfo.Id = uint64(comment.ID)
	commentInfo.UserInfo.Id = userInfo.Id
	commentInfo.UserInfo.Name = userInfo.Name
	commentInfo.Content = comment.CommentContent
	fmt.Println(commentInfo)
}
