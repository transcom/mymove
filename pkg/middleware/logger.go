package middleware

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
	With(fields ...zap.Field) *zap.Logger
	WithOptions(options ...zap.Option) *zap.Logger
}

// InfoLogger is a logger interface with the Info method.
type InfoLogger interface {
	Info(msg string, fields ...zap.Field)
}
