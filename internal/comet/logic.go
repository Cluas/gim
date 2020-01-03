package comet

import (
	"context"

	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
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
		log.Bg().Info("创建LogicRPC客户端, 绑定地址: ", zap.String("addr", bind.Addr))

	}
	d := client.NewMultipleServersDiscovery(LogicAddr)

	logicRPCClient = client.NewXClient("LogicRPC", client.Failover, client.RoundRobin, d, client.DefaultOption)

	return
}

func connect(ctx context.Context, c *ConnectArg) (uid string, err error) {

	reply := &ConnectReply{}
	err = logicRPCClient.Call(ctx, "Connect", c, reply)
	if err != nil {
		log.Bg().Error("Comet 调用 Logic Connect(c *ConnectArg) 失败, 原因:", zap.Error(err))

	}

	uid = reply.UID
	log.Bg().Info("Comet 调用 Logic Connect(c *ConnectArg) 返回UID :", zap.String("uid", reply.UID))

	return
}

func disconnect(_ context.Context, d *DisconnectArg) (err error) {

	reply := &DisconnectReply{}
	if err = logicRPCClient.Call(context.Background(), "Disconnect", d, reply); err != nil {
		log.Bg().Error("failed to call:", zap.Error(err))
	}
	return
}
