package internalapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

const cookieExpiresInHours = 1

// SetDPSAuthCookieOKResponder is a custom responder that sets the DPS authentication cookie
// when writing the response
type SetDPSAuthCookieOKResponder struct {
	request     *http.Request
	redirectURL string
}

// NewSetDPSAuthCookieOKResponder creates a new SetDPSAuthCookieOKResponder
func NewSetDPSAuthCookieOKResponder(r *http.Request, url string) *SetDPSAuthCookieOKResponder {
	return &SetDPSAuthCookieOKResponder{request: r, redirectURL: url}
}

// WriteResponse to the client
func (o *SetDPSAuthCookieOKResponder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	//rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	http.Redirect(rw, o.request, o.redirectURL, http.StatusSeeOther)
}

// DPSAuthCookieHandler handles the authentication process for DPS
type DPSAuthCookieHandler struct {
	handlers.HandlerContext
}

// Handle begins the authentication process for DPS
func (h DPSAuthCookieHandler) Handle(params dps_auth.SetDPSAuthCookieParams) middleware.Responder {
	cookieName := "DPS"
	if params.CookieName != nil {
		cookieName = *params.CookieName
	}

	request := params.HTTPRequest
	session := auth.SessionFromRequestContext(request)
	user, err := models.GetUser(h.DB(), session.UserID)
	if err != nil {
		h.Logger().Error("Fetching user", zap.Error(err))
		return dps_auth.NewSetDPSAuthCookieInternalServerError()
	}
	redirectURL, err := url.Parse(fmt.Sprintf("http://%s:8080%s", h.SDDCHostname(), dpsauth.SetCookiePath))
	if err != nil {
		h.Logger().Error("Creating redirect URL", zap.Error(err))
		return dps_auth.NewSetDPSAuthCookieInternalServerError()
	}
	q := redirectURL.Query()
	q.Set("login_gov_id", user.LoginGovUUID.String())
	q.Set("cookie_name", cookieName)
	redirectURL.RawQuery = q.Encode()
	fmt.Println(redirectURL.String())
	return NewSetDPSAuthCookieOKResponder(request, redirectURL.String())
}
