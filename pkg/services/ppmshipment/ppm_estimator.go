package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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
			checkEstimatedWeight(), // add this function to validation.go
		},
	}
}

// EstimateIncentiveWithDefaultChecks func that returns the estimate hard coded to 12K (because it'll be clear that the value is coming from teh service)
func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error) {
	return f.estimateIncentive(appCtx, oldppmshipment, newppmshipment, f.checks...)
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, oldppmshipment models.PPMShipment, newppmshipment *models.PPMShipment, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	// Check that the PPMShipment has an ID
	var err error
	if newppmshipment.ID == uuid.Nil {
		return nil, err
	}
	// PPM Shipment Created:
	// 1.  Check that shipment status is in DRAFT or SUBMITTED
	// 2. Check that all the fields we need are present.
	// 3. Move on to estimate incentive

	// PPM Shipment Updated:
	// 1. Check that shipment status is in DRAFT or SUBMITTED
	// 2. Check that all or some required fields are not present OR have not changed.
	// 3. Skip the estimator and return log message with reason
	// 4. Check that all the required fields are present AND have been updated to a value
	// 5. estimate the incentive again and set old advance to nil, and advance_requested to false, if they existed

	// Check that shipment status is in DRAFT or SUBMITTED
	// Check that all the fields we need are present.
	// Counselor flow: Check if there is no change to the required fields if status is also SUBMITTED
	// move on to estimate incentive
	// Update the value of the Estimated Incentive field
	// update to the hard coded incentive amount
	// Will not do the following:
	// TODO: Call the pricer to calculate the incentive
	// return the entire ppmShipment
	// return an error

	return nil, err
}
