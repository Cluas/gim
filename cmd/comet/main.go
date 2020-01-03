package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/comet"
	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/internal/comet/rpc"
	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/perf"
	"github.com/Cluas/gim/pkg/tracing"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		log.Bg().Panic("初始化配置文件失败", zap.Error(err))
	}

	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)

	if conf.Conf.Log != nil {
		log.Init(conf.Conf.Log)
	}

	_, closer := tracing.Init("comet", nil)
	defer Close(closer)

	perf.Init(conf.Conf.Base.PprofBind)

	if err := comet.InitLogic(); err != nil {
		log.Bg().Panic("初始化LogicRPC客户端失败", zap.Error(err))
	}

	srv := comet.NewServer(conf.Conf)

	if err := rpc.Init(); err != nil {
		log.Bg().Panic("初始化RPC服务端失败n", zap.Error(err))
	}
	log.Bg().Info("开始启动websocket服务", zap.String("bind", conf.Conf.Websocket.Bind))
	if err := comet.InitWebsocket(srv, conf.Conf.Websocket); err != nil {
		log.Bg().Fatal("", zap.Error(err))
	}
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Bg().Info("gim-comet get a signal", zap.String("signal", s.String()))
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
