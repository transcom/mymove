package handlers

import (
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
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
