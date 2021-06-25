package movetaskorder_test

import (
	"time"

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
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden: false,
	}

	actualMTO, err := mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID, &searchParams)
	suite.NoError(err)

	suite.NotZero(expectedMTO.ID, actualMTO.ID)
	suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
	suite.NotZero(actualMTO.Orders)
	suite.NotNil(expectedMTO.ReferenceID)
	suite.NotNil(expectedMTO.Locator)
	suite.Nil(expectedMTO.AvailableToPrimeAt)
	suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)
}

func (suite *MoveTaskOrderServiceSuite) TestListMoveTaskOrdersFetcher() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	hide := false
	hiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
		Move: models.Move{
			Show: &hide,
		},
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	suite.RunWithRollback("implicitly non-hidden move task orders", func() {
		searchParams := services.MoveTaskOrderFetcherParams{} // should default to IncludeHidden being false
		moveTaskOrders, err := mtoFetcher.ListMoveTaskOrders(expectedOrder.ID, &searchParams)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		actualMTO := moveTaskOrders[0]

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(actualMTO.Locator)
		suite.NotNil(actualMTO.ReferenceID)
		suite.Nil(actualMTO.AvailableToPrimeAt)
		suite.NotEqual(actualMTO.Status, models.MoveStatusCANCELED)
	})

	suite.RunWithRollback("include hidden move task orders", func() {
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden: true,
		}
		moveTaskOrders, err := mtoFetcher.ListMoveTaskOrders(expectedOrder.ID, &searchParams)
		suite.NoError(err)

		// The hidden move be in this output list since we weren't excluding hidden MTOs:
		found := false
		for _, move := range moveTaskOrders {
			if move.ID == hiddenMTO.ID {
				found = true
				break
			}
		}
		suite.True(found)
		suite.Equal(len(moveTaskOrders), 2)
	})

	suite.RunWithRollback("default search - excludes hidden move task orders", func() {
		moveTaskOrders, err := mtoFetcher.ListMoveTaskOrders(expectedOrder.ID, nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}
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
		now := time.Now()

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
		since := now.Unix()
		searchParams.Since = &since
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
		suite.NoError(err)
		suite.Equal(1, len(mtosWithSince))
	})
}
