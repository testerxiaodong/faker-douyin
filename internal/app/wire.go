//go:build wireinject
// +build wireinject

package app

import (
	v1 "faker-douyin/internal/app/api/v1"
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/dao"
	log "faker-douyin/internal/app/log"
	"faker-douyin/internal/app/middleware/rabbitmq"
	"faker-douyin/internal/app/router"
	"faker-douyin/internal/app/service"
	"faker-douyin/internal/pkg/utils"
	"github.com/google/wire"
)

func CreateApp() (*App, error) {
	wire.Build(
		config.NewConfig,
		log.NewLogger,
		log.NewGormLogger,
		dao.ProviderSet,
		utils.NewFilter,
		utils.NewFtpClient,
		utils.NewFfmpegClient,
		rabbitmq.ProviderSet,
		service.ProviderSet,
		v1.ProviderSet,
		router.ProviderSet,
		InitGinEngine,
		InjectorSet,
	)
	return new(App), nil
}
