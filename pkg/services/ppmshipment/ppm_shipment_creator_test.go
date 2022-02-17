package ppmshipment

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"

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

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moverouter.NewMoveRouter()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	subtestData.ppmShipmentCreator = NewPPMShipmentCreator(mtoShipmentCreator)

	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())

	// Create a valid ppm shipment associated with a move
	subtestData.newPPMShipment = &models.PPMShipment{
		Shipment: models.MTOShipment{
			MoveTaskOrderID: subtestData.move.ID,
		},
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
		subtestData.newPPMShipment.ExpectedDepartureDate = models.TimePointer(time.Now())
		subtestData.newPPMShipment.PickupPostalCode = models.StringPointer("90909")
		subtestData.newPPMShipment.DestinationPostalCode = models.StringPointer("90905")
		subtestData.newPPMShipment.SitExpected = models.BoolPointer(false)

		createdPPMShipment, err := subtestData.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
	})

	// InvalidInputError
	suite.T().Run("A PPM shipment with validation errors returns an InvalidInputError with a bad UUID", func(t *testing.T) {
		blankPPMShipment := models.PPMShipment{}
		subtestData := suite.createSubtestData()
		createdPPMShipment, err := subtestData.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &blankPPMShipment)

		suite.Error(err)
		suite.Nil(createdPPMShipment)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
