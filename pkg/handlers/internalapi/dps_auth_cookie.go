package internalapi

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// DPSAuthGetCookieURLHandler generates the URL to redirect to that begins the authentication process for DPS
type DPSAuthGetCookieURLHandler struct {
	handlers.HandlerContext
}

// Handle generates the URL to redirect to that begins the authentication process for DPS
func (h DPSAuthGetCookieURLHandler) Handle(params dps_auth.GetCookieURLParams) middleware.Responder {
	cookieName := "DPS"
	if params.CookieName != nil {
		cookieName = *params.CookieName
	}

	request := params.HTTPRequest
	session := auth.SessionFromRequestContext(request)
	user, err := models.GetUser(h.DB(), session.UserID)
	if err != nil {
		h.Logger().Error("Fetching user", zap.Error(err))
		return dps_auth.NewGetCookieURLInternalServerError()
	}

	// TODO: pass in port and protocol
	url, err := url.Parse(fmt.Sprintf("https://%s%s", h.SDDCHostname(), dpsauth.SetCookiePath))
	if err != nil {
		h.Logger().Error("Creating redirect URL", zap.Error(err))
		return dps_auth.NewGetCookieURLInternalServerError()
	}

	q := url.Query()
	q.Set("login_gov_id", user.LoginGovUUID.String())
	q.Set("cookie_name", cookieName)
	url.RawQuery = q.Encode()

	payload := internalmessages.DPSAuthCookieURLPayload{CookieURL: strfmt.URI(url.String())}
	return dps_auth.NewGetCookieURLOK().WithPayload(&payload)
}
