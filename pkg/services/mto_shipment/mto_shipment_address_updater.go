package mtoshipment

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// mtoShipmentAddressUpdater handles the db connection
type mtoShipmentAddressUpdater struct {
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater() services.MTOShipmentAddressUpdater {
	return mtoShipmentAddressUpdater{}
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

func UpdateSITServiceItemDestinationAddressToMTOShipmentAddress(mtoServiceItems *models.MTOServiceItems, newAddress *models.Address, appCtx appcontext.AppContext) (*models.MTOServiceItems, error) {
	// Change the address ID of destination SIT service items to match the address of the shipment address ID
	var updatedMtoServiceItems models.MTOServiceItems
	for _, s := range *mtoServiceItems {
		serviceItem := s
		reServiceCode := serviceItem.ReService.Code
		if reServiceCode == models.ReServiceCodeDDDSIT ||
			reServiceCode == models.ReServiceCodeDDFSIT ||
			reServiceCode == models.ReServiceCodeDDASIT ||
			reServiceCode == models.ReServiceCodeDDSFSC {

			serviceItem.SITDestinationFinalAddressID = &newAddress.ID
			updatedMtoServiceItems = append(updatedMtoServiceItems, serviceItem)
			transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
				// update service item final destination address ID to match shipment address ID
				verrs, err := txnCtx.DB().ValidateAndUpdate(&serviceItem)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(newAddress.ID, err, verrs, "invalid input found while updating final destination address of service item")
				} else if err != nil {
					return apperror.NewQueryError("Service item", err, "")
				}

				return nil
			})

			// if there was a transaction error, we'll return nothing but the error
			if transactionError != nil {
				return nil, transactionError
			}
		}
	}

	return &updatedMtoServiceItems, nil
}

// UpdateMTOShipmentAddress updates an address on an MTO shipment.
// Since address records have no parent id, caller must supply the mtoShipmentID associated with this address.
// Function will check that the etag matches before making the update.
// If mustBeAvailableToPrime is set, update will not happen unless the mto with which the address + shipment is associated
// is also availableToPrime and not an external vendor shipment.
func (f mtoShipmentAddressUpdater) UpdateMTOShipmentAddress(appCtx appcontext.AppContext, newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error) {

	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	oldAddress := models.Address{}

	// Find the shipment, return error if not found.  If this shipment must be available to the prime,
	// do not include any shipments assigned to an external vendor.
	query := appCtx.DB().Q()
	if mustBeAvailableToPrime {
		query.Where("uses_external_vendor = FALSE")
	}
	err := query.Scope(utilities.ExcludeDeletedScope()).Eager("MTOServiceItems", "MTOServiceItems.ReService").Find(&mtoShipment, mtoShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "looking for mtoShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	if mustBeAvailableToPrime {
		// Make sure the associated move is available to the prime
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
		mtoAvailableToPrime, _ := mtoChecker.MTOAvailableToPrime(appCtx, mtoShipment.MoveTaskOrderID)
		if !mtoAvailableToPrime {
			return nil, apperror.NewNotFoundError(mtoShipment.MoveTaskOrderID, "looking for moveTaskOrder")
		}
	}

	// Find the address, return error if not found
	err = appCtx.DB().Find(&oldAddress, newAddress.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(newAddress.ID, "looking for address")
		default:
			return nil, apperror.NewQueryError("Address", err, "")
		}
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAddress.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(newAddress.ID, err)
	}

	// Check that address is associated with this shipment
	if !isAddressOnShipment(newAddress, &mtoShipment) {
		return nil, apperror.NewConflictError(newAddress.ID, ": Address is not associated with the provided MTOShipmentID.")
	}

	// Make the update and create a InvalidInput Error if there were validation issues
	verrs, err := appCtx.DB().ValidateAndSave(newAddress)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(newAddress.ID, err, verrs, "")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("Address", err, "")
	}

	_, err = UpdateSITServiceItemDestinationAddressToMTOShipmentAddress(&mtoShipment.MTOServiceItems, newAddress, appCtx)
	if err != nil {
		return nil, apperror.NewQueryError("No updated service items on shipment address change", err, "")
	}

	// Get the updated address and return
	updatedAddress := models.Address{}
	err = appCtx.DB().Find(&updatedAddress, newAddress.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(newAddress.ID, "looking for Address")
		default:
			return nil, apperror.NewQueryError("Address", err, fmt.Sprintf("Unexpected error after saving: %v", err))
		}
	}
	return &updatedAddress, nil
}
