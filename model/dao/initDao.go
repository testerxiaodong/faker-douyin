package dao

import (
	"faker-douyin/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	Db *gorm.DB
)

func Init() {
	mysqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Config.Mysql.User, global.Config.Mysql.Password, global.Config.Mysql.Host, global.Config.Mysql.Port,
		global.Config.Mysql.Name)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: mysqlLogger})
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	fmt.Println("数据库链接成功")
}
