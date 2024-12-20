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

	flagMH, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHome, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHome), zap.Error(err))
	}

	flagDOP, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDOPEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDOPEnabled), zap.Error(err))
	}

	flagDDP, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeDDPEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeDDPEnabled), zap.Error(err))
	}

	flagDPK, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomePackingEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomePackingEnabled), zap.Error(err))
	}

	flagDUPK, err := featureFlagFetcher.GetBooleanFlag(appCtx.DB().Context(), appCtx.Logger(), "", featureflag.DomesticMobileHomeUnpackingEnabled, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureflag.DomesticMobileHomeUnpackingEnabled), zap.Error(err))
	}

	featureFlagValues[featureflag.DomesticMobileHome] = flagMH.Match
	featureFlagValues[featureflag.DomesticMobileHomeDOPEnabled] = flagDOP.Match
	featureFlagValues[featureflag.DomesticMobileHomeDDPEnabled] = flagDDP.Match
	featureFlagValues[featureflag.DomesticMobileHomePackingEnabled] = flagDPK.Match
	featureFlagValues[featureflag.DomesticMobileHomeUnpackingEnabled] = flagDUPK.Match
	return featureFlagValues, nil
}
