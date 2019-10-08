package job

import (
	"context"
	"encoding/json"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/client"
)

var RPCClientList map[int8]client.XClient

type CometRPC int

type PushMsgArg struct {
	UID string
	P   Proto
}

type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type NoReply struct {
}
type RoomMsgArg struct {
	RoomID int32
	P      Proto
}

type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)

}

func InitComets() (err error) {
	CometAddress := make([]*client.KVPair, len(conf.Conf.Comet))
	RPCClientList = make(map[int8]client.XClient, len(conf.Conf.Comet))

	for i, bind := range conf.Conf.Comet {
		b := new(client.KVPair)
		b.Key = bind.Addr
		CometAddress[i] = b
		d := client.NewPeer2PeerDiscovery(bind.Addr, "")
		RPCClientList[bind.Key] = client.NewXClient("CometRPC", client.Failtry, client.RandomSelect, d, client.DefaultOption)
		log.Infof("CometRPC client %s", bind.Addr)
	}
	return

}

func PushSingle(serverId int8, userID string, msg []byte) {

	pushMsgArg := &PushMsgArg{UID: userID, P: Proto{Ver: 1, Operation: 2, Body: msg}}
	reply := &SuccessReply{}
	err := RPCClientList[serverId].Call(context.Background(), "PushSingleMsg", pushMsgArg, reply)
	if err != nil {
		log.Infof(" PushSingle Call err %v", err)
	}
	log.Infof("reply %s", reply.Msg)
}
func broadcastRoom(RoomId int32, msg []byte) {
	pushMsgArg := &RoomMsgArg{RoomID: RoomId, P: Proto{Ver: 1, Operation: 2, Body: msg}}
	reply := &SuccessReply{}
	log.Infof("broadcastRoom room_id %d", RoomId)
	for _, rpc := range RPCClientList {
		log.Infof("broadcastRoom rpc  %v", rpc)
		_ = rpc.Call(context.Background(), "PushRoomMsg", pushMsgArg, reply)
	}
}
