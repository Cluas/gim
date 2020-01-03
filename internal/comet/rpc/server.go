package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/comet"
	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
)

const (
	successCode = 0
	successMsg  = "success"
	splitString = "@"
)

// PushMsgArg is struct of push msg arg
type PushMsgArg struct {
	UID string
	P   comet.Proto
}

// SuccessReply if struct of success reply
type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Server is struct of Comet RPC Server
type Server int

// ParseNetwork is util func used to parse network
func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, splitString); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
	} else {
		network = str[:idx]
		addr = str[idx+1:]
	}
	return
}

// Init is func to initial comet rpc server
func Init() (err error) {
	var (
		network, addr string
	)
	binds := conf.Conf.RPC.CometAddr
	for _, bind := range binds {
		if network, addr, err = ParseNetwork(bind.Addr); err != nil {
			log.Bg().Panic("InitLogicRpc ParseNetwork error : ", zap.Error(err))
		}
		log.Bg().Info("创建comet RPC", zap.String("bind", bind.Addr))
		go createServer(network, addr)
	}
	return
}

func createServer(network string, addr string) {
	s := server.NewServer()
	p := serverplugin.OpenTracingPlugin{}
	s.Plugins.Add(p)
	_ = s.RegisterName("CometRPC", new(Server), "")
	_ = s.Serve(network, addr)
}

// PushSingleMsg is rpc func used to push single msg
func (rpc *Server) PushSingleMsg(ctx context.Context, args *PushMsgArg, SuccessReply *SuccessReply) (err error) {
	var (
		bucket  *comet.Bucket
		channel *comet.Channel
	)

	if args == nil {
		log.For(ctx).Error("rpc Server() error", zap.Error(err))
		return
	}

	bucket = comet.CurrentServer.Bucket(ctx, args.UID)
	if channel = bucket.Channel(args.UID); channel != nil {
		err = channel.Push(&args.P)
		return
	}

	SuccessReply.Code = 1
	SuccessReply.Msg = "success"
	return
}

// PushRoomMsg is func used tp push room msg
func (rpc *Server) PushRoomMsg(ctx context.Context, args *comet.RoomMsgArg, SuccessReply *SuccessReply) (err error) {

	SuccessReply.Code = successCode
	SuccessReply.Msg = successMsg

	for _, bucket := range comet.CurrentServer.Buckets {
		bucket.BroadcastRoom(args)
	}

	return
}
