package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightOriginalLookup does lookup on original weight
type WeightOriginalLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightOriginalLookup) lookup(_ appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var originalWeight *unit.Pound

	switch keyData.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDOSHUT,
		models.ReServiceCodeDDSHUT,
		models.ReServiceCodeIOSHUT,
		models.ReServiceCodeIDSHUT:
		// Attempt to grab the service item's actual weight. If it can't be found, default to shipment
		if keyData.MTOServiceItem.ActualWeight == nil {
			originalWeight = r.MTOShipment.PrimeActualWeight
			if originalWeight == nil {
				return "", fmt.Errorf("could not find actual weight for MTOServiceItemID [%s] or for MTOShipmentID [%s]", keyData.MTOServiceItem.ID, r.MTOShipment.ID)
			}
		} else {
			originalWeight = keyData.MTOServiceItem.ActualWeight
		}
	default:
		// Make sure there's an actual weight since that's nullable
		originalWeight = r.MTOShipment.PrimeActualWeight
		if originalWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
	}

	value := fmt.Sprintf("%d", int(*originalWeight))
	return value, nil
}
