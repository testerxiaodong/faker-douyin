package main

import (
	"faker-douyin/global"
	"faker-douyin/model/dao"
)

func main() {
	global.LoadConfig()
	dao.Init()
}
