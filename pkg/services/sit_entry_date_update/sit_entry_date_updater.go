package sitentrydateupdate

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
	var serviceItem models.MTOServiceItem
	findServiceItemQuery := appCtx.DB().Q()

	// finding the current service item
	err := findServiceItemQuery.Find(&serviceItem, s.ID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(s.ID, "while looking for service item")
		default:
			return nil, apperror.NewQueryError("ServiceItem", err, "")
		}
	}

	// updating service item struct with the new SIT entry date
	if s.SITEntryDate != nil {
		serviceItem.SITEntryDate = s.SITEntryDate
	}

	// Make the update and create a InvalidInputError if there were validation issues
	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		verrs, err := txnCtx.DB().ValidateAndUpdate(&serviceItem)
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
