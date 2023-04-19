package main

import (
	"faker-douyin/internal/app"
	"faker-douyin/internal/pkg/utils"
)

func main() {
	// 启用敏感词过滤工具
	utils.InitFilter()
	// 启动应用程序
	err := app.Init()
	if err != nil {
		panic(err)
	}
	// 监听信号
	app.Run()
}
