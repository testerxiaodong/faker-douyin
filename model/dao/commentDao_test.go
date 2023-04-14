package dao

import (
	"faker-douyin/global"
	"fmt"
	"testing"
)

func TestGetCommentById(t *testing.T) {
	global.LoadConfig()
	Init()
	comment, err := GetCommentById(2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(comment)
}

func TestInsertComment(t *testing.T) {
	global.LoadConfig()
	Init()
	comment, err := InsertComment(1, 9, "测试评论")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(comment)
}

func TestDeleteCommentById(t *testing.T) {
	global.LoadConfig()
	Init()
	err := DeleteCommentById(1)
	if err != nil {
		t.Error(err)
	}
}

func TestGetCommentIdList(t *testing.T) {
	global.LoadConfig()
	Init()
	idList, err := GetCommentIdList(9)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(idList)
}

func TestGetCommentList(t *testing.T) {
	global.LoadConfig()
	Init()
	comments, err := GetCommentList(9)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(comments)
}

func TestCount(t *testing.T) {
	global.LoadConfig()
	Init()
	count, err := Count(9)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}
