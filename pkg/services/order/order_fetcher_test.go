package order

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	moveservice "github.com/transcom/mymove/pkg/services/move"
	officeuserservice "github.com/transcom/mymove/pkg/services/office_user"
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
	setupTestData := func() (models.OfficeUser, models.Move, auth.Session) {

		// Make an office user → GBLOC X
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// Create a move with a shipment → GBLOC X
		move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		// Make a postal code and GBLOC → AGFM
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), agfmPostalCode, "AGFM")
		return officeUser, move, session
	}
	orderFetcher := NewOrderFetcher()

	suite.Run("returns moves", func() {
		// Under test: ListOrders
		// Mocked:           None
		// Set up:           Make 2 moves, one with a shipment and one without.
		//                   The shipment should have a pickup GBLOC that matches the office users transportation GBLOC
		//                   In other words, shipment should originate from same GBLOC as the office user
		// Expected outcome: Only the move with a shipment should be returned by ListOrders
		officeUser, expectedMove, session := setupTestData()

		// Create a Move without a shipment
		factory.BuildMove(suite.DB(), nil, nil)

		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{})

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
		officeUser, expectedMove, session := setupTestData()

		// This move's pickup GBLOC of the office user's GBLOC, so it should not be returned
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: agfmPostalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		}, nil)

		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{Page: models.Int64Pointer(1)})

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
		officeUser, expectedMove, session := setupTestData()

		params := services.ListOrderParams{}
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: models.BoolPointer(false),
				},
			},
		}, nil)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

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
		officeUser, expectedMove, session := setupTestData()
		expectedComboMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{})

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
		officeUser, _, session := setupTestData()

		// Create the airforce move
		airForce := models.AffiliationAIRFORCE
		airForceString := "AIR_FORCE"
		airForceMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &airForce,
				},
			},
		}, nil)
		// Filter by airforce move
		params := services.ListOrderParams{Branch: &airForceString, Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(airForceMove.ID, move.ID)

	})

	suite.Run("returns moves filtered submitted at", func() {
		// Under test: ListOrders
		// Set up:           Make 3 moves, with different submitted_at times, and search for a specific move
		// Expected outcome: Only the one move with the right date should be returned
		officeUser, _, session := setupTestData()

		// Move with specified timestamp
		submittedAt := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		expectedMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt: &submittedAt,
				},
			},
		}, nil)
		// Test edge cases (one day later)
		submittedAt2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt: &submittedAt2,
				},
			},
		}, nil)
		// Test edge cases (one second earlier)
		submittedAt3 := time.Date(2022, 03, 31, 23, 59, 59, 59, time.UTC)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt: &submittedAt3,
				},
			},
		}, nil)

		// Filter by submittedAt timestamp
		params := services.ListOrderParams{SubmittedAt: &submittedAt}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("returns moves filtered appeared in TOO at", func() {
		// Under test: ListOrders
		// Expected outcome: Only the three move with the right date should be returned
		officeUser, _, session := setupTestData()

		// Moves with specified timestamp
		specifiedDay := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		specifiedTimestamp1 := time.Date(2022, 04, 01, 1, 0, 0, 0, time.UTC)
		specifiedTimestamp2 := time.Date(2022, 04, 01, 23, 59, 59, 999999000, time.UTC) // the upper bound is 999999499 nanoseconds but the DB only stores microseconds

		matchingSubmittedAt := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt: &specifiedDay,
				},
			},
		}, nil)
		matchingSCCompletedAt := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ServiceCounselingCompletedAt: &specifiedTimestamp1,
				},
			},
		}, nil)
		matchingApprovalsRequestedAt := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ApprovalsRequestedAt: &specifiedTimestamp2,
				},
			},
		}, nil)
		// Test non dates matching
		nonMatchingDate1 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		nonMatchingDate2 := time.Date(2022, 03, 31, 23, 59, 59, 999999000, time.UTC) // the upper bound is 999999499 nanoseconds but the DB only stores microseconds
		nonMatchingDate3 := time.Date(2023, 04, 01, 0, 0, 0, 0, time.UTC)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt:                  &nonMatchingDate1,
					ServiceCounselingCompletedAt: &nonMatchingDate2,
					ApprovalsRequestedAt:         &nonMatchingDate3,
				},
			},
		}, nil)
		// Filter by AppearedInTOOAt timestamp
		params := services.ListOrderParams{AppearedInTOOAt: &specifiedDay}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

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
		officeUser, _, session := setupTestData()

		requestedPickupDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
				},
			},
		}, nil)
		requestedMoveDateString := createdMove.MTOShipments[0].RequestedPickupDate.Format("2006-01-02")
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			RequestedMoveDate: &requestedMoveDateString,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.Run("returns moves filtered by ppm type", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, with different ppm types, and search for both types
		// Expected outcome: search results should only include the move with the PPM type that was searched for
		postalCode := "50309"
		officeUser, partialPPMMove, session := setupTestData()
		suite.Equal("PARTIAL", *partialPPMMove.PPMType)
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Move{
					PPMType: models.StringPointer("FULL"),
					Locator: "FULLLL",
				},
			},
			{
				Model: models.Address{
					PostalCode: postalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		})
		fullPPMMove := ppmShipment.Shipment.MoveTaskOrder

		// Search for PARTIAL PPM moves
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			PPMType: models.StringPointer("PARTIAL"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(partialPPMMove.Locator, moves[0].Locator)

		// Search for FULL PPM moves
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			PPMType: models.StringPointer("FULL"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(fullPPMMove.Locator, moves[0].Locator)
	})

	suite.Run("returns moves filtered by ppm status", func() {
		// Under test: ListOrders
		// Set up:           Make 2 moves, with different ppm status, and search for both statues
		// Expected outcome: search results should only include the move with the PPM status that was searched for
		officeUser, partialPPMMove, session := setupTestData()
		suite.Equal("PARTIAL", *partialPPMMove.PPMType)
		postalCode := "50309"

		ppmShipmentNeedsCloseout := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Address{
					PostalCode: postalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		})
		// Search for PARTIAL PPM moves
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			PPMStatus: models.StringPointer("NEEDS_CLOSEOUT"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(moves[0].MTOShipments[0].PPMShipment.Status, ppmShipmentNeedsCloseout.Shipment.PPMShipment.Status)

		ppmShipmentWaiting := factory.BuildPPMShipmentThatNeedsToBeResubmitted(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Address{
					PostalCode: postalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		})
		// Search for FULL PPM moves
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			PPMStatus: models.StringPointer("WAITING_ON_CUSTOMER"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(moves[0].MTOShipments[0].PPMShipment.Status, ppmShipmentWaiting.Shipment.PPMShipment.Status)
	})

	suite.Run("returns moves filtered by closeout location", func() {
		// Under test: ListOrders
		// Set up:           Make a move with a closeout office. Search for that closeout office.
		// Expected outcome: Only the one ppmShipment with the right closeout office should be returned
		officeUser, _, session := setupTestData()

		ftBragg := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "Ft Bragg",
				},
			},
		}, nil)
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Move{
					CloseoutOfficeID: &ftBragg.ID,
				},
			},
		})

		// Search should be case insensitive and allow partial matches
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
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
		postalCode := "50309"
		officeUser, _, session := setupTestData()

		// Create a PPM submitted on April 1st
		closeoutInitiatedDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdPPM := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate,
				},
			},
			{
				Model: models.Address{
					PostalCode: postalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		})

		// Create a PPM submitted on April 2nd
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		createdPPM2 := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate2,
				},
			},
		})

		// Search for PPMs submitted on April 1st
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
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
		postalCode := "50309"
		officeUser, _, session := setupTestData()

		// Create a PPM submitted on April 1st
		closeoutInitiatedDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		createdPPM := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate,
				},
			},
			{
				Model: models.Address{
					PostalCode: postalCode,
				},
				Type: &factory.Addresses.PickupAddress,
			},
		})
		// Add another PPM for the same move submitted on April 1st
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)

		factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate2,
					Status:      models.PPMShipmentStatusNeedsCloseout,
				},
			},
			{
				Model:    createdPPM.Shipment.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		// Search for PPMs submitted on April 1st
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate,
		})
		suite.Empty(moves)
		suite.FatalNoError(err)

		// Search for PPMs submitted on April 2nd
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate2,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(createdPPM.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})
}
func (suite *OrderServiceSuite) TestListOrderWithAssignedUserSingle() {
	// Under test: ListOrders
	// Set up:           Make a move, assign one to an SC office user
	// Expected outcome: Only the one move with the assigned user should be returned
	assignedOfficeUserUpdater := moveservice.NewAssignedOfficeUserUpdater(moveservice.NewMoveFetcher())
	scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
	var orderFetcherTest orderFetcher
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           scUser.User.Roles,
		OfficeUserID:    scUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	appCtx := suite.AppContextWithSessionForTest(&session)

	createdMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	createdMove.SCAssignedID = &scUser.ID
	createdMove.SCAssignedUser = &scUser
	_, updateError := assignedOfficeUserUpdater.UpdateAssignedOfficeUser(appCtx, createdMove.ID, &scUser, roles.RoleTypeServicesCounselor)

	moves, _, err := orderFetcherTest.ListOrders(suite.AppContextWithSessionForTest(&session), scUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
		SCAssignedUser: &scUser.LastName,
	})

	suite.FatalNoError(err)
	suite.FatalNoError(updateError)
	suite.Equal(1, len(moves))
	suite.Equal(moves[0].SCAssignedID, createdMove.SCAssignedID)
	suite.Equal(createdMove.SCAssignedUser.ID, moves[0].SCAssignedUser.ID)
	suite.Equal(createdMove.SCAssignedUser.FirstName, moves[0].SCAssignedUser.FirstName)
	suite.Equal(createdMove.SCAssignedUser.LastName, moves[0].SCAssignedUser.LastName)
}
func (suite *OrderServiceSuite) TestListOrdersUSMCGBLOC() {
	orderFetcher := NewOrderFetcher()

	suite.Run("returns USMC order for USMC office user", func() {
		marines := models.AffiliationMARINES
		// It doesn't matter what the Origin GBLOC is for the move. Only the Marines
		// affiliation matters for office users who are tied to the USMC GBLOC.
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &marines,
				},
			},
		}, nil)
		// Create move where service member has the default ARMY affiliation
		factory.BuildMoveWithShipment(suite.DB(), nil, nil)

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
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserOooRah.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationMARINES, *moves[0].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationARMY, *moves[0].Orders.ServiceMember.Affiliation)
	})
}

