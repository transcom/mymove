package audit

import (
	"context"

	"github.com/gofrs/uuid"
)

type contextKey string

var loggerEventNameContextKey = contextKey("operationID")
var loggerAuditIDContextKey = contextKey("AuditID")

// WithEventName returns a new context with the operationId.
func WithEventName(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, loggerEventNameContextKey, operationID)
}

// RetrieveEventNameFromContext returns the operationId associated with a context. If the
// context has no operationId, the empty string is returned
func RetrieveEventNameFromContext(ctx context.Context) string {
	operationID, ok := ctx.Value(loggerEventNameContextKey).(string)
	if !ok {
		return ""
	}
	return operationID
}

// WithAuditUserID returns a new context with the auditUser.
func WithAuditUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, loggerAuditIDContextKey, userID)
}

// RetrieveAuditUserIDFromContext returns the auditUser associated with a context. If the
// context has no auditUser, the empty user model is returned
func RetrieveAuditUserIDFromContext(ctx context.Context) uuid.UUID {
	var userID uuid.UUID
	auditUserID, ok := ctx.Value(loggerAuditIDContextKey).(uuid.UUID)
	if !ok {
		return userID
	}
	return auditUserID
}
