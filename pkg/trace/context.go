package trace

import (
	"context"
)

// the contextKey is typed so as not to conflict between similar keys from different pkgs
type contextKey string

var traceIDContextKey = contextKey("milmove_trace_id")

// NewContext adds the traceID string into the context and returns the new context
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey, traceID)
}

// FromContext returns a traceID that was previously added into the context, if any.
func FromContext(ctx context.Context) string {
	// This is a recursive call that checks the nested contexts for this key
	str, ok := ctx.Value(traceIDContextKey).(string)
	if !ok {
		return ""
	}
	return str
}
