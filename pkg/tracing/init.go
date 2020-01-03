package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	"go.uber.org/zap"

	"github.com/Cluas/gim/pkg/log"
)

// Init creates a new instance of Jaeger tracer.
func Init(serviceName string, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Bg().Fatal("cannot parse Jaeger env vars", zap.Error(err))
	}
	if metricsFactory == nil {
		metricsFactory = jexpvar.NewFactory(10)
	}
	cfg.ServiceName = serviceName

	// TODO(cluas): 后期转移成可配置选项
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	jaegerLogger := jaegerLoggerAdapter{log.Bg()}

	metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: serviceName, Tags: nil})
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
		config.Gen128Bit(true),
	)
	if err != nil {
		log.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
