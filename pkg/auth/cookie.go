package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// UserSessionCookieName is the key at which we're storing our token cookie
const UserSessionCookieName = "session_token"

// SessionExpiryInMinutes is the number of minutes before a fallow session is harvested
const SessionExpiryInMinutes = 15
const sessionExpiryInSeconds = 15 * 60

// A representable date far in the future.  The trouble with something like https://stackoverflow.com/a/32620397
// is that it produces a date which may not marshall well into JSON which makes logging problematic
var likeForever = time.Date(9999, 1, 1, 12, 0, 0, 0, time.UTC)
var likeForeverInSeconds = 99999999

// GetExpiryTimeFromMinutes returns 'min' minutes from now
func GetExpiryTimeFromMinutes(min int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(min))
}

// SessionClaims wraps StandardClaims with some Session info
type SessionClaims struct {
	jwt.StandardClaims
	SessionValue Session
}

func signTokenStringWithUserInfo(expiry time.Time, session *Session, secret string) (string, error) {
	claims := SessionClaims{
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiry.Unix()},
		SessionValue:   *session,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		err = errors.Wrap(err, "Parsing RSA key from PEM")
		return "", err
	}

	ss, err := token.SignedString(rsaKey)
	if err != nil {
		err = errors.Wrap(err, "Signing string with token")
		return "", err
	}
	return ss, err
}

func sessionClaimsFromRequest(logger *zap.Logger, secret string, r *http.Request) (claims *SessionClaims, ok bool) {
	cookie, err := r.Cookie(UserSessionCookieName)
	if err != nil {
		// No cookie set on client
		return
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
		return &rsaKey.PublicKey, err
	})

	if err != nil || token == nil || !token.Valid {
		logger.Error("Failed token validation", zap.Error(err))
		return
	}

	// The token actually just stores a Claims interface, so we need to explicitly cast back to UserClaims
	claims, ok = token.Claims.(*SessionClaims)
	if !ok {
		logger.Error("Failed getting claims from token")
		return
	}
	return claims, ok
}

// WriteSessionCookie update the cookie for the session
func WriteSessionCookie(w http.ResponseWriter, session *Session, secret string, noSessionTimeout bool, logger *zap.Logger) {

	// Delete the cookie
	cookie := http.Cookie{
		Name:    UserSessionCookieName,
		Value:   "blank",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}

	// unless we have a valid session
	if session.IDToken != "" && session.UserID != uuid.Nil {
		expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
		maxAge := sessionExpiryInSeconds
		// Never expire token if in development
		if noSessionTimeout {
			expiry = likeForever
			maxAge = likeForeverInSeconds
		}

		ss, err := signTokenStringWithUserInfo(expiry, session, secret)
		if err != nil {
			logger.Error("Generating signed token string", zap.Error(err))
		} else {
			cookie.Value = ss
			cookie.Expires = expiry
			cookie.MaxAge = maxAge
		}
	}
	http.SetCookie(w, &cookie)
}

// SessionCookieConfig contains secret and other flags for setting session Cookies
type SessionCookieConfig struct {
	Secret    string
	NoTimeout bool
}

// SessionCookieMiddleware handle serializing and de-serializing the session betweem the user_session cookie and the request context
type SessionCookieMiddleware func(next http.Handler) http.Handler

// NewSessionCookieMiddleware is the DI provider for constructing SessionCookieMiddleware
func NewSessionCookieMiddleware(config *SessionCookieConfig, logger *zap.Logger) SessionCookieMiddleware {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			session := Session{}
			claims, ok := sessionClaimsFromRequest(logger, config.Secret, r)
			if ok {
				session = claims.SessionValue
			}

			// And put the session info into the request context
			ctx := SetSessionInRequestContext(r, &session)
			// And update the cookie. May get over-ridden later
			WriteSessionCookie(w, &session, config.Secret, config.NoTimeout, logger)
			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(mw)
	}
}
