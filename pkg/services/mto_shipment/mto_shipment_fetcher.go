package mtoshipment

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

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

func (f mtoShipmentFetcher) ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) ([]models.MTOShipment, error) {
	var move models.Move
	err := appCtx.DB().Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "move not found")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	var shipments []models.MTOShipment
	err = appCtx.DB().
		Eager(
			"MTOServiceItems.ReService",
			"MTOAgents",
			"PickupAddress",
			"SecondaryPickupAddress",
			"DestinationAddress",
			"SecondaryDeliveryAddress",
			"MTOServiceItems.Dimensions",
			"Reweigh",
			"SITExtensions",
			"StorageFacility.Address",
		).
		Where("move_id = ?", moveID).
		Order("created_at asc").
		All(&shipments)

	if err != nil {
		return nil, err
	}

	return shipments, nil
}
