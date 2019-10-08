package job

import (
	"encoding/json"
	"math/rand"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
)

type pushArg struct {
	ServerID int8
	UserID   string
	Msg      []byte
	RoomID   int32
}

var pushChs []chan *pushArg

func InitPush() {
	pushChs = make([]chan *pushArg, conf.Conf.Base.PushChan)
	for i := 0; i < len(pushChs); i++ {

		pushChs[i] = make(chan *pushArg, conf.Conf.Base.PushChanSize)
		go processPush(pushChs[i])
	}
}

func processPush(ch chan *pushArg) {
	var arg *pushArg
	for {
		arg = <-ch
		PushSingle(arg.ServerID, arg.UserID, arg.Msg)

	}
}
func push(msg string) (err error) {
	m := &RedisMsg{}
	msgByte := []byte(msg)
	if err := json.Unmarshal(msgByte, m); err != nil {
		log.Infof(" json.Unmarshal err:%v ", err)
	}
	log.Infof("push m info %s", m)

	switch m.Op {
	case OP_SINGLE_SEND:
		pushChs[rand.Int()%conf.Conf.Base.PushChan] <- &pushArg{
			ServerID: m.ServerID,
			UserID:   m.UserID,
			Msg:      m.Msg,
		}
		break
	case OP_ROOM_SEND:
		broadcastRoom(m.RoomID, m.Msg)
		break
		//case OP_ROOM_COUNT_SEND:
		//	broadcastRoomCountToComet(m.RoomID, m.Count)
		//case OP_ROOM_INFO_SEND:
		//	broadcastRoomInfoToComet(m.RoomID, m.RoomUserInfo)
	}

	return
}
