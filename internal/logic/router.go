package logic

import (
	"github.com/Cluas/gim/pkg/log"
)

var (
	userInfo map[string]string
)

// Router is struct for user router
type Router struct {
	ServerID int8
	RoomID   int32
	UserID   string
	Username string
}

func getRouter(auth string) (router *Router, err error) {
	userInfo, err = RedisCli.HGetAll(getKey(auth)).Result()
	if err != nil {
		return
	}
	log.Infof("getRouter auth :%s, userId:%s", auth, userInfo["UserID"])
	router = &Router{UserID: userInfo["UserID"], Username: userInfo["Username"]}
	return

}
