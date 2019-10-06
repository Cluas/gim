package logic

import (
	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/go-redis/redis"
)

var (
	RedisCli *redis.Client
)

const (
	REDIS_AUTH_PREFIX       = "gim_auth"
	REDIS_SUB_CHANNEL       = "gim_sub_channel"
	REDIS_MESSAGE_BROADCAST = "broadcast"
)

type RedisMsg struct {
	Op       string `json:"op"`
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
	var redisMsg = &RedisMsg{Op: REDIS_MESSAGE_BROADCAST, ServerId: serverId, UserId: uid, Msg: msg}
	err = RedisCli.Publish(REDIS_SUB_CHANNEL, redisMsg).Err()
	return
}
