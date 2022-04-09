package ppmshipment

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func createNewDefaultPPMShipment() *models.PPMShipment {
	estimatedWeight := unit.Pound(5000)
	ppmShipment := models.PPMShipment{
		PickupPostalCode:      "20636",
		DestinationPostalCode: "94040",
		EstimatedWeight:       &estimatedWeight,
	}
	return &ppmShipment
}

func (suite *PPMShipmentSuite) TestEstimatedIncentive() {
	suite.Run("Estimated Incentive - Success", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()
		//ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)
		// Compare that the required fields have changed
		// Check that we have an estimated incentive
		newPPM := createNewDefaultPPMShipment()
		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, newPPM)
		//updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, oldPPMShipment.ShipmentID)
		// Might need to check that the estimated incentive is saved to the DB. LOOK AT EXAMPLES in other tests
		suite.NilOrNoVerrs(err)
		suite.NotEqualValues(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.NotEqualValues(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(int32(100000), ppmEstimate)
	})
	// Add unhappy path tests
}
