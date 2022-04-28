package ppmshipment

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type createShipmentSubtestData struct {
	move               models.Move
	newPPMShipment     *models.PPMShipment
	ppmShipmentCreator services.PPMShipmentCreator
}

func (suite *PPMShipmentSuite) createSubtestData() (subtestData *createShipmentSubtestData) {
	// Create new move
	subtestData = &createShipmentSubtestData{}

	subtestData.ppmShipmentCreator = NewPPMShipmentCreator()

	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())

	mtoShipment := testdatagen.MakeBaseMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			MoveTaskOrderID: subtestData.move.ID,
			ShipmentType:    models.MTOShipmentTypePPM,
			Status:          models.MTOShipmentStatusDraft,
		},
	})

	// Create a valid ppm shipment associated with a move
	subtestData.newPPMShipment = &models.PPMShipment{
		ShipmentID: mtoShipment.ID,
		Shipment:   mtoShipment,
	}

	return subtestData
}

func (suite *PPMShipmentSuite) TestPPMShipmentCreator() {

	suite.T().Run("CreatePPMShipment - Success", func(t *testing.T) {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created

		// Set required fields to their pointer values:
		subtestData := suite.createSubtestData()
		subtestData.newPPMShipment.ExpectedDepartureDate = time.Now()
		subtestData.newPPMShipment.PickupPostalCode = "90909"
		subtestData.newPPMShipment.DestinationPostalCode = "90905"
		subtestData.newPPMShipment.SitExpected = models.BoolPointer(false)

		createdPPMShipment, err := subtestData.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
	})

	// InvalidInputError
	suite.T().Run("Returns an InvalidInputError if MTOShipment type is not PPM", func(t *testing.T) {
		subtestData := suite.createSubtestData()

		subtestData.newPPMShipment.Shipment.ShipmentType = models.MTOShipmentTypeHHG

		createdPPMShipment, err := subtestData.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), subtestData.newPPMShipment)

		suite.Error(err)
		suite.Nil(createdPPMShipment)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("MTO shipment type must be PPM shipment", err.Error())
	})
}
