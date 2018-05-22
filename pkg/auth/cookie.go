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
const UserSessionCookieName = "user_session"

// SessionExpiryInMinutes is the number of minutes before a fallow session is harvested
const SessionExpiryInMinutes = 15

// Taken from answer here: https://stackoverflow.com/a/32620397
var maxPossibleTimeValue = time.Unix(1<<63-62135596801, 999999999)

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

func getUserClaimsFromRequest(logger *zap.Logger, secret string, r *http.Request) (claims *SessionClaims, ok bool) {
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

// SessionCookieMiddleware handle serializing and de-serializing the session betweem the user_session cookie and the request context
func SessionCookieMiddleware(logger *zap.Logger, secret string, noSessionTimeout bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			session := Session{}
			claims, ok := getUserClaimsFromRequest(logger, secret, r)
			if ok {
				session = claims.SessionValue
			}

			// And put the session info into the request context
			ctx := SetSessionInRequestContext(r, &session)
			next.ServeHTTP(w, r.WithContext(ctx))

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
				// Never expire token if in development
				if noSessionTimeout {
					expiry = maxPossibleTimeValue
				}

				ss, err := signTokenStringWithUserInfo(expiry, &session, secret)
				if err != nil {
					logger.Error("Generating signed token string", zap.Error(err))
				} else {
					cookie.Value = ss
					cookie.Expires = expiry
				}
			}
			http.SetCookie(w, &cookie)
		}
		return http.HandlerFunc(mw)
	}
}
