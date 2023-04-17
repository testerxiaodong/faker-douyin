package main

import (
	"faker-douyin/internal/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type Server struct {
	logger     *zap.Logger
	Config     *config.Config
	HttpServer *http.Server
	Router     *gin.Engine
}

func NewServer(logger *zap.Logger, config *config.Config) {
	gin.SetMode(config.Server.Mode)
	router := gin.New()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s", config.Server.HttpPort),
		Handler: router,
	}
	s := &Server{logger: logger, Config: config, HttpServer: httpServer, Router: router}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// print err info when httpServer start failed
			s.logger.Error("unexpected error from ListenAndServe", zap.Error(err))
			fmt.Printf("http server start error:%s\n", err.Error())
			os.Exit(1)
		}
	}()
}

func main() {
	config.NewConfig()
}
