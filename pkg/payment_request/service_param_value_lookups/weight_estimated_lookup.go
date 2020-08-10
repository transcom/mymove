package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// WeightEstimatedLookup does lookup on actual weight billed
type WeightEstimatedLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightEstimatedLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	estimatedWeight := r.MTOShipment.PrimeEstimatedWeight
	if estimatedWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find estimated weight for MTOShipmentID [%s]", r.MTOShipment.ID)
	}

	value := fmt.Sprintf("%d", int(*estimatedWeight))
	return value, nil
}