func getMoveNeedsServiceCounseling(suite *OrderServiceSuite, showMove bool, affiliation models.ServiceMemberAffiliation) models.Move {
	nonCloseoutMove := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
				Show:   &showMove,
			},
		},
		{
			Model: models.ServiceMember{
				Affiliation: &affiliation,
			},
		},
	}, nil)

	return nonCloseoutMove
}

func getSubmittedMove(suite *OrderServiceSuite, showMove bool, affiliation models.ServiceMemberAffiliation) models.Move {
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
				Show:   &showMove,
			},
		},
		{
			Model: models.ServiceMember{
				Affiliation: &affiliation,
			},
		},
	}, nil)
	return move
}

func buildPPMShipmentNeedsCloseout(suite *OrderServiceSuite, move models.Move) models.PPMShipment {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	return ppm
}

func buildPPMShipmentDraft(suite *OrderServiceSuite, move models.Move) models.PPMShipment {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusDraft,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	return ppm
}

func buildPPMShipmentCloseoutComplete(suite *OrderServiceSuite, move models.Move) models.PPMShipment {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusCloseoutComplete,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	return ppm
}
func (suite *OrderServiceSuite) TestListOrdersPPMCloseoutForArmyAirforce() {
	orderFetcher := NewOrderFetcher()

	var session auth.Session

	suite.Run("office user in normal GBLOC should only see non-Navy/Marines/CoastGuard moves that need closeout in closeout tab", func() {
		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		move := getMoveNeedsServiceCounseling(suite, true, models.AffiliationARMY)
		buildPPMShipmentNeedsCloseout(suite, move)

		afMove := getMoveNeedsServiceCounseling(suite, true, models.AffiliationAIRFORCE)
		buildPPMShipmentDraft(suite, afMove)

		cgMove := getMoveNeedsServiceCounseling(suite, true, models.AffiliationCOASTGUARD)
		buildPPMShipmentNeedsCloseout(suite, cgMove)

		params := services.ListOrderParams{PerPage: models.Int64Pointer(9), Page: models.Int64Pointer(1), NeedsPPMCloseout: models.BoolPointer(true), Status: []string{string(models.MoveStatusNeedsServiceCounseling)}}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(move.Locator, moves[0].Locator)
	})

	suite.Run("office user in normal GBLOC should not see moves that require closeout in counseling tab", func() {
		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		closeoutMove := getMoveNeedsServiceCounseling(suite, true, models.AffiliationARMY)
		buildPPMShipmentCloseoutComplete(suite, closeoutMove)

		// PPM moves that are not in one of the closeout statuses
		nonCloseoutMove := getMoveNeedsServiceCounseling(suite, true, models.AffiliationAIRFORCE)
		buildPPMShipmentDraft(suite, nonCloseoutMove)

		params := services.ListOrderParams{PerPage: models.Int64Pointer(9), Page: models.Int64Pointer(1), NeedsPPMCloseout: models.BoolPointer(false), Status: []string{string(models.MoveStatusNeedsServiceCounseling)}}

		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(nonCloseoutMove.Locator, moves[0].Locator)
	})
}

