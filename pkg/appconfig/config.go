package appconfig

import (
	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"
)

// AppConfig should be the first argument passed to all stateless
// methods and contains all necessary config for the app
//
// This is a separate package so that all parts of the app can import
// it without creating an import cycle
type AppConfig interface {
	DB() *pop.Connection
	Logger() *zap.Logger
	NewTransaction(func(appCfg AppConfig) error) error
}

type appConfig struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewAppConfig creates a new AppConfig
func NewAppConfig(db *pop.Connection, logger *zap.Logger) AppConfig {
	return &appConfig{
		db:     db,
		logger: logger,
	}
}

func (ac *appConfig) DB() *pop.Connection {
	return ac.db
}

func (ac *appConfig) Logger() *zap.Logger {
	return ac.logger
}

func (ac *appConfig) NewTransaction(fn func(appCfg AppConfig) error) error {
	// OMFG I fucking hate go and its lack of nested transactions
	if ac.db.TX != nil {
		return fn(ac)
	}
	return ac.db.Transaction(func(tx *pop.Connection) error {
		txnAppCfg := NewAppConfig(tx, ac.logger)
		return fn(txnAppCfg)
	})
}
