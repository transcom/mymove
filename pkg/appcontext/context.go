package appcontext

import (
	"context"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

// AppContext should be the first argument passed to all stateless
// methods and contains all necessary config for the app
//
// This is a separate package so that all parts of the app can import
// it without creating an import cycle
type AppContext interface {
	DB() *pop.Connection
	Logger() *zap.Logger
	NewTransaction(func(appCtx AppContext) error) error
	Context() context.Context
	TraceID() uuid.UUID
}

type appContext struct {
	ctx context.Context
	db  *pop.Connection
}

// NewAppContext creates a new AppConfig
func NewAppContext(ctx context.Context, db *pop.Connection) AppContext {
	return &appContext{
		ctx: ctx,
		db:  db,
	}
}

func (ac *appContext) DB() *pop.Connection {
	return ac.db
}

// Logger gets the logger from the context
func (ac *appContext) Logger() *zap.Logger {
	return logging.FromContext(ac.ctx)
}

// Context gets the underlying context
func (ac *appContext) Context() context.Context {
	return ac.ctx
}

// TraceID gets the uuid of the trace from the context
func (ac *appContext) TraceID() uuid.UUID {
	return trace.FromContext(ac.ctx)
}

// NewTransaction starts a new transaction if not already in one
func (ac *appContext) NewTransaction(fn func(appCtx AppContext) error) error {
	// We need to make sure we don't start a new transaction since pop
	// doesn't support nested transactions
	if ac.db.TX != nil {
		return fn(ac)
	}
	return ac.db.Transaction(func(tx *pop.Connection) error {
		txnAppCtx := NewAppContext(ac.ctx, tx)
		return fn(txnAppCtx)
	})
}
