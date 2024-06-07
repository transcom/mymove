package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
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

func (f *ppmShipmentUpdater) UpdatePPMShipmentSITEstimatedCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMShipment, error) {
	if ppmShipment == nil {
		return nil, nil
	}

	oldPPMShipment, err := FindPPMShipment(appCtx, ppmShipment.ID)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment, err := mergePPMShipment(*ppmShipment, oldPPMShipment)
	if err != nil {
		return nil, err
	}

	err = validatePPMShipment(appCtx, *updatedPPMShipment, oldPPMShipment, &oldPPMShipment.Shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		contractDate := ppmShipment.ExpectedDepartureDate
		contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
		if err != nil {
			return err
		}

		estimatedSITCost, err := CalculateSITCost(appCtx, updatedPPMShipment, contract)
		if err != nil {
			return err
		}

		updatedPPMShipment.SITEstimatedCost = estimatedSITCost

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

	isPrimeCounseled, err := IsPrimeCounseledPPM(appCtx, mtoShipmentID)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment, err := mergePPMShipment(*ppmShipment, oldPPMShipment)
	if err != nil {
		return nil, err
	}

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
				edited := models.PPMAdvanceStatusEdited
				if oldPPMShipment.HasRequestedAdvance != nil && updatedPPMShipment.HasRequestedAdvance != nil {
					if !*oldPPMShipment.HasRequestedAdvance && *updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &edited
					} else if *oldPPMShipment.HasRequestedAdvance && !*updatedPPMShipment.HasRequestedAdvance {
						updatedPPMShipment.AdvanceStatus = &edited
					}
				}
				if oldPPMShipment.AdvanceAmountRequested != nil && updatedPPMShipment.AdvanceAmountRequested != nil {
					if *oldPPMShipment.AdvanceAmountRequested != *updatedPPMShipment.AdvanceAmountRequested {
						updatedPPMShipment.AdvanceStatus = &edited
					}
				}
			}
			if appCtx.Session().IsMilApp() {
				if isPrimeCounseled && updatedPPMShipment.HasRequestedAdvance != nil {
					received := models.PPMAdvanceStatusReceived
					notReceived := models.PPMAdvanceStatusNotReceived

					if updatedPPMShipment.HasReceivedAdvance != nil && *updatedPPMShipment.HasRequestedAdvance {
						if *updatedPPMShipment.HasReceivedAdvance {
							updatedPPMShipment.AdvanceStatus = &received
						}
						if !*updatedPPMShipment.HasReceivedAdvance {
							updatedPPMShipment.AdvanceStatus = &notReceived
						}
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
			var createOrUpdateErr error
			if updatedPPMShipment.W2Address.ID.IsNil() {
				updatedAddress, createOrUpdateErr = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.W2Address)
			} else {
				updatedAddress, createOrUpdateErr = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.W2Address, etag.GenerateEtag(oldPPMShipment.W2Address.UpdatedAt))
			}
			if createOrUpdateErr != nil {
				return createOrUpdateErr
			}
			updatedPPMShipment.W2AddressID = &updatedAddress.ID
			updatedPPMShipment.W2Address = updatedAddress
		}

		if updatedPPMShipment.PickupAddress != nil {
			var updatedAddress *models.Address
			var createOrUpdateErr error
			if updatedPPMShipment.PickupAddress.ID.IsNil() {
				updatedAddress, createOrUpdateErr = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.PickupAddress)
			} else {
				updatedAddress, createOrUpdateErr = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.PickupAddress, etag.GenerateEtag(oldPPMShipment.PickupAddress.UpdatedAt))
			}
			if createOrUpdateErr != nil {
				return createOrUpdateErr
			}
			updatedPPMShipment.PickupAddressID = &updatedAddress.ID
			updatedPPMShipment.PickupAddress = updatedAddress
		}

		if updatedPPMShipment.SecondaryPickupAddress != nil {
			var updatedAddress *models.Address
			var createOrUpdateErr error
			if updatedPPMShipment.SecondaryPickupAddress.ID.IsNil() {
				updatedAddress, createOrUpdateErr = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.SecondaryPickupAddress)
			} else {
				updatedAddress, createOrUpdateErr = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.SecondaryPickupAddress, etag.GenerateEtag(oldPPMShipment.SecondaryPickupAddress.UpdatedAt))
			}
			if createOrUpdateErr != nil {
				return createOrUpdateErr
			}
			updatedPPMShipment.SecondaryPickupAddressID = &updatedAddress.ID
			updatedPPMShipment.SecondaryPickupAddress = updatedAddress
		}

		if updatedPPMShipment.DestinationAddress != nil {
			var updatedAddress *models.Address
			var createOrUpdateErr error
			if updatedPPMShipment.DestinationAddress.ID.IsNil() {
				updatedAddress, createOrUpdateErr = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.DestinationAddress)
			} else {
				updatedAddress, createOrUpdateErr = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.DestinationAddress, etag.GenerateEtag(oldPPMShipment.DestinationAddress.UpdatedAt))
			}
			if createOrUpdateErr != nil {
				return createOrUpdateErr
			}
			updatedPPMShipment.DestinationAddressID = &updatedAddress.ID
			updatedPPMShipment.DestinationAddress = updatedAddress
		}

		if updatedPPMShipment.SecondaryDestinationAddress != nil {
			var updatedAddress *models.Address
			var createOrUpdateErr error
			if updatedPPMShipment.SecondaryDestinationAddress.ID.IsNil() {
				updatedAddress, createOrUpdateErr = f.addressCreator.CreateAddress(txnAppCtx, updatedPPMShipment.SecondaryDestinationAddress)
			} else {
				updatedAddress, createOrUpdateErr = f.addressUpdater.UpdateAddress(txnAppCtx, updatedPPMShipment.SecondaryDestinationAddress, etag.GenerateEtag(oldPPMShipment.SecondaryDestinationAddress.UpdatedAt))
			}
			if createOrUpdateErr != nil {
				return createOrUpdateErr
			}
			updatedPPMShipment.SecondaryDestinationAddressID = &updatedAddress.ID
			updatedPPMShipment.SecondaryDestinationAddress = updatedAddress
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
