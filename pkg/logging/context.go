package logging

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

var loggerContextKey = contextKey("logger")

// NewContext returns a new context with a logger.
func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext returns a logger associated with a context. If the
// context has no logger, the global logger is returned
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*zap.Logger)
	if !ok {
		// globally replaced in serve.go, so if this context doesn't
		// have a logger, use the global logger
		return zap.L()
	}
	return logger
}

// FromContextWithoutDefault returns a logger associated with a
// context. If the context has no logger, nil is returned
func FromContextWithoutDefault(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}
