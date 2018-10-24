package handlers

import (
	"github.com/transcom/mymove/pkg/server"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// CookieUpdateResponder wraps a swagger middleware.Responder in code which sets the session_cookie
// See: https://github.com/go-swagger/go-swagger/issues/748
type CookieUpdateResponder struct {
	session          *server.Session
	cookieSecret     string
	noSessionTimeout bool
	logger           *zap.Logger
	Responder        middleware.Responder
}

// NewCookieUpdateResponder constructs a wrapper for the responder which will update cookies
func NewCookieUpdateResponder(request *http.Request, secret string, noSessionTimeout bool, logger *zap.Logger, responder middleware.Responder) middleware.Responder {
	return &CookieUpdateResponder{
		session:          server.SessionFromRequestContext(request),
		cookieSecret:     secret,
		noSessionTimeout: noSessionTimeout,
		logger:           logger,
		Responder:        responder,
	}
}

// WriteResponse updates the session cookie before writing out the details of the response
func (cur *CookieUpdateResponder) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	server.WriteSessionCookie(rw, cur.session, cur.cookieSecret, cur.noSessionTimeout, cur.logger)
	cur.Responder.WriteResponse(rw, p)
}
