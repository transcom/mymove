package telemetry

import "context"

type telemetryKeyType int

const (
	routePatternKey telemetryKeyType = iota
	statusCodeKey
)

// ContextWithRoutePattern adds the route pattern to the context. See
// NewOtelHTTPMiddleware for usage
func ContextWithRoutePattern(ctx context.Context, routePattern *string) context.Context {
	return context.WithValue(ctx, routePatternKey, routePattern)
}

// RoutePatternFromContext returns the route pattern if found, or the
// empty string
func RoutePatternFromContext(ctx context.Context) *string {
	val := ctx.Value(routePatternKey)
	if val == nil {
		return nil
	}
	routePattern, ok := val.(*string)
	if ok {
		return routePattern
	}
	return nil
}

// ContextWithStatusCode adds the status code to the context. See
// NewOtelHTTPMiddleware for usage
func ContextWithStatusCode(ctx context.Context, statusCode *int) context.Context {
	return context.WithValue(ctx, statusCodeKey, statusCode)
}

// StatusCodeFromContext returns the route pattern if found, or the
// empty string
func StatusCodeFromContext(ctx context.Context) *int {
	val := ctx.Value(statusCodeKey)
	if val == nil {
		return nil
	}
	statusCode, ok := val.(*int)
	if ok {
		return statusCode
	}
	return nil
}
