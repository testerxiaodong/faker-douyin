package global

import (
	"fmt"
	"github.com/spf13/viper"
)

type Server struct {
	AppMode  string
	HttpPort string
}

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type RabbitMq struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Jwt struct {
	secret string
}

type Ssh struct {
	Host     string
	User     string
	Password string
}

type Ftp struct {
	Host     string
	Port     string
	User     string
	Password string
}

type ConfigEnv struct {
	Server   Server
	Mysql    Mysql
	Redis    Redis
	RabbitMq RabbitMq
	Jwt      Jwt
	Ssh      Ssh
	Ftp      Ftp
}

var Config ConfigEnv

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	LoadServerConfig()
	LoadMysqlConfig()
	LoadRedisConfig()
	LoadRabbitMqConfig()
	LoadJwtConfig()
	LoadSshConfig()
	LoadFtpConfig()
	fmt.Println("配置加载完毕")
	fmt.Println(Config)
}

func LoadServerConfig() {
	Config.Server.AppMode = viper.GetString("server.AppMode")
	Config.Server.HttpPort = viper.GetString("server.HttpPort")
}

func LoadMysqlConfig() {
	Config.Mysql.Host = viper.GetString("mysql.Host")
	Config.Mysql.Port = viper.GetString("mysql.Port")
	Config.Mysql.User = viper.GetString("mysql.User")
	Config.Mysql.Password = viper.GetString("mysql.Password")
	Config.Mysql.Name = viper.GetString("mysql.Name")
}

func LoadRedisConfig() {
	Config.Redis.Host = viper.GetString("redis.Host")
	Config.Redis.Port = viper.GetString("redis.Port")
	Config.Redis.Password = viper.GetString("redis.Password")
}

func LoadRabbitMqConfig() {
	Config.RabbitMq.Host = viper.GetString("rabbitmq.Host")
	Config.RabbitMq.Port = viper.GetString("rabbitmq.Port")
	Config.RabbitMq.User = viper.GetString("rabbitmq.User")
	Config.RabbitMq.Password = viper.GetString("rabbitmq.Password")
}

func LoadJwtConfig() {
	Config.Jwt.secret = viper.GetString("jwt.secret")
}

func LoadSshConfig() {
	Config.Ssh.Host = viper.GetString("ssh.Host")
	Config.Ssh.User = viper.GetString("ssh.User")
	Config.Ssh.Password = viper.GetString("ssh.Password")
}

func LoadFtpConfig() {
	Config.Ftp.Host = viper.GetString("ftp.Host")
	Config.Ftp.Port = viper.GetString("ftp.Port")
	Config.Ftp.User = viper.GetString("ftp.User")
	Config.Ftp.Password = viper.GetString("ftp.Password")
}
