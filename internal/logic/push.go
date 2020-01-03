package logic

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/tracing"
)

var (
	// RedisCli is default redis client
	RedisCli *redis.Client
)

// Operation defines the type of operation.
type Operation int

const (
	// OpSend 发送
	_ Operation = iota + 1 //
	// OpSingleSend 指定用户发送
	OpSingleSend
	// OpRoomSend 广播到房间操作
	OpRoomSend
)

const (
	// RedisRoomUserPrefix 用户房间
	RedisRoomUserPrefix = "im_room_user_"
	// RedisAuthPrefix 授权 用来寻找serverID
	RedisAuthPrefix = "gim_auth_"
	// RedisSubChannel 消息订阅地址
	RedisSubChannel = "gim_sub_channel"
	// RedisBaseValidTime 过期时间
	RedisBaseValidTime = 86400
)

// RedisMsg is struct of redis msg
type RedisMsg struct {
	Carrier  []byte    `json:"carrier"` // 携带trace信息
	Op       Operation `json:"op"`
	ServerID int8      `json:"serverID,omitempty"`
	RoomID   string    `json:"roomID,omitempty"`
	UserID   string    `json:"userID,omitempty"`
	Msg      []byte    `json:"msg"`
}

// InitRedis is func to initial redis client
func InitRedis() (err error) {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:        conf.Conf.Redis.Address,
		Password:    conf.Conf.Redis.Password,  // no password set
		DB:          conf.Conf.Redis.DefaultDB, // use default DB
		MaxRetries:  conf.Conf.Redis.MaxRetries,
		IdleTimeout: conf.Conf.Redis.IdleTimeout,
	})
	var pong string
	if pong, err = RedisCli.Ping().Result(); err != nil {
		log.Bg().Info("RedisCli Ping Result pong err", zap.String("pong", pong), zap.Error(err))
	}

	return
}

// RedisPublishCh is func to push msg
func RedisPublishCh(ctx context.Context, serverID int8, uid string, msg []byte) (err error) {
	var redisMsg = &RedisMsg{Op: OpSingleSend, ServerID: serverID, UserID: uid, Msg: msg}
	setSpanCarrier(ctx, redisMsg)
	redisMsgStr, err := json.Marshal(redisMsg)
	err = tracing.WrapRedisClient(ctx, RedisCli).Publish(RedisSubChannel, redisMsgStr).Err()
	return
}

// RedisPublishRoom is func to push msg to room
func RedisPublishRoom(ctx context.Context, rid string, msg []byte) (err error) {
	var redisMsg = &RedisMsg{
		Op:     OpRoomSend,
		RoomID: rid,
		Msg:    msg,
	}
	setSpanCarrier(ctx, redisMsg)
	redisMsgStr, err := json.Marshal(redisMsg)
	log.Bg().Info("RedisPublishRoom redisMsg info", zap.Binary("redisMsgStr", redisMsgStr))
	err = tracing.WrapRedisClient(ctx, RedisCli).Publish(RedisSubChannel, redisMsgStr).Err()
	return
}
func setSpanCarrier(ctx context.Context, redisMsg *RedisMsg) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if opentracing.IsGlobalTracerRegistered() {
			carrier := new(bytes.Buffer)
			tracer := opentracing.GlobalTracer()
			span = tracer.StartSpan("Logic.Redis.PublishMessage", opentracing.FollowsFrom(span.Context()))
			_ = tracer.Inject(span.Context(), opentracing.Binary, carrier)
			redisMsg.Carrier = carrier.Bytes()
			span.Finish()
		}
	}
}
