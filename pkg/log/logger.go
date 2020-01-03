package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	Panic(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
	Warn(msg string, fields ...zapcore.Field)
}

// sLogger delegates all calls to the underlying zap.Logger
type logger struct {
	logger *zap.Logger
}

// Info logs an info msg with fields
func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error msg with fields
func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// Panic logs a panic error msg with fields
func (l logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}

// Warn logs a warn error msg with fields
func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// With creates a child sLogger, and optionally adds some context fields to that Logger.
func (l logger) With(fields ...zapcore.Field) Logger {
	return logger{logger: l.logger.With(fields...)}
}
