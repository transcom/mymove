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
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	suite.T().Run("Success with Prime-available move by ID", func(t *testing.T) {
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(&searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)
	})

	suite.T().Run("Success with Prime-available move by Locator", func(t *testing.T) {
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden: false,
			Locator:       expectedMTO.Locator,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(&searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)
	})

	suite.T().Run("Failure - Not Found with Bad ID", func(t *testing.T) {
		badID, _ := uuid.NewV4()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: badID,
		}

		_, err := mtoFetcher.FetchMoveTaskOrder(&searchParams)
		suite.Error(err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListAllMoveTaskOrdersFetcher() {
	// Set up a hidden move so we can check if it's in the output:
	now := time.Now()
	show := false
	hiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Show:               &show,
		},
	})
	testdatagen.MakeDefaultMove(suite.DB())

	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	suite.RunWithRollback("all move task orders", func() {
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: false,
			IncludeHidden:      true,
			Since:              nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
		suite.NoError(err)

		move := moveTaskOrders[0]

		suite.Equal(2, len(moveTaskOrders))
		suite.NotNil(move.Orders)
		suite.NotNil(move.Orders.OriginDutyStation)
		suite.NotNil(move.Orders.NewDutyStation)
	})

	suite.RunWithRollback("default search - excludes hidden move task orders", func() {
		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		suite.Equal(1, len(moveTaskOrders)) // minus the one hidden MTO
	})

	suite.RunWithRollback("all move task orders that are available to prime and using since", func() {
		now = time.Now()

		testdatagen.MakeAvailableMove(suite.DB())
		oldMTO := testdatagen.MakeAvailableMove(suite.DB())

		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: true,
			// IncludeHidden should be false by default
			Since: nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
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
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
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

	fetcher := NewMoveTaskOrderFetcher(suite.DB())
	searchParams := services.MoveTaskOrderFetcherParams{}

	// Run the fetcher without `since` to get all Prime moves:
	primeMoves, err := fetcher.ListPrimeMoveTaskOrders(&searchParams)
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
	sinceMoves, err := fetcher.ListPrimeMoveTaskOrders(&searchParams)
	suite.NoError(err)
	suite.Len(sinceMoves, 2)

	sinceMoveIDs := []uuid.UUID{sinceMoves[0].ID, sinceMoves[1].ID}
	suite.Contains(sinceMoveIDs, primeMove2.ID)
	suite.Contains(sinceMoveIDs, primeMove3.ID)
}
