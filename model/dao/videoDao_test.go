package dao

import (
	"faker-douyin/global"
	"fmt"
	"testing"
	"time"
)

func TestGetVideosByAuthorId(t *testing.T) {
	global.LoadConfig()
	Init()
	data, err := GetVideosByAuthorId(2)
	if err != nil {
		t.Error(err)
	}
	for _, video := range data {
		fmt.Println(video)
	}
}

func TestGetVideoById(t *testing.T) {
	global.LoadConfig()
	Init()
	data, err := GetVideoById(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestGetVideosByLastTime(t *testing.T) {
	global.LoadConfig()
	Init()
	data, err := GetVideosByLastTime(time.Now())
	if err != nil {
		t.Error(err)
	}
	for _, video := range data {
		fmt.Println(video)
	}
}

func TestInsertTableVideo(t *testing.T) {
	global.LoadConfig()
	Init()
	video, err := InsertTableVideo("测试标题", "测试视频名", "测试封面名", 2)
	if err != nil {
		t.Error(err)
	}
	t.Log(video)
}
