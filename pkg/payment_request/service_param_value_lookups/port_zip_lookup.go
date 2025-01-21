package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PortZipLookup does lookup on the shipment and finds the port zip
// The mileage calculated is from port <-> pickup/destination so this value is important
type PortZipLookup struct {
	ServiceItem models.MTOServiceItem
}

func (p PortZipLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var portLocationID *uuid.UUID
	if p.ServiceItem.PODLocationID != nil {
		portLocationID = p.ServiceItem.PODLocationID
	} else if p.ServiceItem.POELocationID != nil {
		portLocationID = p.ServiceItem.POELocationID
	} else {
		// for PPMs we need to send back the ZIP for the Tacoma Port, they are reimbursed for their CONUS <-> Port travel
		shipment, err := models.FetchShipmentByID(appCtx.DB(), *keyData.mtoShipmentID)
		if err != nil {
			return "", fmt.Errorf("unable to find shipment with id %s", keyData.mtoShipmentID)
		}
		if shipment.ShipmentType == models.MTOShipmentTypePPM && shipment.MarketCode == models.MarketCodeInternational {
			portLocation, err := models.FetchPortLocationByCode(appCtx.DB(), "3002")
			if err != nil {
				return "", fmt.Errorf("unable to find port zip with code %s", "3002")
			}
			return portLocation.UsPostRegionCity.UsprZipID, nil
		} else {
			return "", nil
		}
	}
	var portLocation models.PortLocation
	err := appCtx.DB().Q().
		EagerPreload("UsPostRegionCity").
		Where("id = $1", portLocationID).First(&portLocation)
	if err != nil {
		return "", fmt.Errorf("unable to find port zip with id %s", portLocationID)
	}
	return portLocation.UsPostRegionCity.UsprZipID, nil
}
