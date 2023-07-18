package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// inspired by
// https://github.com/go-chi/chi/issues/695#issuecomment-1587508983

func NewOtelHTTPMiddleware(config *Config, name string, globalLogger *zap.Logger) func(http.Handler) http.Handler {
	// support noop middleware if telemetry is not enabled
	if !config.Enabled {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	requestTelemetry := NewRequestTelemetry(globalLogger)

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
				// For gathering metrics, we want to know the route
				// pattern (e.g. /ghc/v1/move/{locator}) and not the
				// explicit http url (e.g. /ghc/v1/move/ABCDEF)
				// because the latter won't help us find patterns in
				// the stats for API endpoints.
				//
				// With go-chi, we can only get the route after the
				// request is served (e.g. after ServeHTTP is called).
				//
				// However, go-chi isn't the only place that knows
				// about the routing, because we also use openapi (aka
				// swagger) routes. We have middleware that can
				// extract the route (see middleware.OpenAPITracing)
				// but that middleware has to be set up and run after
				// this middleware because this is set up once for all
				// routes and the openapi middleware is set up once
				// per swagger config (internal, ghc, prime, etc)
				//
				// Thus, we need a way for middleware that runs
				// *after* this one, to pass information back to this
				// one.
				//
				// We do that by putting a string pointer into the
				// context via ContextWithRoutePattern. The
				// OpenAPITracing middleware can get the string
				// pointer out of the context and update it and then
				// this middleware can look at the contents of the
				// string pointer to see if it has been changed.
				var routePattern string
				r = r.WithContext(ContextWithRoutePattern(r.Context(), &routePattern))

				var statusCode int
				r = r.WithContext(ContextWithStatusCode(r.Context(), &statusCode))
				h.ServeHTTP(w, r)

				// if nothing set up the route pattern already (e.g.
				// the openapi middleware), get the route pattern from chi
				if routePattern == "" {
					routeCtx := chi.RouteContext(r.Context())
					if routeCtx != nil {
						routePattern = routeCtx.RoutePattern()
					}
				}

				span := trace.SpanFromContext(r.Context())
				if span != nil {
					span.SetAttributes(semconv.HTTPTarget(r.URL.String()), semconv.HTTPRoute(routePattern))
				}

				labeler, ok := otelhttp.LabelerFromContext(r.Context())
				if ok {
					labeler.Add(semconv.HTTPRoute(routePattern))
				}

				requestTelemetry.IncrementRequestCount(r, routePattern, statusCode)
			}),
			name,
			otelHTTPOptions...,
		)
	}
}
