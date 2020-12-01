package mtoshipment

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// mtoShipmentAddressUpdater handles the db connection
type mtoShipmentAddressUpdater struct {
	db *pop.Connection
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater(db *pop.Connection) services.MTOShipmentAddressUpdater {
	return mtoShipmentAddressUpdater{
		db: db}
}

// isAddressOnShipment returns true if address is associated with the shipment, false if not
func isAddressOnShipment(address *models.Address, mtoShipment *models.MTOShipment) bool {
	addressIDs := []*uuid.UUID{
		mtoShipment.PickupAddressID,
		mtoShipment.DestinationAddressID,
		mtoShipment.SecondaryDeliveryAddressID,
		mtoShipment.SecondaryPickupAddressID,
	}

	for _, id := range addressIDs {
		if id != nil {
			if *id == address.ID {
				return true
			}
		}
	}
	return false
}

// UpdateMTOShipmentAddress updates an address on an MTO shipment.
// Since address records have no parent id, caller must supply the mtoShipmentID associated with this address.
// Function will check that the etag matches before making the update
// If mustBeAvailableToPrime is set, update will not happen unless the mto with which the address + shipment is associated, is also availableToPrime
func (f mtoShipmentAddressUpdater) UpdateMTOShipmentAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error) {

	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	oldAddress := models.Address{}

	// Find the shipment, return error if not found
	err := f.db.Find(&mtoShipment, mtoShipmentID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(mtoShipmentID, "looking for mtoShipment")
		}
	}

	if mustBeAvailableToPrime == true {
		// Make sure the associated move is available to the prime
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(f.db)
		mtoAvailableToPrime, _ := mtoChecker.MTOAvailableToPrime(mtoShipment.MoveTaskOrderID)
		if !mtoAvailableToPrime {
			return nil, services.NewNotFoundError(mtoShipment.MoveTaskOrderID, "looking for moveTaskOrder")
		}
	}

	// Find the address, return error if not found
	err = f.db.Find(&oldAddress, newAddress.ID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(newAddress.ID, "looking for address")
		}
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAddress.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(newAddress.ID, err)
	}

	// Check that address is associated with this shipment
	if !isAddressOnShipment(newAddress, &mtoShipment) {
		return nil, services.NewConflictError(newAddress.ID, ": Address is not associated with the provided MTOShipmentID.")
	}

	// Make the update and create a InvalidInput Error if there were validation issues
	verrs, err := f.db.ValidateAndSave(newAddress)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(newAddress.ID, err, verrs, "")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("Address", err, "")
	}

	// Get the updated address and return
	updatedAddress := models.Address{}
	err = f.db.Find(&updatedAddress, newAddress.ID)
	if err != nil {
		return nil, services.NewQueryError("Address", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedAddress, nil
}
