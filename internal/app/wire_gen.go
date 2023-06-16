// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"faker-douyin/internal/app/api/v1"
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"faker-douyin/internal/app/middleware/rabbitmq"
	"faker-douyin/internal/app/router"
	"faker-douyin/internal/app/service"
	"faker-douyin/internal/pkg/utils"
)

// Injectors from wire.go:

func CreateApp() (*App, error) {
	configConfig := config.NewConfig()
	logger := log.NewLogger(configConfig)
	gormLogger := log.NewGormLogger(configConfig, logger)
	query := dao.NewGormMysql(configConfig, gormLogger)
	client := dao.NewRedisClient(configConfig)
	dataRepo := dao.NewDataRepo(query, client)
	userServiceImpl := &service.UserServiceImpl{
		DataRepo: dataRepo,
	}
	userController := v1.UserController{
		U: userServiceImpl,
	}
	ftpClient := utils.NewFtpClient(configConfig)
	ffmpegClient := utils.NewFfmpegClient(configConfig)
	rabbitMQ := rabbitmq.NewRabbitMQ(configConfig)
	commentRabbitMQ := rabbitmq.NewCommentRabbitMQ(rabbitMQ, client)
	commentServiceImpl := &service.CommentServiceImpl{
		DataRepo:        dataRepo,
		CommentRabbitMQ: commentRabbitMQ,
		UserService:     userServiceImpl,
	}
	videoServiceImpl := &service.VideoServiceImpl{
		DataRepo:       dataRepo,
		FtpClient:      ftpClient,
		FfmpegClient:   ffmpegClient,
		UserService:    userServiceImpl,
		CommentService: commentServiceImpl,
	}
	videoController := v1.VideoController{
		V: videoServiceImpl,
	}
	commentController := v1.CommentController{
		CommentService: commentServiceImpl,
	}
	routerRouter := &router.Router{
		UserController:    userController,
		VideoController:   videoController,
		CommentController: commentController,
	}
	engine := InitGinEngine(routerRouter, configConfig)
	app := &App{
		Config: configConfig,
		Engine: engine,
	}
	return app, nil
}
