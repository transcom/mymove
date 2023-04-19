package order

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestFetchOrder() {
	expectedMove := factory.BuildMove(suite.DB(), nil, nil)
	expectedOrder := expectedMove.Orders
	orderFetcher := NewOrderFetcher()

	order, err := orderFetcher.FetchOrder(suite.AppContextForTest(), expectedOrder.ID)
	suite.FatalNoError(err)

	suite.Equal(expectedOrder.ID, order.ID)
	suite.Equal(expectedOrder.ServiceMemberID, order.ServiceMemberID)
	suite.NotNil(order.NewDutyLocation)
	suite.Equal(expectedOrder.NewDutyLocationID, order.NewDutyLocation.ID)
	suite.Equal(expectedOrder.NewDutyLocation.AddressID, order.NewDutyLocation.AddressID)
	suite.Equal(expectedOrder.NewDutyLocation.Address.StreetAddress1, order.NewDutyLocation.Address.StreetAddress1)
	suite.NotNil(order.Entitlement)
	suite.Equal(*expectedOrder.EntitlementID, order.Entitlement.ID)
	suite.Equal(expectedOrder.OriginDutyLocation.ID, order.OriginDutyLocation.ID)
	suite.Equal(expectedOrder.OriginDutyLocation.AddressID, order.OriginDutyLocation.AddressID)
	suite.Equal(expectedOrder.OriginDutyLocation.Address.StreetAddress1, order.OriginDutyLocation.Address.StreetAddress1)
	suite.NotZero(order.OriginDutyLocation)
	suite.Equal(expectedMove.Locator, order.Moves[0].Locator)
}

func (suite *OrderServiceSuite) TestFetchOrderWithEmptyFields() {
	// When move_orders and orders were consolidated, we moved the OriginDutyLocation
	// field that used to only exist on the move_orders table into the orders table.
	// This means that existing orders in production won't have any values in the
	// OriginDutyLocation column. To mimic that and to surface any issues, we didn't
	// update the testdatagen MakeOrder function so that new orders would have
	// an empty OriginDutyLocation. During local testing in the office app, we
	// noticed an exception due to trying to load empty OriginDutyLocations.
	// This was not caught by any tests, so we're adding one now.
	expectedOrder := factory.BuildOrder(suite.DB(), nil, nil)

	expectedOrder.Entitlement = nil
	expectedOrder.EntitlementID = nil
	expectedOrder.Grade = nil
	expectedOrder.OriginDutyLocation = nil
	expectedOrder.OriginDutyLocationID = nil
	suite.MustSave(&expectedOrder)

	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    expectedOrder,
			LinkOnly: true,
		},
	}, nil)

	orderFetcher := NewOrderFetcher()
	order, err := orderFetcher.FetchOrder(suite.AppContextForTest(), expectedOrder.ID)

	suite.FatalNoError(err)
	suite.Nil(order.Entitlement)
	suite.Nil(order.OriginDutyLocation)
	suite.Nil(order.Grade)
}

