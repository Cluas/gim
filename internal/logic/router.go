package logic

import (
	"bytes"
	"strconv"
)

var (
	userInfo map[string]string
)

type Router struct {
	ServerId int64
	RoomId   int64
	UserId   int64
}

func getRouter(auth string) (router *Router, err error) {
	var key bytes.Buffer
	key.WriteString(auth)
	key.WriteString(REDIS_AUTH_PREFIX)
	userInfo, err = RedisCli.HGetAll(key.String()).Result()
	if err != nil {
		return
	}
	uid, err := strconv.ParseInt(userInfo["UserId"], 10, 64)
	if err != nil {
		return
	}
	rid, err := strconv.ParseInt(userInfo["RoomId"], 10, 32)
	if err != nil {
		return
	}
	sid, err := strconv.ParseInt(userInfo["ServerId"], 10, 16)
	if err != nil {
		return
	}
	router = &Router{ServerId: sid, RoomId: rid, UserId: uid}
	return

}
