package utils

import (
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/log"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"time"
)

type FtpClient struct {
	Conn *ftp.ServerConn
}

func NewFtpClient(config *config.Config) *FtpClient {
	var ftpClient FtpClient
	dsn := config.Ftp.Host + ":" + config.Ftp.Port
	var err error
	ftpClient.Conn, err = ftp.Dial(dsn, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	log.AppLogger.Debug("ftp服务器连接成功")
	err = ftpClient.Conn.Login(config.Ftp.User, config.Ftp.Password)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	log.AppLogger.Debug("ftp服务器登陆成功")
	//Linux小知识：用户登陆时所处目录/home/$username
	//cwd, err := MyFtp.CurrentDir()
	//if err != nil {
	//	fmt.Println("获取当前目录失败")
	//}
	//fmt.Println(cwd)
	go ftpClient.FtpKeepAlive()
	return &ftpClient

}

func (f *FtpClient) FtpKeepAlive() {
	time.Sleep(time.Duration(consts.HeartbeatTime) * time.Second)
	err := f.Conn.NoOp()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

// VideoFTP
// 通过ftp将视频传入服务器
func (f *FtpClient) VideoFTP(file io.Reader, videoName string) error {
	//转到video相对路线下
	err := f.Conn.ChangeDir("/home/ftpuser/videos")
	if err != nil {
		log.AppLogger.Error(err.Error())
	} else {
		log.AppLogger.Debug("转到路径videos成功")
	}
	err = f.Conn.Stor(videoName+".mp4", file)
	if err != nil {
		log.AppLogger.Error(err.Error())
		return err
	}
	log.AppLogger.Debug("上传视频成功！")
	return nil
}

// ImageFTP
// 将图片传入FTP服务器中，但是这里要注意图片的格式随着名字一起给,同时调用时需要自己结束流
func (f *FtpClient) ImageFTP(file io.Reader, imageName string) error {
	//转到video相对路线下
	err := f.Conn.ChangeDir("images")
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("转到路径images失败：%s", err.Error()))
		return err
	}
	log.AppLogger.Debug("转到路径images成功！")
	if err = f.Conn.Stor(imageName, file); err != nil {
		log.AppLogger.Error(err.Error())
		return err
	}
	log.AppLogger.Debug("上传图片成功")
	return nil
}