func (suite *OrderServiceSuite) TestListOrders() {

	agfmPostalCode := "06001"
	setupTestData := func() (models.OfficeUser, models.Move) {
		// Make an office user → GBLOC X
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		// Create a move with a shipment → GBLOC X
		move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		// Make a postal code and GBLOC → AGFM
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), agfmPostalCode, "AGFM")

		return officeUser, move
	}
	orderFetcher := NewOrderFetcher()

	suite.Run("returns moves", func() {
		// Under test: ListOrders
		// Mocked:           None
		// Set up:           Make 2 moves, one with a shipment and one without.
		//                   The shipment should have a pickup GBLOC that matches the office users transportation GBLOC
		//                   In other words, shipment should originate from same GBLOC as the office user
		// Expected outcome: Only the move with a shipment should be returned by ListOrders
		officeUser, expectedMove := setupTestData()

		// Create a Move without a shipment
		factory.BuildMove(suite.DB(), nil, nil)

		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{})

		// Expect a single move returned
		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)

		// Check that move matches
		move := moves[0]
		suite.NotNil(move.Orders.ServiceMember)
		suite.Equal(expectedMove.Orders.ServiceMember.FirstName, move.Orders.ServiceMember.FirstName)
		suite.Equal(expectedMove.Orders.ServiceMember.LastName, move.Orders.ServiceMember.LastName)
		suite.Equal(expectedMove.Orders.ID, move.Orders.ID)
		suite.Equal(expectedMove.Orders.ServiceMemberID, move.Orders.ServiceMemberID)
		suite.NotNil(move.Orders.NewDutyLocation)
		suite.Equal(expectedMove.Orders.NewDutyLocationID, move.Orders.NewDutyLocation.ID)
		suite.NotNil(move.Orders.Entitlement)
		suite.Equal(*expectedMove.Orders.EntitlementID, move.Orders.Entitlement.ID)
		suite.Equal(expectedMove.Orders.OriginDutyLocation.ID, move.Orders.OriginDutyLocation.ID)
		suite.NotNil(move.Orders.OriginDutyLocation)
		suite.Equal(expectedMove.Orders.OriginDutyLocation.AddressID, move.Orders.OriginDutyLocation.AddressID)
		suite.Equal(expectedMove.Orders.OriginDutyLocation.Address.StreetAddress1, move.Orders.OriginDutyLocation.Address.StreetAddress1)
	})

	suite.Run("returns moves filtered by GBLOC", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, one with a pickup GBLOC that matches the office users transportation GBLOC
		//                   (which is done in setupTestData) and one with a pickup GBLOC that doesn't
		// Expected outcome: Only the move with the correct GBLOC should be returned by ListOrders
		officeUser, expectedMove := setupTestData()

		// This move's pickup GBLOC of the office user's GBLOC, so it should not be returned
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			PickupAddress: models.Address{
				PostalCode: agfmPostalCode,
			},
		})

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{Page: swag.Int64(1)})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("only returns visible moves (where show = True)", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, one correctly setup in setupTestData (show = True)
		//                   and one with show = False
		// Expected outcome: Only the move with show = True should be returned by ListOrders
		officeUser, expectedMove := setupTestData()

		params := services.ListOrderParams{}
		testdatagen.MakeHiddenHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("includes combo hhg and ppm moves", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, one default move setup in setupTestData (show = True)
		//                   and one a combination HHG and PPM move and make sure it's included
		// Expected outcome: Both moves should be returned by ListOrders
		officeUser, expectedMove := setupTestData()
		expectedComboMove := testdatagen.MakeHHGPPMMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(2, moveCount)
		suite.Len(moves, 2)

		var moveIDs []uuid.UUID
		for _, move := range moves {
			moveIDs = append(moveIDs, move.ID)
		}
		suite.Contains(moveIDs, expectedComboMove.ID)
		suite.Contains(moveIDs, expectedMove.ID)
	})

	suite.Run("returns moves filtered by service member affiliation", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, one default move setup in setupTestData (show = True)
		//                   and one specific to Airforce and make sure it's included
		//                   Fetch filtered to Airforce moves.
		// Expected outcome: Only the Airforce move should be returned
		officeUser, _ := setupTestData()

		// Create the airforce move
		airForce := models.AffiliationAIRFORCE
		airForceString := "AIR_FORCE"
		airForceMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				Affiliation: &airForce,
			},
		})
		// Filter by airforce move
		params := services.ListOrderParams{Branch: &airForceString, Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(airForceMove.ID, move.ID)

	})

	suite.Run("returns moves filtered submitted at", func() {
		// Under test: ListOrders
		// Set up:           Make 3 moves, with different submitted_at times, and search for a specific move
		// Expected outcome: Only the one move with the right date should be returned
		officeUser, _ := setupTestData()

		// Move with specified timestamp
		submittedAt := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		expectedMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt,
			},
		})

		// Test edge cases (one day later)
		submittedAt2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt2,
			},
		})

		// Test edge cases (one second earlier)
		submittedAt3 := time.Date(2022, 03, 31, 23, 59, 59, 59, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &submittedAt3,
			},
		})

		// Filter by submittedAt timestamp
		params := services.ListOrderParams{SubmittedAt: &submittedAt}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("returns moves filtered appeared in TOO at", func() {
		// Under test: ListOrders
		// Expected outcome: Only the three move with the right date should be returned
		officeUser, _ := setupTestData()

		// Moves with specified timestamp
		specifiedDay := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		specifiedTimestamp1 := time.Date(2022, 04, 01, 1, 0, 0, 0, time.UTC)
		specifiedTimestamp2 := time.Date(2022, 04, 01, 23, 59, 59, 999999000, time.UTC) // the upper bound is 999999499 nanoseconds but the DB only stores microseconds

		matchingSubmittedAt := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &specifiedDay,
			},
		})

		matchingSCCompletedAt := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ServiceCounselingCompletedAt: &specifiedTimestamp1,
			},
		})

		matchingApprovalsRequestedAt := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ApprovalsRequestedAt: &specifiedTimestamp2,
			},
		})

		// Test non dates matching
		nonMatchingDate1 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		nonMatchingDate2 := time.Date(2022, 03, 31, 23, 59, 59, 999999000, time.UTC) // the upper bound is 999999499 nanoseconds but the DB only stores microseconds
		nonMatchingDate3 := time.Date(2023, 04, 01, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt:                  &nonMatchingDate1,
				ServiceCounselingCompletedAt: &nonMatchingDate2,
				ApprovalsRequestedAt:         &nonMatchingDate3,
			},
		})

		// Filter by AppearedInTOOAt timestamp
		params := services.ListOrderParams{AppearedInTOOAt: &specifiedDay}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(3, len(moves))
		var foundIDs []uuid.UUID
		for _, move := range moves {
			foundIDs = append(foundIDs, move.ID)
		}
		suite.Contains(foundIDs, matchingSubmittedAt.ID)
		suite.Contains(foundIDs, matchingSCCompletedAt.ID)
		suite.Contains(foundIDs, matchingApprovalsRequestedAt.ID)
	})

	suite.Run("returns moves filtered by requested pickup date", func() {
		// Under test: ListOrders
		// Set up:           Make 3 moves, with different submitted_at times, and search for a specific move
		// Expected outcome: Only the one move with the right date should be returned
		officeUser, _ := setupTestData()

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

	suite.Run("returns moves filtered by ppm type", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, with different ppm types, and search for both types
		// Expected outcome: search results should only include the move with the PPM type that was searched for
		officeUser, partialPPMMove := setupTestData()
		suite.Equal("PARTIAL", *partialPPMMove.PPMType)
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				PPMType: swag.String("FULL"),
				Locator: "FULLLL",
			},
		})
		fullPPMMove := ppmShipment.Shipment.MoveTaskOrder

		// Search for PARTIAL PPM moves
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			PPMType: swag.String("PARTIAL"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(partialPPMMove.Locator, moves[0].Locator)

		// Search for FULL PPM moves
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			PPMType: swag.String("FULL"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(fullPPMMove.Locator, moves[0].Locator)
	})

	suite.Run("returns moves filtered by closeout location", func() {
		// Under test: ListOrders
		// Set up:           Make a move with a closeout office. Search for that closeout office.
		// Expected outcome: Only the one ppmShipment with the right closeout office should be returned
		officeUser, _ := setupTestData()

		ftBragg := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "Ft Bragg",
				},
			},
		}, nil)
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				CloseoutOfficeID: &ftBragg.ID,
			},
		})

		// Search should be case insensitive and allow partial matches
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			CloseoutLocation: models.StringPointer("fT bR"),
			NeedsPPMCloseout: models.BoolPointer(true),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(ppmShipment.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})

	suite.Run("returns moves filtered by closeout initiated date", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves with PPM shipments ready for closeout, with different submitted_at times,
		//                   and search for a specific move
		// Expected outcome: Only the one move with the right date should be returned
		officeUser, _ := setupTestData()

		// Create a PPM submitted on April 1st
		closeoutInitiatedDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdPPM := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate,
			},
		})

		// Create a PPM submitted on April 2nd
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		createdPPM2 := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate2,
			},
		})

		// Search for PPMs submitted on April 1st
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(createdPPM.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.NotEqual(createdPPM2.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})

	suite.Run("latest closeout initiated date is used for filter", func() {
		// Under test: ListOrders
		// Set up:           Make one move with multiple ppm shipments with different closeout initiated times, and
		//                   search for multiple different times
		// Expected outcome: Only a search for the latest of the closeout dates should find the move
		officeUser, _ := setupTestData()

		// Create a PPM submitted on April 1st
		closeoutInitiatedDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdPPM := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate,
			},
		})
		// Add another PPM for the same move submitted on April 1st
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate2,
				Status:      models.PPMShipmentStatusNeedsPaymentApproval,
			},
			Move: createdPPM.Shipment.MoveTaskOrder,
		})

		// Search for PPMs submitted on April 1st
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate,
		})
		suite.Empty(moves)
		suite.FatalNoError(err)

		// Search for PPMs submitted on April 2nd
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate2,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(createdPPM.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})
}

