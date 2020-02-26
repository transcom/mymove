package movetaskorder

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: "MS",
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})
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
	suite.NotNil(expectedMTO.ReferenceID)
	suite.False(expectedMTO.IsAvailableToPrime)
	suite.False(expectedMTO.IsCanceled)

	var serviceItems models.MTOServiceItems
	suite.DB().All(&serviceItems)
	suite.Equal(2, len(serviceItems))

	// Make sure we don't create more default service items
	serviceItems = models.MTOServiceItems{}
	actualMTO, err = mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID)
	suite.NoError(err)

	suite.DB().All(&serviceItems)
	suite.Equal(2, len(serviceItems))
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
	suite.NotNil(expectedMTO.ReferenceID)
	suite.False(expectedMTO.IsAvailableToPrime)
	suite.False(expectedMTO.IsCanceled)
}
