package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Cluas/gim/internal/job"
	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
	}

	// 设置cpu 核数
	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)

	// 初始化redis
	if err := job.InitRedis(); err != nil {
		log.Panic(fmt.Errorf("InitRedis() fatal error : %s \n", err))
	}

	// 通过rpc初始化comet对应的 server bucket等
	if err := job.InitComets(); err != nil {
		log.Panic(fmt.Errorf("InitRPC() fatal error : %s \n", err))
	}

	job.InitPush()
	select {}

}