func (suite *OrderServiceSuite) TestListOrdersUSMCGBLOC() {
	orderFetcher := NewOrderFetcher()

	suite.Run("returns USMC order for USMC office user", func() {
		marines := models.AffiliationMARINES
		// It doesn't matter what the Origin GBLOC is for the move. Only the Marines
		// affiliation matters for office users who are tied to the USMC GBLOC.
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})

		// Create move where service member has the default ARMY affiliation
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		tioRole := roles.Role{RoleType: roles.RoleTypeTIO}
		tooRole := roles.Role{RoleType: roles.RoleTypeTOO}
		officeUserOooRah := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USMC",
				},
			},
			{
				Model: models.User{
					Roles: []roles.Role{tioRole, tooRole},
				},
			},
		}, nil)
		// Create office user tied to the default KKFA GBLOC
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

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

func (suite *OrderServiceSuite) TestListOrdersPPMCloseoutForArmyAirforce() {
	orderFetcher := NewOrderFetcher()
	showMove := true

	suite.Run("office user in normal GBLOC should only see non-Navy/Marines/CoastGuard moves that need closeout in closeout tab", func() {
		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		army := models.AffiliationARMY
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: move,
		})

		// Moves that are not ready for closeout should not show in this queue
		af := models.AffiliationAIRFORCE
		afMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &af,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusDraft,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: afMove,
		})

		// Coast guard moves should not show up in our office user's closeout queue
		cg := models.AffiliationCOASTGUARD
		cgMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &cg,
				},
			},
		}, nil)

		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: cgMove,
		})

		params := services.ListOrderParams{PerPage: swag.Int64(9), Page: swag.Int64(1), NeedsPPMCloseout: swag.Bool(true), Status: []string{string(models.MoveStatusNeedsServiceCounseling)}}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(move.Locator, moves[0].Locator)
	})

	suite.Run("office user in normal GBLOC should not see moves that require closeout in counseling tab", func() {
		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		// PPM moves that need closeout should not show up in counseling queue
		army := models.AffiliationARMY
		closeoutMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: closeoutMove,
		})

		// PPM moves that are not in one of the closeout statuses
		airforce := models.AffiliationAIRFORCE
		nonCloseoutMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &airforce,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusDraft,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: nonCloseoutMove,
		})

		params := services.ListOrderParams{PerPage: swag.Int64(9), Page: swag.Int64(1), NeedsPPMCloseout: swag.Bool(false), Status: []string{string(models.MoveStatusNeedsServiceCounseling)}}

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(nonCloseoutMove.Locator, moves[0].Locator)
	})
}

