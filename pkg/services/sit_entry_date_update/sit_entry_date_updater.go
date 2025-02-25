package sitentrydateupdate

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

type sitEntryDateUpdater struct {
}

// NewSitEntryDateUpdater creates a new sitEntryDateUpdater struct
func NewSitEntryDateUpdater() services.SitEntryDateUpdater {
	return &sitEntryDateUpdater{}
}

// finds the current service item & it's sister add'l days service item
// replaces sit entry date for both
// sends back updated service item
func (p sitEntryDateUpdater) UpdateSitEntryDate(appCtx appcontext.AppContext, s *models.SITEntryDateUpdate) (*models.MTOServiceItem, error) {
	// we will need to update not only the target SIT service item, but it's sister service item
	// of additional days since the entry dates can't be the same
	// and the SIT add'l days service item will need to be the NEXT day
	var serviceItem models.MTOServiceItem
	var serviceItemAdditionalDays models.MTOServiceItem

	// finding the service item and populating serviceItem variable
	// passing in relations so we can get ReService codes & add'l info
	err := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder",
		"ReService",
	).Find(&serviceItem, s.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(s.ID, "while looking for service item")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	// eager associations is needed to get data from other tables
	// retrieving the shipment to get the other service items
	eagerAssociations := []string{"MoveTaskOrder", "MTOServiceItems", "MTOServiceItems.ReService"}
	shipment, err := mtoshipment.NewMTOShipmentFetcher().GetShipment(appCtx, *serviceItem.MTOShipmentID, eagerAssociations...)
	if err != nil {
		return nil, apperror.NewQueryError("Shipment", err, "")
	}

	// the service code can either be DOFSIT/DDFSIT or IOFSIT/IDFSIT
	serviceItemCode := serviceItem.ReService.Code
	if serviceItemCode != models.ReServiceCodeDOFSIT && serviceItemCode != models.ReServiceCodeDDFSIT &&
		serviceItemCode != models.ReServiceCodeIOFSIT && serviceItemCode != models.ReServiceCodeIDFSIT {
		return nil, apperror.NewUnprocessableEntityError(string(serviceItemCode) + "You cannot change the SIT entry date of this service item.")
	}

	// looping through each service item in the shipment based on the service item code
	// then looking for the sister service item of add'l days
	// once found, we'll set the value of variable to that service item
	// so now we have the 1st day of SIT service item & the add'l days SIT service item
	if serviceItemCode == models.ReServiceCodeDOFSIT || serviceItemCode == models.ReServiceCodeIOFSIT {
		for _, si := range shipment.MTOServiceItems {
			if si.ReService.Code == models.ReServiceCodeDOASIT || si.ReService.Code == models.ReServiceCodeIOASIT {
				serviceItemAdditionalDays = si
				break
			}
		}
	} else if serviceItemCode == models.ReServiceCodeDDFSIT || serviceItemCode == models.ReServiceCodeIDFSIT {
		for _, si := range shipment.MTOServiceItems {
			if si.ReService.Code == models.ReServiceCodeDDASIT || si.ReService.Code == models.ReServiceCodeIDASIT {
				serviceItemAdditionalDays = si
				break
			}
		}
	} else {
		// if it is not either service codes, then we shouldn't be updating the SIT entry date this way
		return nil, apperror.NewUnprocessableEntityError("This service item's SIT entry date cannot be updated due to being an uneditable service code.")
	}

	// updating service item struct with the new SIT entry date
	// updating sister service item to have the next day for SIT entry date
	if s.SITEntryDate == nil {
		return nil, apperror.NewUnprocessableEntityError("You must provide the SIT entry date in the request")
	}

	// The new SIT entry date must be before SIT departure date
	if serviceItem.SITDepartureDate != nil && !s.SITEntryDate.Before(*serviceItem.SITDepartureDate) {
		return nil, apperror.NewUnprocessableEntityError(fmt.Sprintf("the SIT Entry Date (%s) must be before the SIT Departure Date (%s)",
			s.SITEntryDate.Format("2006-01-02"), serviceItem.SITDepartureDate.Format("2006-01-02")))
	}

	serviceItem.SITEntryDate = s.SITEntryDate
	dayAfter := s.SITEntryDate.Add(24 * time.Hour)
	serviceItemAdditionalDays.SITEntryDate = &dayAfter

	// Make the update to both service items and create a InvalidInputError if there were validation issues
	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		// updating 1st day of SIT service item
		verrs, err := txnCtx.DB().ValidateAndUpdate(&serviceItem)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(s.ID, err, verrs, "invalid input found while updating service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}

		// updating add'l days of SIT service item
		verrs, err = txnCtx.DB().ValidateAndUpdate(&serviceItemAdditionalDays)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(s.ID, err, verrs, "invalid input found while updating service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}
		// Done with updates to service items, will return nil if there were no errors
		return nil
	})

	// if there was a transaction error, we'll return nothing but the error
	if transactionError != nil {
		return nil, transactionError
	}

	// upon successful validation and update - we'll send back the updated service item
	return &serviceItem, nil
}
