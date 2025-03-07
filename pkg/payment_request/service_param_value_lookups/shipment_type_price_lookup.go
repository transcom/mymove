package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type ShipmentTypePriceLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r ShipmentTypePriceLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	// TODO: Get rid of the OCONUS hard code
	factor, err := models.FetchMarketFactor(appCtx, keyData.ContractID, r.ServiceItem.ReServiceID, models.MarketOconus.String())
	if err != nil {
		return "", fmt.Errorf("ShipmentTypePrice error when fetching market factor for ReServiceItem %s: err: %w", r.ServiceItem.ReService.Code, err)
	}

	return strconv.FormatFloat(factor, 'f', -1, 64), nil
}
