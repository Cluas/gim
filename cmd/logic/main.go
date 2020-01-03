package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/logic"
	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/internal/logic/http"
	"github.com/Cluas/gim/internal/logic/rpc"
	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/tracing"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		log.Bg().Panic("初始化配置文件失败", zap.Error(err))
	}
	if conf.Conf.Log != nil {
		log.Init(conf.Conf.Log)
	}

	_, closer := tracing.Init("logic", nil)
	defer Close(closer)

	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)
	if err := rpc.Init(); err != nil {
		log.Bg().Panic("初始化RPC服务端失败", zap.Error(err))
	}
	if err := logic.InitRedis(); err != nil {
		log.Bg().Panic("初始化Redis失败", zap.Error(err))
	}
	if err := http.Init(conf.Conf); err != nil {
		log.Bg().Panic("初始化HTTP服务端失败", zap.Error(err))
	}
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Bg().Info("gim-logic get a signal", zap.String("signal", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// Close id func to close
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Bg().Fatal("", zap.Error(err))
	}
}
