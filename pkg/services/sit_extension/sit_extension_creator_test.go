package sitextension

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SitExtensionServiceSuite) TestSITExtensionCreator() {
	// Create new mtoShipment
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

	// Create a valid SIT Extension for the move
	sit := &models.SITExtension{
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		Status:        models.SITExtensionStatusApproved,
		MTOShipmentID: mtoShipment.ID,
		RequestedDays: 10,
	}

	appCtx := appcontext.NewAppContext(suite.DB(), suite.logger)

	suite.T().Run("CreateSITExtension - Success", func(t *testing.T) {
		// Under test:	CreateSITExtension
		// Set up:		Established valid shipment and valid SIT extension
		// Expected:	New reweigh successfully created
		sitExtensionCreator := NewSitExtensionCreator()
		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		suite.Nil(err)
		suite.NotNil(createdSITExtension)
		suite.Equal(mtoShipment.ID, createdSITExtension.MTOShipmentID)

	})

	// InvalidInputError
	suite.T().Run("SIT Extension with validation errors returns an InvalidInputError", func(t *testing.T) {
		badRequestReason := models.SITExtensionRequestReason("none")
		sit.RequestReason = badRequestReason
		sitExtensionCreator := NewSitExtensionCreator()
		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		suite.Error(err)
		suite.Nil(createdSITExtension)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		sit.MTOShipmentID = notFoundUUID
		sitExtensionCreator := NewSitExtensionCreator()
		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		suite.Nil(createdSITExtension)
		suite.IsType(services.NotFoundError{}, err)
	})
}
