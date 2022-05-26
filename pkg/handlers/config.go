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
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/trace"
)

// HandlerConfig provides access to all the contextual references needed by individual handlers
//go:generate mockery --name HandlerConfig --disable-version-string
type HandlerConfig interface {
	AppContextFromRequest(r *http.Request) appcontext.AppContext
	AuditableAppContextFromRequestWithErrors(
		*http.Request,
		func(appCtx appcontext.AppContext) (middleware.Responder, error),
	) middleware.Responder
	FileStorer() storage.FileStorer
	SetFileStorer(storer storage.FileStorer)
	NotificationSender() notifications.NotificationSender
	SetNotificationSender(sender notifications.NotificationSender)
	Planner() route.Planner
	SetPlanner(planner route.Planner)
	HHGPlanner() route.Planner
	SetHHGPlanner(planner route.Planner)
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
	GetTraceIDFromRequest(r *http.Request) uuid.UUID
	SetSessionManagers(sessionManagers [3]*scs.SessionManager)
	SessionManager(session *auth.Session) *scs.SessionManager
}

// FeatureFlag struct for feature flags
type FeatureFlag struct {
	Name   string
	Active bool
}

// A single handlerConfig is passed to each handler
type handlerConfig struct {
	db                    *pop.Connection
	logger                *zap.Logger
	cookieSecret          string
	planner               route.Planner
	hhgPlanner            route.Planner
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
	sessionManagers       [3]*scs.SessionManager
}

// NewHandlerConfig returns a new handlerConfig with its required private fields set.
func NewHandlerConfig(db *pop.Connection, logger *zap.Logger) HandlerConfig {
	return &handlerConfig{
		db:     db,
		logger: logger,
	}
}

// AppContextFromRequest builds an AppContext from the http request
// TODO: This should eventually go away and all handlers should use AuditableAppContextFromRequestWithErrors
func (h *handlerConfig) AppContextFromRequest(r *http.Request) appcontext.AppContext {
	// use LoggerFromRequest to get the most specific logger
	return appcontext.NewAppContext(
		h.dBFromContext(r.Context()),
		h.loggerFromRequest(r),
		h.sessionFromRequest(r))
}

