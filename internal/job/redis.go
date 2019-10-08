package job

import (
	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/go-redis/redis"
)

var (
	RedisCli *redis.Client
)

func InitRedis() (err error) {
	RedisCli = redis.NewClient(&redis.Options{
		Addr: conf.Conf.
			Redis.Address,
		Password: conf.Conf.Redis.Password,  // no password set
		DB:       conf.Conf.Redis.DefaultDB, // use default DB
	})
	if pong, err := RedisCli.Ping().Result(); err != nil {
		log.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}
	go func() {
		redisSub := RedisCli.Subscribe("gim_sub_channel")
		ch := redisSub.Channel()
		log.Infof("开始订阅Channel, %s", "gim_sub_channel")
		for {
			msg, ok := <-ch
			if !ok {
				log.Debugf("redisSub Channel !ok: %v", ok)
				break
			}

			_ = push(msg.Payload)
			if conf.Conf.Base.IsDebug == true {
				log.Infof("redisSub Subscribe msg : %s", msg.Payload)
			}

		}

	}()

	return
}
