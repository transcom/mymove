package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ZipSITOriginHHGActualAddressLookup does lookup on the postal code HHG shipment's actual (new) pickup address
type ZipSITOriginHHGActualAddressLookup struct {
	ServiceItem models.MTOServiceItem
}

func (z ZipSITOriginHHGActualAddressLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	if z.ServiceItem.SITOriginHHGActualAddressID != nil && *z.ServiceItem.SITOriginHHGActualAddressID != uuid.Nil {
		err := appCtx.DB().Load(&z.ServiceItem, "SITOriginHHGActualAddress")
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("nil SITOriginHHGActualAddressID for service item ID %s", z.ServiceItem.ID.String())
	}

	if z.ServiceItem.SITOriginHHGActualAddress == nil {
		return "", fmt.Errorf("db load for SITOriginHHGActualAddress failed service item ID %s", z.ServiceItem.ID.String())
	}

	value := z.ServiceItem.SITOriginHHGActualAddress.PostalCode
	return value, nil
}
