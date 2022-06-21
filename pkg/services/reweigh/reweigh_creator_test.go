package reweigh

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReweighSuite) TestReweighCreator() {
	// Create new mtoShipment

	suite.Run("CreateReweigh - Success", func() {
		// Under test:	CreateReweigh
		// Set up:		Established valid shipment and valid reweigh
		// Expected:	New reweigh successfully created
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

		// Create a valid reweigh for the move
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  mtoShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Nil(err)
		suite.NotNil(createdReweigh)
		suite.Equal(mtoShipment.ID, createdReweigh.ShipmentID)

	})

	// InvalidInputError
	suite.Run("Reweigh with validation errors returns an InvalidInputError", func() {
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

		// Create a reweigh with a bad requester
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequester("not requested by anyone"),
			ShipmentID:  mtoShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		createReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Error(err)
		suite.Nil(createReweigh)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Not Found Error", func() {

		// Create a reweigh with a shipment that doesn't exist
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  uuid.Must(uuid.NewV4()),
		}
		reweighCreator := NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Nil(createdReweigh)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
