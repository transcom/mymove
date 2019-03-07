package internalapi

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

type errUserMissing struct {
	userID uuid.UUID
}

func (e *errUserMissing) Error() string {
	return fmt.Sprintf("Unable to fetch user: %s", e.userID.String())
}

// DPSAuthGetCookieURLHandler generates the URL to redirect to that begins the authentication process for DPS
type DPSAuthGetCookieURLHandler struct {
	handlers.HandlerContext
}

// Handle generates the URL to redirect to that begins the authentication process for DPS
func (h DPSAuthGetCookieURLHandler) Handle(params dps_auth.GetCookieURLParams) middleware.Responder {
	// TODO: Currently, only whitelisted DPS users can access this endpoint because
	//   1. The /dps_cookie page is ungated on the front-end. The restriction here will prevent
	//      people from actually doing anything useful with that page.
	//   2. This feature is in testing and isn't open to service members yet.
	// However, when we're able to gate the /dps_cookie page on the front end and/or we're ready to
	// launch this feature, all service members should be able to access this endpoint.
	// Important: Only DPS users should ever be allowed to set parameters though (for testing).
	// Service members should never be allowed to set params and only be allowed to use the default params.
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsDpsUser() {
		return dps_auth.NewGetCookieURLForbidden()
	}

	dpsParams := h.DPSAuthParams()
	portSuffix := ""
	if dpsParams.SDDCPort != "" {
		portSuffix = fmt.Sprintf(":%s", dpsParams.SDDCPort)
	}
	url, err := url.Parse(fmt.Sprintf("%s://%s%s%s", dpsParams.SDDCProtocol, dpsParams.SDDCHostname, portSuffix, dpsauth.SetCookiePath))
	if err != nil {
		h.Logger().Error("Parsing cookie URL", zap.Error(err))
		return dps_auth.NewGetCookieURLInternalServerError()
	}

	token, err := h.generateToken(params)
	if err != nil {
		switch e := err.(type) {
		case *errUserMissing:
			h.Logger().Error("Generating token for cookie URL", zap.Error(err), zap.String("user", e.userID.String()))
		default:
			h.Logger().Error("Generating token for cookie URL", zap.Error(err))
		}

		return dps_auth.NewGetCookieURLInternalServerError()
	}

	q := url.Query()
	q.Set("token", token)
	url.RawQuery = q.Encode()

	payload := internalmessages.DPSAuthCookieURLPayload{CookieURL: strfmt.URI(url.String())}
	return dps_auth.NewGetCookieURLOK().WithPayload(&payload)
}

func (h DPSAuthGetCookieURLHandler) generateToken(params dps_auth.GetCookieURLParams) (string, error) {
	dpsParams := h.DPSAuthParams()
	cookieName := dpsParams.CookieName
	if params.CookieName != nil {
		cookieName = *params.CookieName
	}

	dpsRedirectURL := dpsParams.DPSRedirectURL
	if params.DpsRedirectURL != nil {
		dpsRedirectURL = *params.DpsRedirectURL
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	user, err := models.GetUser(h.DB(), session.UserID)
	if err != nil {
		return "", &errUserMissing{userID: session.UserID}
	}

	return dpsauth.GenerateToken(user.LoginGovUUID.String(), cookieName, dpsRedirectURL, h.DPSAuthParams().SecretKey)
}
