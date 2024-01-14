package reweigh

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type reweighFetcher struct{}

func NewReweighFetcher() services.ReweighFetcher {
	return &reweighFetcher{}
}

// Returns a map of reweighs for the provided shipment IDs
func (f *reweighFetcher) ListReweighsByShipmentIDs(appCtx appcontext.AppContext, shipmentIDs []uuid.UUID) (map[uuid.UUID]models.Reweigh, error) {
	var reweighs []models.Reweigh
	err := appCtx.DB().Where("shipment_id IN (?)", shipmentIDs).All(&reweighs)
	if err != nil {
		return nil, err
	}

	reweighMap := make(map[uuid.UUID]models.Reweigh)
	for _, reweigh := range reweighs {
		reweighMap[reweigh.ShipmentID] = reweigh
	}
	return reweighMap, nil
}
