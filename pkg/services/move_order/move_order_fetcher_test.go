package moveorder

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveOrderServiceSuite) TestMoveOrderFetcher() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())

	moveOrder, err := moveOrderFetcher.FetchMoveOrder(expectedMoveOrder.ID)
	suite.FatalNoError(err)

	suite.Equal(expectedMoveOrder.ID, moveOrder.ID)
	suite.Equal(expectedMoveOrder.CustomerID, moveOrder.CustomerID)
	suite.NotNil(moveOrder.DestinationDutyStation)
	suite.Equal(*expectedMoveOrder.DestinationDutyStationID, moveOrder.DestinationDutyStation.ID)
	suite.NotNil(moveOrder.Entitlement)
	suite.Equal(*expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}

func (suite *MoveOrderServiceSuite) TestListMoveOrder() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())
	moveOrders, err := moveOrderFetcher.ListMoveOrders()
	suite.FatalNoError(err)
	suite.Len(moveOrders, 1)

	moveOrder := moveOrders[0]
	suite.NotNil(moveOrder.Customer)
	suite.Equal(expectedMoveOrder.Customer.FirstName, moveOrder.Customer.FirstName)
	suite.Equal(expectedMoveOrder.Customer.LastName, moveOrder.Customer.LastName)
	suite.Equal(expectedMoveOrder.ID, moveOrder.ID)
	suite.Equal(expectedMoveOrder.CustomerID, moveOrder.CustomerID)
	suite.NotNil(moveOrder.DestinationDutyStation)
	suite.Equal(*expectedMoveOrder.DestinationDutyStationID, moveOrder.DestinationDutyStation.ID)
	suite.NotNil(moveOrder.Entitlement)
	suite.Equal(*expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotNil(moveOrder.OriginDutyStation)
}
