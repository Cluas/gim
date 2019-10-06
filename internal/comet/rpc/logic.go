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

type ConnectReply struct {
	UID string
}
type DisconnectReply struct {
	Has bool
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

func connect(c *ConnectArg) (uid string, err error) {

	log.Info("connect logic rpc...")
	reply := &ConnectReply{}
	err = logicRpcClient.Call(context.Background(), "Connect", c, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	uid = reply.UID
	log.Infof("comet logic uid :%s", reply.UID)

	return
}

func disconnect(d *DisconnectArg) (err error) {

	reply := &DisconnectReply{}
	if err = logicRpcClient.Call(context.Background(), "Disconnect", d, reply); err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	return
}