func (suite *OrderServiceSuite) TestListOrdersPPMCloseoutForNavyCoastGuardAndMarines() {
	orderFetcher := NewOrderFetcher()

	suite.Run("returns Navy order for NAVY office user when there's a ppm shipment in closeout", func() {
		// It doesn't matter what the Origin GBLOC is for the move. Only the navy
		// affiliation matters for SC  who are tied to the NAVY GBLOC.
		move := getSubmittedMove(suite, true, models.AffiliationNAVY)
		buildPPMShipmentNeedsCloseout(suite, move)

		cgMove := getSubmittedMove(suite, true, models.AffiliationCOASTGUARD)
		buildPPMShipmentNeedsCloseout(suite, cgMove)

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "NAVY",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(9), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationNAVY, *moves[0].Orders.ServiceMember.Affiliation)

	})

	suite.Run("returns TVCB order for TVCB office user when there's a ppm shipment in closeout", func() {
		// It doesn't matter what the Origin GBLOC is for the move. Only the marines
		// affiliation matters for SC  who are tied to the TVCB GBLOC.
		move := getSubmittedMove(suite, true, models.AffiliationMARINES)
		buildPPMShipmentNeedsCloseout(suite, move)

		nonMarineMove := getSubmittedMove(suite, true, models.AffiliationARMY)
		buildPPMShipmentNeedsCloseout(suite, nonMarineMove)

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "TVCB",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationMARINES, *moves[0].Orders.ServiceMember.Affiliation)

	})

	suite.Run("returns coast guard order for USCG office user when there's a ppm shipment in closeout and filters out non coast guard moves", func() {
		// It doesn't matter what the Origin GBLOC is for the move. Only the coast guard
		// affiliation matters for SC  who are tied to the USCG GBLOC.
		move := getSubmittedMove(suite, true, models.AffiliationCOASTGUARD)
		buildPPMShipmentNeedsCloseout(suite, move)

		armyMove := getSubmittedMove(suite, true, models.AffiliationARMY)
		buildPPMShipmentNeedsCloseout(suite, armyMove)

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USCG",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(models.AffiliationCOASTGUARD, *moves[0].Orders.ServiceMember.Affiliation)
	})

	suite.Run("Filters out moves with PPM shipments not in the status of NeedsApproval", func() {

		cgMoveInWrongStatus := getSubmittedMove(suite, true, models.AffiliationCOASTGUARD)
		buildPPMShipmentCloseoutComplete(suite, cgMoveInWrongStatus)

		officeUserSC := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USCG",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		var session auth.Session
		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})

	suite.Run("Filters out moves with no PPM shipment", func() {

		moveWithHHG := getSubmittedMove(suite, true, models.AffiliationCOASTGUARD)
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

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUserSC.User.Roles,
			OfficeUserID:    officeUserSC.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUserSC.ID, roles.RoleTypeServicesCounselor, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListOrdersMarines() {
	suite.Run("does not return moves where the service member affiliation is Marines for non-USMC office user", func() {
		orderFetcher := NewOrderFetcher()
		marines := models.AffiliationMARINES
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &marines,
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{PerPage: models.Int64Pointer(2), Page: models.Int64Pointer(1)}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

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
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           officeUser.User.Roles,
		OfficeUserID:    officeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	orderFetcher := NewOrderFetcher()
	moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{PerPage: models.Int64Pointer(1), Page: models.Int64Pointer(1)})

	suite.FatalNoError(err)
	suite.Nil(moves)

}

func (suite *OrderServiceSuite) TestListOrdersWithPagination() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           officeUser.User.Roles,
		OfficeUserID:    officeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	for i := 0; i < 2; i++ {
		factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	}

	orderFetcher := NewOrderFetcher()
	params := services.ListOrderParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(1)}
	moves, count, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

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

	setupTestData := func() (models.Move, models.Move, auth.Session) {

		// CREATE EXPECTED MOVES
		expectedMove1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{ // Default New Duty Location name is Fort Eisenhower
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "AA1234",
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate1,
				},
			},
		}, nil)
		expectedMove2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "TTZ123",
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &affiliation,
					FirstName:   &serviceMemberFirstName,
					Edipi:       &edipi,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate2,
				},
			},
		}, nil)
		// Create a second shipment so we can test min() sort
		factory.BuildMTOShipmentWithMove(&expectedMove2, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate3,
				},
			},
		}, nil)
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		return expectedMove1, expectedMove2, session
	}

	orderFetcher := NewOrderFetcher()

	suite.Run("Sort by locator code", func() {
		expectedMove1, expectedMove2, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		params = services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Locator, moves[0].Locator)
		suite.Equal(expectedMove1.Locator, moves[1].Locator)
	})

	suite.Run("Sort by move status", func() {
		expectedMove1, expectedMove2, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Status, moves[0].Status)
		suite.Equal(expectedMove2.Status, moves[1].Status)

		params = services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Status, moves[0].Status)
		suite.Equal(expectedMove1.Status, moves[1].Status)
	})

	suite.Run("Sort by service member affiliations", func() {
		expectedMove1, expectedMove2, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)
	})

	suite.Run("Sort by request move date", func() {
		_, _, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("requestedMoveDate"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(2, len(moves[0].MTOShipments)) // the move with two shipments has the earlier date
		suite.Equal(1, len(moves[1].MTOShipments))
		// NOTE: You have to use Jan 02, 2006 as the example for date/time formatting in Go
		suite.Equal(requestedMoveDate1.Format("2006/01/02"), moves[1].MTOShipments[0].RequestedPickupDate.Format("2006/01/02"))

		params = services.ListOrderParams{Sort: models.StringPointer("requestedMoveDate"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
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
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}
		now := time.Now()
		oneWeekAgo := now.AddDate(0, 0, -7)
		move1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					SubmittedAt: &oneWeekAgo,
				},
			},
		}, nil)
		move2 := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentWithMove(&move2, suite.DB(), nil, nil)
		move3 := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentWithMove(&move3, suite.DB(), nil, nil)

		params := services.ListOrderParams{Sort: models.StringPointer("appearedInTooAt"), Order: models.StringPointer("asc")}

		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal(moves[0].ID, move1.ID)
		suite.Equal(moves[1].ID, move2.ID)
		suite.Equal(moves[2].ID, move3.ID)
	})

	// MUST BE LAST, ADDS EXTRA MOVE
	suite.Run("Sort by service member last name", func() {
		_, _, session := setupTestData()

		// Last name sort is the only one that needs 3 moves for a complete test, so add that here at the end
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{ // Leo Zephyer
					LastName: &serviceMemberLastName,
				},
			},
		}, nil)
		params := services.ListOrderParams{Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Spacemen, Lea", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Zephyer, Leo", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)

		params = services.ListOrderParams{Sort: models.StringPointer("customerName"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Zephyer, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Lea", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)
	})
}

