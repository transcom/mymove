package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/audit"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/trace"
)

// HandlerConfig provides access to all the contextual references
// needed by individual handlers
type HandlerConfig interface {
	AppContextFromRequest(r *http.Request) appcontext.AppContext
	AuditableAppContextFromRequestWithErrors(
		*http.Request,
		func(appCtx appcontext.AppContext) (middleware.Responder, error),
	) middleware.Responder
	FileStorer() storage.FileStorer
	NotificationSender() notifications.NotificationSender
	Planner() route.Planner
	HHGPlanner() route.Planner
	DtodPlanner() route.Planner
	CookieSecret() string
	IWSPersonLookup() iws.PersonLookup
	SendProductionInvoice() bool
	UseSecureCookie() bool
	AppNames() auth.ApplicationServername
	GetFeatureFlag(name string) bool

	GexSender() services.GexSender
	ICNSequencer() sequence.Sequencer
	GetTraceIDFromRequest(r *http.Request) uuid.UUID
	SessionManager(session *auth.Session) *scs.SessionManager
	GetSessionManagers() [3]*scs.SessionManager
	GetMilSessionManager() *scs.SessionManager
	GetAdminSessionManager() *scs.SessionManager
	GetOfficeSessionManager() *scs.SessionManager
}

// FeatureFlag struct for feature flags
type FeatureFlag struct {
	Name   string
	Active bool
}

// A single Config is passed to each handler. This should be
// instantiated by NewHandlerConfig
type Config struct {
	db                    *pop.Connection
	logger                *zap.Logger
	cookieSecret          string
	planner               route.Planner
	hhgPlanner            route.Planner
	dtodPlanner           route.Planner
	storage               storage.FileStorer
	notificationSender    notifications.NotificationSender
	iwsPersonLookup       iws.PersonLookup
	sendProductionInvoice bool
	senderToGex           services.GexSender
	icnSequencer          sequence.Sequencer
	useSecureCookie       bool
	appNames              auth.ApplicationServername
	featureFlags          map[string]bool
	sessionManagers       [3]*scs.SessionManager
}

// NewHandlerConfig returns a new HandlerConfig interface with its
// required private fields set.
func NewHandlerConfig(
	db *pop.Connection,
	logger *zap.Logger,
	cookieSecret string,
	planner route.Planner,
	hhgPlanner route.Planner,
	dtodPlanner route.Planner,
	storage storage.FileStorer,
	notificationSender notifications.NotificationSender,
	iwsPersonLookup iws.PersonLookup,
	sendProductionInvoice bool,
	senderToGex services.GexSender,
	icnSequencer sequence.Sequencer,
	useSecureCookie bool,
	appNames auth.ApplicationServername,
	featureFlags []FeatureFlag,
	sessionManagers [3]*scs.SessionManager,
) HandlerConfig {
	featureFlagMap := make(map[string]bool)
	for _, ff := range featureFlags {
		featureFlagMap[ff.Name] = ff.Active
	}
	return &Config{
		db:                    db,
		logger:                logger,
		cookieSecret:          cookieSecret,
		planner:               planner,
		hhgPlanner:            hhgPlanner,
		dtodPlanner:           dtodPlanner,
		storage:               storage,
		notificationSender:    notificationSender,
		iwsPersonLookup:       iwsPersonLookup,
		sendProductionInvoice: sendProductionInvoice,
		senderToGex:           senderToGex,
		icnSequencer:          icnSequencer,
		useSecureCookie:       useSecureCookie,
		appNames:              appNames,
		featureFlags:          featureFlagMap,
		sessionManagers:       sessionManagers,
	}
}

// AppContextFromRequest builds an AppContext from the http request
// TODO: This should eventually go away and all handlers should use AuditableAppContextFromRequestWithErrors
func (c *Config) AppContextFromRequest(r *http.Request) appcontext.AppContext {
	// use LoggerFromRequest to get the most specific logger
	return appcontext.NewAppContext(
		c.dBFromContext(r.Context()),
		c.loggerFromRequest(r),
		c.sessionFromRequest(r))
}

// AuditableAppContextFromRequestWithErrors creates a transaction and sets local
// variables for use by the auditable trigger and also allows handlers to return errors.
func (c *Config) AuditableAppContextFromRequestWithErrors(
	r *http.Request,
	handler func(appCtx appcontext.AppContext) (middleware.Responder, error),
) middleware.Responder {
	// use LoggerFromRequest to get the most specific logger
	var resp middleware.Responder
	appCtx := c.AppContextFromRequest(r)
	err := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var userID uuid.UUID
		if txnAppCtx.Session() != nil {
			userID = txnAppCtx.Session().UserID
		}
		// not sure why, but using RawQuery("SET LOCAL foo = ?",
		// thing) did not work
		err := txnAppCtx.DB().RawQuery("SET LOCAL audit.current_user_id = '" + userID.String() + "'").Exec()
		if err != nil {
			return err
		}
		eventName := audit.EventNameFromContext(r.Context())
		err = txnAppCtx.DB().RawQuery("SET LOCAL audit.current_event_name = '" + eventName + "'").Exec()
		if err != nil {
			return err
		}
		resp, err = handler(txnAppCtx)
		return err
	})
	if err != nil {
		return resp
	}
	return resp
}

func (c *Config) sessionFromRequest(r *http.Request) *auth.Session {
	return auth.SessionFromContext(r.Context())
}

