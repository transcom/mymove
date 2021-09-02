package reweigh

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// reweighUpdater handles the db connection
type reweighUpdater struct {
	db     *pop.Connection
	checks []reweighValidator
}

// NewReweighUpdater creates a new struct with the service dependencies
func NewReweighUpdater(db *pop.Connection, moveAvailabilityChecker services.MoveTaskOrderChecker) services.ReweighUpdater {
	return &reweighUpdater{
		db: db,
		checks: []reweighValidator{
			checkShipmentID(),
			checkReweighID(),
			checkRequiredFields(),
			checkPrimeAvailability(moveAvailabilityChecker),
		},
	}
}

// UpdateReweighCheck passes the Prime validator key to CreateReweigh
func (f *reweighUpdater) UpdateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error) {
	return f.UpdateReweigh(appCtx, reweigh, eTag, f.checks...)
}

// UpdateReweigh updates the Reweigh table
func (f *reweighUpdater) UpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	oldReweigh := models.Reweigh{}

	// Find the reweigh, return error if not found
	err := appCtx.DB().Find(&oldReweigh, reweigh.ID)
	if err != nil {
		return nil, services.NewNotFoundError(reweigh.ID, "while looking for Reweigh")
	}

	shipment := models.MTOShipment{}
	// Find the shipment, return error if not found
	err = appCtx.DB().Find(&shipment, reweigh.ShipmentID)
	if err != nil {
		return nil, services.NewNotFoundError(reweigh.ID, "while looking for Shipment")
	}
	oldReweigh.Shipment = shipment

	err = validateReweigh(appCtx, *reweigh, &oldReweigh, &oldReweigh.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldReweigh.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(reweigh.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := f.db.ValidateAndSave(reweigh)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(reweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("Reweigh", err, "")
	}

	// Get the updated reweigh and return
	updatedReweigh := models.Reweigh{}
	err = f.db.Find(&updatedReweigh, reweigh.ID)
	if err != nil {
		return nil, services.NewQueryError("Reweigh", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedReweigh, nil
}
