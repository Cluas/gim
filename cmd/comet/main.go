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
		panic(fmt.Errorf("Fatal error conf file: %s \n ", err))
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Init(conf.Conf.Log)
	perf.Init(conf.Conf.Base.PprofBind)
	if err := rpc.InitLogicRpc(); err != nil {
		log.Panic(fmt.Errorf("InitLogicRpc Fatal error: %s \n", err))
	}
	server := comet.NewServer(conf.Conf)

	log.Info("Starting WebSocket...")
	log.Infof("WebSocket server : %s", conf.Conf.Websocket)
	if err := comet.InitWebsocket(server, conf.Conf.Websocket); err != nil {
		log.Fatal(err)
	}

}
