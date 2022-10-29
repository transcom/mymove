package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchPPMShipment() {
	t := suite.T()

	ppm := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

	_, err := suite.DB().ValidateAndSave(&ppm)
	if err != nil {
		t.Errorf("could not save PPM: %v", err)
		return
	}

	retrievedPPM, _ := models.FetchPPMShipmentFromMTOShipmentID(suite.DB(), ppm.ShipmentID)

	suite.Equal(retrievedPPM.ID, ppm.ID)
	suite.Equal(retrievedPPM.ShipmentID, ppm.ShipmentID)
}

func (suite *ModelSuite) TestGetPPMShipment() {
	suite.Run("Can find an existing PPM shipment and loads associations", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentThatNeedsCloseOut(appCtx.DB(), testdatagen.Assertions{})

		assertions := testdatagen.Assertions{
			PPMShipment: existingPPMShipment,
		}

		testdatagen.MakeMovingExpense(appCtx.DB(), assertions)

		testdatagen.MakeProgearWeightTicket(appCtx.DB(), assertions)

		ppmShipment, err := models.GetPPMShipment(appCtx, existingPPMShipment.ID)

		if suite.NoError(err) {
			suite.Equal(existingPPMShipment.ID, ppmShipment.ID)

			suite.NotNil(ppmShipment.Shipment)
			suite.NotNil(ppmShipment.Shipment.MoveTaskOrder.ID)
			suite.True(len(ppmShipment.WeightTickets) > 0, "Expected weight tickets to be loaded")
			suite.True(len(ppmShipment.MovingExpenses) > 0, "Expected moving expenses to be loaded")
			suite.True(len(ppmShipment.ProgearExpenses) > 0, "Expected progear weight tickets to be loaded")
			suite.NotNil(ppmShipment.SignedCertification)
		}
	})

	suite.Run("Returns an error if PPM shipment does not exist", func() {
		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		_, err := models.GetPPMShipment(suite.AppContextForTest(), nonexistentPPMShipmentID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "not found while looking for PPMShipment")
		}
	})
}
