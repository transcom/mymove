package ppmshipment

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMShipmentSuite) TestPPMShipmentCreator() {
	// Create new mtoShipment
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

	// Create a valid ppm shipment associated with a move
	newPPMShipment := &models.PPMShipment{
		CreatedAt:  time.Now(),
		ShipmentID: mtoShipment.ID,
	}

	suite.T().Run("CreatePPMShipment - Success", func(t *testing.T) {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created
		ppmShipmentCreator := NewPPMShipmentCreator()
		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
		suite.Equal(mtoShipment.ID, createdPPMShipment.ShipmentID)

	})

	// InvalidInputError
	suite.T().Run("A PPM shipment with validation errors returns an InvalidInputError", func(t *testing.T) {
		badCreatedAt := models.PPMShipment{CreatedAt: time.Time{}} // createdAt is empty because there is no PPM shipment
		newPPMShipment.CreatedAt = badCreatedAt
		ppmShipmentCreator := NewPPMShipmentCreator()
		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), newPPMShipment)

		suite.Error(err)
		suite.Nil(createdPPMShipment)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		newPPMShipment.ShipmentID = notFoundUUID
		ppmShipmentCreator := NewPPMShipmentCreator()
		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), newPPMShipment)

		suite.Nil(createdPPMShipment)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
