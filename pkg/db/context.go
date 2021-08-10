package db

import (
	"context"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/logging"
)

type contextKey string

var dbContextKey = contextKey("db")

// NewContext returns a new context with a logger.
func NewContext(ctx context.Context, db *pop.Connection) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

// FromContext returns a db associated with a context, if any.
func FromContext(ctx context.Context) *pop.Connection {
	logger := logging.FromContext(ctx)
	db, ok := ctx.Value(dbContextKey).(*pop.Connection)
	if !ok {
		logger.Panic("Cannot get db connection from context")
	}
	return db
}
