package auth

import (
	"net/http"
	"strings"

	"github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Application describes the application name
type Application string

const (
	// TspApp indicates tsp.move.mil
	TspApp Application = "TSP"
	// OfficeApp indicates office.move.mil
	OfficeApp Application = "OFFICE"
	// MyApp indicates my.move.mil
	MyApp Application = "MY"
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

// ApplicationName returns the application name given the hostname
func ApplicationName(hostname, myHostname, officeHostname, tspHostname string) (Application, error) {
	var appName Application
	if strings.EqualFold(hostname, myHostname) {
		return MyApp, nil
	} else if strings.EqualFold(hostname, officeHostname) {
		return OfficeApp, nil
	} else if strings.EqualFold(hostname, tspHostname) {
		return TspApp, nil
	}
	return appName, errors.New("Bad hostname")
}

// DetectorMiddleware detects which application we are serving based on the hostname
func DetectorMiddleware(logger *zap.Logger, myHostname string, officeHostname string, tspHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating host detector", zap.String("myHost", myHostname), zap.String("officeHost", officeHostname), zap.String("tspHost", tspHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := beeline.StartSpan(r.Context(), "DetectorMiddleware")
			defer span.Send()

			session := SessionFromRequestContext(r)

			// Split the hostname from the port
			hostname := strings.Split(r.Host, ":")[0]
			appName, err := ApplicationName(hostname, myHostname, officeHostname, tspHostname)
			if err != nil {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
			}
			session.ApplicationName = appName
			session.Hostname = strings.ToLower(hostname)

			span.AddTraceField("auth.application_name", session.ApplicationName)
			span.AddTraceField("auth.hostname", session.Hostname)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}
