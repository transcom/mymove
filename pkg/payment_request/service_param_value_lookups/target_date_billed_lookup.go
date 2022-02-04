package serviceparamvaluelookups

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// TargetDateBilledLookup determines the target date to use for billing (e.g., determining
// peak vs. non peak, escalations, etc.)
type TargetDateBilledLookup struct {
	MTOShipment models.MTOShipment
}

func (r TargetDateBilledLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var targetDateBilled *time.Time

	// Most shipment types should use RequestedPickupDate, but there are exceptions.
	switch r.MTOShipment.ShipmentType {
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		actualPickupDate := r.MTOShipment.ActualPickupDate
		if actualPickupDate == nil || actualPickupDate.IsZero() {
			return "", fmt.Errorf("could not find a valid actual pickup date for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
		targetDateBilled = actualPickupDate
	default:
		requestedPickupDate := r.MTOShipment.RequestedPickupDate
		if requestedPickupDate == nil || requestedPickupDate.IsZero() {
			return "", fmt.Errorf("could not find a valid requested pickup date for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
		targetDateBilled = r.MTOShipment.RequestedPickupDate
	}

	return targetDateBilled.Format(ghcrateengine.DateParamFormat), nil
}
