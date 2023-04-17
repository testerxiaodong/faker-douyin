package dao

import (
	"context"
	"faker-douyin/internal/pkg/config"
	"fmt"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	config *config.Mysql
	logger logger.Interface
}

func NewGormMysql(config *config.Mysql, gormLogger logger.Interface) *Query {
	// 拼接地址
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port,
		config.Name)

	fmt.Println(dsn)
	// 连接数据库
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		fmt.Println("连接Mysql服务器失败：", zap.Error(err))
	}
	fmt.Println("数据库链接成功")
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("获取数据库连接失败：", zap.Error(err))
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
	Config *config.Redis
}

// NewRedisClient 获取redis连接
func NewRedisClient(config *config.Redis) *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("redis connect ping failed: %s", err))
	}
	return Rdb
}
