package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// RequestedPickupDateLookup does lookup on requested pickup date
type RequestedPickupDateLookup struct {
	MTOShipment models.MTOShipment
}

func (r RequestedPickupDateLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Make sure there's a requested pickup date since that's nullable
	requestedPickupDate := r.MTOShipment.RequestedPickupDate
	if requestedPickupDate == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find a requested pickup date for MTOShipmentID [%s]", r.MTOShipment.ID)
	}

	return requestedPickupDate.Format(ghcrateengine.DateParamFormat), nil
}
