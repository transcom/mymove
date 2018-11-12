package server

import (
	"net/http"
	"strings"

	"github.com/honeycombio/beeline-go"
	"go.uber.org/zap"
)

type application string

const (
	// TspApp indicates tsp.move.mil
	TspApp application = "TSP"
	// OfficeApp indicates office.move.mil
	OfficeApp application = "OFFICE"
	// MyApp indicates my.move.mil
	MyApp application = "MY"
)

// IsTspApp returns true iff the request is for the office.move.mil host
func (s *Session) IsTspApp() bool {
	return s.ApplicationName == TspApp
}

// IsOfficeApp returns true iff the request is for the office.move.mil host
func (s *Session) IsOfficeApp() bool {
	return s.ApplicationName == OfficeApp
}

// IsMyApp returns true iff the request is for the my.move.mil host
func (s *Session) IsMyApp() bool {
	return s.ApplicationName == MyApp
}

// AppDetectorMiddleware is a unique type definition for this middleware, so DI can identify the right object
type AppDetectorMiddleware func(next http.Handler) http.Handler

// NewAppDetectorMiddleware detects which application we are serving based on the hostname
func NewAppDetectorMiddleware(cfg *HostsConfig, l *zap.Logger) AppDetectorMiddleware {
	l.Info("Creating host detector", zap.String("myHost", cfg.MyName), zap.String("officeHost", cfg.OfficeName), zap.String("tspHost", cfg.TspName))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			_, span := beeline.StartSpan(r.Context(), "DetectorMiddleware")
			session := SessionFromRequestContext(r)
			parts := strings.Split(r.Host, ":")
			var appName application
			if strings.EqualFold(parts[0], cfg.MyName) {
				appName = MyApp
			} else if strings.EqualFold(parts[0], cfg.OfficeName) {
				appName = OfficeApp
			} else if strings.EqualFold(parts[0], cfg.TspName) {
				appName = TspApp
			} else {
				l.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			session.ApplicationName = appName
			session.Hostname = strings.ToLower(parts[0])
			span.AddField("auth.application_name", session.ApplicationName)
			span.AddField("auth.hostname", session.Hostname)
			span.Send()
			next.ServeHTTP(w, r)
			return
		}
		return http.HandlerFunc(mw)
	}
}