// AuditableAppContextFromRequestWithErrors creates a transaction and sets local
// variables for use by the auditable trigger and also allows handlers to return errors.
func (h *handlerConfig) AuditableAppContextFromRequestWithErrors(
	r *http.Request,
	handler func(appCtx appcontext.AppContext) (middleware.Responder, error),
) middleware.Responder {
	// use LoggerFromRequest to get the most specific logger
	var resp middleware.Responder
	appCtx := h.AppContextFromRequest(r)
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

func (h *handlerConfig) sessionFromRequest(r *http.Request) *auth.Session {
	return auth.SessionFromContext(r.Context())
}

func (h *handlerConfig) loggerFromRequest(r *http.Request) *zap.Logger {
	return h.loggerFromContext(r.Context())
}

// LoggerFromContext returns the logger from the context. If the
// context has no appCtx.Logger(), the handlerConfig logger is returned
func (h *handlerConfig) loggerFromContext(ctx context.Context) *zap.Logger {
	logger := logging.FromContextWithoutDefault(ctx)
	if logger != nil {
		return logger
	}
	return h.logger
}

// dBFromContext returns a POP db connection for the context
func (h *handlerConfig) dBFromContext(ctx context.Context) *pop.Connection {
	return h.db.WithContext(ctx)
}

// FileStorer returns the storage to use in the current context
func (h *handlerConfig) FileStorer() storage.FileStorer {
	return h.storage
}

// SetFileStorer is a simple setter for storage private field
func (h *handlerConfig) SetFileStorer(storer storage.FileStorer) {
	h.storage = storer
}

// AppNames returns a struct of all the app names for the current environment
func (h *handlerConfig) AppNames() auth.ApplicationServername {
	return h.appNames
}

// SetAppNames is a simple setter for private field
func (h *handlerConfig) SetAppNames(appNames auth.ApplicationServername) {
	h.appNames = appNames
}

// NotificationSender returns the sender to use in the current context
func (h *handlerConfig) NotificationSender() notifications.NotificationSender {
	return h.notificationSender
}

// SetNotificationSender is a simple setter for AWS SES private field
func (h *handlerConfig) SetNotificationSender(sender notifications.NotificationSender) {
	h.notificationSender = sender
}

// Planner returns the planner for the current context
func (h *handlerConfig) Planner() route.Planner {
	return h.planner
}

// SetPlanner is a simple setter for the route.Planner private field
func (h *handlerConfig) SetPlanner(planner route.Planner) {
	h.planner = planner
}

// HHGPlanner returns the HHG planner for the current context
func (h *handlerConfig) HHGPlanner() route.Planner {
	return h.hhgPlanner
}

// SetHHGPlanner is a simple setter for the route.Planner private field
func (h *handlerConfig) SetHHGPlanner(hhgPlanner route.Planner) {
	h.hhgPlanner = hhgPlanner
}

// CookieSecret returns the secret key to use when signing cookies
func (h *handlerConfig) CookieSecret() string {
	return h.cookieSecret
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (h *handlerConfig) SetCookieSecret(cookieSecret string) {
	h.cookieSecret = cookieSecret
}

func (h *handlerConfig) IWSPersonLookup() iws.PersonLookup {
	return h.iwsPersonLookup
}

func (h *handlerConfig) SetIWSPersonLookup(rbs iws.PersonLookup) {
	h.iwsPersonLookup = rbs
}

// SendProductionInvoice is a flag to notify EDI invoice generation whether it should be sent as a test or production transaction
func (h *handlerConfig) SendProductionInvoice() bool {
	return h.sendProductionInvoice
}

// Set UsageIndicator flag for use in EDI invoicing (ediinvoice pkg)
func (h *handlerConfig) SetSendProductionInvoice(sendProductionInvoice bool) {
	h.sendProductionInvoice = sendProductionInvoice
}

func (h *handlerConfig) GexSender() services.GexSender {
	return h.senderToGex
}

func (h *handlerConfig) SetGexSender(sendGexRequest services.GexSender) {
	h.senderToGex = sendGexRequest
}

func (h *handlerConfig) ICNSequencer() sequence.Sequencer {
	return h.icnSequencer
}

func (h *handlerConfig) SetICNSequencer(sequencer sequence.Sequencer) {
	h.icnSequencer = sequencer
}

func (h *handlerConfig) DPSAuthParams() dpsauth.Params {
	return h.dpsAuthParams
}

func (h *handlerConfig) SetDPSAuthParams(params dpsauth.Params) {
	h.dpsAuthParams = params
}

// UseSecureCookie determines if the field "Secure" is set to true or false upon cookie creation
func (h *handlerConfig) UseSecureCookie() bool {
	return h.useSecureCookie
}

// Sets flag for using Secure cookie
func (h *handlerConfig) SetUseSecureCookie(useSecureCookie bool) {
	h.useSecureCookie = useSecureCookie
}

func (h *handlerConfig) SetFeatureFlag(flag FeatureFlag) {
	if h.featureFlags == nil {
		h.featureFlags = make(map[string]bool)
	}

	h.featureFlags[flag.Name] = flag.Active
}

func (h *handlerConfig) GetFeatureFlag(flag string) bool {
	if value, ok := h.featureFlags[flag]; ok {
		return value
	}
	return false
}

// GetTraceIDFromRequest returns the request traceID. It
// returns the Nil UUID if no traceid is found
func (h *handlerConfig) GetTraceIDFromRequest(r *http.Request) uuid.UUID {
	return trace.FromContext(r.Context())
}

func (h *handlerConfig) SetSessionManagers(sessionManagers [3]*scs.SessionManager) {
	h.sessionManagers = sessionManagers
}

// SessionManager returns the session manager corresponding to the current app.
// A user can be signed in at the same time across multiple apps.
func (h *handlerConfig) SessionManager(session *auth.Session) *scs.SessionManager {
	if session.IsMilApp() {
		return h.sessionManagers[0]
	} else if session.IsAdminApp() {
		return h.sessionManagers[1]
	} else if session.IsOfficeApp() {
		return h.sessionManagers[2]
	}

	return nil
}
