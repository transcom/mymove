package sitentrydateupdate

import (
	"database/sql"
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

// finds the current service item
// replaces sit entry date
// sends back updated service item
func (p sitEntryDateUpdater) UpdateSitEntryDate(appCtx appcontext.AppContext, s *models.SITEntryDateUpdate) (*models.MTOServiceItem, error) {
	// we will need to update not only the target SIT service item, but it's sister service item
	// of additional days since the entry dates can't be the same
	// and the SIT add'l days service item will need to be the NEXT day
	var serviceItem models.MTOServiceItem
	var serviceItemAdditionalDays models.MTOServiceItem

	// finding the service item and populating serviceItem variable
	err := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder",
		"SITDestinationFinalAddress",
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

	// the service code can either be DOFSIT or DDFSIT
	serviceItemCode := serviceItem.ReService.Code
	if serviceItemCode != models.ReServiceCodeDOFSIT && serviceItemCode != models.ReServiceCodeDDFSIT {
		return nil, apperror.NewUnprocessableEntityError("You cannot change the SIT entry date of this service item.")
	}

	// looping through each service item in the shipment based on the service item code
	// looking for the sister service item of add'l days
	if serviceItemCode == models.ReServiceCodeDOFSIT {
		for _, si := range shipment.MTOServiceItems {
			if si.ReService.Code == models.ReServiceCodeDOASIT {
				serviceItemAdditionalDays = si
				break
			}
		}
	} else if serviceItemCode == models.ReServiceCodeDDFSIT {
		for _, si := range shipment.MTOServiceItems {
			if si.ReService.Code == models.ReServiceCodeDDASIT {
				serviceItemAdditionalDays = si
				break
			}
		}
	} else {
		return nil, apperror.NewUnprocessableEntityError("This service item's SIT entry date cannot be updated due to being an uneditable service code.")
	}

	// updating service item struct with the new SIT entry date
	// updating sister service item to have the next day for SIT entry date
	if s.SITEntryDate == nil {
		return nil, apperror.NewUnprocessableEntityError("You must provide the SIT entry date in the request")
	} else if s.SITEntryDate != nil {
		serviceItem.SITEntryDate = s.SITEntryDate
		dayAfter := s.SITEntryDate.Add(24 * time.Hour)
		serviceItemAdditionalDays.SITEntryDate = &dayAfter
	}

	// Make the update to both service items and create a InvalidInputError if there were validation issues
	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		verrs, err := txnCtx.DB().ValidateAndUpdate(&serviceItem)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(s.ID, err, verrs, "invalid input found while updating service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}

		verrs, err = txnCtx.DB().ValidateAndUpdate(&serviceItemAdditionalDays)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(s.ID, err, verrs, "invalid input found while updating service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}
		// Done with updates to service item
		return nil
	})

	// if there was a transaction error, we'll return nothing but the error
	if transactionError != nil {
		return nil, transactionError
	}

	// upon successful validation and update - we'll send back the updated service item
	return &serviceItem, nil
}
