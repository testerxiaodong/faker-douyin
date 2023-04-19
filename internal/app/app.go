package app

import (
	"faker-douyin/internal/app/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var InjectorSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	Config *config.Config
	Engine *gin.Engine
}

func Init() error {
	app, err := CreateApp()
	if err != nil {
		return err
	}
	app.InitHttpServer()
	return nil
}

func (a *App) InitHttpServer() {
	addr := fmt.Sprintf("127.0.0.1:%s", a.Config.Server.HttpPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.Engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func Run() {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

EXIT:
	for {
		sig := <-sc
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	time.Sleep(time.Second)
	os.Exit(state)
}
