package moveorder

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveOrderServiceSuite) TestMoveOrderFetcher() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())

	moveOrder, err := moveOrderFetcher.FetchMoveOrder(expectedMoveOrder.ID)
	suite.NoError(err)

	suite.Equal(expectedMoveOrder.ID, moveOrder.ID)
	suite.Equal(expectedMoveOrder.CustomerID, moveOrder.CustomerID)
	suite.Equal(expectedMoveOrder.DestinationDutyStationID, moveOrder.DestinationDutyStation.ID)
	suite.NotZero(moveOrder.DestinationDutyStation)
	suite.Equal(expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.NotZero(moveOrder.Entitlement)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}

func (suite *MoveOrderServiceSuite) TestListMoveOrder() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())

	moveOrders, err := moveOrderFetcher.ListMoveOrders()
	suite.NoError(err)
	suite.Len(moveOrders, 1)

	moveOrder := moveOrders[0]
	suite.Equal(expectedMoveOrder.Customer.FirstName, moveOrder.Customer.FirstName)
	suite.Equal(expectedMoveOrder.Customer.LastName, moveOrder.Customer.LastName)
	suite.Equal(expectedMoveOrder.ID, moveOrder.ID)
	suite.Equal(expectedMoveOrder.CustomerID, moveOrder.CustomerID)
	suite.Equal(expectedMoveOrder.DestinationDutyStationID, moveOrder.DestinationDutyStation.ID)
	suite.NotZero(moveOrder.DestinationDutyStation)
	suite.Equal(expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.NotZero(moveOrder.Entitlement)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}
