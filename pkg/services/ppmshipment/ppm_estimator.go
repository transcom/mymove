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

	if newPPMShipment.Status != models.PPMShipmentStatusDraft {
		return nil, nil
	}
	// Check that all the required fields we need are present.
	err = validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	// If a field does not pass validation return nil as error handling is happening in the validator
	if err != nil {
		switch err.(type) {
		case apperror.InvalidInputError:
			return nil, nil
		default:
			return nil, err
		}
	}

	if oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate) && newPPMShipment.PickupPostalCode == oldPPMShipment.PickupPostalCode && newPPMShipment.DestinationPostalCode == oldPPMShipment.DestinationPostalCode && oldPPMShipment.EstimatedWeight != nil && *newPPMShipment.EstimatedWeight == *oldPPMShipment.EstimatedWeight {
		return oldPPMShipment.EstimatedIncentive, nil
	}
	// Clear out advance and advance requested fields when the estimated incentive is reset.
	newPPMShipment.AdvanceRequested = nil
	newPPMShipment.Advance = nil

	// TODO: Call the pricer to calculate the incentive

	estimatedIncentive := int32(1000000)
	return &estimatedIncentive, nil
}
