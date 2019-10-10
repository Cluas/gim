package comet

import (
	"context"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/client"
)

var (
	logicRPCClient client.XClient
)

// ConnectReply is struct for Connect reply
type ConnectReply struct {
	UID string
}

// DisconnectReply is struct for Disconnect reply
type DisconnectReply struct {
	Has bool
}

// InitLogic is func for initial logic rpc client
func InitLogic() (err error) {

	LogicAddr := make([]*client.KVPair, len(conf.Conf.RPC.LogicAddr))

	for i, bind := range conf.Conf.RPC.LogicAddr {
		b := new(client.KVPair)
		b.Key = bind.Addr
		LogicAddr[i] = b
		log.Infof("创建LogicRPC客户端, 绑定地址: %s", bind.Addr)

	}
	d := client.NewMultipleServersDiscovery(LogicAddr)

	logicRPCClient = client.NewXClient("LogicRPC", client.Failover, client.RoundRobin, d, client.DefaultOption)

	return
}

func connect(c *ConnectArg) (uid string, err error) {

	reply := &ConnectReply{}
	err = logicRPCClient.Call(context.Background(), "Connect", c, reply)
	log.Info("Comet 调用 Logic Connect(c *ConnectArg)...")
	if err != nil {
		log.Errorf("Comet 调用 Logic Connect(c *ConnectArg) 失败, 原因: %v", err)
	}

	uid = reply.UID
	log.Infof("Comet 调用 Logic Connect(c *ConnectArg) 返回UID :%s", reply.UID)

	return
}

func disconnect(d *DisconnectArg) (err error) {

	reply := &DisconnectReply{}
	if err = logicRPCClient.Call(context.Background(), "Disconnect", d, reply); err != nil {
		log.Errorf("failed to call: %v", err)
	}
	return
}
