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
	fmt.Println(dsn)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: mysqlLogger})
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	fmt.Println("数据库链接成功")
	// err = Db.AutoMigrate(&entity.User{}, &entity.TableVideo{}, entity.TableComment{})
	SetDefault(Db)
	if err != nil {
		panic(fmt.Sprintf("数据库迁移失败,%s", err))
	}
	DB, err := Db.DB()
	if err != nil {
		panic(fmt.Sprintf("数据库初始化配置失败,%s", err))
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	DB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	DB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	DB.SetConnMaxLifetime(10 * time.Second)
}
