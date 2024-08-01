package boatshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type boatShipmentUpdater struct {
	checks []boatShipmentValidator
}

var BoatShipmentUpdaterChecks = []boatShipmentValidator{
	checkShipmentType(),
	checkShipmentID(),
	checkBoatShipmentID(),
	checkRequiredFields(),
}

func NewBoatShipmentUpdater() services.BoatShipmentUpdater {
	return &boatShipmentUpdater{
		checks: BoatShipmentUpdaterChecks,
	}
}

func (f *boatShipmentUpdater) UpdateBoatShipmentWithDefaultCheck(appCtx appcontext.AppContext, boatShipment *models.BoatShipment, mtoShipmentID uuid.UUID) (*models.BoatShipment, error) {
	return f.updateBoatShipment(appCtx, boatShipment, mtoShipmentID, f.checks...)
}

func (f *boatShipmentUpdater) updateBoatShipment(appCtx appcontext.AppContext, boatShipment *models.BoatShipment, mtoShipmentID uuid.UUID, checks ...boatShipmentValidator) (*models.BoatShipment, error) {
	if boatShipment == nil {
		return nil, nil
	}

	oldBoatShipment, err := FindBoatShipmentByMTOID(appCtx, mtoShipmentID)
	if err != nil {
		return nil, err
	}

	updatedBoatShipment, err := mergeBoatShipment(*boatShipment, oldBoatShipment)
	if err != nil {
		return nil, err
	}

	err = validateBoatShipment(appCtx, *updatedBoatShipment, oldBoatShipment, &oldBoatShipment.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedBoatShipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedBoatShipment.ID, err, verrs, "Invalid input found while updating the BoatShipments.")
		} else if err != nil {
			return apperror.NewQueryError("BoatShipments", err, "")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return updatedBoatShipment, nil
}
