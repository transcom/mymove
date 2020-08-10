package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// WeightActualLookup does lookup on actual weight billed
type WeightActualLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightActualLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Make sure there's an actual weight since that's nullable
	actualWeight := r.MTOShipment.PrimeActualWeight
	if actualWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", r.MTOShipment.ID)
	}

	value := fmt.Sprintf("%d", int(*actualWeight))
	return value, nil
}
