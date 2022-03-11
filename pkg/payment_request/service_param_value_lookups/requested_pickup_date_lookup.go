package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// RequestedPickupDateLookup does lookup on requested pickup date
type RequestedPickupDateLookup struct {
	MTOShipment models.MTOShipment
}

func (r RequestedPickupDateLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	requestedPickupDate := r.MTOShipment.RequestedPickupDate
	if requestedPickupDate == nil {
		return "", nil
	}

	return requestedPickupDate.Format(ghcrateengine.DateParamFormat), nil
}
