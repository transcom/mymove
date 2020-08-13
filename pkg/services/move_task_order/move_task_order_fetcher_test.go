package movetaskorder_test

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	actualMTO, err := mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID)
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
	expectedMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	moveTaskOrders, err := mtoFetcher.ListMoveTaskOrders(expectedOrder.ID)
	suite.NoError(err)

	actualMTO := moveTaskOrders[0]

	suite.NotZero(expectedMTO.ID, actualMTO.ID)
	suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
	suite.NotZero(actualMTO.Orders)
	suite.NotNil(expectedMTO.Locator)
	suite.NotNil(expectedMTO.ReferenceID)
	suite.Nil(expectedMTO.AvailableToPrimeAt)
	suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)
}

func (suite *MoveTaskOrderServiceSuite) TestListAllMoveTaskOrdersFetcher() {
	suite.T().Run("all move task orders", func(t *testing.T) {
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

		mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(false, nil)
		suite.NoError(err)

		mto := moveTaskOrders[0]

		suite.Equal(len(moveTaskOrders), 3)
		suite.Nil(mto.AvailableToPrimeAt)
	})

	suite.T().Run("all move task orders that are available to prime and using since", func(t *testing.T) {
		time1 := time.Now()
		time2 := time.Now()
		time3 := time.Now()
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
			MoveTaskOrder: models.Move{
				AvailableToPrimeAt: &time1,
			},
		})
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
			MoveTaskOrder: models.Move{
				AvailableToPrimeAt: &time2,
			},
		})

		oldMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
			MoveTaskOrder: models.Move{
				AvailableToPrimeAt: &time3,
			},
		})
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

		mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

		moveTaskOrders, err := mtoFetcher.ListAllMoveTaskOrders(true, nil)
		suite.NoError(err)
		suite.Equal(len(moveTaskOrders), 3)

		// Put 1 Move updatedAt in the past
		suite.NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
			time3.Add(-2*time.Second), oldMTO.ID).Exec())
		since := time3.Unix()
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(true, &since)
		suite.NoError(err)
		suite.Equal(len(mtosWithSince), 2)
	})
}
