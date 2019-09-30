package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Config is log config
type Config struct {
	LogPath  string `mapstructure:"logpath"`
	LogLevel string `mapstructure:"loglevel"`
}

// NewLogger is func to new Logger
func NewLogger(c *Config) *zap.Logger {

	hook := lumberjack.Logger{
		Filename:   c.LogPath, // 日志文件路径
		MaxSize:    128,       // megabytes
		MaxBackups: 30,        // 最多保留300个备份
		MaxAge:     7,         // days
		Compress:   true,      // 是否压缩 disabled by default
	}

	var level zapcore.Level
	switch c.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level,
	)

	logger := zap.New(core)
	return logger
}

func init() {

	logger = NewLogger(&Config{LogPath: "./log.log", LogLevel: "info"})

}

//Init is initial func
func Init(c *Config) {
	logger = NewLogger(c)
}

//Debug 最低等级的，主要用于开发过程中打印一些运行/调试信息，不允许生产环境打开debug级别
func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

// Info 打印一些你感兴趣的或者重要的信息，这个可以用于生产环境中输出程序运行的一些重要信息
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// Warn 表明会出现潜在错误的情形，有些信息不是错误信息，但是也要给程序员的一些提示
func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

// Error 指出虽然发生错误事件，但仍然不影响系统的继续运行。打印错误和异常信息
func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

// DPanic is logs of DPanicLevel
func DPanic(msg string, fields ...zapcore.Field) {
	logger.DPanic(msg, fields...)
}

// Panic is logs of PanicLevel
func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

// Fatal 指出每个严重的错误事件将会导致应用程序的退出。这个级别比较高了。重大错误
func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}
