package dpsauth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// SetCookiePath is the path for this resource
const SetCookiePath = "/dps_auth/set_cookie"

// Claims contains information passed to the endpoint that sets the DPS auth cookie
type Claims struct {
	jwt.RegisteredClaims
	CookieName     string
	DPSRedirectURL string
}

// SetCookieHandler handles setting the DPS auth cookie and redirecting to DPS
type SetCookieHandler struct {
	secretKey     string
	cookieDomain  string
	cookieSecret  []byte
	cookieExpires int
}

// NewSetCookieHandler creates a new SetCookieHandler
func NewSetCookieHandler(secretKey string, cookieDomain string, cookieSecret []byte, cookieExpires int) SetCookieHandler {
	return SetCookieHandler{
		secretKey:     secretKey,
		cookieDomain:  cookieDomain,
		cookieSecret:  cookieSecret,
		cookieExpires: cookieExpires}
}

func (h SetCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	claims, err := ParseToken(r.URL.Query().Get("token"), h.secretKey)
	if err != nil {
		logger.Error("Parsing token", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	cookie, err := LoginGovIDToCookie(claims.RegisteredClaims.Subject, h.cookieSecret, h.cookieExpires)
	if err != nil {
		logger.Error("Converting user ID to cookie value", zap.Error(err))
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

// InitDPSAuthParams initializes the DPS Auth Params
func InitDPSAuthParams(v *viper.Viper, appnames auth.ApplicationServername) Params {
	return Params{
		SDDCProtocol:   v.GetString(cli.HTTPSDDCProtocolFlag),
		SDDCHostname:   appnames.SddcServername,
		SDDCPort:       v.GetInt(cli.HTTPSDDCPortFlag),
		SecretKey:      v.GetString(cli.DPSAuthSecretKeyFlag),
		DPSRedirectURL: v.GetString(cli.DPSRedirectURLFlag),
		CookieName:     v.GetString(cli.DPSCookieNameFlag),
		CookieDomain:   v.GetString(cli.DPSCookieDomainFlag),
		CookieSecret:   []byte(v.GetString(cli.DPSAuthCookieSecretKeyFlag)),
		CookieExpires:  v.GetInt(cli.DPSCookieExpiresInMinutesFlag),
	}
}
