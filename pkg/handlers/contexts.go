package handlers

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging/hnyzap"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext provides access to all the contextual references needed by individual handlers
type HandlerContext interface {
	DB() *pop.Connection
	Logger() Logger
	HoneyZapLogger() *hnyzap.Logger
	FileStorer() storage.FileStorer
	SetFileStorer(storer storage.FileStorer)
	NotificationSender() notifications.NotificationSender
	SetNotificationSender(sender notifications.NotificationSender)
	Planner() route.Planner
	SetPlanner(planner route.Planner)
	CookieSecret() string
	SetCookieSecret(secret string)
	NoSessionTimeout() bool
	SetNoSessionTimeout()
	IWSPersonLookup() iws.PersonLookup
	SetIWSPersonLookup(rbs iws.PersonLookup)
	SendProductionInvoice() bool
	SetSendProductionInvoice(sendProductionInvoice bool)
	UseSecureCookie() bool
	SetUseSecureCookie(useSecureCookie bool)

	GexSender() services.GexSender
	SetGexSender(gexSender services.GexSender)
	ICNSequencer() sequence.Sequencer
	SetICNSequencer(sequencer sequence.Sequencer)
	DPSAuthParams() dpsauth.Params
	SetDPSAuthParams(params dpsauth.Params)
	RespondAndTraceError(ctx context.Context, err error, msg string, fields ...zap.Field) middleware.Responder
	RespondAndTraceVErrors(ctx context.Context, verrs *validate.Errors, err error, msg string, fields ...zap.Field) middleware.Responder
}

// A single handlerContext is passed to each handler
type handlerContext struct {
	db                    *pop.Connection
	logger                Logger
	cookieSecret          string
	noSessionTimeout      bool
	planner               route.Planner
	storage               storage.FileStorer
	notificationSender    notifications.NotificationSender
	iwsPersonLookup       iws.PersonLookup
	sendProductionInvoice bool
	dpsAuthParams         dpsauth.Params
	senderToGex           services.GexSender
	icnSequencer          sequence.Sequencer
	useSecureCookie       bool
}

// NewHandlerContext returns a new handlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger Logger) HandlerContext {
	return &handlerContext{
		db:     db,
		logger: logger,
	}
}

// DB returns a POP db connection for the context
func (hctx *handlerContext) DB() *pop.Connection {
	return hctx.db
}

// Logger returns the logger to use in this context
func (hctx *handlerContext) Logger() Logger {
	return hctx.logger
}

// HoneyZapLogger returns the logger capable of writing to Honeycomb to use in this context
func (hctx *handlerContext) HoneyZapLogger() *hnyzap.Logger {
	if zapLogger, ok := hctx.logger.(*zap.Logger); ok {
		return &hnyzap.Logger{Logger: zapLogger}
	}
	return nil
}

// RespondAndTraceError uses Honeycomb to trace errors and then passes response to the standard ResponseForError
func (hctx *handlerContext) RespondAndTraceError(ctx context.Context, err error, msg string, fields ...zap.Field) middleware.Responder {
	hctx.HoneyZapLogger().TraceError(ctx, msg, fields...)
	return ResponseForError(hctx.Logger(), err)
}

// RespondAndTraceVErrors uses Honeycomb to trace errors and then passes response to the standard ResponseForVErrors
func (hctx *handlerContext) RespondAndTraceVErrors(ctx context.Context, verrs *validate.Errors, err error, msg string, fields ...zap.Field) middleware.Responder {
	hctx.HoneyZapLogger().TraceError(ctx, msg, fields...)
	return ResponseForVErrors(hctx.Logger(), verrs, err)
}

// FileStorer returns the storage to use in the current context
func (hctx *handlerContext) FileStorer() storage.FileStorer {
	return hctx.storage
}

// SetFileStorer is a simple setter for storage private field
func (hctx *handlerContext) SetFileStorer(storer storage.FileStorer) {
	hctx.storage = storer
}

// NotificationSender returns the sender to use in the current context
func (hctx *handlerContext) NotificationSender() notifications.NotificationSender {
	return hctx.notificationSender
}

// SetNotificationSender is a simple setter for AWS SES private field
func (hctx *handlerContext) SetNotificationSender(sender notifications.NotificationSender) {
	hctx.notificationSender = sender
}

// Planner is a simple setter for the route.Planner private field
func (hctx *handlerContext) Planner() route.Planner {
	return hctx.planner
}

// SetPlanner is a simple setter for the route.Planner private field
func (hctx *handlerContext) SetPlanner(planner route.Planner) {
	hctx.planner = planner
}

// CookieSecret returns the secret key to use when signing cookies
func (hctx *handlerContext) CookieSecret() string {
	return hctx.cookieSecret
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (hctx *handlerContext) SetCookieSecret(cookieSecret string) {
	hctx.cookieSecret = cookieSecret
}

// NoSessionTimeout is a flag which, when true, indicates that sessions should not timeout. Used in dev.
func (hctx *handlerContext) NoSessionTimeout() bool {
	return hctx.noSessionTimeout
}

// SetNoSessionTimeout is a simple setter for the noSessionTimeout private Field
func (hctx *handlerContext) SetNoSessionTimeout() {
	hctx.noSessionTimeout = true
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
