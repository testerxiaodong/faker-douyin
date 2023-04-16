package utils

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

const WordDictPath = "/Users/cengdong/GolandProjects/faker-douyin/document/sensitiveDict.txt"

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadNetWordDict("https://raw.githubusercontent.com/importcjj/sensitive/master/dict/dict.txt")
	if err != nil {
		log.Println("Load network sensitive word fail：", err.Error())
	}
	err = Filter.LoadWordDict(WordDictPath)
	if err != nil {
		log.Println("InitFilter Fail,Err=" + err.Error())
	}
}
