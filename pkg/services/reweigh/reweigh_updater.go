package reweigh

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// reweighUpdater needs to be updates to have checks for validation
type reweighUpdater struct {
	checks       []reweighValidator
	recalculator services.PaymentRequestShipmentRecalculator
}

// NewReweighUpdater creates a new struct with the service dependencies
func NewReweighUpdater(moveAvailabilityChecker services.MoveTaskOrderChecker, recalculator services.PaymentRequestShipmentRecalculator) services.ReweighUpdater {
	return &reweighUpdater{
		checks: []reweighValidator{
			checkShipmentID(),
			checkReweighID(),
			checkRequiredFields(),
			checkPrimeAvailability(moveAvailabilityChecker),
		},
		recalculator: recalculator,
	}
}

// UpdateReweighCheck passes the Prime validator key to CreateReweigh
func (f *reweighUpdater) UpdateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error) {
	return f.UpdateReweigh(appCtx, reweigh, eTag, f.checks...)
}

// UpdateReweigh updates the Reweigh table
func (f *reweighUpdater) UpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	var updatedReweigh *models.Reweigh

	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		updatedReweigh, err = f.doUpdateReweigh(txnAppCtx, reweigh, eTag, checks...)
		return err
	})
	if transactionError != nil {
		return nil, transactionError
	}

	return updatedReweigh, nil
}

func (f *reweighUpdater) doUpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	oldReweigh := models.Reweigh{}

	// Find the reweigh, return error if not found
	err := appCtx.DB().Find(&oldReweigh, reweigh.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Reweigh")
		default:
			return nil, apperror.NewQueryError("Reweigh", err, "")
		}
	}

	shipment := models.MTOShipment{}
	// Find the shipment, return error if not found
	err = appCtx.DB().Find(&shipment, reweigh.ShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}
	oldReweigh.Shipment = shipment

	err = validateReweigh(appCtx, *reweigh, &oldReweigh, &oldReweigh.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	if reweigh.VerificationReason != nil {
		now := time.Now()
		reweigh.VerificationProvidedAt = &now
	}

	newReweigh := mergeReweigh(*reweigh, &oldReweigh)

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldReweigh.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(reweigh.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := appCtx.DB().ValidateAndSave(newReweigh)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(reweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("Reweigh", err, "")
	}

	// Get the updated reweigh and return
	updatedReweigh := models.Reweigh{}
	err = appCtx.DB().Find(&updatedReweigh, reweigh.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "looking for Reweigh")
		default:
			return nil, apperror.NewQueryError("Reweigh", err, fmt.Sprintf("Unexpected error after saving: %v", err))
		}
	}

	// Need to pull out this common code. It is from reweigh requester
	err = appCtx.DB().Q().
		Eager("Reweigh").
		Find(&shipment, reweigh.ShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ShipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("Shipment", err, "")
		}
	}

	// Recalculate payment request for the shipment, if the reweigh weight changed
	if reweighChanged(oldReweigh, updatedReweigh) {
		_, err = f.recalculator.ShipmentRecalculatePaymentRequest(appCtx, reweigh.ShipmentID)
		if err != nil {
			return nil, err
		}
	}

	return &updatedReweigh, nil
}
