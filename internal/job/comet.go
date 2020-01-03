package job

import (
	"context"
	"encoding/json"

	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
)

// RPCClientList is list of CometRPC
var RPCClientList map[int8]client.XClient

// CometRPC is comet rpc Client
type CometRPC int

// PushMsgArg is struct of push msg
type PushMsgArg struct {
	UID string
	P   Proto
}

// SuccessReply is struct of success reply
type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NoReply is struct of no reply
type NoReply struct {
}

// RoomMsgArg is struct of room msg arg
type RoomMsgArg struct {
	RoomID string
	P      Proto
}

// Proto is struct of msg protocol
type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation Operation       `json:"op"`   // operation for request
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)

}

// InitComets is func to initial CometRPC client
func InitComets() (err error) {
	CometAddress := make([]*client.KVPair, len(conf.Conf.Comet))
	RPCClientList = make(map[int8]client.XClient, len(conf.Conf.Comet))

	for i, bind := range conf.Conf.Comet {
		b := new(client.KVPair)
		b.Key = bind.Addr
		CometAddress[i] = b
		d := client.NewPeer2PeerDiscovery(bind.Addr, "")
		c := client.NewXClient("CometRPC", client.Failtry, client.RandomSelect, d, client.DefaultOption)
		p := &client.OpenTracingPlugin{}
		pc := client.NewPluginContainer()
		pc.Add(p)
		c.SetPlugins(pc)

		RPCClientList[bind.Key] = c
		log.Bg().Info("CometRPC client:", zap.String("bind", bind.Addr))
	}
	return

}

// pushSingle is func to pushSingle msg
func pushSingle(ctx context.Context, serverID int8, userID string, msg []byte) {

	pushMsgArg := &PushMsgArg{UID: userID, P: Proto{Ver: 1, Operation: OpSingleSend, Body: msg}}
	reply := &SuccessReply{}
	err := RPCClientList[serverID].Call(ctx, "PushSingleMsg", pushMsgArg, reply)
	if err != nil {
		log.For(ctx).Info("pushSingle Call err:", zap.Error(err))
	}
}

// broadcastRoom is func to broadcast room msg
func broadcastRoom(ctx context.Context, RoomID string, msg []byte) {
	pushMsgArg := &RoomMsgArg{RoomID: RoomID, P: Proto{Ver: 1, Operation: OpRoomSend, Body: msg}}
	reply := &SuccessReply{}
	for _, rpc := range RPCClientList {
		log.For(ctx).Info("job call comet PushRoomMsg")
		err := rpc.Call(ctx, "PushRoomMsg", pushMsgArg, reply)
		if err != nil {
			log.For(ctx).Info("PushRoomMsg Call err:", zap.Error(err))
		}
	}
}
