package movetaskorder

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	expectedMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveOrder: expectedMoveOrder,
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	actualMTO, err := mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID)
	suite.NoError(err)

	suite.NotZero(expectedMTO.ID, actualMTO.ID)
	suite.Equal(expectedMTO.MoveOrder.ID, actualMTO.MoveOrder.ID)
	suite.NotZero(actualMTO.MoveOrder)
	suite.Nil(expectedMTO.ReferenceID)
	suite.False(expectedMTO.IsAvailableToPrime)
	suite.False(expectedMTO.IsCanceled)
}

func (suite *MoveTaskOrderServiceSuite) TestListMoveTaskOrdersFetcher() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	expectedMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveOrder: expectedMoveOrder,
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	moveTaskOrders, err := mtoFetcher.ListMoveTaskOrders(expectedMoveOrder.ID)
	suite.NoError(err)

	actualMTO := moveTaskOrders[0]

	suite.NotZero(expectedMTO.ID, actualMTO.ID)
	suite.Equal(expectedMTO.MoveOrder.ID, actualMTO.MoveOrder.ID)
	suite.NotZero(actualMTO.MoveOrder)
	suite.Nil(expectedMTO.ReferenceID)
	suite.False(expectedMTO.IsAvailableToPrime)
	suite.False(expectedMTO.IsCanceled)
}