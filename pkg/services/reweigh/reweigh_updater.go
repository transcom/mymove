package reweigh

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// reweighUpdater needs to be updates to have checks for validation
type reweighUpdater struct {
}

// NewReweighUpdater creates a new struct with the service dependencies
func NewReweighUpdater() services.ReweighUpdater {
	return &reweighUpdater{}
}

// UpdateReweigh updates the Reweigh table
func (f *reweighUpdater) UpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error) {
	oldReweigh := models.Reweigh{}

	// Find the reweigh, return error if not found
	err := appCtx.DB().Find(&oldReweigh, reweigh.ID)
	if err != nil {
		return nil, services.NewNotFoundError(reweigh.ID, "while looking for a reweigh")
	}

	if reweigh.VerificationReason != nil {
		now := time.Now()
		reweigh.VerificationProvidedAt = &now
	}

	newReweigh := mergeReweigh(*reweigh, &oldReweigh)

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldReweigh.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(reweigh.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := appCtx.DB().ValidateAndSave(newReweigh)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(reweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("Reweigh", err, "")
	}

	// Get the updated reweigh and return
	updatedReweigh := models.Reweigh{}
	err = appCtx.DB().Find(&updatedReweigh, reweigh.ID)
	if err != nil {
		return nil, services.NewQueryError("Reweigh", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	// Need to pull out this common code. It is from reweigh requester
	var shipment models.MTOShipment
	err = appCtx.DB().Q().
		Eager("Reweigh").
		Find(&shipment, reweigh.ShipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(reweigh.ShipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	updatedReweigh.Shipment = shipment

	return &updatedReweigh, nil
}