func getTransportationOffice(suite *OrderServiceSuite, name string) models.TransportationOffice {
	trasportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name: name,
			},
		}}, nil)
	return trasportationOffice
}

func getPPMShipmentWithCloseoutOfficeNeedsCloseout(suite *OrderServiceSuite, closeoutOffice models.TransportationOffice) models.PPMShipment {
	ppm := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
		{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		},
	})
	return ppm
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

	var session auth.Session

	suite.Run("Sort by PPM closeout initiated", func() {
		officeUser := setupTestData()
		// Create a PPM submitted on April 1st
		closeoutInitiatedDate1 := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Gbloc: "KKFA"},
			},
		}, nil)

		ppm1 := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate1,
				},
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})

		// Create a PPM submitted on April 2nd
		closeoutInitiatedDate2 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		ppm2 := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.PPMShipment{
					SubmittedAt: &closeoutInitiatedDate2,
				},
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})

		// Sort by closeout initiated date (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutInitiated"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppm1.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppm2.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout initiated date (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
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

		locationA := getTransportationOffice(suite, "A")
		ppmShipmentA := getPPMShipmentWithCloseoutOfficeNeedsCloseout(suite, locationA)

		locationB := getTransportationOffice(suite, "B")
		ppmShipmentB := getPPMShipmentWithCloseoutOfficeNeedsCloseout(suite, locationB)

		// Sort by closeout location (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
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

		ppmShipmentA := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model:    dutyLocationA,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation,
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})
		dutyLocationB := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "B",
				},
			},
		}, nil)
		ppmShipmentB := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model:    dutyLocationB,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation,
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})

		// Sort by destination duty location (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("destinationDutyLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by destination duty location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
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
		ppmShipmentPartial := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Move{
					PPMType: models.StringPointer("Partial"),
				},
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})
		ppmShipmentFull := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
			{
				Model: models.Move{
					PPMType: models.StringPointer("FULL"),
				},
			},
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		})

		// Sort by PPM type (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmType"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentFull.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentPartial.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by PPM type (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmType"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentPartial.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentFull.Shipment.MoveTaskOrder.Locator, moves[1].Locator)
	})
	suite.Run("Sort by PPM status", func() {
		officeUser := setupTestData()
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Gbloc: "KKFA"},
			},
		}, nil)
		ppmShipmentNeedsCloseout := getPPMShipmentWithCloseoutOfficeNeedsCloseout(suite, closeoutOffice)

		// Sort by PPM type (ascending)
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmStatus"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(ppmShipmentNeedsCloseout.Status, moves[0].MTOShipments[0].PPMShipment.Status)

		// Sort by PPM type (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmStatus"),
			Order:            models.StringPointer("desc"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(ppmShipmentNeedsCloseout.Status, moves[0].MTOShipments[0].PPMShipment.Status)
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
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// Create a move with Origin KKFA, needs service couseling
		kkfaMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
		}, nil)
		// Create data for a second Origin ZANY
		dutyLocationAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Anchor 1212",
					City:           "Fort Eisenhower",
					State:          "GA",
					PostalCode:     "89898",
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
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusNeedsServiceCounseling,
					Locator: "ZZ1234",
				},
			},
			{
				Model:    originDutyLocation2,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		// Setup and run the function under test requesting status NEEDS SERVICE COUNSELING
		orderFetcher := NewOrderFetcher()
		statuses := []string{"NEEDS SERVICE COUNSELING"}
		// Sort by origin GBLOC, filter by status
		params := services.ListOrderParams{Sort: models.StringPointer("originGBLOC"), Order: models.StringPointer("asc"), Status: statuses}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &params)

		// Expect only LKNQ move to be returned
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(kkfaMove.ID, moves[0].ID)
	})
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithNTSRelease() {
	// Make an NTS-Release shipment (and a move).  Should not have a pickup address.
	factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
		},
	}, nil)
	// Make a TOO user and the postal code to GBLOC link.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           tooOfficeUser.User.Roles,
		OfficeUserID:    tooOfficeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{})

	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPM() {
	postalCode := "50309"
	partialPPMType := models.MovePPMTypePARTIAL

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				ID: uuid.UUID{uuid.V4},
			},
		},
		{
			Model: models.Move{
				Status:  models.MoveStatusAPPROVED,
				PPMType: &partialPPMType,
			},
		},
		{
			Model: models.Address{
				PostalCode: postalCode,
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)
	// Make a TOO user.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           tooOfficeUser.User.Roles,
		OfficeUserID:    tooOfficeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	// GBLOC for the below doesn't really matter, it just means the query for the moves passes the inner join in ListOrders
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), ppmShipment.PickupAddress.PostalCode, tooOfficeUser.TransportationOffice.Gbloc)

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersWithViewAsGBLOCParam() {
	var hqOfficeUser models.OfficeUser
	var hqOfficeUserAGFM models.OfficeUser

	requestedMoveDate1 := time.Date(testdatagen.GHCTestYear, 02, 20, 0, 0, 0, 0, time.UTC)
	requestedMoveDate2 := time.Date(testdatagen.GHCTestYear, 03, 03, 0, 0, 0, 0, time.UTC)

	setupTestData := func() (models.Move, models.Move, models.MTOShipment, auth.Session, auth.Session) {
		// CREATE EXPECTED MOVES
		expectedMove1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{ // Default New Duty Location name is Fort Eisenhower
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "AA1234",
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate1,
				},
			},
		}, nil)
		expectedMove2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "TTZ123",
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate2,
				},
			},
		}, nil)

		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "06001", "AGFM")

		expectedShipment3 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Fort Punxsutawney",
					Gbloc: "AGFM",
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model: models.Address{
					PostalCode: "06001",
				},
				Type: &factory.Addresses.PickupAddress,
			},
		}, nil)

		hqOfficeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeHQ})
		hqSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           hqOfficeUser.User.Roles,
			OfficeUserID:    hqOfficeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		hqOfficeUserAGFM = factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Scott AFB",
					Gbloc: "AGFM",
				},
			},
		}, []roles.RoleType{roles.RoleTypeHQ})
		hqSessionAGFM := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           hqOfficeUserAGFM.User.Roles,
			OfficeUserID:    hqOfficeUserAGFM.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		return expectedMove1, expectedMove2, expectedShipment3, hqSession, hqSessionAGFM
	}

	orderFetcher := NewOrderFetcher()

	suite.Run("Sort by locator code", func() {
		expectedMove1, expectedMove2, expectedShipment3, hqSession, hqSessionAGFM := setupTestData()

		// Request as an HQ user with their default GBLOC, KKFA
		params := services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		// Expect the same results with a ViewAsGBLOC that equals the user's default GBLOC, KKFA
		params = services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc"), ViewAsGBLOC: models.StringPointer("KKFA")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		// Expect the AGFM move when using the ViewAsGBLOC param set to AGFM
		params = services.ListOrderParams{ViewAsGBLOC: models.StringPointer("AGFM")}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(expectedShipment3.ID, moves[0].MTOShipments[0].ID)

		// Expect the same results without a ViewAsGBLOC for a user whose default GBLOC is AGFM
		params = services.ListOrderParams{}
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&hqSessionAGFM), hqOfficeUserAGFM.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(expectedShipment3.ID, moves[0].MTOShipments[0].ID)
	})
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithDeletedShipment() {
	postalCode := "50309"
	deletedAt := time.Now()
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				PostalCode: postalCode,
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)
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
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           tooOfficeUser.User.Roles,
		OfficeUserID:    tooOfficeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{Status: []string{string(models.MoveStatusSUBMITTED)}})
	suite.FatalNoError(err)
	suite.Equal(0, moveCount)
	suite.Len(moves, 0)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithOneDeletedShipmentButOtherExists() {
	postalCode := "50309"
	deletedAt := time.Now()
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	// This shipment is created first, but later deleted
	ppmShipment1 := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				CreatedAt: time.Now(),
			},
		},
		{
			Model: models.Address{
				PostalCode: postalCode,
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)
	// This shipment is created after the first one, but not deleted
	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				CreatedAt: time.Now().Add(time.Minute * time.Duration(1)),
			},
		},
		{
			Model: models.Address{
				PostalCode: postalCode,
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)
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
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           tooOfficeUser.User.Roles,
		OfficeUserID:    tooOfficeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	orderFetcher := NewOrderFetcher()
	moves, moveCount, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListAllOrderLocations() {
	suite.Run("returns a list of all order locations in the current users queue", func() {
		orderFetcher := NewOrderFetcher()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		params := services.ListOrderParams{}
		moves, err := orderFetcher.ListAllOrderLocations(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListOrdersFilteredByCustomerName() {
	serviceMemberFirstName := "Margaret"
	serviceMemberLastName := "Starlight"
	edipi := "9999999998"
	var officeUser models.OfficeUser
	var session auth.Session

	requestedMoveDate1 := time.Date(testdatagen.GHCTestYear, 05, 20, 0, 0, 0, 0, time.UTC)
	requestedMoveDate2 := time.Date(testdatagen.GHCTestYear, 07, 03, 0, 0, 0, 0, time.UTC)

	suite.PreloadData(func() {
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "AA1235",
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate1,
				},
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "TTZ125",
				},
			},
			{
				Model: models.ServiceMember{
					FirstName: &serviceMemberFirstName,
					Edipi:     &edipi,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedMoveDate2,
				},
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{ // Leo Zephyer
					LastName: &serviceMemberLastName,
				},
			},
		}, nil)
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}
	})

	orderFetcher := NewOrderFetcher()

	suite.Run("list moves by customer name - full name (last, first)", func() {
		// Search "Spacemen, Margaret"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Spacemen, Margaret"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - full name (first last)", func() {
		// Search "Margaret Spacemen"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Margaret Spacemen"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial last (multiple)", func() {
		// Search "space"
		params := services.ListOrderParams{CustomerName: models.StringPointer("space"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Margaret", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial last (single)", func() {
		// Search "Light"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Light"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Starlight, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial first", func() {
		// Search "leo"
		params := services.ListOrderParams{CustomerName: models.StringPointer("leo"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Starlight, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial matching within first or last", func() {
		// Search "ar"
		params := services.ListOrderParams{CustomerName: models.StringPointer("ar"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Starlight, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - empty", func() {
		// Search "johnny"
		params := services.ListOrderParams{CustomerName: models.StringPointer("johnny"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params)
		suite.NoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListAllOrderLocationsWithViewAsGBLOCParam() {
	suite.Run("returns a list of all order locations in the current users queue", func() {
		orderFetcher := NewOrderFetcher()
		officeUserFetcher := officeuserservice.NewOfficeUserFetcherPop()
		movesContainOriginDutyLocation := func(moves models.Moves, keyword string) func() (success bool) {
			return func() (success bool) {
				for _, record := range moves {
					if strings.Contains(record.Orders.OriginDutyLocation.Name, keyword) {
						return true
					}
				}
				return false
			}
		}

		// Create SC office user with a default transportation office in the AGFM GBLOC
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Fort Punxsutawney",
					Gbloc: "AGFM",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		// Add a secondary GBLOC to the above office user, this should default to KKFA
		factory.BuildAlternateTransportationOfficeAssignment(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					ID: officeUser.ID,
				},
				LinkOnly: true,
			},
		}, nil)
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// Create three default moves with shipment, should be in KKFA and have the status SUBMITTED
		KKFAMove1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		KKFAMove2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		KKFAMove3 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		// Create fourth move with the same origin duty location as one of the above
		KKFAMove4 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ID: KKFAMove3.Orders.OriginDutyLocation.ID,
				},
				Type:     &factory.DutyLocations.OriginDutyLocation,
				LinkOnly: true,
			},
		}, nil)

		// Create AGFM Move
		AGFM := "AGFM"
		AGFMTransportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Fort Punxsutawney",
					ID:    uuid.Must(uuid.NewV4()),
					Gbloc: AGFM,
				},
			},
			{
				Model: models.Address{
					PostalCode: "15767",
				},
				Type: &factory.Addresses.DutyLocationAddress,
			},
		}, nil)
		AGFMDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    AGFMTransportationOffice,
				LinkOnly: true,
			},
		}, nil)
		AGFMOrders := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    AGFMDutyLocation,
				Type:     &factory.DutyLocations.OriginDutyLocation,
				LinkOnly: true,
			},
			{
				Model: models.Order{
					OriginDutyLocationGBLOC: &AGFM,
				},
			},
		}, nil)
		AGFMMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    AGFMOrders,
				LinkOnly: true,
			},
		}, nil)
		// Create one AGFM shipment, this should result in a move as well
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Fort Punxsutawney",
					Gbloc: AGFM,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model: models.Address{
					PostalCode: "15767",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model:    AGFMMove,
				LinkOnly: true,
				Type:     &factory.Move,
			},
		}, nil)

		officeUser, _ = officeUserFetcher.FetchOfficeUserByIDWithTransportationOfficeAssignments(suite.AppContextForTest(), officeUser.ID)

		// Confirm office user has the desired transportation office assignments
		suite.Equal("AGFM", officeUser.TransportationOffice.Gbloc)
		if officeUser.TransportationOfficeAssignments[0].TransportationOffice.Gbloc == "AGFM" {
			suite.Equal("AGFM", officeUser.TransportationOfficeAssignments[0].TransportationOffice.Gbloc)
			suite.Equal(true, *officeUser.TransportationOfficeAssignments[0].PrimaryOffice)
			suite.Equal("KKFA", officeUser.TransportationOfficeAssignments[1].TransportationOffice.Gbloc)
			suite.Equal(false, *officeUser.TransportationOfficeAssignments[1].PrimaryOffice)
		} else {
			suite.Equal("KKFA", officeUser.TransportationOfficeAssignments[0].TransportationOffice.Gbloc)
			suite.Equal(false, *officeUser.TransportationOfficeAssignments[0].PrimaryOffice)
			suite.Equal("AGFM", officeUser.TransportationOfficeAssignments[1].TransportationOffice.Gbloc)
			suite.Equal(true, *officeUser.TransportationOfficeAssignments[1].PrimaryOffice)
		}

		// Confirm the factory created moves have the desired GBLOCS, 4x KKFA, 1x AGFM
		suite.Equal("AGFM", *AGFMMove.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove1.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove2.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove3.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove4.Orders.OriginDutyLocationGBLOC)

		// Fetch and check default GBLOC
		params := services.ListOrderParams{}
		AGFMmoves, err := orderFetcher.ListAllOrderLocations(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(AGFMmoves))
		suite.Equal("AGFM", *AGFMmoves[0].Orders.OriginDutyLocationGBLOC)
		suite.Condition(movesContainOriginDutyLocation(AGFMmoves, AGFMMove.Orders.OriginDutyLocation.Name), "Should contain first AGFM move's origin duty location")

		// Fetch and check secondary GBLOC
		KKFA := "KKFA"
		params = services.ListOrderParams{
			ViewAsGBLOC: &KKFA,
		}
		KKFAmoves, err := orderFetcher.ListAllOrderLocations(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		// This value should be updated to 3 if ListAllOrderLocations is updated to return distinct locations
		suite.Equal(4, len(KKFAmoves))

		suite.Equal("KKFA", *KKFAmoves[0].Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAmoves[1].Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAmoves[2].Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAmoves[3].Orders.OriginDutyLocationGBLOC)

		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove1.Orders.OriginDutyLocation.Name), "Should contain first KKFA move's origin duty location")
		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove2.Orders.OriginDutyLocation.Name), "Should contain second KKFA move's origin duty location")
		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove3.Orders.OriginDutyLocation.Name), "Should contain third KKFA move's origin duty location")
		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove4.Orders.OriginDutyLocation.Name), "Should contain third KKFA move's origin duty location")
	})
}

