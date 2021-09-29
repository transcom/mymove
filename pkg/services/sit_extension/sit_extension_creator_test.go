package sitextension

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SitExtensionServiceSuite) TestSITExtensionCreator() {
	// Create new mtoShipment
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &move, testdatagen.Assertions{})

	// Create a valid SIT Extension for the move
	sit := &models.SITExtension{
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		MTOShipmentID: shipment.ID,
		RequestedDays: 10,
	}

	appCtx := appcontext.NewAppContext(suite.DB(), suite.logger)

	// Create move router for SitExtension Createor
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionCreator := NewSitExtensionCreator(moveRouter)
	movefetcher := movefetcher.NewMoveTaskOrderFetcher()

	suite.T().Run("Success - CreateSITExtension with no status passed in", func(t *testing.T) {
		// Under test:	CreateSITExtension
		// Set up:		Established valid shipment and valid SIT extension
		// Expected:	New reweigh successfully created
		createdSITExtension, sitErr := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		// Retrieve updated move
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: move.ID,
		}
		updatedMove, moveErr := movefetcher.FetchMoveTaskOrder(appCtx, &searchParams)

		suite.Nil(sitErr)
		suite.Nil(moveErr)
		suite.NotNil(createdSITExtension)
		suite.Equal(models.SITExtensionStatusPending, createdSITExtension.Status)
		suite.Equal(sit.RequestedDays, createdSITExtension.RequestedDays)
		suite.Equal(sit.RequestReason, createdSITExtension.RequestReason)
		suite.Equal(shipment.ID, createdSITExtension.MTOShipmentID)
		suite.Equal(move.ID, updatedMove.ID)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	// InvalidInputError
	suite.T().Run("Failure - SIT Extension with validation errors returns an InvalidInputError", func(t *testing.T) {
		badRequestReason := models.SITExtensionRequestReason("none")
		sit.RequestReason = badRequestReason
		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		suite.Error(err)
		suite.Nil(createdSITExtension)
		suite.IsType(services.InvalidInputError{}, err)

		// Reset request reason to correct reason
		sit.RequestReason = models.SITExtensionRequestReasonAwaitingCompletionOfResidence
	})

	suite.T().Run("Failure - Not Found Error", func(t *testing.T) {
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		sit.MTOShipmentID = notFoundUUID
		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(appCtx, sit)

		suite.Nil(createdSITExtension)
		suite.IsType(services.NotFoundError{}, err)

		// Reset shipment id to correct id
		sit.MTOShipmentID = shipment.ID
	})

	suite.T().Run("Success - CreateSITExtension with status passed in ", func(t *testing.T) {
		// Create new mtoShipment
		move2 := testdatagen.MakeAvailableMove(suite.DB())
		shipment2 := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &move, testdatagen.Assertions{})

		// Create a valid SIT Extension for the move
		sit2 := &models.SITExtension{
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			MTOShipmentID: shipment2.ID,
			Status:        models.SITExtensionStatusApproved,
			RequestedDays: 10,
		}
		createdSITExtension, sitErr2 := sitExtensionCreator.CreateSITExtension(appCtx, sit2)

		// Retrieve updated move
		searchParams2 := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: move2.ID,
		}
		updatedMove, moveErr2 := movefetcher.FetchMoveTaskOrder(appCtx, &searchParams2)

		suite.Nil(sitErr2)
		suite.Nil(moveErr2)
		suite.NotNil(createdSITExtension)
		suite.Equal(shipment2.ID, createdSITExtension.MTOShipmentID)
		suite.Equal(models.SITExtensionStatusApproved, createdSITExtension.Status)
		suite.Equal(sit2.RequestedDays, createdSITExtension.RequestedDays)
		suite.Equal(sit2.RequestReason, createdSITExtension.RequestReason)
		suite.Equal(move2.ID, updatedMove.ID)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)
	})
}
