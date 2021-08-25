package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightActualLookup does lookup on actual weight billed
type WeightActualLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightActualLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var actualWeight *unit.Pound

	switch keyData.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDOSHUT,
		models.ReServiceCodeDDSHUT,
		models.ReServiceCodeIOSHUT,
		models.ReServiceCodeIDSHUT:
		actualWeight = keyData.MTOServiceItem.ActualWeight
		if actualWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find actual weight for MTOServiceItemID [%s]", keyData.MTOServiceItem.ID)
		}
	default:
		// Make sure there's an actual weight since that's nullable
		actualWeight = r.MTOShipment.PrimeActualWeight
		if actualWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", r.MTOShipment.ID)
		}
	}

	value := fmt.Sprintf("%d", int(*actualWeight))
	return value, nil
}
