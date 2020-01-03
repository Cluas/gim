package job

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/job/conf"
	"github.com/Cluas/gim/pkg/log"
)

type pushArg struct {
	ServerID int8
	UserID   string
	Msg      []byte
	RoomID   int32
	Context  context.Context
}

var pushChs []chan *pushArg

// InitPush is func to initial push msg
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
		pushSingle(arg.Context, arg.ServerID, arg.UserID, arg.Msg)

	}
}

func push(msg string) (err error) {
	m := &RedisMsg{}
	msgByte := []byte(msg)
	if err = json.Unmarshal(msgByte, m); err != nil {
		log.Bg().Info(" json.Unmarshal err:%v ", zap.Error(err))
	}

	ctx, spanCtx := genContext(m)
	ctx, span := contextWithSpan(ctx, spanCtx, "Job.Redis.ReceiveMessage")
	ctx = opentracing.ContextWithSpan(ctx, span)

	switch m.Op {
	case OpSingleSend:
		pushChs[rand.Int()%conf.Conf.Base.PushChan] <- &pushArg{
			ServerID: m.ServerID,
			UserID:   m.UserID,
			Msg:      m.Msg,
			Context:  ctx,
		}
		break
	case OpRoomSend:
		broadcastRoom(ctx, m.RoomID, m.Msg)
		break
	}
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()

	return
}

func contextWithSpan(ctx context.Context, spanCtx opentracing.SpanContext, spanName string) (context.Context, opentracing.Span) {
	if spanCtx != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName, opentracing.FollowsFrom(spanCtx))
		ctx = opentracing.ContextWithSpan(ctx, span)
		return ctx, span
	}
	return ctx, nil
}

func genContext(m *RedisMsg) (context.Context, opentracing.SpanContext) {
	ctx := context.Background()
	tracer := opentracing.GlobalTracer()
	if opentracing.IsGlobalTracerRegistered() {
		carrier := new(bytes.Buffer)
		carrier.Write(m.Carrier)
		spanCtx, _ := tracer.Extract(opentracing.Binary, carrier)
		return ctx, spanCtx
	}
	return ctx, nil
}
