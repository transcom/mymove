package dpsauth

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// SetCookiePath is the path for this resource
const SetCookiePath = "/dps_auth/set_cookie"

// Claims contains information passed to the endpoint that sets the DPS auth cookie
type Claims struct {
	jwt.StandardClaims
	CookieName     string
	DPSRedirectURL string
}

// SetCookieHandler handles setting the DPS auth cookie and redirecting to DPS
type SetCookieHandler struct {
	logger       *zap.Logger
	secretKey    string
	cookieDomain string
}

// NewSetCookieHandler creates a new SetCookieHandler
func NewSetCookieHandler(logger *zap.Logger, secretKey string, cookieDomain string) SetCookieHandler {
	return SetCookieHandler{logger: logger, secretKey: secretKey, cookieDomain: cookieDomain}
}

func (h SetCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := ParseToken(r.URL.Query().Get("token"), h.secretKey)
	if err != nil {
		h.logger.Error("Parsing token", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	cookie, err := LoginGovIDToCookie(claims.StandardClaims.Subject)
	if err != nil {
		h.logger.Error("Converting user ID to cookie value", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	cookie.Name = claims.CookieName
	cookie.Domain = h.cookieDomain
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	roleCookie := http.Cookie{
		Domain:  h.cookieDomain,
		Expires: cookie.Expires,
		Name:    "DPSETAROLE",
		Path:    "/",
		Value:   "dodcustomer",
	}
	http.SetCookie(w, &roleCookie)

	http.Redirect(w, r, claims.DPSRedirectURL, http.StatusSeeOther)
}
