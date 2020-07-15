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
	suite.Equal(expectedMoveOrder.ServiceMemberID, moveOrder.ServiceMemberID)
	suite.NotNil(moveOrder.NewDutyStation)
	suite.Equal(expectedMoveOrder.NewDutyStationID, moveOrder.NewDutyStation.ID)
	suite.NotNil(moveOrder.Entitlement)
	suite.Equal(*expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}

func (suite *MoveOrderServiceSuite) TestMoveOrderFetcherWithEmptyFields() {
	// When move_orders and orders were consolidated, we moved the OriginDutyStation
	// field that used to only exist on the move_orders table into the orders table.
	// This means that existing orders in production won't have any values in the
	// OriginDutyStation column. To mimic that and to surface any issues, we didn't
	// update the testdatagen MakeOrder function so that new orders would have
	// an empty OriginDutyStation. During local testing in the office app, we
	// noticed an exception due to trying to load empty OriginDutyStations.
	// This was not caught by any tests, so we're adding one now.
	expectedMoveOrder := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())
	moveOrder, err := moveOrderFetcher.FetchMoveOrder(expectedMoveOrder.ID)

	suite.FatalNoError(err)
	suite.Nil(moveOrder.Entitlement)
	suite.Nil(moveOrder.OriginDutyStation)
	suite.Nil(moveOrder.Grade)
}

func (suite *MoveOrderServiceSuite) TestListMoveOrder() {
	expectedMoveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())
	moveOrders, err := moveOrderFetcher.ListMoveOrders()
	suite.FatalNoError(err)
	suite.Len(moveOrders, 1)

	moveOrder := moveOrders[0]
	suite.NotNil(moveOrder.ServiceMember)
	suite.Equal(expectedMoveOrder.ServiceMember.FirstName, moveOrder.ServiceMember.FirstName)
	suite.Equal(expectedMoveOrder.ServiceMember.LastName, moveOrder.ServiceMember.LastName)
	suite.Equal(expectedMoveOrder.ID, moveOrder.ID)
	suite.Equal(expectedMoveOrder.ServiceMemberID, moveOrder.ServiceMemberID)
	suite.NotNil(moveOrder.NewDutyStation)
	suite.Equal(expectedMoveOrder.NewDutyStationID, moveOrder.NewDutyStation.ID)
	suite.NotNil(moveOrder.Entitlement)
	suite.Equal(*expectedMoveOrder.EntitlementID, moveOrder.Entitlement.ID)
	suite.Equal(expectedMoveOrder.OriginDutyStation.ID, moveOrder.OriginDutyStation.ID)
	suite.NotNil(moveOrder.OriginDutyStation)
}

func (suite *MoveOrderServiceSuite) TestListMoveOrderWithEmptyFields() {
	testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{})
	moveOrderFetcher := NewMoveOrderFetcher(suite.DB())
	moveOrders, err := moveOrderFetcher.ListMoveOrders()
	moveOrder := moveOrders[0]

	suite.FatalNoError(err)
	suite.Nil(moveOrder.Entitlement)
	suite.Nil(moveOrder.OriginDutyStation)
	suite.Nil(moveOrder.Grade)
}
