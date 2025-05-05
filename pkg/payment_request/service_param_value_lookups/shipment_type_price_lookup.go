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

var optionalLookupsReServiceCodes = []models.ReServiceCode{
	models.ReServiceCodeIHPK,
}

func (r ShipmentTypePriceLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	// Skip if optional
	if optional := r.isShipmentTypePriceLookupOptional(); optional {
		return "", nil
	}
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
	// If you are debugging this and you wonder why you might not see your ReServiceCode in the supported ShipmentTypePrices,
	// double check that you aren't encountering the scenario of ShipmentTypePrice being looked up when it shouldn't.
	// You should only be looking this up if you have special market factors for your ReService, like INPK, DNPK, IBTF, etc.
	// This is not just a normal CONUS/OCONUS checker. These items have special factors depending on conus or oconus in the same db
	// row, we have no method of currently dynamically fetching that so we have to hard code it until the enhancement is introduced
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

// Returns whether or not ShipmentTypePrice is an optional lookup for the provided
// ReServiceCode.
// This is used because IHPK is used to price INPK, meaning ShipmentTypePrice is only
// present on IHPK when pricing INPK. When not INPK, ShipmentTypePrice can be ignored on IHPK
func (r ShipmentTypePriceLookup) isShipmentTypePriceLookupOptional() bool {
	return containsReServiceCode(optionalLookupsReServiceCodes, r.ServiceItem.ReService.Code)
}
