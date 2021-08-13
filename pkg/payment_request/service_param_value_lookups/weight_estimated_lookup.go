package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightEstimatedLookup does lookup on estimated weight billed
type WeightEstimatedLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightEstimatedLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	var estimatedWeight *unit.Pound

	switch keyData.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDOSHUT,
		models.ReServiceCodeDDSHUT,
		models.ReServiceCodeIOSHUT,
		models.ReServiceCodeIDSHUT:
		estimatedWeight = keyData.MTOServiceItem.EstimatedWeight
		if estimatedWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find estimated weight for MTOServiceItemID [%s]", keyData.MTOServiceItem.ID)
		}
	default:
		// Make sure there's an estimated weight since that's nullable
		estimatedWeight = r.MTOShipment.PrimeEstimatedWeight
		if estimatedWeight == nil {
			return "", nil
		}
	}

	value := fmt.Sprintf("%d", int(*estimatedWeight))
	return value, nil
}
