package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/server"
)

// Server is struct of logic rpc server
type Server int

const (
	split = "@"
)

// ConnectArg is struct of connect arg
type ConnectArg struct {
	Auth     string
	RoomID   int32
	ServerID int8
}

// ConnectReply is struct of connect reply
type ConnectReply struct {
	UID string
}

// ParseNetwork is func to parse network string
func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, split); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
	} else {
		network = str[:idx]
		addr = str[idx+1:]
	}
	return
}

// Init is func to initial LogicRPC server
func Init() (err error) {

	var (
		network, addr string
	)
	for _, bind := range conf.Conf.RPC.Address {
		if network, addr, err = ParseNetwork(bind); err != nil {
			log.Panicf("InitLogicRpc ParseNetwork error : %s", err)
		}
		go createServer(network, addr)
	}
	return
}

func createServer(network string, addr string) {

	s := server.NewServer()
	_ = s.RegisterName("LogicRPC", new(Server), "")
	_ = s.Serve(network, addr)

}

// Connect is api for connect
func (rpc *Server) Connect(ctx context.Context, args *ConnectArg, reply *ConnectReply) (err error) {
	log.Info("rpc logic 2  rpc uid ")
	if args == nil {
		log.Errorf("Connect() error(%v)", err)
		return
	}
	// TODO now is mock
	reply.UID = "555"
	log.Infof("logic rpc uid:%s", reply.UID)

	return
}
