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

// OfficeApp describes an office app
const OfficeApp = "OFFICE"

// MyApp describes a my move app
const MyApp = "MY"

func fromContext(r *http.Request, key appCtxKey) string {
	if app, ok := r.Context().Value(key).(string); ok {
		return app
	}
	return ""
}

// GetAppFromContext returns an app string for a request
func GetAppFromContext(r *http.Request) string {
	return fromContext(r, appKey)
}

// IsOfficeApp returns true iff the request is for the office.move.mil host
func IsOfficeApp(r *http.Request) bool {
	return GetAppFromContext(r) == OfficeApp
}

// IsMyApp returns true iff the request is for the my.move.mil host
func IsMyApp(r *http.Request) bool {
	return GetAppFromContext(r) == MyApp
}

// GetHostname returns the hostname used to hit this server
func GetHostname(r *http.Request) string {
	return fromContext(r, hostnameKey)
}

// PopulateAppContext adds the app onto request context
func PopulateAppContext(ctx context.Context, app string) context.Context {
	ctx = context.WithValue(ctx, appKey, app)
	return ctx
}

// PopulateHostnameContext adds the app onto request context
func PopulateHostnameContext(ctx context.Context, hostname string) context.Context {
	ctx = context.WithValue(ctx, hostnameKey, hostname)
	return ctx
}

// DetectorMiddleware detects which application we are serving based on the hostname
func DetectorMiddleware(logger *zap.Logger, myHostname string, officeHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating host detector", zap.String("myHost", myHostname), zap.String("officeHost", officeHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Host, ":")
			var app string
			if strings.EqualFold(parts[0], myHostname) {
				app = MyApp
			} else if strings.EqualFold(parts[0], officeHostname) {
				app = OfficeApp
			} else {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			ctx := PopulateAppContext(r.Context(), app)
			ctx = PopulateHostnameContext(ctx, strings.ToLower(parts[0]))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}
