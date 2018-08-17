package utils

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

// SetnotificationSender is a simple setter for AWS SES private field
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

// CookieUpdateResponder wraps a swagger middleware.Responder in code which sets the session_cookie
// See: https://github.com/go-swagger/go-swagger/issues/748
type CookieUpdateResponder struct {
	session          *auth.Session
	cookieSecret     string
	noSessionTimeout bool
	logger           *zap.Logger
	Responder        middleware.Responder
}

// NewCookieUpdateResponder constructs a wrapper for the responder which will update cookies
func NewCookieUpdateResponder(request *http.Request, secret string, noSessionTimeout bool, logger *zap.Logger, responder middleware.Responder) middleware.Responder {
	return &CookieUpdateResponder{
		session:          auth.SessionFromRequestContext(request),
		cookieSecret:     secret,
		noSessionTimeout: noSessionTimeout,
		logger:           logger,
		Responder:        responder,
	}
}

// WriteResponse updates the session cookie before writing out the details of the response
func (cur *CookieUpdateResponder) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	auth.WriteSessionCookie(rw, cur.session, cur.cookieSecret, cur.noSessionTimeout, cur.logger)
	cur.Responder.WriteResponse(rw, p)
}

// CreateFailedValidationPayload Converts the value returned by Pop's ValidateAnd* methods into a payload that can
// be returned to clients. This payload contains an object with a key,  `errors`, the
// value of which is a name -> validation error object.
func CreateFailedValidationPayload(verrs *validate.Errors) *internalmessages.InvalidRequestResponsePayload {
	errs := make(map[string]string)
	for _, key := range verrs.Keys() {
		errs[key] = strings.Join(verrs.Get(key), " ")
	}
	return &internalmessages.InvalidRequestResponsePayload{
		Errors: errs,
	}
}
