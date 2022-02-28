package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ActualPickupDateLookup does lookup on actual pickup date
type ActualPickupDateLookup struct {
	MTOShipment models.MTOShipment
}

func (r ActualPickupDateLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	actualPickupDate := r.MTOShipment.ActualPickupDate
	if actualPickupDate == nil {
		return "", nil
	}

	return actualPickupDate.Format(ghcrateengine.DateParamFormat), nil
}
