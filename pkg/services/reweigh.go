package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ReweighCreator creates a reweigh
type ReweighCreator interface {
	CreateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh) (*models.Reweigh, error)
}

// ReweighUpdater updates a reweigh
type ReweighUpdater interface {
	UpdateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error)
}

// ReweighFetcher allows us to fetch reweights. This is primarily used during the reweigh update of a
// diverted shipment
type ReweighFetcher interface {
	ListReweighsByShipmentIDs(appCtx appcontext.AppContext, shipmentIDs []uuid.UUID) (map[uuid.UUID]models.Reweigh, error)
}
