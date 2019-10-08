package logic

import (
	"bytes"
	"encoding/json"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/go-redis/redis"
)

var (
	RedisCli *redis.Client
)

const (
	REDIS_PREFIX           = "gim_"
	REDIS_ROOM_USER_PREFIX = "im_room_user_"
	REDIS_AUTH_PREFIX      = "gim_auth_"
	REDIS_SUB_CHANNEL      = "gim_sub_channel"
	OP_SINGLE_SEND         = int32(2)
	REDIS_BASE_VALID_TIME  = 86400
)

type RedisMsg struct {
	Op       int32  `json:"op"`
	ServerId int8   `json:"serverId,omitempty"`
	RoomId   int32  `json:"roomId,omitempty"`
	UserId   string `json:"userId,omitempty"`
	Msg      []byte `json:"msg"`
}

func InitRedis() (err error) {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Address,
		Password: conf.Conf.Redis.Password,  // no password set
		DB:       conf.Conf.Redis.DefaultDB, // use default DB
	})
	if pong, err := RedisCli.Ping().Result(); err != nil {
		log.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}

	return
}

func RedisPublishCh(serverId int8, uid string, msg []byte) (err error) {
	var redisMsg = &RedisMsg{Op: OP_SINGLE_SEND, ServerId: serverId, UserId: uid, Msg: msg}
	redisMsgStr, err := json.Marshal(redisMsg)
	err = RedisCli.Publish(REDIS_SUB_CHANNEL, redisMsgStr).Err()
	return
}

func getKey(key string) string {

	var returnKey bytes.Buffer
	returnKey.WriteString(REDIS_AUTH_PREFIX)
	returnKey.WriteString(key)
	return returnKey.String()
}

func getRoomUserKey(key string) string {

	var returnKey bytes.Buffer
	returnKey.WriteString(REDIS_ROOM_USER_PREFIX)
	returnKey.WriteString(key)
	return returnKey.String()
}
