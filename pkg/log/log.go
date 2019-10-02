package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// Config is log conf
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

	logger = NewLogger(&Config{LogPath: "./log.log", LogLevel: "info"}).Sugar()

}

//Init is initial func
func Init(c *Config) {
	logger = NewLogger(c).Sugar()
}

//Debug 最低等级的，主要用于开发过程中打印一些运行/调试信息，不允许生产环境打开debug级别
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf 支持参数格式化
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

// Info 打印一些你感兴趣的或者重要的信息，这个可以用于生产环境中输出程序运行的一些重要信息
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof 支持参数格式化
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Warn 表明会出现潜在错误的情形，有些信息不是错误信息，但是也要给程序员的一些提示
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf 支持参数格式化
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Error 指出虽然发生错误事件，但仍然不影响系统的继续运行。打印错误和异常信息
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf 支持参数格式化
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// DPanic is logs of DPanicLevel
func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

// DPanicf 支持参数格式化
func DPanicf(msg string, args ...interface{}) {
	logger.DPanicf(msg, args...)
}

// Panic is logs of PanicLevel
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf 支持参数格式化
func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

// Fatal 指出每个严重的错误事件将会导致应用程序的退出。这个级别比较高了。重大错误
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf 支持参数格式化
func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}
