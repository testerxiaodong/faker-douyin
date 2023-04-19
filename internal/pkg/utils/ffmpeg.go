package utils

import (
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/log"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

type FfmpegClient struct {
	SshClient *ssh.Client
	Ffchan    chan Ffmsg
}

// NewFfmpegClient  建立SSH客户端，但是会不会超时导致无法链接，这个需要做一些措施
func NewFfmpegClient(config *config.Config) *FfmpegClient {
	var ffmpegClient FfmpegClient
	var err error
	//创建ssh登陆配置
	SshConfig := &ssh.ClientConfig{
		Timeout:         5 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            config.Ssh.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全

		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if config.Ssh.TypeSsh == "password" {
		SshConfig.Auth = []ssh.AuthMethod{ssh.Password(config.Ssh.Password)}
	}
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%s", config.Ssh.Host, config.Ssh.Port)
	ffmpegClient.SshClient, err = ssh.Dial("tcp", addr, SshConfig)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	//建立通道，作为队列使用,并且确立缓冲区大小
	ffmpegClient.Ffchan = make(chan Ffmsg, consts.MaxMsgCount)
	//建立携程用于派遣
	go ffmpegClient.dispatcher()
	go ffmpegClient.SshKeepAlive()
	return &ffmpegClient
}

// 通过增加携程，将获取的信息进行派遣，当信息处理失败之后，还会将处理方式放入通道形成的队列中
func (f *FfmpegClient) dispatcher() {
	for ffmsg := range f.Ffchan {
		go func(fs Ffmsg) {
			err := f.Ffmpeg(fs.VideoName, fs.ImageName)
			if err != nil {
				f.Ffchan <- fs
				log.AppLogger.Error(err.Error())
			}
			log.AppLogger.Debug(fmt.Sprintf("视频处理成功，video name: %s", fs.VideoName))
		}(ffmsg)
	}
}

// Ffmpeg 通过远程调用ffmpeg命令来创建视频截图
func (f *FfmpegClient) Ffmpeg(videoName string, imageName string) error {
	session, err := f.SshClient.NewSession()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	//执行远程命令 ffmpeg -ss 00:00:01 -i /home/ftpuser/video/1.mp4 -vframes 1 /home/ftpuser/images/4.jpg
	combo, err := session.CombinedOutput("ls;/usr/local/ffmpeg/bin/ffmpeg -ss 00:00:01 -i /home/ftpuser/videos/" + videoName + ".mp4 -vframes 1 /home/ftpuser/images/" + imageName + ".jpg")
	if err != nil {
		log.AppLogger.Error(err.Error())
		return err
	}
	log.AppLogger.Debug(fmt.Sprintf("命令输出：%s", combo))
	return nil
}

// SshKeepAlive 维持长链接
func (f *FfmpegClient) SshKeepAlive() {
	time.Sleep(time.Duration(consts.SSHHeartbeatTime) * time.Second)
	session, _ := f.SshClient.NewSession()
	err := session.Close()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}
