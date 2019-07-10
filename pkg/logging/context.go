package logging

import (
	"context"
)

type contextKey int

var loggerContextKey = contextKey(1)

// NewContext returns a new context with a logger.
func NewContext(ctx context.Context, logger interface{}) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext returns a logger associated with a context, if any.
func FromContext(ctx context.Context) interface{} {
	return ctx.Value(loggerContextKey)
}
