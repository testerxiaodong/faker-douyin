package entity

import "gorm.io/gorm"

type TableComment struct {
	gorm.Model
	UserId         uint64
	VideoId        uint64
	CommentContent string
}

func (tableComment TableComment) TableName() string {
	return "comments"
}