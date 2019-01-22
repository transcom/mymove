package internalapi

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
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
	// Only DPS users can set the cookie name and redirect URL for testing purposes
	if params.CookieName != nil || params.DpsRedirectURL != nil {
		session := auth.SessionFromRequestContext(params.HTTPRequest)
		if !session.CanAccessFeature(auth.FeatureDPS) {
			return dps_auth.NewGetCookieURLForbidden()
		}
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
