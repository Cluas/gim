package job

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
)

var (
	// RedisCli is Client of redis
	RedisCli *redis.Client
)

// InitRedis is func to initial redis
func InitRedis() (err error) {
	RedisCli = redis.NewClient(&redis.Options{
		Addr: conf.Conf.
			Redis.Address,
		Password: conf.Conf.Redis.Password,  // no password set
		DB:       conf.Conf.Redis.DefaultDB, // use default DB
	})
	var pong string
	if pong, err = RedisCli.Ping().Result(); err != nil {
		log.Bg().Info("RedisCli Ping Result pong err:", zap.String("pong", pong), zap.Error(err))
	}
	go func() {
		redisSub := RedisCli.Subscribe("gim_sub_channel")
		ch := redisSub.Channel()
		log.Bg().Info("开始订阅Channel:", zap.String("channel", "gim_sub_channel"))
		for {
			msg, ok := <-ch
			if !ok {
				break
			}
			_ = push(msg.Payload)
			if conf.Conf.Base.IsDebug == true {
				log.Bg().Info("redisSub Subscribe msg :", zap.String("payload", msg.Payload))
			}
		}
	}()

	return
}
