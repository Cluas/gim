package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Factory is the default logging wrapper that can create
// sLogger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) Factory {
	return Factory{logger: logger}
}

// Bg creates a context-unaware sLogger.
func (b Factory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		spanCtx, ok := sp.Context().(jaeger.SpanContext)
		if ok {
			return spanLogger{span: sp, logger: b.logger}.With(zap.String("trace_id", spanCtx.TraceID().String()))
		}
		return spanLogger{span: sp, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child sLogger, and optionally adds some context fields to that sLogger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: b.logger.With(fields...)}
}
