package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/job"
	"github.com/Cluas/gim/internal/job/conf"
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

	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)

	_, closer := tracing.Init("job", nil)
	defer Close(closer)

	if err := job.InitRedis(); err != nil {
		log.Bg().Panic("初始化Redis失败", zap.Error(err))
	}

	if err := job.InitComets(); err != nil {
		log.Bg().Panic("初始化Comet客户端失败", zap.Error(err))
	}

	job.InitPush()

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Bg().Info("gim-job get a signal", zap.String("signal", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			err := job.RedisCli.Close()
			if err != nil {
				log.Bg().Error("Redis关闭失败", zap.Error(err))
			}
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