func (c *Config) loggerFromRequest(r *http.Request) *zap.Logger {
	return c.loggerFromContext(r.Context())
}

// loggerFromContext returns the logger from the context. If the
// context has no appCtx.Logger(), the handlerConfig logger is returned
func (c *Config) loggerFromContext(ctx context.Context) *zap.Logger {
	logger := logging.FromContextWithoutDefault(ctx)
	if logger != nil {
		return logger
	}
	return c.logger
}

// dBFromContext returns a POP db connection for the context
func (c *Config) dBFromContext(ctx context.Context) *pop.Connection {
	return c.db.WithContext(ctx)
}

// FileStorer returns the storage to use in the current context
func (c *Config) FileStorer() storage.FileStorer {
	return c.storage
}

// SetFileStorer is a simple setter for storage private field
func (c *Config) SetFileStorer(storer storage.FileStorer) {
	c.storage = storer
}

// AppNames returns a struct of all the app names for the current environment
func (c *Config) AppNames() auth.ApplicationServername {
	return c.appNames
}

// SetAppNames is a simple setter for private field
func (c *Config) SetAppNames(appNames auth.ApplicationServername) {
	c.appNames = appNames
}

// NotificationSender returns the sender to use in the current context
func (c *Config) NotificationSender() notifications.NotificationSender {
	return c.notificationSender
}

// SetNotificationSender is a simple setter for AWS SES private field
func (c *Config) SetNotificationSender(sender notifications.NotificationSender) {
	c.notificationSender = sender
}

// Planner returns the planner for the current context
func (c *Config) Planner() route.Planner {
	return c.planner
}

// SetPlanner is a simple setter for the route.Planner private field
func (c *Config) SetPlanner(planner route.Planner) {
	c.planner = planner
}

// HHGPlanner returns the HHG planner for the current context
func (c *Config) HHGPlanner() route.Planner {
	return c.hhgPlanner
}

// DtodPlanner returns the DTOD planner for the current context
func (c *Config) DtodPlanner() route.Planner {
	return c.dtodPlanner
}

// CookieSecret returns the secret key to use when signing cookies
func (c *Config) CookieSecret() string {
	return c.cookieSecret
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (c *Config) SetCookieSecret(cookieSecret string) {
	c.cookieSecret = cookieSecret
}

func (c *Config) IWSPersonLookup() iws.PersonLookup {
	return c.iwsPersonLookup
}

func (c *Config) SetIWSPersonLookup(rbs iws.PersonLookup) {
	c.iwsPersonLookup = rbs
}

// SendProductionInvoice is a flag to notify EDI invoice generation whether it should be sent as a test or production transaction
func (c *Config) SendProductionInvoice() bool {
	return c.sendProductionInvoice
}

// Set UsageIndicator flag for use in EDI invoicing (ediinvoice pkg)
func (c *Config) SetSendProductionInvoice(sendProductionInvoice bool) {
	c.sendProductionInvoice = sendProductionInvoice
}

func (c *Config) GexSender() services.GexSender {
	return c.senderToGex
}

func (c *Config) SetGexSender(sendGexRequest services.GexSender) {
	c.senderToGex = sendGexRequest
}

func (c *Config) ICNSequencer() sequence.Sequencer {
	return c.icnSequencer
}

func (c *Config) SetICNSequencer(sequencer sequence.Sequencer) {
	c.icnSequencer = sequencer
}

// UseSecureCookie determines if the field "Secure" is set to true or false upon cookie creation
func (c *Config) UseSecureCookie() bool {
	return c.useSecureCookie
}

// Sets flag for using Secure cookie
func (c *Config) SetUseSecureCookie(useSecureCookie bool) {
	c.useSecureCookie = useSecureCookie
}

func (c *Config) SetFeatureFlag(flag FeatureFlag) {
	if c.featureFlags == nil {
		c.featureFlags = make(map[string]bool)
	}

	c.featureFlags[flag.Name] = flag.Active
}

func (c *Config) GetFeatureFlag(flag string) bool {
	if value, ok := c.featureFlags[flag]; ok {
		return value
	}
	return false
}

// GetTraceIDFromRequest returns the request traceID. It
// returns the Nil UUID if no traceid is found
func (c *Config) GetTraceIDFromRequest(r *http.Request) uuid.UUID {
	return trace.FromContext(r.Context())
}

func (c *Config) SetSessionManagers(sessionManagers [3]*scs.SessionManager) {
	c.sessionManagers = sessionManagers
}

func (c *Config) GetMilSessionManager() *scs.SessionManager {
	return c.sessionManagers[0]
}

func (c *Config) GetAdminSessionManager() *scs.SessionManager {
	return c.sessionManagers[1]
}

func (c *Config) GetOfficeSessionManager() *scs.SessionManager {
	return c.sessionManagers[2]
}

// SessionManager returns the session manager corresponding to the current app.
// A user can be signed in at the same time across multiple apps.
func (c *Config) SessionManager(session *auth.Session) *scs.SessionManager {
	if session.IsMilApp() {
		return c.GetMilSessionManager()
	} else if session.IsAdminApp() {
		return c.GetAdminSessionManager()
	} else if session.IsOfficeApp() {
		return c.GetOfficeSessionManager()
	}

	return nil
}

// GetSessionManagers returns all session managers
func (c *Config) GetSessionManagers() [3]*scs.SessionManager {
	return c.sessionManagers
}
