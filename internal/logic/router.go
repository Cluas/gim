package logic

import (
	"bytes"
	"strconv"

	"github.com/smallnest/rpcx/log"
)

var (
	userInfo map[string]string
)

type Router struct {
	ServerID int8
	RoomID   int32
	UserID   string
}

func getRouter(auth string) (router *Router, err error) {
	var key bytes.Buffer
	key.WriteString(auth)
	key.WriteString(REDIS_AUTH_PREFIX)
	log.Infof("key %s", key.String())

	userInfo, err = RedisCli.HGetAll(key.String()).Result()
	if err != nil {
		return
	}
	log.Infof("user_id %v", userInfo)
	uid, err := strconv.ParseInt(userInfo["UserID"], 10, 64)
	if err != nil {
		return
	}
	rid, err := strconv.ParseInt(userInfo["RoomID"], 10, 32)
	if err != nil {
		return
	}
	sid, err := strconv.ParseInt(userInfo["ServerID"], 10, 16)
	if err != nil {
		return
	}
	router = &Router{ServerID: int8(sid), RoomID: int32(rid), UserID: string(uid)}
	return

}
