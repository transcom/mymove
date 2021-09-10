package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type mtoShipmentFetcher struct {
}

// NewMTOShipmentFetcher creates a new MTOShipmentFetcher struct that supports ListMTOShipments
func NewMTOShipmentFetcher() services.MTOShipmentFetcher {
	return &mtoShipmentFetcher{}
}

func (f mtoShipmentFetcher) ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.MTOShipments, error) {
	var move models.Move
	err := appCtx.DB().Find(&move, moveID)
	if err != nil {
		return nil, services.NewNotFoundError(moveID, "move not found")
	}

	var shipments models.MTOShipments
	// TODO: These associations could be preloaded, but it will require Pop 5.3.4 to land first as it
	//   has a fix for using a "has_many" association that has a pointer-based foreign key (like the
	//   case with "MTOServiceItems.ReService"). There appear to be other changes that will need to be
	//   made for Pop 5.3.4 though (see https://ustcdp3.slack.com/archives/CP497TGAU/p1620421441217700).
	err = appCtx.DB().Eager("MTOServiceItems.ReService", "MTOAgents", "PickupAddress", "SecondaryPickupAddress", "DestinationAddress", "SecondaryDeliveryAddress", "MTOServiceItems.Dimensions", "Reweigh", "SITExtensions").
		Where("move_id = ?", moveID).
		Order("created_at asc").
		All(&shipments)

	if err != nil {
		return nil, err
	}

	return &shipments, nil
}
