package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Cluas/gim/pkg/comet/config"
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := config.InitConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n ", err))
	}
	fmt.Printf("%v\n", config.Conf)
}
