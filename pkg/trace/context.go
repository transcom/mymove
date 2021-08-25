package trace

import (
	"context"

	"github.com/gofrs/uuid"
)

// the contextKey is typed so as not to conflict between similar keys from different pkgs
type contextKey string

var traceIDContextKey = contextKey("milmove_trace_id")

// NewContext adds the traceID string into the context and returns the new context
func NewContext(ctx context.Context, traceID uuid.UUID) context.Context {
	return context.WithValue(ctx, traceIDContextKey, traceID)
}

// FromContext returns a traceID that was previously added into the context, if any.
func FromContext(ctx context.Context) uuid.UUID {
	// This is a recursive call that checks the nested contexts for this key
	traceID, ok := ctx.Value(traceIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return traceID
}