func (suite *OrderServiceSuite) TestOriginDutyLocationFilter() {
	var session auth.Session
	var expectedMove models.Move
	var officeUser models.OfficeUser
	orderFetcher := NewOrderFetcher()
	suite.PreloadData(func() {
		setupTestData := func() (models.OfficeUser, models.Move, auth.Session) {
			officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
			session := auth.Session{
				ApplicationName: auth.OfficeApp,
				Roles:           officeUser.User.Roles,
				OfficeUserID:    officeUser.ID,
				IDToken:         "fake_token",
				AccessToken:     "fakeAccessToken",
			}
			move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
			return officeUser, move, session
		}
		officeUser, expectedMove, session = setupTestData()
	})
	locationName := expectedMove.Orders.OriginDutyLocation.Name
	suite.Run("Returns orders matching full originDutyLocation name filter", func() {
		expectedMoves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{OriginDutyLocation: strings.Split(locationName, " ")})
		suite.NoError(err)
		suite.Equal(1, len(expectedMoves))
		suite.Equal(locationName, string(expectedMoves[0].Orders.OriginDutyLocation.Name))
	})

	suite.Run("Returns orders matching partial originDutyLocation name filter", func() {
		//Split the location name and retrieve a substring (first string) for the search param
		partialParamSearch := strings.Split(locationName, " ")[0]
		expectedMoves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{OriginDutyLocation: strings.Split(partialParamSearch, " ")})
		suite.NoError(err)
		suite.Equal(1, len(expectedMoves))
		suite.Equal(locationName, string(expectedMoves[0].Orders.OriginDutyLocation.Name))
	})
}
