package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Cluas/gim/pkg/comet/config"
	"github.com/Cluas/gim/pkg/comet/log"
)

func main() {
	flag.Parse()
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n ", err))
	}
	log.Info("测试初始化日志")
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Init(config.Conf.Log)
	log.Info("success")
}