func (suite *OrderServiceSuite) TestListOrdersPPMCloseoutForNavyCoastGuardAndMarines() {
	orderFetcher := NewOrderFetcher()
	showMove := true

	suite.Run("returns Navy order for NAVY office user when there's a ppm shipment in closeout", func() {
		navy := models.AffiliationNAVY
		// It doesn't matter what the Origin GBLOC is for the move. Only the navy
		// affiliation matters for SC  who are tied to the NAVY GBLOC.
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &navy,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: move,
		})

		cg := models.AffiliationCOASTGUARD
		cgMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &cg,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: cgMove,
		})

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "NAVY",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		params := services.ListOrderParams{PerPage: swag.Int64(9), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationNAVY, *moves[0].Orders.ServiceMember.Affiliation)

	})

	suite.Run("returns TVCB order for TVCB office user when there's a ppm shipment in closeout", func() {
		marines := models.AffiliationMARINES
		// It doesn't matter what the Origin GBLOC is for the move. Only the marines
		// affiliation matters for SC  who are tied to the TVCB GBLOC.
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &marines,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: move,
		})

		army := models.AffiliationARMY
		nonMarineMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: nonMarineMove,
		})

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "TVCB",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationMARINES, *moves[0].Orders.ServiceMember.Affiliation)

	})

	suite.Run("returns coast guard order for USCG office user when there's a ppm shipment in closeout and filters out non coast guard moves", func() {
		cg := models.AffiliationCOASTGUARD
		// It doesn't matter what the Origin GBLOC is for the move. Only the coast guard
		// affiliation matters for SC  who are tied to the USCG GBLOC.
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &cg,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: move,
		})

		army := models.AffiliationARMY
		armyMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: armyMove,
		})

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USCG",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationCOASTGUARD, *moves[0].Orders.ServiceMember.Affiliation)
	})

	suite.Run("Filters out moves with PPM shipments not in the status of NeedsApproval", func() {
		cg := models.AffiliationCOASTGUARD

		cgMoveInWrongStatus := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &cg,
				},
			},
		}, nil)
		testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusPaymentApproved,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			Move: cgMoveInWrongStatus,
		})

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USCG",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})

	suite.Run("Filters out moves with no PPM shipment", func() {
		cg := models.AffiliationCOASTGUARD

		moveWithHHG := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
					Show:   &showMove,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &cg,
				},
			},
		}, nil)
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    moveWithHHG,
				LinkOnly: true,
			},
		}, nil)

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USCG",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUserSC.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListOrdersMarines() {
	suite.Run("does not return moves where the service member affiliation is Marines for non-USMC office user", func() {
		orderFetcher := NewOrderFetcher()
		marines := models.AffiliationMARINES
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		params := services.ListOrderParams{PerPage: swag.Int64(2), Page: swag.Int64(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListOrdersWithEmptyFields() {
	expectedOrder := factory.BuildOrder(suite.DB(), nil, nil)

	expectedOrder.Entitlement = nil
	expectedOrder.EntitlementID = nil
	expectedOrder.Grade = nil
	expectedOrder.OriginDutyLocation = nil
	expectedOrder.OriginDutyLocationID = nil
	suite.MustSave(&expectedOrder)

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    expectedOrder,
			LinkOnly: true,
		},
	}, nil)
	// Only orders with shipments are returned, so we need to add a shipment
	// to the move we just created
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
	// Add a second shipment to make sure we only return 1 order even if its
	// move has more than one shipment
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	orderFetcher := NewOrderFetcher()
	moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{PerPage: swag.Int64(1), Page: swag.Int64(1)})

	suite.FatalNoError(err)
	suite.Nil(moves)

}

