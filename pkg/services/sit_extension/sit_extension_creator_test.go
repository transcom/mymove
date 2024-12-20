package sitextension

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
)

func (suite *SitExtensionServiceSuite) TestSITExtensionCreator() {

	// Create move router for SitExtension Creator
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionCreator := NewSitExtensionCreator(moveRouter)
	movefetcher := movefetcher.NewMoveTaskOrderFetcher()

	suite.Run("Success - CreateSITExtension with no status passed in", func() {
		// Under test:	CreateSITExtension
		// Set up:		Established valid shipment and valid SIT extension
		// Expected:	New sit successfully created
		// Create new mtoShipment
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// Create a valid SIT Extension for the move
		sit := &models.SITDurationUpdate{
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
	suite.Run("Failure - SIT Extension with validation errors returns an InvalidInputError", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// Create a SIT Extension for the move
		sit := &models.SITDurationUpdate{
			RequestReason: models.SITDurationUpdateRequestReason("none"),
			MTOShipmentID: shipment.ID,
			RequestedDays: 10,
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Error(err)
		suite.Nil(createdSITExtension)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Failure - Not Found Error because shipment not found", func() {
		// Create a SIT Extension for the move
		sit := &models.SITDurationUpdate{
			MTOShipmentID: uuid.Must(uuid.NewV4()),
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Nil(createdSITExtension)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Failure - Not Found Error because shipment uses external vendor", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		externalShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)

		// Create a SIT Extension for the move
		sit := &models.SITDurationUpdate{
			MTOShipmentID: externalShipment.ID,
		}

		createdSITExtension, err := sitExtensionCreator.CreateSITExtension(suite.AppContextForTest(), sit)

		suite.Nil(createdSITExtension)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Success - CreateSITExtension with status passed in ", func() {
		// Create new mtoShipment
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		move2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment2 := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// Create a valid SIT Extension for the move
		sit2 := &models.SITDurationUpdate{
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
