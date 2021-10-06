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

func (r WeightOriginalLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var originalWeight *unit.Pound

	switch keyData.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDOSHUT,
		models.ReServiceCodeDDSHUT,
		models.ReServiceCodeIOSHUT,
		models.ReServiceCodeIDSHUT:
		originalWeight = keyData.MTOServiceItem.ActualWeight
		if originalWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find actual weight for MTOServiceItemID [%s]", keyData.MTOServiceItem.ID)
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