func (suite *OrderServiceSuite) TestListOrdersWithPagination() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

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

func (suite *OrderServiceSuite) TestListOrdersWithSortOrder() {

	// SET UP: Service Members for sorting by Service Member Last Name and Branch
	// - We'll need two other service members to test the last name sort, Lea Spacemen and Leo Zephyer
	serviceMemberFirstName := "Lea"
	serviceMemberLastName := "Zephyer"
	affiliation := models.AffiliationNAVY
	edipi := "9999999999"
	var officeUser models.OfficeUser

	// SET UP: Dates for sorting by Requested Move Date
	// - We want dates 2 and 3 to sandwich requestedMoveDate1 so we can test that the min() query is working
	requestedMoveDate1 := time.Date(testdatagen.GHCTestYear, 02, 20, 0, 0, 0, 0, time.UTC)
	requestedMoveDate2 := time.Date(testdatagen.GHCTestYear, 03, 03, 0, 0, 0, 0, time.UTC)
	requestedMoveDate3 := time.Date(testdatagen.GHCTestYear, 01, 15, 0, 0, 0, 0, time.UTC)

	setupTestData := func() (models.Move, models.Move) {

		// CREATE EXPECTED MOVES
		expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			// Default New Duty Location name is Fort Gordon
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
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedMoveDate2,
			},
		})
		// Create a second shipment so we can test min() sort
		factory.BuildMTOShipmentWithMove(&expectedMove2, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate3,
				},
			},
		}, nil)
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		return expectedMove1, expectedMove2
	}

	orderFetcher := NewOrderFetcher()

	suite.Run("Sort by locator code", func() {
		expectedMove1, expectedMove2 := setupTestData()
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

	suite.Run("Sort by move status", func() {
		expectedMove1, expectedMove2 := setupTestData()
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

	suite.Run("Sort by service member affiliations", func() {
		expectedMove1, expectedMove2 := setupTestData()
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

	suite.Run("Sort by request move date", func() {
		setupTestData()
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

	suite.Run("Sort by submitted date (appearedInTooAt) in TOO queue ", func() {
		// Scenario: In order to sort the moves the submitted_at, service_counseling_completed_at, and approvals_requested_at are checked to which are the minimum
		// Expected: The moves appear in the order they are created below
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		now := time.Now()
		oneWeekAgo := now.AddDate(0, 0, -7)
		move1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				SubmittedAt: &oneWeekAgo,
			},
		})

		move2 := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		factory.BuildMTOShipmentWithMove(&move2, suite.DB(), nil, nil)
		move3 := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{})
		factory.BuildMTOShipmentWithMove(&move3, suite.DB(), nil, nil)

		params := services.ListOrderParams{Sort: models.StringPointer("appearedInTooAt"), Order: models.StringPointer("asc")}

		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal(moves[0].ID, move1.ID)
		suite.Equal(moves[1].ID, move2.ID)
		suite.Equal(moves[2].ID, move3.ID)
	})

	// MUST BE LAST, ADDS EXTRA MOVE
	suite.Run("Sort by service member last name", func() {
		setupTestData()

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

func (suite *OrderServiceSuite) TestListOrdersNeedingServicesCounselingWithPPMCloseoutColumnsSort() {
	defaultShipmentPickupPostalCode := "90210"
	setupTestData := func() models.OfficeUser {
		// Make an office user → GBLOC X
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

		// Ensure there's an entry connecting the default shipment pickup postal code with the office user's gbloc
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(),
			defaultShipmentPickupPostalCode,
			officeUser.TransportationOffice.Gbloc)

		return officeUser
	}
	orderFetcher := NewOrderFetcher()

	suite.Run("Sort by PPM closeout initiated", func() {
		officeUser := setupTestData()
		// Create a PPM submitted on April 1st
		closeoutInitiatedDate1 := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Gbloc: "KKFA"},
			},
		}, nil)

		ppm1 := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate1,
			},
			Move: models.Move{
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})

		// Create a PPM submitted on April 2nd
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		ppm2 := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SubmittedAt: &closeoutInitiatedDate2,
			},
			Move: models.Move{
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})

		// Sort by closeout initiated date (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutInitiated"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppm1.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppm2.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout initiated date (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutInitiated"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppm2.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppm1.Shipment.MoveTaskOrder.Locator, moves[1].Locator)
	})

	suite.Run("Sort by PPM closeout location", func() {
		officeUser := setupTestData()

		locationA := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "A",
				},
			}}, nil)
		ppmShipmentA := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				CloseoutOfficeID: &locationA.ID,
			},
		})
		locationB := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "B",
				},
			}}, nil)
		ppmShipmentB := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				CloseoutOfficeID: &locationB.ID,
			},
		})

		// Sort by closeout location (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutLocation"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[1].Locator)
	})

	suite.Run("Sort by destination duty location", func() {
		officeUser := setupTestData()

		dutyLocationA := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "A",
				},
			},
		}, nil)
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Gbloc: "KKFA"},
			},
		}, nil)

		ppmShipmentA := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				NewDutyLocationID: dutyLocationA.ID,
				NewDutyLocation:   dutyLocationA,
			},
			Move: models.Move{
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})
		dutyLocationB := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "B",
				},
			},
		}, nil)
		ppmShipmentB := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				NewDutyLocationID: dutyLocationB.ID,
				NewDutyLocation:   dutyLocationB,
			},
			Move: models.Move{
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})

		// Sort by destination duty location (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("destinationDutyLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by destination duty location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("destinationDutyLocation"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[1].Locator)
	})

	suite.Run("Sort by PPM type (full or partial)", func() {
		officeUser := setupTestData()
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Gbloc: "KKFA"},
			},
		}, nil)
		ppmShipmentPartial := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				PPMType:          swag.String("Partial"),
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})
		ppmShipmentFull := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				PPMType:          swag.String("FULL"),
				CloseoutOffice:   &closeoutOffice,
				CloseoutOfficeID: &closeoutOffice.ID,
			},
		})

		// Sort by PPM type (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmType"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentFull.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentPartial.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by PPM type (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmType"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentPartial.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentFull.Shipment.MoveTaskOrder.Locator, moves[1].Locator)
	})
}

