package ppmshipment

import (
	"testing"

	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type createShipmentSubtestData struct {
	move           models.Move
	newPPMShipment *models.PPMShipment
}

func (suite *PPMShipmentSuite) createSubtestData(assertions testdatagen.Assertions) (subtestData *createShipmentSubtestData) {
	// Create new move
	subtestData = &createShipmentSubtestData{}

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
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moverouter.NewMoveRouter()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	suite.T().Run("CreatePPMShipment - Success", func(t *testing.T) {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		ppmShipmentCreator := NewPPMShipmentCreator(mtoShipmentCreator)
		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
	})

	// InvalidInputError
	suite.T().Run("A PPM shipment with validation errors returns an InvalidInputError with a bad UUID", func(t *testing.T) {
		//badCreatedAt := models.PPMShipment{CreatedAt: time.Time{}} // createdAt is empty because there is no PPM shipment
		//newPPMShipment.CreatedAt = badCreatedAt
		blankPPMShipment := models.PPMShipment{}
		ppmShipmentCreator := NewPPMShipmentCreator(mtoShipmentCreator)
		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), &blankPPMShipment)

		suite.Error(err)
		suite.Nil(createdPPMShipment)
		suite.IsType(apperror.QueryError{}, err)
	})
}
