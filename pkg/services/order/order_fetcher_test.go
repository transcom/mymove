package order

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestOrderFetcher() {
	expectedMove := testdatagen.MakeDefaultMove(suite.DB())
	expectedOrder := expectedMove.Orders
	orderFetcher := NewOrderFetcher()

	order, err := orderFetcher.FetchOrder(suite.AppContextForTest(), expectedOrder.ID)
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
	orderFetcher := NewOrderFetcher()
	order, err := orderFetcher.FetchOrder(suite.AppContextForTest(), expectedOrder.ID)

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
	//"30813"
	// May have to create postalcodetogbolc for office user
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	agfmPostalCode := "06001"
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), agfmPostalCode, "AGFM")
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	orderFetcher := NewOrderFetcher()

	suite.T().Run("returns moves", func(t *testing.T) {
		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{})

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
			PickupAddress: models.Address{
				PostalCode: agfmPostalCode,
			},
		})

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{Page: swag.Int64(1)})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("only returns visible moves (where show = True)", func(t *testing.T) {
		params := services.ListOrderParams{}
		testdatagen.MakeHiddenHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("includes combo hhg and ppm moves", func(t *testing.T) {
		// Create a combination HHG and PPM move and make sure it's included
		expectedComboMove := testdatagen.MakeHHGPPMMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{})

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

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("returns moves filtered submitted at", func(t *testing.T) {

		submittedAt := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt,
			},
		})

		// Test edge cases
		submittedAt2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt2,
			},
		})

		// Test edge cases
		submittedAt3 := time.Date(2022, 03, 31, 23, 59, 59, 59, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt3,
			},
		})

		params := services.ListOrderParams{SubmittedAt: &submittedAt}

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.T().Run("returns moves filtered by requested pickup date", func(t *testing.T) {
		requestedPickupDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
			},
		})
		requestedMoveDateString := createdMove.MTOShipments[0].RequestedPickupDate.Format("2006-01-02")
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			RequestedMoveDate: &requestedMoveDateString,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListMovesUSMCGBLOC() {
	orderFetcher := NewOrderFetcher()

	suite.T().Run("returns USMC order for USMC office user", func(t *testing.T) {
		// Map default shipment ZIP code to default office user GBLOC
		testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", "KKFA")
		testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", "KKFA")

		marines := models.AffiliationMARINES
		// It doesn't matter what the Origin GBLOC is for the move. Only the Marines
		// affiliation matters for office users who are tied to the USMC GBLOC.
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})

		// Create move where service member has the default ARMY affiliation
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		officeUserOooRah := testdatagen.MakeOfficeUserWithUSMCGBLOC(suite.DB())
		// Create office user tied to the default KKFA GBLOC
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserOooRah.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationMARINES, *moves[0].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationARMY, *moves[0].Orders.ServiceMember.Affiliation)
	})
}

func (suite *OrderServiceSuite) TestListMovesMarines() {
	suite.T().Run("does not return moves where the service member affiliation is Marines for non-USMC office user", func(t *testing.T) {
		orderFetcher := NewOrderFetcher()
		marines := models.AffiliationMARINES
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		// Map default shipment ZIP code to default office user GBLOC
		testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
		testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

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
	orderFetcher := NewOrderFetcher()
	moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{PerPage: swag.Int64(1), Page: swag.Int64(1)})

	suite.FatalNoError(err)
	suite.Nil(moves)

}

func (suite *OrderServiceSuite) TestListMovesWithPagination() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Map default shipment postal code to office user's GBLOC
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	for i := 0; i < 2; i++ {
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	}

	orderFetcher := NewOrderFetcher()
	params := services.ListOrderParams{Page: swag.Int64(1), PerPage: swag.Int64(1)}
	moves, count, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

	suite.NoError(err)
	suite.Equal(1, len(moves))
	suite.Equal(2, count)

}

