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

const (
	// Checks if the mobile home FF is enabled
	DomesticMobileHome string = "mobile_home"

	// Toggles service items on/off completely for mobile home shipments
	DomesticMobileHomeDOPEnabled       string = "domestic_mobile_home_origin_price_enabled"
	DomesticMobileHomeDDPEnabled       string = "domestic_mobile_home_destination_price_enabled"
	DomesticMobileHomePackingEnabled   string = "domestic_mobile_home_packing_enabled"
	DomesticMobileHomeUnpackingEnabled string = "domestic_mobile_home_unpacking_enabled"

	// Toggles whether or not the DMHF is applied to these service items for Mobile Home shipments (if they are not toggled off by the above flags)
	DomesticMobileHomeDOPFactor       string = "domestic_mobile_home_factor_origin_price"
	DomesticMobileHomeDDPFactor       string = "domestic_mobile_home_factor_destination_price"
	DomesticMobileHomePackingFactor   string = "domestic_mobile_home_factor_packing"
	DomesticMobileHomeUnpackingFactor string = "domestic_mobile_home_factor_unpacking"
)
