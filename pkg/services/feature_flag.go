package services

import (
	"context"

	"github.com/transcom/mymove/pkg/appcontext"
)

// Simplifed struct based on
// https://pkg.go.dev/go.flipt.io/flipt/rpc/flipt#EvaluationResponse
type FeatureFlag struct {
	Entity    string
	Key       string
	Enabled   bool
	Value     string
	Namespace string
}

// FeatureFlagFetcher is the exported interface for feature flags
//
//go:generate mockery --name FeatureFlagFetcher
type FeatureFlagFetcher interface {
	GetFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (FeatureFlag, error)
	IsEnabledForUser(ctx context.Context, appCtx appcontext.AppContext, key string) (bool, error)
	GetFlag(ctx context.Context, entityID string, key string, flagContext map[string]string) (FeatureFlag, error)
}
