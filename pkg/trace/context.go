package trace

import (
	"context"
)

type contextKey int

var traceIDContextKey = contextKey(1)

// NewContext returns a new context with a logger.
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey, traceID)
}

// FromContext returns a logger associated with a context, if any.
func FromContext(ctx context.Context) string {
	str, ok := ctx.Value(traceIDContextKey).(string)
	if !ok {
		return ""
	}
	return str
}
