package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightReweighLookup does lookup on estimated weight billed
type WeightReweighLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightReweighLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {

	err := appCtx.DB().Load(&r.MTOShipment, "Reweigh")
	if err != nil {
		return "", err
	}

	var reweighWeight *unit.Pound
	// Make sure there's a reweigh weight since that's nullable
	if r.MTOShipment.Reweigh != nil {
		reweighWeight = r.MTOShipment.Reweigh.Weight
	}

	if reweighWeight == nil {
		return "", nil
	}

	value := fmt.Sprintf("%d", int(*reweighWeight))
	return value, nil
}
