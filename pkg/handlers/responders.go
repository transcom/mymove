package handlers

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/auth"
)

// CookieUpdateResponder wraps a swagger middleware.Responder in code which sets the session_cookie
// See: https://github.com/go-swagger/go-swagger/issues/748
type CookieUpdateResponder struct {
	session        *auth.Session
	Responder      middleware.Responder
	sessionManager auth.SessionManager
	ctx            context.Context
}

// NewCookieUpdateResponder constructs a wrapper for the responder which will update cookies
func NewCookieUpdateResponder(request *http.Request, responder middleware.Responder, sessionManager auth.SessionManager, session *auth.Session) middleware.Responder {
	return &CookieUpdateResponder{
		session:        session,
		Responder:      responder,
		sessionManager: sessionManager,
		ctx:            request.Context(),
	}
}

// WriteResponse updates the session cookie before writing out the details of the response
func (cur *CookieUpdateResponder) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	cur.sessionManager.Put(cur.ctx, "session", cur.session)
	cur.Responder.WriteResponse(rw, p)
}
