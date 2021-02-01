package movetaskorder_test

import (
	"testing"
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
	searchParams := services.FetchMoveTaskOrderParams{
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

	suite.T().Run("implicitly non-hidden move task orders", func(t *testing.T) {
		searchParams := services.ListMoveTaskOrderParams{} // should default to IncludeHidden being false
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
		suite.NotNil(expectedMTO.Locator)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)
	})

	suite.T().Run("include hidden move task orders", func(t *testing.T) {
		searchParams := services.ListMoveTaskOrderParams{
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

	suite.T().Run("default search - excludes hidden move task orders", func(t *testing.T) {
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
	hide := false
	hiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Show:               &hide,
		},
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	suite.T().Run("all move task orders", func(t *testing.T) {
		testdatagen.MakeDefaultMove(suite.DB())
		testdatagen.MakeDefaultMove(suite.DB())
		testdatagen.MakeDefaultMove(suite.DB())

		searchParams := services.ListMoveTaskOrderParams{
			IsAvailableToPrime: false,
			IncludeHidden:      true,
			Since:              nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
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
		suite.Equal(len(moveTaskOrders), 4)
	})

	suite.T().Run("default search - excludes hidden move task orders", func(t *testing.T) {
		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moveTaskOrders {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		suite.Equal(len(moveTaskOrders), 3) // minus the one hidden MTO
	})

	suite.T().Run("all move task orders that are available to prime and using since", func(t *testing.T) {
		now := time.Now()

		testdatagen.MakeAvailableMove(suite.DB())
		testdatagen.MakeAvailableMove(suite.DB())
		oldMTO := testdatagen.MakeAvailableMove(suite.DB())
		testdatagen.MakeDefaultMove(suite.DB())
		testdatagen.MakeDefaultMove(suite.DB())

		searchParams := services.ListMoveTaskOrderParams{
			IsAvailableToPrime: true,
			// IncludeHidden should be false by default
			Since: nil,
		}

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(&searchParams)
		suite.NoError(err)
		suite.Equal(len(moveTaskOrders), 3)

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
		suite.Equal(len(mtosWithSince), 2)
	})
}