func (suite *OrderServiceSuite) TestListOrdersNeedingServicesCounselingWithGBLOCSortFilter() {

	suite.Run("Filter by origin GBLOC", func() {

		// TESTCASE SCENARIO
		// Under test: OrderFetcher.ListOrders function
		// Mocked:     None
		// Set up:     We create 2 moves with different GBLOCs, KKFA and ZANY. Both moves require service counseling
		//             We create an office user with the GBLOC KKFA
		//             Then we request a list of moves sorted by GBLOC, ascending for service counseling
		// Expected outcome:
		//             We expect only the move that matches the counselors GBLOC - aka the KKFA move.

		// Create a services counselor (default GBLOC is KKFA)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		// Create a move with Origin KKFA, needs service couseling
		kkfaMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
		})

		// Create data for a second Origin ZANY
		dutyLocationAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Anchor 1212",
					City:           "Augusta",
					State:          "GA",
					PostalCode:     "89898",
					Country:        models.StringPointer("United States"),
				},
			},
		}, nil)

		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), dutyLocationAddress2.PostalCode, "ZANY")
		originDutyLocation2 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "Fort Sam Snap",
				},
			},
			{
				Model:    dutyLocationAddress2,
				LinkOnly: true,
			},
		}, nil)

		// Create a second move from the ZANY gbloc
		testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:  models.MoveStatusNeedsServiceCounseling,
				Locator: "ZZ1234",
			},
			Order: models.Order{
				OriginDutyLocation:   &originDutyLocation2,
				OriginDutyLocationID: &originDutyLocation2.ID,
			},
			OriginDutyLocation: originDutyLocation2,
		})

		// Setup and run the function under test requesting status NEEDS SERVICE COUNSELING
		orderFetcher := NewOrderFetcher()
		statuses := []string{"NEEDS SERVICE COUNSELING"}
		// Sort by origin GBLOC, filter by status
		params := services.ListOrderParams{Sort: swag.String("originGBLOC"), Order: swag.String("asc"), Status: statuses}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextForTest(), officeUser.ID, &params)

		// Expect only LKNQ move to be returned
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(kkfaMove.ID, moves[0].ID)
	})
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithNTSRelease() {
	// Make an NTS-Release shipment (and a move).  Should not have a pickup address.
	testdatagen.MakeNTSRMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	// Make a TOO user and the postal code to GBLOC link.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), tooOfficeUser.ID, &services.ListOrderParams{})

	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPM() {
	postalCode := "90210"
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		PPMShipment: models.PPMShipment{
			PickupPostalCode: postalCode,
		},
	})

	// Make a TOO user.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	// GBLOC for the below doesn't really matter, it just means the query for the moves passes the inner join in ListOrders
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), ppmShipment.PickupPostalCode, tooOfficeUser.TransportationOffice.Gbloc)

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), tooOfficeUser.ID, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithDeletedShipment() {
	postalCode := "90210"
	deletedAt := time.Now()
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		PPMShipment: models.PPMShipment{
			PickupPostalCode: postalCode,
		},
	})
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				DeletedAt: &deletedAt,
			},
		},
		{
			Model:    ppmShipment,
			LinkOnly: true,
		},
	}, nil)

	// Make a TOO user.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), tooOfficeUser.ID, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(0, moveCount)
	suite.Len(moves, 0)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithOneDeletedShipmentButOtherExists() {
	postalCode := "90210"
	deletedAt := time.Now()
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	// This shipment is created first, but later deleted
	ppmShipment1 := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		PPMShipment: models.PPMShipment{
			PickupPostalCode: postalCode,
			CreatedAt:        time.Now(),
		},
	})
	// This shipment is created after the first one, but not deleted
	testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		PPMShipment: models.PPMShipment{
			PickupPostalCode: postalCode,
			CreatedAt:        time.Now().Add(time.Minute * time.Duration(1)),
		},
	})
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				DeletedAt: &deletedAt,
			},
		},
		{
			Model:    ppmShipment1,
			LinkOnly: true,
		},
	}, nil)

	// Make a TOO user and the postal code to GBLOC link.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextForTest(), tooOfficeUser.ID, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}
