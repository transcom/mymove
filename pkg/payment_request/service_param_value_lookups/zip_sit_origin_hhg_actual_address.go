package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// ZipSITOriginHHGActualAddressLookup does lookup on the postal code HHG shipment's actual (new) pickup address
type ZipSITOriginHHGActualAddressLookup struct {
	ServiceItem models.MTOServiceItem
}

func (z ZipSITOriginHHGActualAddressLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	if z.ServiceItem.SITOriginHHGActualAddressID != nil && *z.ServiceItem.SITOriginHHGActualAddressID != uuid.Nil {
		err := keyData.db.Load(&z.ServiceItem, "SITOriginHHGActualAddress")
		if err != nil {
			return "", err
		}
	}

	if z.ServiceItem.SITOriginHHGActualAddress != nil {
		return "", fmt.Errorf("db load for SITOriginHHGActualAddress failed service item ID %s", z.ServiceItem.ID.String())
	}

	value := fmt.Sprintf("%s", z.ServiceItem.SITOriginHHGActualAddress.PostalCode)
	return value, nil
}