func (suite *OrderServiceSuite) TestListMovesWithSortOrder() {
	// SET UP: Dates for sorting by Requested Move Date
	// - We want dates 2 and 3 to sandwich requestedMoveDate1 so we can test that the min() query is working
	requestedMoveDate1 := time.Date(testdatagen.GHCTestYear, 02, 20, 0, 0, 0, 0, time.UTC)
	requestedMoveDate2 := time.Date(testdatagen.GHCTestYear, 03, 03, 0, 0, 0, 0, time.UTC)
	requestedMoveDate3 := time.Date(testdatagen.GHCTestYear, 01, 15, 0, 0, 0, 0, time.UTC)

	// SET UP: Service Members for sorting by Service Member Last Name and Branch
	// - We'll need two other service members to test the last name sort, Lea Spacemen and Leo Zephyer
	serviceMemberFirstName := "Lea"
	serviceMemberLastName := "Zephyer"
	affiliation := models.AffiliationNAVY
	edipi := "9999999999"

	// SET UP: New Duty Station for sorting by destination duty station
	newDutyStationName := "Ze Duty Station"
	newDutyStation2 := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: newDutyStationName,
		},
	})

	// CREATE EXPECTED MOVES
	expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		// Default New Duty Station name is Fort Gordon
		Move: models.Move{
			Status:  models.MoveStatusAPPROVED,
			Locator: "AA1234",
		},
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &requestedMoveDate1,
		},
	})
	expectedMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Locator: "TTZ123",
		},
		// Lea Spacemen
		ServiceMember: models.ServiceMember{Affiliation: &affiliation, FirstName: &serviceMemberFirstName, Edipi: &edipi},
		Order: models.Order{
			NewDutyStation:   newDutyStation2,
			NewDutyStationID: newDutyStation2.ID,
		},
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &requestedMoveDate2,
		},
	})
	// Create a second shipment so we can test min() sort
	testdatagen.MakeMTOShipmentWithMove(suite.DB(), &expectedMove2, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &requestedMoveDate3,
		},
	})

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)
	orderFetcher := NewOrderFetcher()

	suite.T().Run("Sort by locator code", func(t *testing.T) {
		params := services.ListOrderParams{Sort: swag.String("locator"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		params = services.ListOrderParams{Sort: swag.String("locator"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Locator, moves[0].Locator)
		suite.Equal(expectedMove1.Locator, moves[1].Locator)
	})

	suite.T().Run("Sort by move status", func(t *testing.T) {
		params := services.ListOrderParams{Sort: swag.String("status"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Status, moves[0].Status)
		suite.Equal(expectedMove2.Status, moves[1].Status)

		params = services.ListOrderParams{Sort: swag.String("status"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Status, moves[0].Status)
		suite.Equal(expectedMove1.Status, moves[1].Status)
	})

	suite.T().Run("Sort by service member affiliations", func(t *testing.T) {
		params := services.ListOrderParams{Sort: swag.String("branch"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{Sort: swag.String("branch"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)
	})

	suite.T().Run("Sort by destination duty station", func(t *testing.T) {
		params := services.ListOrderParams{Sort: swag.String("destinationDutyStation"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Orders.NewDutyStation.Name, moves[0].Orders.NewDutyStation.Name)
		suite.Equal(expectedMove2.Orders.NewDutyStation.Name, moves[1].Orders.NewDutyStation.Name)

		params = services.ListOrderParams{Sort: swag.String("destinationDutyStation"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Orders.NewDutyStation.Name, moves[0].Orders.NewDutyStation.Name)
		suite.Equal(expectedMove1.Orders.NewDutyStation.Name, moves[1].Orders.NewDutyStation.Name)
	})

	suite.T().Run("Sort by request move date", func(t *testing.T) {
		params := services.ListOrderParams{Sort: swag.String("requestedMoveDate"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(2, len(moves[0].MTOShipments)) // the move with two shipments has the earlier date
		suite.Equal(1, len(moves[1].MTOShipments))
		// NOTE: You have to use Jan 02, 2006 as the example for date/time formatting in Go
		suite.Equal(requestedMoveDate1.Format("2006/01/02"), moves[1].MTOShipments[0].RequestedPickupDate.Format("2006/01/02"))

		params = services.ListOrderParams{Sort: swag.String("requestedMoveDate"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(1, len(moves[0].MTOShipments)) // the move with one shipment should be first
		suite.Equal(2, len(moves[1].MTOShipments))
		suite.Equal(requestedMoveDate1.Format("2006/01/02"), moves[0].MTOShipments[0].RequestedPickupDate.Format("2006/01/02"))
	})

	// MUST BE LAST, ADDS EXTRA MOVE
	suite.T().Run("Sort by service member last name", func(t *testing.T) {
		// Last name sort is the only one that needs 3 moves for a complete test, so add that here at the end
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			// Leo Zephyer
			ServiceMember: models.ServiceMember{LastName: &serviceMemberLastName},
		})

		params := services.ListOrderParams{Sort: swag.String("lastName"), Order: swag.String("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Spacemen, Lea", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Zephyer, Leo", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)

		params = services.ListOrderParams{Sort: swag.String("lastName"), Order: swag.String("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Zephyer, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Lea", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)
	})
}

func (suite *OrderServiceSuite) TestListMovesNeedingServicesCounselingWithGBLOCSortFilter() {

	// TESTCASE SCENARIO
	// Under test: OrderFetcher.ListOrders function
	// Mocked:     None
	// Set up:     We create 2 USMC moves with different GBLOCs, LKNQ and ZANY
	//             We create an office user with the LKNQ GBLOC

	// Create a services counselor with defauult GBLOC, LKNQ
	officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})

	// Default Origin Duty Station GBLOC is LKNQ
	hhgMoveType := models.SelectedMoveTypeHHG
	submittedAt := time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC)
	lknqMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusNeedsServiceCounseling,
			SubmittedAt:      &submittedAt,
		},
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	// Create a dutystation with ZANY GBLOC
	dutyStationAddress2 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Anchor 1212",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "89898",
			Country:        swag.String("United States"),
		},
	})

	originDutyStation2 := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:      "Fort Sam Snap",
			AddressID: dutyStationAddress2.ID,
			Address:   dutyStationAddress2,
		},
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(), dutyStationAddress2.PostalCode, "ZANY")

	// Create a second move from the ZANY gbloc
	testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:  models.MoveStatusNeedsServiceCounseling,
			Locator: "ZZ1234",
		},
		Order: models.Order{
			OriginDutyStation:   &originDutyStation2,
			OriginDutyStationID: &originDutyStation2.ID,
		},
	})

	suite.T().Run("Filter by origin GBLOC", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: OrderFetcher.ListOrders function
		// Mocked:     None
		// Set up:     We create 2 USMC moves with different GBLOCs, ACME and ZANY
		//             We create an office user with the ACME GBLOC
		//             Then we request a list of moves sorted by GBLOC, first ascending then descending
		// Expected outcome:
		//             We expect one move to be returned: the the ACME move

		// Setup and run the function under test sorting GBLOC with ascending mode
		orderFetcher := NewOrderFetcher()
		statuses := []string{"NEEDS SERVICE COUNSELING"}
		// Sort by service member name
		params := services.ListOrderParams{Sort: swag.String("originGBLOC"), Order: swag.String("asc"), Status: statuses}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		// Check the results
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(lknqMove.ID, moves[0].ID)
	})
}
