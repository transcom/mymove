package sitextension

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SitExtensionServiceSuite) TestSITExtensionCreator() {
	move := testdatagen.MakeAvailableMove(suite.DB())

	// Create move router for SitExtension Createor
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionCreator := NewSitExtensionCreator(moveRouter)
	movefetcher := movefetcher.NewMoveTaskOrderFetcher()

	suite.T().Run("Success - CreateSITExtension with no status passed in", func(t *testing.T) {
		// Under test:	CreateSITExtension
		// Set up:		Established valid shipment and valid SIT extension
		// Expected:	New sit successfully created
		// Create new mtoShipment
		shipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &move, testdatagen.Assertions{})

		// Create a valid SIT Extension for the move
		sit := &models.SITExtension{
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			MTOShipmentID: shipment.ID,
			RequestedDays: 10,
		}

		createdSITExtension, sitErr := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		// Retrieve updated move
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: move.ID,
		}
		updatedMove, moveErr := movefetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)

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
		shipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &move, testdatagen.Assertions{})

		// Create a SIT Extension for the move
		sit := &models.SITExtension{
			RequestReason: models.SITExtensionRequestReason("none"),
			MTOShipmentID: shipment.ID,
			RequestedDays: 10,
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Error(err)
		suite.Nil(createdSITExtension)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Failure - Not Found Error because shipment not found", func(t *testing.T) {
		// Create a SIT Extension for the move
		sit := &models.SITExtension{
			MTOShipmentID: uuid.Must(uuid.NewV4()),
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Nil(createdSITExtension)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Failure - Not Found Error because shipment uses external vendor", func(t *testing.T) {
		externalShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})

		// Create a SIT Extension for the move
		sit := &models.SITExtension{
			MTOShipmentID: externalShipment.ID,
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Nil(createdSITExtension)
		suite.IsType(apperror.NotFoundError{}, err)
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
		createdSITExtension, sitErr2 := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit2)

		// Retrieve updated move
		searchParams2 := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: move2.ID,
		}
		updatedMove, moveErr2 := movefetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams2)

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
