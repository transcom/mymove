package order

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestOrderFetcher() {
	expectedMove := testdatagen.MakeDefaultMove(suite.DB())
	expectedOrder := expectedMove.Orders
	orderFetcher := NewOrderFetcher(suite.DB())

	order, err := orderFetcher.FetchOrder(expectedOrder.ID)
	suite.FatalNoError(err)

	suite.Equal(expectedOrder.ID, order.ID)
	suite.Equal(expectedOrder.ServiceMemberID, order.ServiceMemberID)
	suite.NotNil(order.NewDutyStation)
	suite.Equal(expectedOrder.NewDutyStationID, order.NewDutyStation.ID)
	suite.Equal(expectedOrder.NewDutyStation.AddressID, order.NewDutyStation.AddressID)
	suite.Equal(expectedOrder.NewDutyStation.Address.StreetAddress1, order.NewDutyStation.Address.StreetAddress1)
	suite.NotNil(order.Entitlement)
	suite.Equal(*expectedOrder.EntitlementID, order.Entitlement.ID)
	suite.Equal(expectedOrder.OriginDutyStation.ID, order.OriginDutyStation.ID)
	suite.Equal(expectedOrder.OriginDutyStation.AddressID, order.OriginDutyStation.AddressID)
	suite.Equal(expectedOrder.OriginDutyStation.Address.StreetAddress1, order.OriginDutyStation.Address.StreetAddress1)
	suite.NotZero(order.OriginDutyStation)
	suite.Equal(expectedMove.Locator, order.Moves[0].Locator)
}

func (suite *OrderServiceSuite) TestOrderFetcherWithEmptyFields() {
	// When move_orders and orders were consolidated, we moved the OriginDutyStation
	// field that used to only exist on the move_orders table into the orders table.
	// This means that existing orders in production won't have any values in the
	// OriginDutyStation column. To mimic that and to surface any issues, we didn't
	// update the testdatagen MakeOrder function so that new orders would have
	// an empty OriginDutyStation. During local testing in the office app, we
	// noticed an exception due to trying to load empty OriginDutyStations.
	// This was not caught by any tests, so we're adding one now.
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())

	expectedOrder.Entitlement = nil
	expectedOrder.EntitlementID = nil
	expectedOrder.Grade = nil
	expectedOrder.OriginDutyStation = nil
	expectedOrder.OriginDutyStationID = nil
	suite.MustSave(&expectedOrder)

	testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	orderFetcher := NewOrderFetcher(suite.DB())
	order, err := orderFetcher.FetchOrder(expectedOrder.ID)

	suite.FatalNoError(err)
	suite.Nil(order.Entitlement)
	suite.Nil(order.OriginDutyStation)
	suite.Nil(order.Grade)
}

func (suite *OrderServiceSuite) TestListMoves() {
	// Create a Move without a shipment to test that only Orders with shipments
	// are displayed to the TOO
	testdatagen.MakeDefaultMove(suite.DB())

	expectedMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	orderFetcher := NewOrderFetcher(suite.DB())

	suite.T().Run("returns moves", func(t *testing.T) {
		moves, moveCount, err := orderFetcher.ListOrders(officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)

		move := moves[0]

		suite.NotNil(move.Orders.ServiceMember)
		suite.Equal(expectedMove.Orders.ServiceMember.FirstName, move.Orders.ServiceMember.FirstName)
		suite.Equal(expectedMove.Orders.ServiceMember.LastName, move.Orders.ServiceMember.LastName)
		suite.Equal(expectedMove.Orders.ID, move.Orders.ID)
		suite.Equal(expectedMove.Orders.ServiceMemberID, move.Orders.ServiceMemberID)
		suite.NotNil(move.Orders.NewDutyStation)
		suite.Equal(expectedMove.Orders.NewDutyStationID, move.Orders.NewDutyStation.ID)
		suite.NotNil(move.Orders.Entitlement)
		suite.Equal(*expectedMove.Orders.EntitlementID, move.Orders.Entitlement.ID)
		suite.Equal(expectedMove.Orders.OriginDutyStation.ID, move.Orders.OriginDutyStation.ID)
		suite.NotNil(move.Orders.OriginDutyStation)
		suite.Equal(expectedMove.Orders.OriginDutyStation.AddressID, move.Orders.OriginDutyStation.AddressID)
		suite.Equal(expectedMove.Orders.OriginDutyStation.Address.StreetAddress1, move.Orders.OriginDutyStation.Address.StreetAddress1)
	})

	suite.T().Run("returns moves filtered by GBLOC", func(t *testing.T) {
		// This move is outside of the office user's GBLOC, so it should not be returned
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				Gbloc: "AGFM",
			},
		})

		moves, _, err := orderFetcher.ListOrders(officeUser.ID, &services.ListOrderParams{Page: swag.Int64(1)})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("only returns visible moves (where show = True)", func(t *testing.T) {
		params := services.ListOrderParams{}
		testdatagen.MakeHiddenHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, _, err := orderFetcher.ListOrders(officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("includes combo hhg and ppm moves", func(t *testing.T) {
		// Create a combination HHG and PPM move and make sure it's included
		expectedComboMove := testdatagen.MakeHHGPPMMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, moveCount, err := orderFetcher.ListOrders(officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(2, moveCount)
		suite.Len(moves, 2)

		moveIDs := []uuid.UUID{moves[0].ID, moves[1].ID}

		suite.Contains(moveIDs, expectedComboMove.ID)
	})

	suite.T().Run("returns moves filtered by service member affiliation", func(t *testing.T) {
		airForce := models.AffiliationAIRFORCE
		airForceString := "AIR_FORCE"
		params := services.ListOrderParams{Branch: &airForceString, Page: swag.Int64(1)}
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				Affiliation: &airForce,
			},
		})

		moves, _, err := orderFetcher.ListOrders(officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListMovesUSMCGBLOC() {
	orderFetcher := NewOrderFetcher(suite.DB())

	suite.T().Run("returns USMC order for USMC office user", func(t *testing.T) {
		marines := models.AffiliationMARINES
		// It doesn't matter what the Origin GBLOC is for the move. Only the Marines
		// affiliation matters for office users who are tied to the USMC GBLOC.
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})

		// Create move where service member has the default ARMY affiliation
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		officeUserOooRah := testdatagen.MakeOfficeUserWithUSMCGBLOC(suite.DB())
		// Create office user tied to the default LKNQ GBLOC
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(officeUserOooRah.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationMARINES, *moves[0].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationARMY, *moves[0].Orders.ServiceMember.Affiliation)
	})
}

func (suite *OrderServiceSuite) TestListMovesMarines() {
	suite.T().Run("does not return moves where the service member affiliation is Marines for non-USMC office user", func(t *testing.T) {
		orderFetcher := NewOrderFetcher(suite.DB())
		marines := models.AffiliationMARINES
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListMovesWithEmptyFields() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())

	expectedOrder.Entitlement = nil
	expectedOrder.EntitlementID = nil
	expectedOrder.Grade = nil
	expectedOrder.OriginDutyStation = nil
	expectedOrder.OriginDutyStationID = nil
	suite.MustSave(&expectedOrder)

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})
	// Only orders with shipments are returned, so we need to add a shipment
	// to the move we just created
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	// Add a second shipment to make sure we only return 1 order even if its
	// move has more than one shipment
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
	orderFetcher := NewOrderFetcher(suite.DB())
	moves, _, err := orderFetcher.ListOrders(officeUser.ID, &services.ListOrderParams{PerPage: swag.Int64(1), Page: swag.Int64(1)})

	suite.FatalNoError(err)
	suite.Nil(moves)

}

func (suite *OrderServiceSuite) TestListMovesWithPagination() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	for i := 0; i < 2; i++ {
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	}

	orderFetcher := NewOrderFetcher(suite.DB())
	params := services.ListOrderParams{Page: swag.Int64(1), PerPage: swag.Int64(1)}
	moves, count, err := orderFetcher.ListOrders(officeUser.ID, &params)

	suite.NoError(err)
	suite.Equal(1, len(moves))
	suite.Equal(2, count)

}

