package reweigh

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
)

// reweighCreator sets up the service object
type reweighCreator struct {
	checks []reweighValidator
}

// NewReweighCreator creates a new struct with the service dependencies
func NewReweighCreator(db *pop.Connection) services.ReweighCreator {
	return &reweighCreator{
		checks: []reweighValidator{
			checkShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreateReweighCheck passes the Prime validator key to CreateReweigh
func (f *reweighCreator) CreateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh) (*models.Reweigh, error) {
	return f.CreateReweigh(appCtx, reweigh, f.checks...)
}

// CreateReweigh creates a reweigh
func (f *reweighCreator) CreateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, checks ...reweighValidator) (*models.Reweigh, error) {
	// Get existing shipment information for validation
	mtoShipment := &models.MTOShipment{}
	// Find the shipment, return error if not found
	err := appCtx.DB().Find(mtoShipment, reweigh.ShipmentID)

	if err != nil {
		return nil, apperror.NewNotFoundError(reweigh.ShipmentID, "while looking for MTOShipment")
	}

	err = validateReweigh(appCtx, *reweigh, nil, mtoShipment, checks...)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndCreate(reweigh)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("Reweigh", err, "")
	}

	return reweigh, nil
}
