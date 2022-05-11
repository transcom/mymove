package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	prhelpermocks "github.com/transcom/mymove/pkg/payment_request/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func createDefaultPPMShipment() *models.PPMShipment {
	ppmShipment := models.PPMShipment{
		PickupPostalCode: "20636",
		// SitExpected: true,
	}
	return &ppmShipment
}

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {
	ppmEstimator := NewEstimatePPM(&mocks.Planner{}, &prhelpermocks.Helper{})

	suite.Run("UpdatePPMShipment - Success", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)

		newPPM := createDefaultPPMShipment()
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, oldPPMShipment.ShipmentID)

		suite.NilOrNoVerrs(err)
		suite.Equal(newPPM.PickupPostalCode, updatedPPMShipment.PickupPostalCode)
		suite.Equal(*oldPPMShipment.ProGearWeight, *updatedPPMShipment.ProGearWeight)
	})

	suite.Run("Not Found Error", func() {
		ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)

		newPPM := createDefaultPPMShipment()
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, uuid.Nil)

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
