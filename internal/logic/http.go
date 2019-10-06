package logic

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
)

const (
	split = "@"
)

var router *Router

func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, split); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx]
		addr = str[idx+1:]
		return
	}
}
func InitHTTP() (err error) {
	var network, addr string

	for i := 0; i < len(conf.Conf.HTTP.Address); i++ {

		httpServeMux := http.NewServeMux()
		httpServeMux.HandleFunc("/api/v1/push", Push)

		if network, addr, err = ParseNetwork(conf.Conf.HTTP.Address[i]); err != nil {
			log.Errorf("ParseNetwork() error(%v)", err)
			return
		}

		log.Infof("start http listen:\"%s\"", conf.Conf.HTTP.Address[i])

		go httpListen(httpServeMux, network, addr)
		select {}

	}
	return
}
func Push(w http.ResponseWriter, r *http.Request) {
	// log.Info("yes")
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
	}

	var (
		auth      = r.URL.Query().Get("auth")
		err       error
		bodyBytes []byte
	)

	if router, err = getRouter(auth); err != nil {

		log.Errorf("get router error : %s", err)
		return
	}

	log.Infof("router info %v", router)

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		log.Errorf("get router error : %s", err)
	}
	log.Infof("get bodyBytes : %s", bodyBytes)

	if err := RedisPublishCh(router.ServerID, router.UserID, bodyBytes); err != nil {
		log.Errorf("redis Publish err: %s", err)
	}

}

func httpListen(mux *http.ServeMux, network, addr string) {

	httpServer := &http.Server{Handler: mux, ReadTimeout: conf.Conf.HTTP.ReadTimeout, WriteTimeout: conf.Conf.HTTP.WriteTimeout}
	httpServer.SetKeepAlivesEnabled(true)

	l, err := net.Listen(network, addr)
	if err != nil {
		log.Errorf("net.Listen(\"%s\", \"%s\") error(%v)", network, addr, err)
		panic(err)
	}
	if err := httpServer.Serve(l); err != nil {
		log.Errorf("server.Serve() error(%v)", err)
		panic(err)
	}
}
