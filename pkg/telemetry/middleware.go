package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// inspired by
// https://github.com/go-chi/chi/issues/695#issuecomment-1587508983

func NewOtelHTTPMiddleware(config *Config, name string) func(http.Handler) http.Handler {
	meterProvider := otel.GetMeterProvider()
	// set up telemetry options for the server
	otelHTTPOptions := []otelhttp.Option{
		otelhttp.WithMeterProvider(meterProvider),
	}
	if config.ReadEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.ReadEvents))
	}
	if config.WriteEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.WriteEvents))
	}
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)

				routePattern := chi.RouteContext(r.Context()).RoutePattern()

				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(semconv.HTTPTarget(r.URL.String()), semconv.HTTPRoute(routePattern))

				labeler, ok := otelhttp.LabelerFromContext(r.Context())
				if ok {
					labeler.Add(semconv.HTTPRoute(routePattern))
				}
			}),
			name,
			otelHTTPOptions...,
		)
	}
}
