package rpc

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/logic"
	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
)

// Server is struct of logic rpc server
type Server int

const (
	split = "@"
)

// ConnectArg is struct of connect arg
type ConnectArg struct {
	Auth     string
	RoomID   string
	ServerID string
}

// ConnectReply is struct of connect reply
type ConnectReply struct {
	UID string
}

// DisconnectArg is struct of disconnect arg
type DisconnectArg struct {
	RoomID string
	UID    string
}

//DisconnectReply is struct of disconnect reply
type DisconnectReply struct {
	Has bool
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
			log.Bg().Panic("InitLogicRpc ParseNetwork error", zap.Error(err))
		}
		log.Bg().Info("创建Logic RPC", zap.String("bind", bind))
		go createServer(network, addr)
	}
	return
}

func createServer(network string, addr string) {

	s := server.NewServer()
	p := serverplugin.OpenTracingPlugin{}
	s.Plugins.Add(p)
	_ = s.RegisterName("LogicRPC", new(Server), "")
	_ = s.Serve(network, addr)

}

// Connect is api for connect
func (rpc *Server) Connect(_ context.Context, args *ConnectArg, reply *ConnectReply) (err error) {

	if len(args.Auth) == 0 {
		return
	}
	var member *Member
	member, err = JwtParseMember(args.Auth)
	if err != nil {

	}

	if member != nil {
		reply.UID = strconv.Itoa(member.ID)
		logic.RedisCli.Set(logic.RedisAuthPrefix+reply.UID, args.ServerID, logic.RedisBaseValidTime*time.Second)
		logic.RedisCli.HSet(logic.RedisRoomUserPrefix+args.RoomID, reply.UID, member.Nickname)
	}

	return
}

// Disconnect if func to remove room user
func (rpc *Server) Disconnect(_ context.Context, args *DisconnectArg, reply *DisconnectReply) (err error) {
	if args.UID != "" {
		err = logic.RedisCli.HDel(logic.RedisRoomUserPrefix+args.RoomID, args.UID).Err()
		reply.Has = true

	}

	return
}
