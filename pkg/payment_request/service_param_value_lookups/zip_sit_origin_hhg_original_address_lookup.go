package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ZipSITOriginHHGOriginalAddressLookup does lookup on the postal code HHG shipment's original pickup address
type ZipSITOriginHHGOriginalAddressLookup struct {
	ServiceItem models.MTOServiceItem
}

func (z ZipSITOriginHHGOriginalAddressLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {

	// load updated origin SIT addresses from service item
	if z.ServiceItem.SITOriginHHGOriginalAddressID != nil && *z.ServiceItem.SITOriginHHGOriginalAddressID != uuid.Nil {
		err := appCtx.DB().Load(&z.ServiceItem, "SITOriginHHGOriginalAddress")
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("nil SITOriginHHGOriginalAddressID for service item ID %s", z.ServiceItem.ID.String())
	}

	if z.ServiceItem.SITOriginHHGOriginalAddress == nil {
		return "", fmt.Errorf("db load for SITOriginHHGOriginalAddress failed service item ID %s", z.ServiceItem.ID.String())
	}

	value := z.ServiceItem.SITOriginHHGOriginalAddress.PostalCode
	return value, nil
}
