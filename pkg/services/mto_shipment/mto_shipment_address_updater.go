package mtoshipment

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// mtoShipmentAddressUpdater handles the db connection
type mtoShipmentAddressUpdater struct {
	planner        route.Planner
	addressCreator services.AddressCreator
	addressUpdater services.AddressUpdater
	checks         []addressUpdateValidator
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater(planner route.Planner, addressCreator services.AddressCreator, addressUpdater services.AddressUpdater) services.MTOShipmentAddressUpdater {
	return mtoShipmentAddressUpdater{
		planner:        planner,
		addressCreator: addressCreator,
		addressUpdater: addressUpdater,
		checks:         []addressUpdateValidator{checkAddressUpdateAllowed()},
	}
}

// isAddressOnShipment returns true if address is associated with the shipment, false if not
func isAddressOnShipment(address *models.Address, mtoShipment *models.MTOShipment) bool {
	addressIDs := []*uuid.UUID{
		mtoShipment.PickupAddressID,
		mtoShipment.DestinationAddressID,
		mtoShipment.SecondaryDeliveryAddressID,
		mtoShipment.SecondaryPickupAddressID,
		mtoShipment.TertiaryDeliveryAddressID,
		mtoShipment.TertiaryPickupAddressID,
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

func UpdateOriginSITServiceItemSITDeliveryMiles(planner route.Planner, addressCreator services.AddressCreator, shipment *models.MTOShipment, newAddress *models.Address, oldAddress *models.Address, appCtx appcontext.AppContext) (*models.MTOServiceItems, error) {
	// Change the SITDeliveryMiles of origin SIT service items
	var updatedMtoServiceItems models.MTOServiceItems

	eagerAssociations := []string{"MTOServiceItems.ReService.Code", "MTOServiceItems.SITOriginHHGActualAddress", "MTOServiceItems"}
	mtoShipment, err := FindShipment(appCtx, shipment.ID, eagerAssociations...)
	if err != nil {
		return &updatedMtoServiceItems, err
	}

	mtoServiceItems := mtoShipment.MTOServiceItems
	for _, s := range mtoServiceItems {
		serviceItem := s
		reServiceCode := serviceItem.ReService.Code
		if reServiceCode == models.ReServiceCodeDOPSIT ||
			reServiceCode == models.ReServiceCodeDOSFSC ||
			reServiceCode == models.ReServiceCodeIOPSIT ||
			reServiceCode == models.ReServiceCodeIOSFSC {

			if serviceItem.SITOriginHHGActualAddress != nil {
				milesCalculated, err := planner.ZipTransitDistance(appCtx, newAddress.PostalCode, serviceItem.SITOriginHHGActualAddress.PostalCode)
				if err != nil {
					return nil, err
				}
				serviceItem.SITDeliveryMiles = &milesCalculated

				updatedMtoServiceItems = append(updatedMtoServiceItems, serviceItem)

			} else {
				milesCalculated, err := planner.ZipTransitDistance(appCtx, oldAddress.PostalCode, newAddress.PostalCode)
				if err != nil {
					return nil, err
				}
				serviceItem.SITDeliveryMiles = &milesCalculated
				updatedMtoServiceItems = append(updatedMtoServiceItems, serviceItem)
			}
		}
	}

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		// update service item SITDeliveryMiles
		verrs, err := txnCtx.DB().ValidateAndUpdate(&updatedMtoServiceItems)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(newAddress.ID, err, verrs, "invalid input found while updating SIT delivery miles for service items")
		} else if err != nil {
			return apperror.NewQueryError("Service items", err, "")
		}

		return nil
	})

	// if there was a transaction error, we'll return nothing but the error
	if transactionError != nil {
		return nil, transactionError
	}

	return &updatedMtoServiceItems, nil
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
	var verrs *validate.Errors
	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	oldAddress := models.Address{}

	// Find the shipment, return error if not found.  If this shipment must be available to the prime,
	// do not include any shipments assigned to an external vendor.
	query := appCtx.DB().Q()
	if mustBeAvailableToPrime {
		query.Where("uses_external_vendor = FALSE")
	}
	err := query.Scope(utilities.ExcludeDeletedScope()).Eager("MTOServiceItems", "MTOServiceItems.ReService", "PickupAddress").Find(&mtoShipment, mtoShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "looking for mtoShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	if newAddress != nil {
		for _, check := range f.checks {
			if err = check.Validate(appCtx, newAddress, &mtoShipment); err != nil {
				switch e := err.(type) {
				case *validate.Errors:
					// Accumulate all validation errors
					verrs.Append(e)
				default:
					// Non-validation errors have priority and short-circuit doing any further checks
					return nil, err
				}
			}
		}
	}
	if verrs.HasAny() {
		return nil, verrs
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

	var address *models.Address
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		if newAddress.ID == uuid.Nil {
			// New address doesn't have an ID yet, it should be created
			address, err = f.addressCreator.CreateAddress(appCtx, newAddress)
		} else {
			// It has an ID, it should be updated
			address, err = f.addressUpdater.UpdateAddress(appCtx, newAddress, etag.GenerateEtag(oldAddress.UpdatedAt))
		}
		if err != nil {
			return apperror.NewQueryError("Address", err, "")
		}

		if mtoShipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			var shipment models.MTOShipment
			err := query.Scope(utilities.ExcludeDeletedScope()).Eager("DestinationAddress", "PickupAddress").Find(&shipment, mtoShipment.ID)
			if err != nil {
				return apperror.NewQueryError("Error when querying UB shipment", err, "")
			}
			// check that one of the addresses is OCONUS
			isShipmentOCONUS := models.IsShipmentOCONUS(shipment)
			if isShipmentOCONUS != nil && !*isShipmentOCONUS {
				return apperror.NewConflictError(shipment.ID, "At least one address for a UB shipment must be OCONUS")
			}
		}

		_, err = UpdateOriginSITServiceItemSITDeliveryMiles(f.planner, f.addressCreator, &mtoShipment, newAddress, &oldAddress, appCtx)
		if err != nil {
			return apperror.NewQueryError("No updated service items on shipment address change", err, "")
		}

		// update the service item pricing if relevant fields have changed..ie mileage
		err = models.UpdateEstimatedPricingForShipmentBasicServiceItems(appCtx.DB(), &mtoShipment, nil)
		if err != nil {
			return err
		}

		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return address, nil
}
