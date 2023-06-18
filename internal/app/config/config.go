package config

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	Mode     string `mapstructure:"mode"`
	HttpPort string `mapstructure:"port"`
	WorkDir  string `mapstructure:"work_dir"`
	LogDir   string `mapstructure:"log_dir"`
}

type Log struct {
	FileName string `mapstructure:"filename"`
	Levels   Levels `mapstructure:"level"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
	Compress bool   `mapstructure:"compress"`
}

type Levels struct {
	App  string `mapstructure:"app"`
	Gorm string `mapstructure:"gorm"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMq struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Ssh struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	TypeSsh  string `mapstructure:"type_ssh"`
}

type Ftp struct {
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	HeartbeatTime int64  `mapstructure:"heartbeat_time"`
}

type Config struct {
	Sensitive string   `mapstructure:"sensitive"`
	Server    Server   `mapstructure:"server"`
	Log       Log      `mapstructure:"log"`
	Mysql     Mysql    `mapstructure:"mysql"`
	Redis     Redis    `mapstructure:"redis"`
	RabbitMq  RabbitMq `mapstructure:"rabbitmq"`
	Ssh       Ssh      `mapstructure:"ssh"`
	Ftp       Ftp      `mapstructure:"ftp"`
}

func NewConfig() *Config {
	var configFile string
	flag.StringVar(&configFile, "config", "", "")
	flag.Parse()
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("yaml")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("/Users/cengdong/GolandProjects/faker-douyin/config")
		viper.SetConfigName("config")
	}
	conf := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}
	if conf.Server.WorkDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "init config: get current dir"))
		}
		conf.Server.WorkDir, _ = filepath.Abs(pwd)
	} else {
		workDir, err := filepath.Abs(conf.Server.WorkDir)
		if err != nil {
			panic(err)
		}
		conf.Server.WorkDir = workDir
	}
	normalizeDir := func(path *string, subDir string) {
		if *path == "" {
			*path = filepath.Join(conf.Server.WorkDir, subDir)
		} else {
			temp, err := filepath.Abs(*path)
			if err != nil {
				panic(err)
			}
			*path = temp
		}
	}
	normalizeDir(&conf.Server.LogDir, "logs")

	initDirectory(conf)
	mode = conf.Server.Mode
	return conf
}

func initDirectory(conf *Config) {
	mkdirFunc := func(dir string, err error) error {
		if err == nil {
			if _, err = os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm)
			}
		}
		return err
	}
	err := mkdirFunc(conf.Server.LogDir, nil)
	if err != nil {
		panic(err)
	}
}

var mode string

// IsDev 用于日志器的配置
func IsDev() bool {
	return mode == "debug"
}
