package rpc

import (
	"context"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/client"
)

var (
	logicRpcClient client.XClient
)

type ConnArg struct {
	Auth   string
	RoomId int32
	Server int32
}

type ConnReply struct {
	Uid string
}

func InitLogicRpc() (err error) {

	LogicAddr := make([]*client.KVPair, len(conf.Conf.RPC.LogicAddr))

	for i, bind := range conf.Conf.RPC.LogicAddr {
		log.Infof("logic rpc bind %s", bind)
		b := new(client.KVPair)
		b.Key = bind
		LogicAddr[i] = b

	}
	d := client.NewMultipleServersDiscovery(LogicAddr)

	logicRpcClient = client.NewXClient("LogicRpc", client.Failover, client.RoundRobin, d, client.DefaultOption)
	return
}

func connect(connArg *ConnArg) (uid string, err error) {

	log.Infof("comet logic rpc logicRpcClient %s:", logicRpcClient)
	reply := &ConnReply{}
	err = logicRpcClient.Call(context.Background(), "Connect", connArg, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	uid = reply.Uid
	log.Infof("comet logic uid :%s", reply.Uid)

	return
}
