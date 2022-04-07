package ppmshipment

//import (
//	"github.com/transcom/mymove/pkg/apperror"
//	"github.com/transcom/mymove/pkg/models"
//	"github.com/transcom/mymove/pkg/testdatagen"
//	"github.com/transcom/mymove/pkg/unit"
//)
//
//func createDefaultPPMShipment() *models.PPMShipment {
//	ppmShipment := models.PPMShipment{
//		PickupPostalCode:      "20636",
//		DestinationPostalCode: "94040",
//	}
//	return &ppmShipment
//}

//func (suite *PPMShipmentSuite) TestEstimatedIncentive() {
//	suite.Run("Estimated Incentive - Success", func() {
//		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
//		ppmEstimator := NewEstimatePPM()
//		ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)
//
//		newPPM := createDefaultPPMShipment()
//		ppmEstimate := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, newPPM)
//		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, oldPPMShipment.ShipmentID)
//		// Might need to check that the estimated incentive is saved to the DB. LOOK AT EXAMPLES in other tests
//		suite.NilOrNoVerrs(err)
//		suite.Equal(int32(100000), ppmEstimator)
//	})
//}
