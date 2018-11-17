package handlers

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging/hnyzap"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"go.uber.org/zap"
)

// SendProdInvoice is a type marker for dependency injection
type SendProdInvoice bool

// HandlerContext provides access to all the contextual references needed by individual handlers
type HandlerContext interface {
	DB() *pop.Connection
	Logger() *zap.Logger
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
	FetchServiceMember() services.FetchServiceMember
	SetFetchServiceMember(services.FetchServiceMember)
	FetchDocument() services.FetchDocument
	SetFetchDocument(services.FetchDocument)
	FetchUpload() services.FetchUpload
	SetFetchUpload(services.FetchUpload)
	IWSRealTimeBrokerService() iws.RealTimeBrokerService
	SetIWSRealTimeBrokerService(rbs iws.RealTimeBrokerService)
	SendProductionInvoice() bool
	SetSendProductionInvoice(send SendProdInvoice)
	DPSAuthParams() *dpsauth.Params
	SetDPSAuthParams(params *dpsauth.Params)
}

// A single handlerContext is passed to each handler
type handlerContext struct {
	db                       *pop.Connection
	logger                   *zap.Logger
	cookieSecret             string
	noSessionTimeout         bool
	planner                  route.Planner
	storage                  storage.FileStorer
	notificationSender       notifications.NotificationSender
	fetchServiceMember       services.FetchServiceMember
	fetchDocument            services.FetchDocument
	fetchUpload              services.FetchUpload
	iwsRealTimeBrokerService iws.RealTimeBrokerService
	sendProductionInvoice    bool
	dpsAuthParams            *dpsauth.Params
}

// NewHandlerContext returns a new handlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return &handlerContext{
		db:     db,
		logger: logger,
	}
}

// DB returns a POP db connection for the context
func (context *handlerContext) DB() *pop.Connection {
	return context.db
}

// Logger returns the logger to use in this context
func (context *handlerContext) Logger() *zap.Logger {
	return context.logger
}

// HoneyZapLogger returns the logger capable of writing to Honeycomb to use in this context
func (context *handlerContext) HoneyZapLogger() *hnyzap.Logger {
	return &hnyzap.Logger{Logger: context.logger}
}

// FileStorer returns the storage to use in the current context
func (context *handlerContext) FileStorer() storage.FileStorer {
	return context.storage
}

// SetFileStorer is a simple setter for storage private field
func (context *handlerContext) SetFileStorer(storer storage.FileStorer) {
	context.storage = storer
}

// NotificationSender returns the sender to use in the current context
func (context *handlerContext) NotificationSender() notifications.NotificationSender {
	return context.notificationSender
}

// SetNotificationSender is a simple setter for AWS SES private field
func (context *handlerContext) SetNotificationSender(sender notifications.NotificationSender) {
	context.notificationSender = sender
}

// FetchServiceMember is a simple setter for FetchServiceMember service private field
func (context *handlerContext) FetchServiceMember() services.FetchServiceMember {
	return context.fetchServiceMember
}

// SetFetchServiceMember is a simple setter for FetchServiceMember service private field
func (context *handlerContext) SetFetchServiceMember(fetcher services.FetchServiceMember) {
	context.fetchServiceMember = fetcher
}

// FetchDocument is a simple setter for FetchDocument service private field
func (context *handlerContext) FetchDocument() services.FetchDocument {
	return context.fetchDocument
}

// SetFetchDocument is a simple setter for FetchDocument service private field
func (context *handlerContext) SetFetchDocument(fetcher services.FetchDocument) {
	context.fetchDocument = fetcher
}

// FetchUpload is a simple setter for FetchUpload service private field
func (context *handlerContext) FetchUpload() services.FetchUpload {
	return context.fetchUpload
}

// SetFetchUpload is a simple setter for FetchUpload service private field
func (context *handlerContext) SetFetchUpload(fetcher services.FetchUpload) {
	context.fetchUpload = fetcher
}

// Planner is a simple setter for the route.Planner private field
func (context *handlerContext) Planner() route.Planner {
	return context.planner
}

// SetPlanner is a simple setter for the route.Planner private field
func (context *handlerContext) SetPlanner(planner route.Planner) {
	context.planner = planner
}

// CookieSecret returns the secret key to use when signing cookies
func (context *handlerContext) CookieSecret() string {
	return context.cookieSecret
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (context *handlerContext) SetCookieSecret(cookieSecret string) {
	context.cookieSecret = cookieSecret
}

// NoSessionTimeout is a flag which, when true, indicates that sessions should not timeout. Used in dev.
func (context *handlerContext) NoSessionTimeout() bool {
	return context.noSessionTimeout
}

// SetNoSessionTimeout is a simple setter for the noSessionTimeout private Field
func (context *handlerContext) SetNoSessionTimeout() {
	context.noSessionTimeout = true
}

func (context *handlerContext) IWSRealTimeBrokerService() iws.RealTimeBrokerService {
	return context.iwsRealTimeBrokerService
}

func (context *handlerContext) SetIWSRealTimeBrokerService(rbs iws.RealTimeBrokerService) {
	context.iwsRealTimeBrokerService = rbs
}

// InvoiceIsATest is a flag to notify EDI invoice generation whether it should be sent as a test transaction
func (context *handlerContext) SendProductionInvoice() bool {
	return context.sendProductionInvoice
}

// Set UsageIndicator flag for use in EDI invoicing (ediinvoice pkg)
func (context *handlerContext) SetSendProductionInvoice(send SendProdInvoice) {
	context.sendProductionInvoice = bool(send)
}

func (context *handlerContext) DPSAuthParams() *dpsauth.Params {
	return context.dpsAuthParams
}

func (context *handlerContext) SetDPSAuthParams(params *dpsauth.Params) {
	context.dpsAuthParams = params
}
