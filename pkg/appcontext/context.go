package appcontext

import (
	"context"

	"github.com/gobuffalo/pop/v6"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	apiversion "github.com/transcom/mymove/pkg/handlers/routing/api_version"
	"github.com/transcom/mymove/pkg/logging"
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
	Session() *auth.Session
	// Adding setter and getter methods for the API version flags.
	// There is probably a better way to flag which api version we are using, however this is just a quick example.
	// We need to pass the flag to the services and they have access to the AppContext.
	SetAPIVersion(flag apiversion.Flag)
	GetAPIVersion() *apiversion.Flag
}

type appContext struct {
	db         *pop.Connection
	logger     *zap.Logger
	session    *auth.Session
	apiVersion *apiversion.Flag
}

// NewAppContext creates a new AppContext
func NewAppContext(db *pop.Connection, logger *zap.Logger, session *auth.Session) AppContext {
	return &appContext{
		db:      db,
		logger:  logger,
		session: session,
	}
}

// NewAppContextFromContext creates a new AppContext taking context.Context into account and pulling some values from it
func NewAppContextFromContext(ctx context.Context, appCtx AppContext) AppContext {
	return &appContext{
		db:      appCtx.DB().WithContext(ctx),
		logger:  logging.FromContext(ctx),
		session: auth.SessionFromContext(ctx),
	}
}

// DB returns the db connection
func (ac *appContext) DB() *pop.Connection {
	return ac.db
}

// Logger returns the logger
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
		txnAppCtx := NewAppContext(tx, ac.logger, ac.session)
		return fn(txnAppCtx)
	})
}

func (ac *appContext) Session() *auth.Session {
	return ac.session
}

func (ac *appContext) SetAPIVersion(flag apiversion.Flag) {
	ac.apiVersion = &flag
}

func (ac *appContext) GetAPIVersion() *apiversion.Flag {
	return ac.apiVersion
}
