package auth

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

// DetectorMiddleware detects which application we are serving based on the hostname
func DetectorMiddleware(logger *zap.Logger, myHostname string, officeHostname string, tspHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating host detector", zap.String("myHost", myHostname), zap.String("officeHost", officeHostname), zap.String("tspHost", tspHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			session := SessionFromRequestContext(r)
			parts := strings.Split(r.Host, ":")
			var appName application
			if strings.EqualFold(parts[0], myHostname) {
				appName = MyApp
			} else if strings.EqualFold(parts[0], officeHostname) {
				appName = OfficeApp
			} else if strings.EqualFold(parts[0], tspHostname) {
				appName = TspApp
			} else {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			session.ApplicationName = appName
			session.Hostname = strings.ToLower(parts[0])
			beeline.AddField(r.Context(), "application_name", appName)
			next.ServeHTTP(w, r)
			return
		}
		return http.HandlerFunc(mw)
	}
}
