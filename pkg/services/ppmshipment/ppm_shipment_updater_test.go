package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func createDefaultPPMShipment() *models.PPMShipment {
	ppmShipment := models.PPMShipment{
		PickupPostalCode: "20636",
		// SitExpected: true,
	}
	return &ppmShipment
}

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {
	suite.Run("UpdatePPMShipment - Success", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()
		ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)

		newPPM := createDefaultPPMShipment()
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, oldPPMShipment.ShipmentID)

		suite.NilOrNoVerrs(err)
		suite.Equal(newPPM.PickupPostalCode, updatedPPMShipment.PickupPostalCode)
		// suite.True(updatedPPMShipment.SitExpected)
		suite.Equal(unit.Pound(1150), *updatedPPMShipment.ProGearWeight)
	})

	suite.Run("Not Found Error", func() {
		ppmEstimator := NewEstimatePPM()
		ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)

		newPPM := createDefaultPPMShipment()
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), newPPM, uuid.Nil)

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
