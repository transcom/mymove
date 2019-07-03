package internalapi

import (
	"go.uber.org/zap"
)

// Logger is a logger interface for middleware
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	WithOptions(options ...zap.Option) *zap.Logger
}
