package internalapi

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/handlers"
	"go.uber.org/zap"
)

const cookieExpiresInHours = 1

// SetDPSAuthCookieOKResponder is a custom responder that sets the DPS authentication cookie
// when writing the response
type SetDPSAuthCookieOKResponder struct {
	cookie http.Cookie
}

// NewSetDPSAuthCookieOKResponder creates a new SetDPSAuthCookieOKResponder
func NewSetDPSAuthCookieOKResponder(cookie http.Cookie) *SetDPSAuthCookieOKResponder {
	return &SetDPSAuthCookieOKResponder{cookie: cookie}
}

// WriteResponse to the client
func (o *SetDPSAuthCookieOKResponder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	http.SetCookie(rw, &o.cookie)

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses
	rw.WriteHeader(200)
}

// DPSAuthCookieHandler handles the authentication process for DPS
type DPSAuthCookieHandler struct {
	handlers.HandlerContext
}

// Handle sets the cookie necessary for beginning the authentication process for DPS
func (h DPSAuthCookieHandler) Handle(params dps_auth.SetDPSAuthCookieParams) middleware.Responder {
	cookieName := "DPS"
	if params.CookieName != nil {
		cookieName = *params.CookieName
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	cookie, err := dpsauth.UserIDToCookie(session.ServiceMemberID.String())
	if err != nil {
		h.Logger().Error("Converting user ID to cookie value", zap.Error(err))
		return dps_auth.NewSetDPSAuthCookieInternalServerError()
	}

	cookie.Name = cookieName
	cookie.Domain = ".sddc.army.mil"
	return NewSetDPSAuthCookieOKResponder(*cookie)
}
