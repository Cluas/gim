package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Cluas/gim/internal/logic"
	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/internal/logic/rpc"
	"github.com/Cluas/gim/pkg/log"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {

	}
	// 设置cpu 核数
	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)
	if err := rpc.Init(); err != nil {
		log.Panic(fmt.Errorf("Init() fatal error : %s \n", err))
	}
	if err := logic.InitRedis(); err != nil {
		log.Panic(fmt.Errorf("InitRedis() fatal error : %s \n", err))
	}
	if err := logic.InitHTTP(); err != nil {
		log.Panic(fmt.Errorf("InitHTTP() fatal error : %s \n", err))
	}

}
