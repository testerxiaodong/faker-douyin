package dao

import (
	"context"
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/log"
	"fmt"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormMysql, NewRedisClient, NewDataRepo)

// DataRepo 数据获取源
type DataRepo struct {
	Db  *Query
	Rdb *redis.Client
}

func NewDataRepo(db *Query, rdb *redis.Client) *DataRepo {
	return &DataRepo{
		Db:  db,
		Rdb: rdb,
	}
}

type GormClient struct {
	config *config.Config
	logger log.GormLogger
}

func NewGormMysql(config *config.Config, gormLogger *log.GormLogger) *Query {
	// 拼接地址
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port,
		config.Mysql.Name)

	log.AppLogger.Debug(dsn)
	// 连接数据库
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("连接Mysql服务器失败：%s", err.Error()))
	}
	log.AppLogger.Debug("Mysql服务器连接成功")
	sqlDB, err := DB.DB()
	if err != nil {
		log.AppLogger.Error(fmt.Sprintf("获取数据库连接失败：%s", err.Error()))
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	SetDefault(DB)

	return Use(DB)
}

type RedisClient struct {
	Config *config.Config
}

// NewRedisClient 获取redis连接
func NewRedisClient(config *config.Config) *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return Rdb
}

func GetRedisKeyByPrefix(prefix string, videoId int64) string {
	return fmt.Sprintf("%s%s", prefix, strconv.FormatInt(videoId, 10))
}
