package mtoshipment

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

//ErrNotFound is returned when a given mto shipment is not found
type ErrNotFound struct {
	id uuid.UUID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("mto shipment id: %s not found", e.id.String())
}

type mtoShipmentFetcher struct {
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMTOShipmentFetcher(db *pop.Connection) services.MTOShipmentFetcher {
	return &mtoShipmentFetcher{db}
}

//FetchMTOShipment retrieves a MTOShipment for a given UUID
func (f mtoShipmentFetcher) FetchMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error) {
	shipment := &models.MTOShipment{}
	if err := f.db.Eager().Find(shipment, mtoShipmentID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MTOShipment{}, ErrNotFound{mtoShipmentID}
		default:
			return &models.MTOShipment{}, err
		}
	}

	return shipment, nil
}

// type mtoShipmentUpdater struct {
// 	db *pop.Connection
// 	mtoShipmentFetcher
// }

// // NewMTOShipmentUpdater creates a new struct with the service dependencies
// func NewMTOShipmentUpdater(db *pop.Connection) services.TOShipmentUpdater {
// 	return &mtoShipmentUpdater{db, mtoShipmentFetcher{db}}
// }

// //UpdateMTOShipment updates the mto shipment
// func (f moveTaskOrderFetcher) UpdateMTOShipment((mtoShipment *models.MTOShipment) (*models.MTOShipment, error) {
// 	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
// 	if err != nil {
// 		return &models.MoveTaskOrder{}, err
// 	}
// 	mto.IsAvailableToPrime = true
// 	vErrors, err := f.db.ValidateAndUpdate(mto)
// 	if vErrors.HasAny() {
// 		return &models.MoveTaskOrder{}, ErrInvalidInput{}
// 	}
// 	if err != nil {
// 		return &models.MoveTaskOrder{}, err
// 	}
// 	return mto, nil
// }
