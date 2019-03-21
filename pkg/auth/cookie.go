package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/gorilla/csrf"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type errInvalidHostname struct {
	Hostname  string
	MilApp    string
	OfficeApp string
	TspApp    string
}

func (e *errInvalidHostname) Error() string {
	return fmt.Sprintf("invalid hostname %s, must be one of %s, %s, or %s", e.Hostname, e.MilApp, e.OfficeApp, e.TspApp)
}

// UserSessionCookieName is the key suffix at which we're storing our token cookie
const UserSessionCookieName = "session_token"

// MaskedGorillaCSRFToken is the masked CSRF token used to send back in the 'X-CSRF-Token' request header
const MaskedGorillaCSRFToken = "masked_gorilla_csrf"

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

// GetCookie returns a cookie from a request
func GetCookie(name string, r *http.Request) (*http.Cookie, error) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return cookie, nil
		}
	}
	return nil, errors.Errorf("Unable to find cookie: %s", name)
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

func sessionClaimsFromRequest(logger Logger, secret string, appName Application, r *http.Request) (claims *SessionClaims, ok bool) {
	// Name the cookie with the app name
	cookieName := fmt.Sprintf("%s_%s", string(appName), UserSessionCookieName)
	cookie, err := r.Cookie(cookieName)
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

// WriteMaskedCSRFCookie update the masked_gorilla_csrf cookie value
func WriteMaskedCSRFCookie(w http.ResponseWriter, csrfToken string, noSessionTimeout bool, logger Logger, useSecureCookie bool) {

	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
	maxAge := sessionExpiryInSeconds
	// Never expire token if in development
	if noSessionTimeout {
		expiry = likeForever
		maxAge = likeForeverInSeconds
	}

	// New cookie
	cookie := http.Cookie{
		Name:     MaskedGorillaCSRFToken,
		Value:    csrfToken,
		Path:     "/",
		Expires:  expiry,
		MaxAge:   maxAge,
		HttpOnly: false,                // must be false to be read by client for use in POST/PUT/PATCH/DELETE requests
		SameSite: http.SameSiteLaxMode, // Using lax mode for now since strict is causing issues with Firefox/Safari
		Secure:   useSecureCookie,
	}

	http.SetCookie(w, &cookie)
}

// MaskedCSRFMiddleware handles setting the CSRF Token cookie
func MaskedCSRFMiddleware(logger Logger, noSessionTimeout bool, useSecureCookie bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			// Write a CSRF cookie if none exists
			if _, err := GetCookie(MaskedGorillaCSRFToken, r); err != nil {
				WriteMaskedCSRFCookie(w, csrf.Token(r), noSessionTimeout, logger, useSecureCookie)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

// WriteSessionCookie update the cookie for the session
func WriteSessionCookie(w http.ResponseWriter, session *Session, secret string, noSessionTimeout bool, logger Logger, useSecureCookie bool) {

	// Delete the cookie
	cookieName := fmt.Sprintf("%s_%s", string(session.ApplicationName), UserSessionCookieName)
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    "blank",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode, // Using 'strict' breaks the use of the login.gov redirect
		Secure:   useSecureCookie,
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
			logger.Info("Cookie", zap.Int("Size", len(ss)))
			cookie.Value = ss
			cookie.Expires = expiry
			cookie.MaxAge = maxAge
		}
	}
	// http.SetCookie calls Header().Add() instead of .Set(), which can result in duplicate cookies
	// It's ok to use this here because we want to delete and rewrite `Set-Cookie` on login or if the
	// session token is lost.  However, we would normally use http.SetCookie for any other cookie operations
	// so as not to delete the session token.
	w.Header().Set("Set-Cookie", cookie.String())
}

// ApplicationName returns the application name given the hostname
func ApplicationName(hostname, milHostname, officeHostname, tspHostname string) (Application, error) {
	var appName Application
	if strings.EqualFold(hostname, milHostname) {
		return MilApp, nil
	} else if strings.EqualFold(hostname, officeHostname) {
		return OfficeApp, nil
	} else if strings.EqualFold(hostname, tspHostname) {
		return TspApp, nil
	}
	return appName, errors.Wrap(
		&errInvalidHostname{
			Hostname:  hostname,
			MilApp:    milHostname,
			OfficeApp: officeHostname,
			TspApp:    tspHostname}, fmt.Sprintf("%s is invalid", hostname))
}

// SessionCookieMiddleware handle serializing and de-serializing the session between the user_session cookie and the request context
func SessionCookieMiddleware(logger Logger, secret string, noSessionTimeout bool, milHostname, officeHostname, tspHostname string, useSecureCookie bool) func(next http.Handler) http.Handler {
	logger.Info("Creating session", zap.String("milHost", milHostname), zap.String("officeHost", officeHostname), zap.String("tspHost", tspHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := beeline.StartSpan(r.Context(), "SessionCookieMiddleware")
			defer span.Send()

			// Set up the new session object
			session := Session{}

			// Split the hostname from the port
			hostname := strings.Split(r.Host, ":")[0]
			appName, err := ApplicationName(hostname, milHostname, officeHostname, tspHostname)
			if err != nil {
				logger.Error("Bad Hostname", zap.Error(err))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			claims, ok := sessionClaimsFromRequest(logger, secret, appName, r)
			if ok {
				session = claims.SessionValue
			}

			// Set more information on the session
			session.ApplicationName = appName
			session.Hostname = strings.ToLower(hostname)

			// And put the session info into the request context
			ctx = SetSessionInRequestContext(r.WithContext(ctx), &session)

			// And update the cookie. May get over-ridden later
			WriteSessionCookie(w, &session, secret, noSessionTimeout, logger, useSecureCookie)

			span.AddTraceField("auth.application_name", session.ApplicationName)
			span.AddTraceField("auth.hostname", session.Hostname)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(mw)
	}
}
