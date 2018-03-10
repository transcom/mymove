package context

import (
	"context"

	"github.com/satori/go.uuid"
)

type authCtxKey string

var userIDKey = authCtxKey("user_id")
var idTokenKey = authCtxKey("id_token")

// PopulateAuthContext sets the values that the auth package wants in the context
func PopulateAuthContext(ctx context.Context, userID uuid.UUID, idToken string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, idTokenKey, idToken)

	return ctx
}

// GetUserID retrieves the UserID from the context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// GetIDToken retrieves the IDToken from the context
func GetIDToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(idTokenKey).(string)
	return token, ok
}
