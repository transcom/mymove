package scenario

import (
	"go.uber.org/zap"
)

// Logger is an interface that describes the logging requirements of this package.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	WithOptions(options ...zap.Option) *zap.Logger
}
