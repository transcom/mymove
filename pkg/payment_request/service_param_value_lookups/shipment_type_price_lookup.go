package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type ShipmentTypePriceLookup struct {
	ServiceItem models.MTOServiceItem
}

func containsReServiceCode(validCodes []models.ReServiceCode, code models.ReServiceCode) bool {
	for _, validCode := range validCodes {
		if validCode == code {
			return true
		}
	}
	return false
}

var internationalSupportedShipmentTypePriceServices = []models.ReServiceCode{
	// IMHF does not exist yet
	models.ReServiceCodeINPK,
	models.ReServiceCodeIBTF,
	models.ReServiceCodeIBHF,
}

var domesticSupportedShipmentTypePriceServices = []models.ReServiceCode{
	models.ReServiceCodeDNPK,
	models.ReServiceCodeDMHF,
	models.ReServiceCodeDBTF,
	models.ReServiceCodeDBHF,
}

func (r ShipmentTypePriceLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	// Until the new `service_location` table is fully in place and not an optional field, we do not have a 100%
	// certain method of determining the market code dynamically.
	// Until then, we will just lock down the shipment type price lookup to known items
	// as well as determine "O" vs "D" based on the incoming code
	var market models.Market
	if containsReServiceCode(internationalSupportedShipmentTypePriceServices, r.ServiceItem.ReService.Code) {
		market = models.MarketOconus
	}
	if containsReServiceCode(domesticSupportedShipmentTypePriceServices, r.ServiceItem.ReService.Code) {
		market = models.MarketConus
	}
	if market.String() == "" {
		return "", errors.New(`service param value lookup package failed on
		ShipmentTypePriceLookup due to lookup initialized service item not having its
		ReServiceCode joined properly or it is currently not supported by the package`)
	}

	factor, err := models.FetchMarketFactor(appCtx, keyData.ContractID, r.ServiceItem.ReServiceID, market.String())
	if err != nil {
		return "", fmt.Errorf("ShipmentTypePrice error when fetching market factor for ReServiceItem %s: err: %w", r.ServiceItem.ReService.Code, err)
	}

	return strconv.FormatFloat(factor, 'f', -1, 64), nil
}
