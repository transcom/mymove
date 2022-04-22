package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ppmShipmentUpdater struct {
	checks    []ppmShipmentValidator
	estimator services.PPMEstimator
}

func NewPPMShipmentUpdater(ppmEstimator services.PPMEstimator) services.PPMShipmentUpdater {
	return &ppmShipmentUpdater{
		checks: []ppmShipmentValidator{
			checkShipmentType(),
			checkShipmentID(),
			checkPPMShipmentID(),
			checkRequiredFields(),
			checkAdvance(),
		},
		estimator: ppmEstimator,
	}
}

func (f *ppmShipmentUpdater) UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, mtoShipmentID uuid.UUID) (*models.PPMShipment, error) {
	return f.updatePPMShipment(appCtx, ppmShipment, mtoShipmentID, f.checks...)
}

func (f *ppmShipmentUpdater) updatePPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, mtoShipmentID uuid.UUID, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	if ppmShipment == nil {
		return nil, nil
	}

	oldPPMShipment, err := models.FetchPPMShipmentFromMTOShipmentID(appCtx.DB(), mtoShipmentID)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment := mergePPMShipment(*ppmShipment, oldPPMShipment)

	err = validatePPMShipment(appCtx, *updatedPPMShipment, oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	estimatedIncentive, err := f.estimator.EstimateIncentiveWithDefaultChecks(appCtx, *oldPPMShipment, updatedPPMShipment)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment.EstimatedIncentive = estimatedIncentive
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedPPMShipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedPPMShipment.ID, err, verrs, "Invalid input found while updating the PPMShipments.")
		} else if err != nil {
			return apperror.NewQueryError("PPMShipments", err, "")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return updatedPPMShipment, nil
}
