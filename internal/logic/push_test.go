package logic_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/Cluas/gim/internal/logic"
	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/tracing"
)

// ExampleRedisPublishCh example used to push single msg
func ExampleRedisPublishCh() {
	c := conf.NewConfig()
	conf.Conf = c
	if err := logic.InitRedis(); err != nil {
		log.Panic(fmt.Errorf("InitRedis() fatal error : %s \n", err))
	}

	_ = os.Setenv("JAEGER_REPORTER_LOG_SPANS", "true")
	_ = os.Setenv("JAEGER_AGENT_HOST", "118.25.5.127")
	_ = os.Setenv("JAEGER_AGENT_PORT", "6831")
	tracer, _ := tracing.Init("logic", nil)
	span := tracer.StartSpan("logic")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	_ = logic.RedisPublishCh(ctx, 1, "2850516", []byte(`这是一条发送到 "个人" 的消息`))

	span.Finish()

	time.Sleep(1 * time.Second)

	// Output:

}

// ExampleRedisPublishRoom example used to push room msg
func ExampleRedisPublishRoom() {
	c := conf.NewConfig()
	conf.Conf = c
	if err := logic.InitRedis(); err != nil {
		log.Panic(fmt.Errorf("InitRedis() fatal error : %s \n", err))
	}
	err := logic.RedisPublishRoom(context.Background(), "123", []byte(`这是一条发送到 "房间" 的消息`))
	log.Print(err)

	// Output:
}
