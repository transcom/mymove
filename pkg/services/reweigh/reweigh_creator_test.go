package reweigh

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReweighSuite) TestReweighCreator() {
	// Create new mtoShipment
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

	// Create a valid reweigh for the move
	newReweigh := &models.Reweigh{
		RequestedAt: time.Now(),
		RequestedBy: models.ReweighRequesterPrime,
		ShipmentID:  mtoShipment.ID,
	}

	appCtx := appcontext.NewAppContext(suite.DB(), suite.logger)

	suite.T().Run("CreateReweigh - Success", func(t *testing.T) {
		// Under test:	CreateReweigh
		// Set up:		Established valid shipment and valid reweigh
		// Expected:	New reweigh successfully created
		reweighCreator := NewReweighCreator(suite.DB())
		createdReweigh, err := reweighCreator.CreateReweighCheck(appCtx, newReweigh)

		suite.Nil(err)
		suite.NotNil(createdReweigh)
		suite.Equal(mtoShipment.ID, createdReweigh.ShipmentID)

	})

	// InvalidInputError
	suite.T().Run("Reweigh with validation errors returns an InvalidInputError", func(t *testing.T) {
		badRequestedby := models.ReweighRequester("not requested by anyone")
		newReweigh.RequestedBy = badRequestedby
		reweighCreator := NewReweighCreator(suite.DB())
		createReweigh, err := reweighCreator.CreateReweighCheck(appCtx, newReweigh)

		suite.Error(err)
		suite.Nil(createReweigh)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		newReweigh.ShipmentID = notFoundUUID
		reweighCreator := NewReweighCreator(suite.DB())
		createdReweigh, err := reweighCreator.CreateReweighCheck(appCtx, newReweigh)

		suite.Nil(createdReweigh)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
