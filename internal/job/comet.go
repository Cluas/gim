package job

import (
	"context"
	"encoding/json"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/client"
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
	RoomID int32
	P      Proto
}

// Proto is struct of msg protocol
type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
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
		RPCClientList[bind.Key] = client.NewXClient("CometRPC", client.Failtry, client.RandomSelect, d, client.DefaultOption)
		log.Infof("CometRPC client %s", bind.Addr)
	}
	return

}

// PushSingle is func to PushSingle msg
func PushSingle(serverID int8, userID string, msg []byte) {

	pushMsgArg := &PushMsgArg{UID: userID, P: Proto{Ver: 1, Operation: 2, Body: msg}}
	reply := &SuccessReply{}
	err := RPCClientList[serverID].Call(context.Background(), "PushSingleMsg", pushMsgArg, reply)
	if err != nil {
		log.Infof(" PushSingle Call err %v", err)
	}
	log.Infof("reply %s", reply.Msg)
}

// broadcastRoom is func to broadcast room msg
func broadcastRoom(RoomID int32, msg []byte) {
	pushMsgArg := &RoomMsgArg{RoomID: RoomID, P: Proto{Ver: 1, Operation: 2, Body: msg}}
	reply := &SuccessReply{}
	log.Infof("broadcastRoom room_id %d", RoomID)
	for _, rpc := range RPCClientList {
		log.Infof("broadcastRoom rpc  %v", rpc)
		_ = rpc.Call(context.Background(), "PushRoomMsg", pushMsgArg, reply)
	}
}
