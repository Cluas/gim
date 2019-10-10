package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Cluas/gim/internal/comet"
	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/internal/comet/rpc"
	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/perf"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Panic(fmt.Errorf("初始化配置文件失败: %s \n ", err))
	}
	runtime.GOMAXPROCS(conf.Conf.Base.MaxProc)
	log.Init(conf.Conf.Log)
	perf.Init(conf.Conf.Base.PprofBind)
	if err := comet.InitLogic(); err != nil {
		log.Panic(fmt.Errorf("初始化LogicRPC客户端失败: %s \n", err))
	}
	srv := comet.NewServer(conf.Conf)
	if err := rpc.Init(); err != nil {
		log.Panic(fmt.Errorf("初始化RPC服务端失败: %s \n", err))
	}
	log.Infof("开始启动websocket服务: %s", conf.Conf.Websocket.Bind)
	if err := comet.InitWebsocket(srv, conf.Conf.Websocket); err != nil {
		log.Fatal(err)
	}

}
