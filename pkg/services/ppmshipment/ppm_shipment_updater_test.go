package ppmshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {
	ppmShipmentUpdater := NewPPMShipmentUpdater()

	oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

	newPPM := models.PPMShipment{
		PickupPostalCode: "20636",
		// SitExpected: true,
	}

	// eTag := etag.GenerateEtag(oldPPMShipment.UpdatedAt)
	suite.T().Run("UpdatePPMShipment - Success", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, oldPPMShipment.ShipmentID)

		suite.NilOrNoVerrs(err)
		suite.Equal(newPPM.PickupPostalCode, updatedPPMShipment.PickupPostalCode)
		// suite.True(updatedPPMShipment.SitExpected)
		suite.Equal(unit.Pound(1150), *updatedPPMShipment.ProGearWeight)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, uuid.Nil)

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	// suite.T().Run("Precondition Failed", func(t *testing.T) {
	// 	updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, oldPPMShipment.ShipmentID)

	// 	suite.Nil(updatedPPMShipment)
	// 	suite.Error(err)
	// 	suite.IsType(apperror.PreconditionFailedError{}, err)
	// })
}
