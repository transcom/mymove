package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ppmShipmentUpdater struct {
	checks         []ppmShipmentValidator
	estimator      services.PPMEstimator
	addressCreator services.AddressCreator
	addressUpdater services.AddressUpdater
}

var PPMShipmentUpdaterChecks = []ppmShipmentValidator{
	checkShipmentType(),
	checkShipmentID(),
	checkPPMShipmentID(),
	checkRequiredFields(),
	checkAdvanceAmountRequested(),
}

func NewPPMShipmentUpdater(ppmEstimator services.PPMEstimator, addressCreator services.AddressCreator, addressUpdater services.AddressUpdater) services.PPMShipmentUpdater {
	return &ppmShipmentUpdater{
		checks:         PPMShipmentUpdaterChecks,
		estimator:      ppmEstimator,
		addressCreator: addressCreator,
		addressUpdater: addressUpdater,
	}
}

func (f *ppmShipmentUpdater) UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, mtoShipmentID uuid.UUID) (*models.PPMShipment, error) {
	return f.updatePPMShipment(appCtx, ppmShipment, mtoShipmentID, f.checks...)
}

func (f *ppmShipmentUpdater) updatePPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, mtoShipmentID uuid.UUID, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	if ppmShipment == nil {
		return nil, nil
	}

	oldPPMShipment, err := FindPPMShipmentByMTOID(appCtx, mtoShipmentID)
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
				// rejected := models.PPMAdvanceStatusRejected
				edited := models.PPMAdvanceStatusEdited
				if oldPPMShipment.HasRequestedAdvance != nil && updatedPPMShipment.HasRequestedAdvance != nil {
					if !*oldPPMShipment.HasRequestedAdvance && *updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &edited
					} else if *oldPPMShipment.HasRequestedAdvance && !*updatedPPMShipment.HasRequestedAdvance {
						// If a SC edits HasRequestedAdvance to be changed from true to false and the advanceRequest
						// was already previously approved this will ensure advanceStatus is reset to null
						updatedPPMShipment.AdvanceStatus = nil
					}
				}
				if oldPPMShipment.AdvanceAmountRequested != nil && updatedPPMShipment.AdvanceAmountRequested != nil {
					if *oldPPMShipment.AdvanceAmountRequested != *updatedPPMShipment.AdvanceAmountRequested {
						updatedPPMShipment.AdvanceStatus = &edited
					}
				}
			}
		}

		finalIncentive, err := f.estimator.FinalIncentiveWithDefaultChecks(appCtx, *oldPPMShipment, updatedPPMShipment)
		if err != nil {
			return err
		}
		updatedPPMShipment.FinalIncentive = finalIncentive

		if updatedPPMShipment.W2Address != nil {
			var updatedAddress *models.Address
			var error error
			if updatedPPMShipment.W2Address.ID.IsNil() {
				updatedAddress, error = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.W2Address)
			} else {
				updatedAddress, error = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.W2Address, etag.GenerateEtag(oldPPMShipment.W2Address.UpdatedAt))
			}
			if error != nil {
				return error
			}
			updatedPPMShipment.W2AddressID = &updatedAddress.ID
			updatedPPMShipment.W2Address = updatedAddress
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
