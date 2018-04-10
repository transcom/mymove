package context

import (
	"context"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// Keys for values stored in the context are supposed to be of a custom type, not string
// so as to avoid collisions with other packages.
type moveCtxKey string

var userIDKey = moveCtxKey("user_id")
var userKey = moveCtxKey("user")
var idTokenKey = moveCtxKey("id_token")

// PopulateAuthContext sets the values that the auth package wants in the context
func PopulateAuthContext(ctx context.Context, userID uuid.UUID, idToken string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, idTokenKey, idToken)
	return ctx
}

// PopulateUserModel sets a user model onto a request context
func PopulateUserModel(ctx context.Context, user models.User) context.Context {
	ctx = context.WithValue(ctx, userKey, user)
	return ctx
}

// GetUserID retrieves the UserID from the context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// GetUser retrieves the User from the context
func GetUser(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userKey).(models.User)
	return user, ok
}

// GetIDToken retrieves the IDToken from the context
func GetIDToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(idTokenKey).(string)
	return token, ok
}
