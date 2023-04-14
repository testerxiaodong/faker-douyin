package dao

import (
	"faker-douyin/model/entity"
	"fmt"
)

func GetCommentById(commentId uint64) (entity.TableComment, error) {
	var comment entity.TableComment
	if err := Db.Find(&comment, commentId).Error; err != nil {
		fmt.Println("GetCommentById failed, err：", err)
		return comment, err
	}
	return comment, nil
}

func DeleteCommentById(commentId uint64) error {
	err := Db.Where("id = ?", commentId).Delete(&entity.TableComment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertComment(userId uint64, videoId uint64, commentContent string) (entity.TableComment, error) {
	var comment entity.TableComment
	var user entity.TableUser
	var video entity.TableVideo
	if err := Db.Where("id = ?", userId).First(&user).Error; err != nil {
		fmt.Println("user not found, user_id：", userId)
		return comment, err
	}
	if err := Db.Where("id = ?", videoId).First(&video).Error; err != nil {
		fmt.Println("video not found, video_id：", videoId)
		return comment, err
	}
	comment.UserId = userId
	comment.VideoId = videoId
	comment.CommentContent = commentContent
	Db.Create(&comment)
	return comment, nil
}

func GetCommentIdList(videoId uint64) ([]uint64, error) {
	var idList []uint64
	if err := Db.Model(&entity.TableComment{}).Select("id").Where("video_id = ?", videoId).Find(&idList).Error; err != nil {
		return idList, err
	}
	return idList, nil
}

func GetCommentList(videoId uint64) ([]entity.TableComment, error) {
	var comments []entity.TableComment
	if err := Db.Where("video_id = ?", videoId).Find(&comments).Error; err != nil {
		return comments, err
	}
	return comments, nil
}

func Count(videoId uint64) (int64, error) {
	var count int64
	err := Db.Model(&entity.TableComment{}).Where("video_id = ?", videoId).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}
