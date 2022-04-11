package ppmshipment

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// estimatePPM Struct
type estimatePPM struct {
	checks []ppmShipmentValidator
}

// NewEstimatePPM returns the estimatePPM (pass in checkRequiredFields() and checkEstimatedWeight)
func NewEstimatePPM() services.PPMEstimator {
	return &estimatePPM{
		checks: []ppmShipmentValidator{
			checkRequiredFields(),
			checkEstimatedWeight(),
		},
	}
}

// EstimateIncentiveWithDefaultChecks func that returns the estimate hard coded to 12K (because it'll be clear that the value is coming from teh service)
func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*int32, error) {
	return f.estimateIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*int32, error) {
	// Check that the PPMShipment has an ID
	var err error
	if newPPMShipment.ID.IsNil() {
		return nil, err
	}

	// Use Case: Calculating the final incentive in the closeout flow, what is the status then?
	// Might not need to block on teh PPM status.
	if newPPMShipment.Status != models.PPMShipmentStatusDraft {
		return nil, err
	}
	// 2. Check that all the fields we need are present.
	err = validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	if newPPMShipment.ExpectedDepartureDate == oldPPMShipment.ExpectedDepartureDate && newPPMShipment.PickupPostalCode == oldPPMShipment.PickupPostalCode && newPPMShipment.DestinationPostalCode == oldPPMShipment.DestinationPostalCode && newPPMShipment.AdvanceRequested == oldPPMShipment.AdvanceRequested && newPPMShipment.EstimatedWeight == oldPPMShipment.EstimatedWeight {
		return oldPPMShipment.EstimatedIncentive, nil
	}
	newPPMShipment.AdvanceRequested = nil
	newPPMShipment.Advance = nil
	// TODO: Call the pricer to calculate the incentive
	newPPMShipment.EstimatedIncentive = models.Int32Pointer(int32(1000000))
	// Saving Esimtated INcentive to DB
	verrs, err := appCtx.DB().ValidateAndSave(newPPMShipment)
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(newPPMShipment.ID, err, verrs, "Invalid input found while creating the Estimated Incentive.")
	} else if err != nil {
		return nil, apperror.NewQueryError("EstimatedIncentive", err, "")
	}

	return newPPMShipment.EstimatedIncentive, nil
}
