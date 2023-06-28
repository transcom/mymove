package services

import (
	"context"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

const (
	// use a convention where the name of the variant is enabled or
	// disabled for booleans
	FeatureFlagEnabledVariant  = "enabled"
	FeatureFlagDisabledVariant = "disabled"
)

// Simplifed struct based on
// https://pkg.go.dev/go.flipt.io/flipt/rpc/flipt#EvaluationResponse
type FeatureFlag struct {
	Entity    string
	Key       string
	Match     bool
	Value     string
	Namespace string
}

func (ff FeatureFlag) IsEnabledVariant() bool {
	return ff.Match && ff.Value == FeatureFlagEnabledVariant
}

func (ff FeatureFlag) IsVariant(variant string) bool {
	return ff.Match && ff.Value == variant
}

// FeatureFlagFetcher is the exported interface for feature flags
//
//go:generate mockery --name FeatureFlagFetcher
type FeatureFlagFetcher interface {
	GetFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (FeatureFlag, error)
	GetFlag(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (FeatureFlag, error)
}
