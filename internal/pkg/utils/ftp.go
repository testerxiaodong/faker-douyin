package utils

import (
	"faker-douyin/internal/pkg/config"
	"faker-douyin/internal/pkg/const"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"time"
)

var MyFtp *ftp.ServerConn

func InitFtp(config *config.Config) {
	dsn := config.Ftp.Host + ":" + config.Ftp.Port
	fmt.Println(dsn)
	var err error
	MyFtp, err = ftp.Dial(dsn, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ftp服务器链接成功")
	err = MyFtp.Login(config.Ftp.User, config.Ftp.Password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("登陆ftp服务器成功")
	//Linux小知识：用户登陆时所处目录/home/$username
	//cwd, err := MyFtp.CurrentDir()
	//if err != nil {
	//	fmt.Println("获取当前目录失败")
	//}
	//fmt.Println(cwd)
	go FtpKeepAlive()
}

func FtpKeepAlive() {
	time.Sleep(time.Duration(_const.HeartbeatTime) * time.Second)
	err := MyFtp.NoOp()
	if err != nil {
		log.Fatal("维持ftp长链接失败")
	}
}

// VideoFTP
// 通过ftp将视频传入服务器
func VideoFTP(file io.Reader, videoName string) error {
	//转到video相对路线下
	err := MyFtp.ChangeDir("/home/ftpuser/videos")
	if err != nil {
		log.Println("转到路径videos失败！！！")
	} else {
		log.Println("转到路径videos成功！！！")
	}
	err = MyFtp.Stor(videoName+".mp4", file)
	if err != nil {
		log.Println("上传视频失败！！！！！")
		return err
	}
	log.Println("上传视频成功！！！！！")
	return nil
}

// ImageFTP
// 将图片传入FTP服务器中，但是这里要注意图片的格式随着名字一起给,同时调用时需要自己结束流
func ImageFTP(file io.Reader, imageName string) error {
	//转到video相对路线下
	err := MyFtp.ChangeDir("images")
	if err != nil {
		log.Println("转到路径images失败！！！")
		return err
	}
	log.Println("转到路径images成功！！！")
	if err = MyFtp.Stor(imageName, file); err != nil {
		log.Println("上传图片失败！！！！！")
		return err
	}
	log.Println("上传图片成功！！！！！")
	return nil
}
