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
		return "", nil
	}

	value := fmt.Sprintf("%d", (*r.MTOShipment.BillableWeightCap).Int())
	return value, nil
}
