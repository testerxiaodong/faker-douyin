package entity

import "gorm.io/gorm"

type TableVideo struct {
	gorm.Model
	AuthorId uint64 // 视频作者id
	Title    string // 视频标题
	PlayUrl  string // 播放地址
	CoverUrl string // 封面地址
}

func (tableVideo TableVideo) TableName() string {
	return "videos"
}
