package apiversion

import (
	"context"
)

type contextKey string

var apiVersionContextKey = contextKey("apiVersion")

type Flag byte

const (
	NoneSpecified Flag = iota
	PrimeVersion1
	PrimeVersion2
)

// WithAPIVersion returns a new context with the API Version.
func WithAPIVersion(ctx context.Context, apiVersion Flag) context.Context {
	return context.WithValue(ctx, apiVersionContextKey, apiVersion)
}

// RetrieveAPIVersionFromContext returns the Flag associated with a context. If the
// context has no Flag, the noneSpecified flag is returned
func RetrieveAPIVersionFromContext(ctx context.Context) Flag {
	apiVersion, ok := ctx.Value(apiVersionContextKey).(Flag)
	if !ok {
		return NoneSpecified
	}
	return apiVersion
}
