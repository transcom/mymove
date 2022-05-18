package movetaskorder_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	setupTestData := func() (models.Move, models.MTOShipment) {
		expectedMTO := testdatagen.MakeDefaultMove(suite.DB())

		// Make a couple of shipments for the move; one prime, one external
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: expectedMTO,
			MTOShipment: models.MTOShipment{
				UsesExternalVendor: false,
			},
		})
		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: expectedMTO,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})
		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: expectedMTO,
			MTOShipment: models.MTOShipment{
				DeletedAt: models.TimePointer(time.Now()),
			},
		})

		return expectedMTO, primeShipment
	}

	mtoFetcher := NewMoveTaskOrderFetcher()

	suite.T().Run("Success with Prime-available move by ID, fetch all non-deleted shipments", func(t *testing.T) {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get two shipments back since we didn't set searchParams to exclude external ones.
		suite.Len(actualMTO.MTOShipments, 2)
	})

	suite.T().Run("Success with Prime-available move by Locator, no deleted or external shipments", func(t *testing.T) {
		expectedMTO, primeShipment := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			Locator:                  expectedMTO.Locator,
			ExcludeExternalShipments: true,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get one shipment back since we requested no external shipments.
		if suite.Len(actualMTO.MTOShipments, 1) {
			suite.Equal(expectedMTO.ID.String(), actualMTO.ID.String())
			suite.Equal(primeShipment.ID.String(), actualMTO.MTOShipments[0].ID.String())
		}
	})

	suite.T().Run("Success with move that has only deleted shipments", func(t *testing.T) {
		mtoWithAllShipmentsDelted := testdatagen.MakeDefaultMove(suite.DB())
		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: mtoWithAllShipmentsDelted,
			MTOShipment: models.MTOShipment{
				DeletedAt: models.TimePointer(time.Now()),
			},
		})

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: mtoWithAllShipmentsDelted.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Equal(mtoWithAllShipmentsDelted.ID, actualMTO.ID)
		suite.Len(actualMTO.MTOShipments, 0)
	})

	suite.T().Run("Failure - nil searchParams", func(t *testing.T) {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), nil)
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should not be nil")
	})

	suite.T().Run("Failure - searchParams with no ID/locator set", func(t *testing.T) {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &services.MoveTaskOrderFetcherParams{})
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should have either a move ID or locator set")
	})

	suite.T().Run("Failure - Not Found with Bad ID", func(t *testing.T) {
		badID, _ := uuid.NewV4()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: badID,
		}

		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.Error(err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListAllMoveTaskOrdersFetcher() {
	setupTestData := func() (models.Move, models.Move, models.MTOShipment) {
		// Set up a hidden move so we can check if it's in the output:
		now := time.Now()
		show := false
		hiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Show:               &show,
			},
		})

		mto := testdatagen.MakeDefaultMove(suite.DB())

		// Make a couple of shipments for the default move; one prime, one external
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				UsesExternalVendor: false,
			},
		})
		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})
		return mto, hiddenMTO, primeShipment
	}

	mtoFetcher := NewMoveTaskOrderFetcher()

	suite.Run("all move task orders", func() {
		setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: false,
			IncludeHidden:      true,
			Since:              nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		move := moveTaskOrders[0]

		suite.Equal(2, len(moveTaskOrders))
		suite.NotNil(move.Orders)
		suite.NotNil(move.Orders.OriginDutyLocation)
		suite.NotNil(move.Orders.NewDutyLocation)
	})

	suite.Run("default search - excludes hidden move task orders", func() {
		mto, hiddenMTO, _ := setupTestData()
		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		if suite.Equal(1, len(moveTaskOrders)) { // minus the one hidden MTO
			// We should get back all shipments on that move since we did not exclude external vendor shipments.
			suite.Equal(mto.ID.String(), moveTaskOrders[0].ID.String())
			suite.Len(moveTaskOrders[0].MTOShipments, 2)
		}
	})

	suite.Run("returns shipments that respect the external searchParams flag", func() {
		mto, _, primeShipment := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			ExcludeExternalShipments: true,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// We should only get back the one move that's not hidden.
		if suite.Len(moveTaskOrders, 1) {
			if suite.Equal(mto.ID.String(), moveTaskOrders[0].ID.String()) {
				// That move should get one shipment back since we requested no external shipments.
				if suite.Len(moveTaskOrders[0].MTOShipments, 1) {
					suite.Equal(primeShipment.ID.String(), moveTaskOrders[0].MTOShipments[0].ID.String())
				}
			}
		}
	})

	suite.Run("all move task orders that are available to prime and using since", func() {
		_, hiddenMTO, _ := setupTestData()
		now := time.Now()

		testdatagen.MakeAvailableMove(suite.DB())
		oldMTO := testdatagen.MakeAvailableMove(suite.DB())

		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: true,
			// IncludeHidden should be false by default
			Since: nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Equal(2, len(moveTaskOrders))

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		// Put 1 Move updatedAt in the past
		suite.NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
			now.Add(-2*time.Second), oldMTO.ID).Exec())
		searchParams.Since = &now
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Equal(1, len(mtosWithSince))
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListPrimeMoveTaskOrdersFetcher() {
	// Set up a hidden move so we can check if it's in the output:
	now := time.Now()
	show := false
	hiddenMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Show:               &show,
		},
	})
	// Make a default, not Prime-available move:
	nonPrimeMove := testdatagen.MakeDefaultMove(suite.DB())
	// Make some Prime moves:
	primeMove1 := testdatagen.MakeAvailableMove(suite.DB())
	primeMove2 := testdatagen.MakeAvailableMove(suite.DB())
	primeMove3 := testdatagen.MakeAvailableMove(suite.DB())
	testdatagen.MakeMTOShipmentWithMove(suite.DB(), &primeMove3, testdatagen.Assertions{})

	// Move primeMove1 and primeMove3 into the past so we can exclude them:
	suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=$1 WHERE id IN ($2, $3);",
		now.Add(-10*time.Second), primeMove1.ID, primeMove3.ID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id=$2;",
		now.Add(-10*time.Second), primeMove1.OrdersID).Exec())

	fetcher := NewMoveTaskOrderFetcher()
	searchParams := services.MoveTaskOrderFetcherParams{}

	// Run the fetcher without `since` to get all Prime moves:
	primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(primeMoves, 3)

	moveIDs := []uuid.UUID{primeMoves[0].ID, primeMoves[1].ID, primeMoves[2].ID}
	suite.NotContains(moveIDs, hiddenMove.ID)
	suite.NotContains(moveIDs, nonPrimeMove.ID)
	suite.Contains(moveIDs, primeMove1.ID)
	suite.Contains(moveIDs, primeMove2.ID)
	suite.Contains(moveIDs, primeMove3.ID)

	// Run the fetcher with `since` to get primeMove2 and primeMove3 (because of the shipment)
	since := now.Add(-5 * time.Second)
	searchParams.Since = &since
	sinceMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(sinceMoves, 2)

	sinceMoveIDs := []uuid.UUID{sinceMoves[0].ID, sinceMoves[1].ID}
	suite.Contains(sinceMoveIDs, primeMove2.ID)
	suite.Contains(sinceMoveIDs, primeMove3.ID)
}
