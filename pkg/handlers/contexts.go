package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext provides access to all the contextual references needed by individual handlers
//go:generate mockery -name HandlerContext
type HandlerContext interface {
	DB() *pop.Connection
	SessionAndLoggerFromContext(ctx context.Context) (*auth.Session, Logger)
	SessionAndLoggerFromRequest(r *http.Request) (*auth.Session, Logger)
	SessionFromRequest(r *http.Request) *auth.Session
	SessionFromContext(ctx context.Context) *auth.Session
	LoggerFromContext(ctx context.Context) Logger
	LoggerFromRequest(r *http.Request) Logger
	FileStorer() storage.FileStorer
	SetFileStorer(storer storage.FileStorer)
	NotificationSender() notifications.NotificationSender
	SetNotificationSender(sender notifications.NotificationSender)
	Planner() route.Planner
	SetPlanner(planner route.Planner)
	GHCPlanner() route.Planner
	SetGHCPlanner(planner route.Planner)
	CookieSecret() string
	SetCookieSecret(secret string)
	IWSPersonLookup() iws.PersonLookup
	SetIWSPersonLookup(rbs iws.PersonLookup)
	SendProductionInvoice() bool
	SetSendProductionInvoice(sendProductionInvoice bool)
	UseSecureCookie() bool
	SetUseSecureCookie(useSecureCookie bool)
	SetAppNames(appNames auth.ApplicationServername)
	AppNames() auth.ApplicationServername
	SetFeatureFlag(flags FeatureFlag)
	GetFeatureFlag(name string) bool

	GexSender() services.GexSender
	SetGexSender(gexSender services.GexSender)
	ICNSequencer() sequence.Sequencer
	SetICNSequencer(sequencer sequence.Sequencer)
	DPSAuthParams() dpsauth.Params
	SetDPSAuthParams(params dpsauth.Params)
	SetTraceID(traceID uuid.UUID)
	GetTraceID() uuid.UUID
	SetSessionManagers(sessionManagers [3]*scs.SessionManager)
	SessionManager(session *auth.Session) *scs.SessionManager
}

// FeatureFlag struct for feature flags
type FeatureFlag struct {
	Name   string
	Active bool
}

// A single handlerContext is passed to each handler
type handlerContext struct {
	db                    *pop.Connection
	logger                Logger
	cookieSecret          string
	planner               route.Planner
	ghcPlanner            route.Planner
	storage               storage.FileStorer
	notificationSender    notifications.NotificationSender
	iwsPersonLookup       iws.PersonLookup
	sendProductionInvoice bool
	dpsAuthParams         dpsauth.Params
	senderToGex           services.GexSender
	icnSequencer          sequence.Sequencer
	useSecureCookie       bool
	appNames              auth.ApplicationServername
	featureFlags          map[string]bool
	traceID               uuid.UUID
	sessionManagers       [3]*scs.SessionManager
}

// NewHandlerContext returns a new handlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger Logger) HandlerContext {
	return &handlerContext{
		db:     db,
		logger: logger,
	}
}

func (hctx *handlerContext) SessionAndLoggerFromRequest(r *http.Request) (*auth.Session, Logger) {
	return hctx.SessionAndLoggerFromContext(r.Context())
}

func (hctx *handlerContext) SessionAndLoggerFromContext(ctx context.Context) (*auth.Session, Logger) {
	return auth.SessionFromContext(ctx), hctx.LoggerFromContext(ctx)
}

func (hctx *handlerContext) SessionFromRequest(r *http.Request) *auth.Session {
	return auth.SessionFromContext(r.Context())
}

func (hctx *handlerContext) SessionFromContext(ctx context.Context) *auth.Session {
	return auth.SessionFromContext(ctx)
}

func (hctx *handlerContext) LoggerFromRequest(r *http.Request) Logger {
	return hctx.LoggerFromContext(r.Context())
}

func (hctx *handlerContext) LoggerFromContext(ctx context.Context) Logger {
	if logger, ok := logging.FromContext(ctx).(Logger); ok {
		return logger
	}
	return hctx.logger
}

// DB returns a POP db connection for the context
func (hctx *handlerContext) DB() *pop.Connection {
	return hctx.db
}

// FileStorer returns the storage to use in the current context
func (hctx *handlerContext) FileStorer() storage.FileStorer {
	return hctx.storage
}

// SetFileStorer is a simple setter for storage private field
func (hctx *handlerContext) SetFileStorer(storer storage.FileStorer) {
	hctx.storage = storer
}

// AppNames returns a struct of all the app names for the current environment
func (hctx *handlerContext) AppNames() auth.ApplicationServername {
	return hctx.appNames
}

