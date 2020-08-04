package mtoshipment

import (
	"database/sql"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type mtoShipmentFetcher struct {
	db *pop.Connection
}

func (f mtoShipmentFetcher) ListMTOShipments(moveTaskOrderID uuid.UUID) ([]models.MTOShipment, error) {
	var mtoShipments []models.MTOShipment
	err := f.db.Where("move_order_id = $1", moveTaskOrderID).Eager().All(&mtoShipments)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MTOShipment{}, services.NotFoundError{}
		default:
			return []models.MTOShipment{}, err
		}
	}
	return mtoShipments, nil
}

// NewMTOShipmentFetcher creates a new struct with the service dependencies
func NewMTOShipmentFetcher(db *pop.Connection) services.MTOShipmentFetcher {
	return &mtoShipmentFetcher{db}
}

//FetchMTOShipment retrieves an MTOShipment for a given UUID
func (f mtoShipmentFetcher) FetchMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error) {
	mto := &models.MTOShipment{}
	if err := f.db.Eager("DestinationAddress",
		"PickupAddress",
		"MTOAgents",
	).Find(mto, mtoShipmentID); err != nil {

		switch err {
		case sql.ErrNoRows:
			return &models.MTOShipment{}, services.NewNotFoundError(mtoShipmentID, "")
		default:
			return &models.MTOShipment{}, err
		}
	}

	return mto, nil
}
