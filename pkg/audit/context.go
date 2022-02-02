package audit

import (
	"context"
)

type contextKey string

var loggerContextKey = contextKey("operationID")

// NewContext returns a new context with the operationId.
func WithEventName(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, loggerContextKey, operationID)
}

// FromContext returns the operationId associated with a context. If the
// context has no operationId, the empty string is returned
func EventNameFromContext(ctx context.Context) string {
	operationID, ok := ctx.Value(loggerContextKey).(string)
	if !ok {
		return ""
	}
	return operationID
}
