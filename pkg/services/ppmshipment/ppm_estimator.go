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
	//oldPPMShipment, err := models.FetchPPMShipmentFromMTOShipmentID(appCtx.DB(), mtoShipmentID)
	//if err != nil {
	//	return nil, err
	//}
	// if etag.GenerateEtag(oldPPMShipment.UpdatedAt) != eTag {
	// 	return nil, apperror.NewPreconditionFailedError(ppmShipment.ID, nil)
	// }

	// PPM Shipment Created:
	// 1.  Check that shipment status is in DRAFT
	// Use Case: Calculating the final incentive in the closeout flow, what is the status then?
	// Might not need to block on teh PPM status.
	if newPPMShipment.Status != models.PPMShipmentStatusDraft {
		return nil, err
	}
	// 2. Check that all the fields we need are present.
	err = validatePPMShipment(appCtx, *newPPMShipment, nil, &newPPMShipment.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	//if transactionError != nil {
	//	return nil, transactionError
	//}

	// PPM Shipment Updated:
	// 2. Check that all or some required fields are not present OR have not changed.
	// 3. Skip the estimator and return log message with reason
	// 4. Check that all the required fields are present AND have been updated to a value
	//verrs := validate.NewErrors()
	// 3. Move on to estimate incentive
	if newPPMShipment.ExpectedDepartureDate == oldPPMShipment.ExpectedDepartureDate && newPPMShipment.PickupPostalCode == oldPPMShipment.PickupPostalCode && newPPMShipment.DestinationPostalCode == oldPPMShipment.DestinationPostalCode && newPPMShipment.AdvanceRequested == oldPPMShipment.AdvanceRequested && newPPMShipment.EstimatedWeight == oldPPMShipment.EstimatedWeight {
		return oldPPMShipment.EstimatedIncentive, nil
	}
	newPPMShipment.AdvanceRequested = nil
	newPPMShipment.Advance = nil
	// TODO: Call the pricer to calculate the incentive
	newPPMShipment.EstimatedIncentive = models.Int32Pointer(int32(1000000))
	verrs, err := appCtx.DB().ValidateAndSave(newPPMShipment)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(newPPMShipment.ID, err, verrs, "Invalid input found while creating the Estimated Incentive.")
	} else if err != nil {
		return nil, apperror.NewQueryError("EstimatedIncentive", err, "")
	}
	return newPPMShipment.EstimatedIncentive, nil
	// 4. Check that all the required fields are present AND have been updated to a value
	// 5. estimate the incentive again and set old advance to nil, and advance_requested to false, if they existed:
}