// SetAppNames is a simple setter for private field
func (hctx *handlerContext) SetAppNames(appNames auth.ApplicationServername) {
	hctx.appNames = appNames
}

// NotificationSender returns the sender to use in the current context
func (hctx *handlerContext) NotificationSender() notifications.NotificationSender {
	return hctx.notificationSender
}

// SetNotificationSender is a simple setter for AWS SES private field
func (hctx *handlerContext) SetNotificationSender(sender notifications.NotificationSender) {
	hctx.notificationSender = sender
}

// Planner returns the planner for the current context
func (hctx *handlerContext) Planner() route.Planner {
	return hctx.planner
}

// SetPlanner is a simple setter for the route.Planner private field
func (hctx *handlerContext) SetPlanner(planner route.Planner) {
	hctx.planner = planner
}

// GHCPlanner returns the GHC planner for the current context
func (hctx *handlerContext) GHCPlanner() route.Planner {
	return hctx.ghcPlanner
}

// SetGHCPlanner is a simple setter for the route.Planner private field
func (hctx *handlerContext) SetGHCPlanner(ghcPlanner route.Planner) {
	hctx.ghcPlanner = ghcPlanner
}

// CookieSecret returns the secret key to use when signing cookies
func (hctx *handlerContext) CookieSecret() string {
	return hctx.cookieSecret
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (hctx *handlerContext) SetCookieSecret(cookieSecret string) {
	hctx.cookieSecret = cookieSecret
}

func (hctx *handlerContext) IWSPersonLookup() iws.PersonLookup {
	return hctx.iwsPersonLookup
}

func (hctx *handlerContext) SetIWSPersonLookup(rbs iws.PersonLookup) {
	hctx.iwsPersonLookup = rbs
}

// SendProductionInvoice is a flag to notify EDI invoice generation whether it should be sent as a test or production transaction
func (hctx *handlerContext) SendProductionInvoice() bool {
	return hctx.sendProductionInvoice
}

// Set UsageIndicator flag for use in EDI invoicing (ediinvoice pkg)
func (hctx *handlerContext) SetSendProductionInvoice(sendProductionInvoice bool) {
	hctx.sendProductionInvoice = sendProductionInvoice
}

func (hctx *handlerContext) GexSender() services.GexSender {
	return hctx.senderToGex
}

func (hctx *handlerContext) SetGexSender(sendGexRequest services.GexSender) {
	hctx.senderToGex = sendGexRequest
}

func (hctx *handlerContext) ICNSequencer() sequence.Sequencer {
	return hctx.icnSequencer
}

func (hctx *handlerContext) SetICNSequencer(sequencer sequence.Sequencer) {
	hctx.icnSequencer = sequencer
}

func (hctx *handlerContext) DPSAuthParams() dpsauth.Params {
	return hctx.dpsAuthParams
}

func (hctx *handlerContext) SetDPSAuthParams(params dpsauth.Params) {
	hctx.dpsAuthParams = params
}

// UseSecureCookie determines if the field "Secure" is set to true or false upon cookie creation
func (hctx *handlerContext) UseSecureCookie() bool {
	return hctx.useSecureCookie
}

// Sets flag for using Secure cookie
func (hctx *handlerContext) SetUseSecureCookie(useSecureCookie bool) {
	hctx.useSecureCookie = useSecureCookie
}

func (hctx *handlerContext) SetFeatureFlag(flag FeatureFlag) {
	if hctx.featureFlags == nil {
		hctx.featureFlags = make(map[string]bool)
	}

	hctx.featureFlags[flag.Name] = flag.Active
}

func (hctx *handlerContext) GetFeatureFlag(flag string) bool {
	if value, ok := hctx.featureFlags[flag]; ok {
		return value
	}
	return false
}

func (hctx *handlerContext) SetTraceID(traceID uuid.UUID) {
	hctx.traceID = traceID
}

func (hctx *handlerContext) GetTraceID() uuid.UUID {
	return hctx.traceID
}

func (hctx *handlerContext) SetSessionManagers(sessionManagers [3]*scs.SessionManager) {
	hctx.sessionManagers = sessionManagers
}

// SessionManager returns the session manager corresponding to the current app.
// A user can be signed in at the same time across multiple apps.
func (hctx *handlerContext) SessionManager(session *auth.Session) *scs.SessionManager {
	if session.IsMilApp() {
		return hctx.sessionManagers[0]
	} else if session.IsAdminApp() {
		return hctx.sessionManagers[1]
	} else if session.IsOfficeApp() {
		return hctx.sessionManagers[2]
	}

	return nil
}
