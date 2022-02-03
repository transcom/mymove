package serviceparamvaluelookups

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// DateBilledLookup determines the date to use for billing (e.g., determining peak vs. non peak,
// escalations, etc.)
type DateBilledLookup struct {
	MTOShipment models.MTOShipment
}

func (r DateBilledLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var dateBilled *time.Time

	// Most shipment types should use RequestedPickupDate, but there are exceptions.
	switch r.MTOShipment.ShipmentType {
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		actualPickupDate := r.MTOShipment.ActualPickupDate
		if actualPickupDate == nil || actualPickupDate.IsZero() {
			return "", fmt.Errorf("could not find a valid actual pickup date for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
		dateBilled = actualPickupDate
	default:
		requestedPickupDate := r.MTOShipment.RequestedPickupDate
		if requestedPickupDate == nil || requestedPickupDate.IsZero() {
			return "", fmt.Errorf("could not find a valid requested pickup date for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
		dateBilled = r.MTOShipment.RequestedPickupDate
	}

	return dateBilled.Format(ghcrateengine.DateParamFormat), nil
}
