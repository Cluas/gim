package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cluas/gim/internal/comet"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/smallnest/rpcx/server"
)

const (
	Success    = 0
	SuccessMsg = "success"
)

type PushMsgArg struct {
	UID string
	P   comet.Proto
}

type NoReply struct {
}

type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CometRPC int

const (
	split = "@"
)

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
func Init() (err error) {
	var (
		network, addr string
	)
	binds := conf.Conf.RPC.CometAddr
	for _, bind := range binds {
		if network, addr, err = ParseNetwork(bind.Addr); err != nil {
			log.Panicf("InitLogicRpc ParseNetwork error : %s", err)
		}
		log.Infof("创建comet RPC, %s", bind.Addr)
		go createServer(network, addr)
	}
	return
}

func createServer(network string, addr string) {
	s := server.NewServer()
	_ = s.RegisterName("CometRPC", new(CometRPC), "")
	_ = s.Serve(network, addr)
}

//func (rpc *CometRPC) MPushMsg(ctx context.Context, args *PushMsgArg, noReply *NoReply) (err error) {
//
//	log.Info("rpc PushMsg :%v ", args)
//	if args == nil {
//		log.Errorf("rpc CometRPC() error(%v)", err)
//		return
//	}
//
//	return
//}

func (rpc *CometRPC) PushSingleMsg(ctx context.Context, args *PushMsgArg, SuccessReply *SuccessReply) (err error) {
	var (
		bucket  *comet.Bucket
		channel *comet.Channel
	)

	log.Info("rpc PushMsg :%v ", args)
	if args == nil {
		log.Errorf("rpc CometRPC() error(%v)", err)
		return
	}
	bucket = comet.CurrentServer.Bucket(args.UID)
	if channel = bucket.Channel(args.UID); channel != nil {
		err = channel.Push(&args.P)

		log.Infof("DefaultServer Channel err nil : %v", err)
		return
	}

	SuccessReply.Code = 1
	SuccessReply.Msg = "success"
	log.Infof("SuccessReply v :%v", SuccessReply)
	return
}

func (rpc *CometRPC) PushRoomMsg(ctx context.Context, args *comet.RoomMsgArg, SuccessReply *SuccessReply) (err error) {

	SuccessReply.Code = Success
	SuccessReply.Msg = SuccessMsg
	log.Infof("PushRoomMsg msg %v", args)
	for _, bucket := range comet.CurrentServer.Buckets {
		bucket.BroadcastRoom(args)
		// room.next

	}
	return
}
