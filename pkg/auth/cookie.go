package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
)

// ApplicationServername is a collection of all the servernames for the application
type ApplicationServername struct {
	MilServername    string
	OfficeServername string
	AdminServername  string
	OrdersServername string
	PrimeServername  string
}

type errInvalidHostname struct {
	Hostname  string
	MilApp    string
	OfficeApp string
	AdminApp  string
}

func (e *errInvalidHostname) Error() string {
	return fmt.Sprintf("invalid hostname %s, must be one of %s, %s, or %s", e.Hostname, e.MilApp, e.OfficeApp, e.AdminApp)
}

// GorillaCSRFToken is the name of the base CSRF token
// RA Summary: gosec - G101 - Password Management: Hardcoded Password
// RA: This line was flagged because it detected use of the word "token"
// RA: This line is used to identify the name of the token. GorillaCSRFToken is the name of the base CSRF token.
// RA: This variable does not store an application token.
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Validator: jneuner@mitre.org
// RA Modified Severity: CAT III
// #nosec G101
const GorillaCSRFToken = "_gorilla_csrf"

// MaskedGorillaCSRFToken is the masked CSRF token used to send back in the 'X-CSRF-Token' request header
const MaskedGorillaCSRFToken = "masked_gorilla_csrf"

// SessionExpiryInMinutes is the number of minutes before a fallow session is harvested
const SessionExpiryInMinutes = 15

// GetExpiryTimeFromMinutes returns 'min' minutes from now
func GetExpiryTimeFromMinutes(min int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(min))
}

// DeleteCookie sends a delete request for the named cookie
func DeleteCookie(w http.ResponseWriter, name string) {
	c := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, &c)
}

// WriteMaskedCSRFCookie update the masked_gorilla_csrf cookie value
func WriteMaskedCSRFCookie(w http.ResponseWriter, csrfToken string, useSecureCookie bool) {
	// Match expiration settings of the _gorilla_csrf cookie (a session cookie); don't set Expires or MaxAge.
	cookie := http.Cookie{
		Name:     MaskedGorillaCSRFToken,
		Value:    csrfToken,
		Path:     "/",
		HttpOnly: false,                // must be false to be read by client for use in POST/PUT/PATCH/DELETE requests
		SameSite: http.SameSiteLaxMode, // Using 'lax' mode for now since 'strict' is causing issues with Firefox/Safari
		Secure:   useSecureCookie,
	}

	http.SetCookie(w, &cookie)
}

// DeleteCSRFCookies deletes the base and masked CSRF cookies
func DeleteCSRFCookies(w http.ResponseWriter) {
	DeleteCookie(w, MaskedGorillaCSRFToken)
	DeleteCookie(w, GorillaCSRFToken)
}

// MaskedCSRFMiddleware handles setting the CSRF Token cookie
func MaskedCSRFMiddleware(globalLogger *zap.Logger, useSecureCookie bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			// Write a masked CSRF cookie (creates a new one with each request).  Per the gorilla/csrf docs:
			// "This library generates unique-per-request (masked) tokens as a mitigation against the BREACH attack."
			// https://github.com/gorilla/csrf#design-notes
			WriteMaskedCSRFCookie(w, csrf.Token(r), useSecureCookie)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

// ApplicationName returns the application name given the hostname
func ApplicationName(hostname string, appnames ApplicationServername) (Application, error) {
	var appName Application
	if strings.EqualFold(hostname, appnames.MilServername) {
		return MilApp, nil
	} else if strings.EqualFold(hostname, appnames.OfficeServername) {
		return OfficeApp, nil
	} else if strings.EqualFold(hostname, appnames.AdminServername) {
		return AdminApp, nil
	}
	return appName, errors.Wrap(
		&errInvalidHostname{
			Hostname:  hostname,
			MilApp:    appnames.MilServername,
			OfficeApp: appnames.OfficeServername,
			AdminApp:  appnames.AdminServername,
		}, fmt.Sprintf("%s is invalid", hostname))
}

// SessionCookieMiddleware handle serializing and de-serializing the session between the user_session cookie and the request context
func SessionCookieMiddleware(globalLogger *zap.Logger, appnames ApplicationServername, sessionManagers AppSessionManagers) func(next http.Handler) http.Handler {
	globalLogger.Info("Creating session middleware",
		zap.String("milServername", appnames.MilServername),
		zap.String("officeServername", appnames.OfficeServername),
		zap.String("adminServername", appnames.AdminServername))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			logger := logging.FromContext(ctx)

			// Split the hostname from the port
			hostname := strings.Split(r.Host, ":")[0]
			app, err := ApplicationName(hostname, appnames)
			if err != nil {
				logger.Error("Bad Hostname", zap.Error(err))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}

			sessionManager := sessionManagers.SessionManagerForApplication(app)

			// The scs session manager Get call will return an empty
			// Session if an existing one is not found in the store
			obj := sessionManager.Get(r.Context(), "session")
			session, ok := obj.(Session)
			if ok {
				logger.Info("Existing session", zap.Any("session.user_id", session.UserID),
					zap.Any("session.appname", session.ApplicationName))
			} else {
				session = Session{
					ApplicationName: app,
					Hostname:        strings.ToLower(hostname),
				}
				logger.Info("Creating new session", zap.Any("session.user_id", session.UserID),
					zap.Any("session.appname", session.ApplicationName))
			}

			// And update the cookie. May get over-ridden later
			sessionManager.Put(r.Context(), "session", session)

			// And put the session info into the request context
			next.ServeHTTP(w, r.WithContext(SetSessionInContext(ctx, &session)))

		})
	}
}
