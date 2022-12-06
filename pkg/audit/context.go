package audit

import (
	"context"

	"github.com/transcom/mymove/pkg/models"
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

// WithAuditUser returns a new context with the auditUser.
func WithAuditUser(ctx context.Context, auditUser models.User) context.Context {
	return context.WithValue(ctx, loggerAuditIDContextKey, auditUser)
}

// RetrieveAuditUserFromContext returns the auditUser associated with a context. If the
// context has no auditUser, the empty user model is returned
func RetrieveAuditUserFromContext(ctx context.Context) models.User {
	var user models.User
	auditUser, ok := ctx.Value(loggerAuditIDContextKey).(models.User)
	if !ok {
		return user
	}
	return auditUser
}
