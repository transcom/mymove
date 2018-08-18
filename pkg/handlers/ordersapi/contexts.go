package ordersapi

import (
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db                 *pop.Connection
	logger             *zap.Logger
	cookieSecret       string
	noSessionTimeout   bool
	planner            route.Planner
	storage            storage.FileStorer
	notificationSender notifications.NotificationSender
}

// NewHandlerContext returns a new HandlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

// SetFileStorer is a simple setter for storage private field
func (context *HandlerContext) SetFileStorer(storer storage.FileStorer) {
	context.storage = storer
}

// SetNotificationSender is a simple setter for AWS SES private field
func (context *HandlerContext) SetNotificationSender(sender notifications.NotificationSender) {
	context.notificationSender = sender
}

// SetPlanner is a simple setter for the route.Planner private field
func (context *HandlerContext) SetPlanner(planner route.Planner) {
	context.planner = planner
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (context *HandlerContext) SetCookieSecret(cookieSecret string) {
	context.cookieSecret = cookieSecret
}

// SetNoSessionTimeout is a simple setter for the noSessionTimeout private Field
func (context *HandlerContext) SetNoSessionTimeout() {
	context.noSessionTimeout = true
}
