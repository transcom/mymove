package context

import (
    "context"

    "github.com/satori/go.uuid"
)

type authCtxKey string

var userIDKey = authCtxKey("user_id")
var idTokenKey = authCtxKey("id_token")

func PopulateAuthContext(ctx context.Context, user_id uuid.UUID, id_token string) context.Context {
    ctx = context.WithValue(ctx, userIDKey, user_id)
    ctx = context.WithValue(ctx, idTokenKey, id_token)

    return ctx
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
    id, ok := ctx.Value(userIDKey).(uuid.UUID)
    return id, ok
}

func GetIDToken(ctx context.Context) (string, bool) {
    token, ok := ctx.Value(idTokenKey).(string)
    return token, ok
}