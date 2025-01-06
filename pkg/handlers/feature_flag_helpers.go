package handlers

import (
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
)

func GetFeatureFlagValue(appCtx appcontext.AppContext, featureFlagFetcher services.FeatureFlagFetcher, featureFlagName string) (bool, error) {
	flagValue := false
	flag, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureFlagName, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
		return flagValue, err
	} else {
		flagValue = flag.Match
	}

	return flagValue, nil
}

func GetAllDomesticMHFlags(appCtx appcontext.AppContext, featureFlagFetcher services.FeatureFlagFetcher) (map[string]bool, error) {
	featureFlagValues := make(map[string]bool)

	// Flags to determine if these service item types should be included at all with mobile home shipments
	flagMH, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHome, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHome), zap.Error(err))
	}

	flagDOPEnabled, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDOPEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDOPEnabled), zap.Error(err))
	}

	flagDDPEnabled, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDDPEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDDPEnabled), zap.Error(err))
	}

	flagDPKEnabled, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomePackingEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomePackingEnabled), zap.Error(err))
	}

	flagDUPKEnabled, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeUnpackingEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeUnpackingEnabled), zap.Error(err))
	}

	// Flags for whether or not the item type is affected by the mobile home factor
	flagDOPFactor, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDOPFactor, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDOPEnabled), zap.Error(err))
	}

	flagDDPFactor, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDDPFactor, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDDPEnabled), zap.Error(err))
	}

	flagDPKFactor, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomePackingFactor, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomePackingEnabled), zap.Error(err))
	}

	flagDUPKFactor, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeUnpackingFactor, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeUnpackingEnabled), zap.Error(err))
	}

	featureFlagValues[featureflag.DomesticMobileHome] = flagMH.Match
	featureFlagValues[featureflag.DomesticMobileHomeDOPEnabled] = flagDOPEnabled.Match
	featureFlagValues[featureflag.DomesticMobileHomeDDPEnabled] = flagDDPEnabled.Match
	featureFlagValues[featureflag.DomesticMobileHomePackingEnabled] = flagDPKEnabled.Match
	featureFlagValues[featureflag.DomesticMobileHomeUnpackingEnabled] = flagDUPKEnabled.Match

	featureFlagValues[featureflag.DomesticMobileHomeDOPFactor] = flagDOPFactor.Match
	featureFlagValues[featureflag.DomesticMobileHomeDDPFactor] = flagDDPFactor.Match
	featureFlagValues[featureflag.DomesticMobileHomePackingFactor] = flagDPKFactor.Match
	featureFlagValues[featureflag.DomesticMobileHomeUnpackingFactor] = flagDUPKFactor.Match

	return featureFlagValues, nil
}
