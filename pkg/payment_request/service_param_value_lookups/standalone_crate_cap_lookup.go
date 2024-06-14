package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// StandaloneCrateCapLookup does lookup on application parameters
type StandaloneCrateCapLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r StandaloneCrateCapLookup) lookup(appCtx appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	applicationParam, _ := models.FetchParameterValueByName(appCtx.DB(), "standaloneCrateCap")

	return *applicationParam.ParameterValue, nil
}
