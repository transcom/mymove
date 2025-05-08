package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// RateAreaLookup does lookup of rate area based on address postal code
type RateAreaLookup struct {
	Address models.Address
}

func (r RateAreaLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	rateArea, err := models.FetchRateArea(appCtx.DB(), r.Address.ID, keyData.MTOServiceItemID, keyData.ContractID)
	if err != nil {
		return "", err
	}
	return rateArea.Code, nil
}
