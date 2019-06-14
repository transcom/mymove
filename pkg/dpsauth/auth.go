package dpsauth

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
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
	logger        Logger
	secretKey     string
	cookieDomain  string
	cookieSecret  []byte
	cookieExpires int
}

// NewSetCookieHandler creates a new SetCookieHandler
func NewSetCookieHandler(logger Logger, secretKey string, cookieDomain string, cookieSecret []byte, cookieExpires int) SetCookieHandler {
	return SetCookieHandler{
		logger:        logger,
		secretKey:     secretKey,
		cookieDomain:  cookieDomain,
		cookieSecret:  cookieSecret,
		cookieExpires: cookieExpires}
}

func (h SetCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := ParseToken(r.URL.Query().Get("token"), h.secretKey)
	if err != nil {
		h.logger.Error("Parsing token", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	cookie, err := LoginGovIDToCookie(claims.StandardClaims.Subject, h.cookieSecret, h.cookieExpires)
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