func (suite *OrderServiceSuite) TestListMovesWithSortOrder() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	// Default New Duty Station name is Fort Gordon
	expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:  models.MoveStatusAPPROVED,
			Locator: "AA1234",
		},
	})

	serviceMemberLastName := "Zephyer"
	affiliation := models.AffiliationNAVY
	edipi := "9999999999"
	newDutyStationName := "Ze Duty Station"
	newDutyStation2 := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: newDutyStationName,
		},
	})

	expectedMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Locator: "TTZ123",
		},
		ServiceMember: models.ServiceMember{Affiliation: &affiliation, LastName: &serviceMemberLastName, Edipi: &edipi},
		Order: models.Order{
			NewDutyStation:   newDutyStation2,
			NewDutyStationID: newDutyStation2.ID,
		},
	})

	orderFetcher := NewOrderFetcher(suite.DB())
	// Sort by service member name
	params := services.ListOrderParams{Sort: swag.String("lastName"), Order: swag.String("asc")}
	moves, _, err := orderFetcher.ListOrders(officeUser.ID, &params)

	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal("Spacemen, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	suite.Equal("Zephyer, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)

	params = services.ListOrderParams{Sort: swag.String("lastName"), Order: swag.String("desc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal("Zephyer, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)

	// Sort by locator
	params = services.ListOrderParams{Sort: swag.String("locator"), Order: swag.String("asc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove1.Locator, moves[0].Locator)
	suite.Equal(expectedMove2.Locator, moves[1].Locator)

	params = services.ListOrderParams{Sort: swag.String("locator"), Order: swag.String("desc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove2.Locator, moves[0].Locator)
	suite.Equal(expectedMove1.Locator, moves[1].Locator)

	// sort by move statuses
	params = services.ListOrderParams{Sort: swag.String("status"), Order: swag.String("asc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove1.Status, moves[0].Status)
	suite.Equal(expectedMove2.Status, moves[1].Status)

	params = services.ListOrderParams{Sort: swag.String("status"), Order: swag.String("desc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove2.Status, moves[0].Status)
	suite.Equal(expectedMove1.Status, moves[1].Status)

	// Sort by service member affiliations
	params = services.ListOrderParams{Sort: swag.String("branch"), Order: swag.String("asc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
	suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)

	params = services.ListOrderParams{Sort: swag.String("branch"), Order: swag.String("desc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
	suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)

	// Sort by destination duty station
	params = services.ListOrderParams{Sort: swag.String("destinationDutyStation"), Order: swag.String("asc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove1.Orders.NewDutyStation.Name, moves[0].Orders.NewDutyStation.Name)
	suite.Equal(expectedMove2.Orders.NewDutyStation.Name, moves[1].Orders.NewDutyStation.Name)

	params = services.ListOrderParams{Sort: swag.String("destinationDutyStation"), Order: swag.String("desc")}
	moves, _, err = orderFetcher.ListOrders(officeUser.ID, &params)
	suite.NoError(err)
	suite.Equal(2, len(moves))
	suite.Equal(expectedMove2.Orders.NewDutyStation.Name, moves[0].Orders.NewDutyStation.Name)
	suite.Equal(expectedMove1.Orders.NewDutyStation.Name, moves[1].Orders.NewDutyStation.Name)

}
