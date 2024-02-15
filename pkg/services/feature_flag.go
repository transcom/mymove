package services

import (
	"context"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// Simplifed struct based on
// https://pkg.go.dev/go.flipt.io/flipt/rpc/flipt@v1.25.0/evaluation#VariantEvaluationResponse
// For boolean responses, Variant is the empty string and Match is the
// flag value
type FeatureFlag struct {
	Entity    string
	Key       string
	Match     bool
	Variant   string
	Namespace string
}

func (ff FeatureFlag) IsVariant(variant string) bool {
	return ff.Match && ff.Variant == variant
}

// FeatureFlagFetcher is the exported interface for feature flags
//
// This service is a thin wrapper around flipt, so it doesn't expose
// all of flipt's API. We can change/expand the API as we get
// experience.
//
//go:generate mockery --name FeatureFlagFetcher
type FeatureFlagFetcher interface {
	GetBooleanFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (FeatureFlag, error)
	GetBooleanFlag(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (FeatureFlag, error)
	GetVariantFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (FeatureFlag, error)
	GetVariantFlag(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (FeatureFlag, error)
}
