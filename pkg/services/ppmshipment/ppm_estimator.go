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
	return f.estimateIncentive(appCtx, ppmshipment, f.checks...)
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	// Check that the PPMShipment has an ID
	var err error
	if ppmshipment.ID == uuid.Nil {
		return nil, err
	}
	// Check if there is an update to any of the required fields
	// IF the required fields values have changed then store the new values on the model
	// Update the value of the Estimated Incentive field
	// update to the hard coded incentive amount
	// Will not do the following:
	// TODO: Call the pricer to calculate the incentive
	// return the entire ppmShipment
	// return an error

	return nil, err
}
