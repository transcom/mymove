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

var PPMShipmentUpdaterChecks = []ppmShipmentValidator{
	checkShipmentType(),
	checkShipmentID(),
	checkPPMShipmentID(),
	checkRequiredFields(),
	checkAdvanceAmountRequested(),
}

func NewPPMShipmentUpdater(ppmEstimator services.PPMEstimator) services.PPMShipmentUpdater {
	return &ppmShipmentUpdater{
		checks:    PPMShipmentUpdaterChecks,
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

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// This potentially updates the MTOShipment.Distance field so include it in the transaction
		estimatedIncentive, estimatedSITCost, err := f.estimator.EstimateIncentiveWithDefaultChecks(appCtx, *oldPPMShipment, updatedPPMShipment)
		if err != nil {
			return err
		}

		updatedPPMShipment.EstimatedIncentive = estimatedIncentive
		updatedPPMShipment.SITEstimatedCost = estimatedSITCost

		if appCtx.Session() != nil {
			if appCtx.Session().IsOfficeUser() {
				rejected := models.PPMAdvanceStatusRejected
				edited := models.PPMAdvanceStatusEdited
				approved := models.PPMAdvanceStatusApproved
				if oldPPMShipment.HasRequestedAdvance != nil && updatedPPMShipment.HasRequestedAdvance != nil {
					if !*oldPPMShipment.HasRequestedAdvance && *updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &edited
					} else if *oldPPMShipment.HasRequestedAdvance && !*updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &rejected
					}
				}
				if oldPPMShipment.AdvanceAmountRequested != nil && updatedPPMShipment.AdvanceAmountRequested != nil {
					if *oldPPMShipment.AdvanceAmountRequested != *updatedPPMShipment.AdvanceAmountRequested {
						updatedPPMShipment.AdvanceStatus = &edited
					}
					if *oldPPMShipment.AdvanceAmountRequested == *updatedPPMShipment.AdvanceAmountRequested && *oldPPMShipment.HasRequestedAdvance == *updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &approved
					}
				}
			}
		}

		if updatedPPMShipment.W2Address != nil {
			if verrs, errors := txnAppCtx.DB().ValidateAndSave(updatedPPMShipment.W2Address); verrs != nil && verrs.HasAny() {
				var id uuid.UUID
				if updatedPPMShipment.W2AddressID != nil {
					id = *updatedPPMShipment.W2AddressID
				}
				return apperror.NewInvalidInputError(id, errors, verrs, "Invalid input found while updating the W2 address for a PPMShipment.")
			} else if errors != nil {
				return apperror.NewQueryError("W2 address for ppmShipment", errors, "")
			}
			updatedPPMShipment.W2AddressID = &updatedPPMShipment.W2Address.ID
		}

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
