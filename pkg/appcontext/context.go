package appcontext

import (
	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"
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
}

type appContext struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewAppContext creates a new AppConfig
func NewAppContext(db *pop.Connection, logger *zap.Logger) AppContext {
	return &appContext{
		db:     db,
		logger: logger,
	}
}

func (ac *appContext) DB() *pop.Connection {
	return ac.db
}

func (ac *appContext) Logger() *zap.Logger {
	return ac.logger
}

func (ac *appContext) NewTransaction(fn func(appCtx AppContext) error) error {
	// We need to make sure we don't start a new transaction since pop
	// doesn't support nested transactions
	if ac.db.TX != nil {
		return fn(ac)
	}
	return ac.db.Transaction(func(tx *pop.Connection) error {
		txnAppCtx := NewAppContext(tx, ac.logger)
		return fn(txnAppCtx)
	})
}
