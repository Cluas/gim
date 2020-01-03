package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _logger *zap.Logger

// XFactory is default factory
var XFactory Factory

// Config is log conf
type Config struct {
	LogPath     string `mapstructure:"log_path"`
	LogLevel    string `mapstructure:"log_level"`
	ServiceName string `mapstructure:"service_name"`
}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)
	_logger = zap.New(core, zap.AddStacktrace(zapcore.FatalLevel), zap.AddCaller())
	XFactory = NewFactory(_logger)

}

//Init is initial func
func Init(c *Config) {
	_logger.With(zap.String("service", c.ServiceName))
}

// Bg creates a context-unaware sLogger.
func Bg() Logger {
	return XFactory.Bg()
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func For(ctx context.Context) Logger {
	return XFactory.For(ctx)
}

// With creates a child sLogger, and optionally adds some context fields to that sLogger.
func With(fields ...zapcore.Field) Factory {
	return XFactory.With(fields...)
}
