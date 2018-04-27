package app

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type appCtxKey string

const appKey = appCtxKey("app")
const hostnameKey = appCtxKey("hostname")
const officeApp = "OFFICE"
const myApp = "MY"

func fromContext(r *http.Request, key appCtxKey) string {
	if app, ok := r.Context().Value(key).(string); ok {
		return app
	}
	return ""
}

func appFromContext(r *http.Request) string {
	return fromContext(r, appKey)
}

// IsOfficeApp returns true iff the request is for the office.move.mil host
func IsOfficeApp(r *http.Request) bool {
	return appFromContext(r) == officeApp
}

// IsMyApp returns true iff the request is for the my.move.mil host
func IsMyApp(r *http.Request) bool {
	return appFromContext(r) == myApp
}

// GetHostname returns the hostname used to hit this server
func GetHostname(r *http.Request) string {
	return fromContext(r, hostnameKey)
}

// DetectorMiddleware detects which application we are serving based on the hostname
func DetectorMiddleware(logger *zap.Logger, myHostname string, officeHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating host detector", zap.String("myHost", myHostname), zap.String("officeHost", officeHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Host, ":")
			var app string
			if strings.EqualFold(parts[0], myHostname) {
				app = myApp
			} else if strings.EqualFold(parts[0], officeHostname) {
				app = officeApp
			} else {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), appKey, app)
			ctx = context.WithValue(ctx, hostnameKey, strings.ToLower(parts[0]))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}
