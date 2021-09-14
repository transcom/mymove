package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// WeightAdjustedLookup does lookup on adjusted weight billed
type WeightAdjustedLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightAdjustedLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	if r.MTOShipment.BillableWeightCap == nil {
		return "", fmt.Errorf("could not find adjusted weight for MTOServiceItemID [%s]", keyData.MTOServiceItem.ID)
	}

	value := fmt.Sprintf("%d", int(*r.MTOShipment.BillableWeightCap))
	return value, nil
}
