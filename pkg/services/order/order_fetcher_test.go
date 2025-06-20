package order

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservice "github.com/transcom/mymove/pkg/services/move"
	officeuserservice "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestFetchOrder() {
	expectedMove := factory.BuildMove(suite.DB(), nil, nil)
	expectedOrder := expectedMove.Orders
	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

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
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)
	order, err := orderFetcher.FetchOrder(suite.AppContextForTest(), expectedOrder.ID)

	suite.FatalNoError(err)
	suite.Nil(order.Entitlement)
	suite.Nil(order.OriginDutyLocation)
	suite.Nil(order.Grade)
}

func (suite *OrderServiceSuite) TestListOrders() {
	waf := entitlements.NewWeightAllotmentFetcher()
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
	orderFetcher := NewOrderFetcher(waf)

	suite.Run("returns moves", func() {
		// Under test: ListOriginRequestsOrders
		// Mocked:           None
		// Set up:           Make 2 moves, one with a shipment and one without.
		//                   The shipment should have a pickup GBLOC that matches the office users transportation GBLOC
		//                   In other words, shipment should originate from same GBLOC as the office user
		// Expected outcome: Only the move with a shipment should be returned by ListOriginRequestsOrders
		officeUser, expectedMove, session := setupTestData()

		// Create a Move without a shipment
		factory.BuildMove(suite.DB(), nil, nil)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

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
		suite.Equal(expectedMove.Orders.ServiceMemberID, move.Orders.ServiceMember.ID)
		suite.NotNil(move.Orders.NewDutyLocation)
		suite.Equal(expectedMove.Orders.OriginDutyLocation.ID, move.Orders.OriginDutyLocation.ID)
		suite.NotNil(move.Orders.OriginDutyLocation)
		suite.Equal(expectedMove.Orders.OriginDutyLocation.Address.StreetAddress1, move.Orders.OriginDutyLocation.Address.StreetAddress1)
	})

	suite.Run("returns moves with all required locked information", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), agfmPostalCode, "AGFM")
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		now := time.Now()

		// build a move that's locked
		lockedMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					LockedByOfficeUserID: &tooUser.ID,
					LockExpiresAt:        &now,
					Show:                 models.BoolPointer(true),
				},
			},
		}, nil)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{Page: models.Int64Pointer(1)})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)

		// Check that move matches
		move := moves[0]
		suite.Equal(move.Locator, lockedMove.Locator)
		suite.NotNil(move.LockedByOfficeUserID)
		suite.NotNil(move.LockExpiresAt)
	})

	suite.Run("returns moves filtered by GBLOC", func() {
		// Under test: ListOriginRequestsOrders
		// Set up:           Make 2 moves, one with a pickup GBLOC that matches the office users transportation GBLOC
		//                   (which is done in setupTestData) and one with a pickup GBLOC that doesn't
		// Expected outcome: Only the move with the correct GBLOC should be returned by ListOriginRequestsOrders
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

		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{Page: models.Int64Pointer(1)})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("only returns visible moves (where show = True)", func() {
		// Under test: ListOriginRequestsOrders
		// Set up:           Make 2 moves, one correctly setup in setupTestData (show = True)
		//                   and one with show = False
		// Expected outcome: Only the move with show = True should be returned by ListOriginRequestsOrders
		officeUser, expectedMove, session := setupTestData()

		params := services.ListOrderParams{}
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: models.BoolPointer(false),
				},
			},
		}, nil)
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(expectedMove.ID, move.ID)

	})

	suite.Run("includes combo hhg and ppm moves", func() {
		// Under test: ListOriginRequestsOrders
		// Set up:           Make 2 moves, one default move setup in setupTestData (show = True)
		//                   and one a combination HHG and PPM move and make sure it's included
		// Expected outcome: Both moves should be returned by ListOriginRequestsOrders
		officeUser, expectedMove, session := setupTestData()
		expectedComboMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

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

	suite.Run("returns moves filtered by ppm_type", func() {
		// Set up:       Make 2 moves, one with ppm_type 'PARTIAL' and another with 'FULL'
		//               The 'FULL' type should be filtered out if the origin_duty_location provides services counseling
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// move with ppm_type 'PARTIAL'
		partial := models.MovePPMTypePARTIAL
		movePartial := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PPMType: &partial,
				},
			},
		}, nil)

		// move with ppm_type 'FULL' and origin_duty_location that provides services counseling
		full := models.MovePPMTypeFULL
		moveFull := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PPMType: &full,
				},
			},
			{
				Model: models.DutyLocation{
					ProvidesServicesCounseling: true,
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})
		suite.FatalNoError(err)

		// the returned moves should not include the 'FULL' type because the origin duty location provides services counseling
		suite.Equal(1, len(moves)) // only the 'PARTIAL' move should be returned
		suite.Equal(movePartial.Locator, moves[0].Locator)
		suite.NotEqual(moveFull.Locator, moves[0].Locator)
	})

	suite.Run("returns moves filtered by service member affiliation", func() {
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		move := moves[0]
		suite.Equal(airForceMove.ID, move.ID)

	})

	suite.Run("returns moves filtered appeared in TOO at", func() {
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

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
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{
			RequestedMoveDate: &requestedMoveDateString,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
	})

	suite.Run("returns moves filtered by ppm expected_departure_date", func() {
		officeUser, _, session := setupTestData()

		move2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		ppmDate := time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC)
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move2,
				LinkOnly: true,
			},
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: ppmDate,
				},
			},
		}, nil)

		dateStr := ppmDate.Format("2006-01-02")
		moves, _, err := orderFetcher.ListOriginRequestsOrders(
			suite.AppContextWithSessionForTest(&session),
			officeUser.ID,
			&services.ListOrderParams{RequestedMoveDate: &dateStr},
		)
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(move2.ID, moves[0].ID)
	})

	suite.Run("returns moves filtered by NTS requested_delivery_date", func() {
		officeUser, _, session := setupTestData()

		ntsrDate := time.Date(2024, 4, 3, 0, 0, 0, 0, time.UTC)
		move3 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType:          models.MTOShipmentTypeHHGOutOfNTS,
					RequestedDeliveryDate: &ntsrDate,
				},
			},
		}, nil)

		dateStr := ntsrDate.Format("2006-01-02")
		moves, _, err := orderFetcher.ListOriginRequestsOrders(
			suite.AppContextWithSessionForTest(&session),
			officeUser.ID,
			&services.ListOrderParams{RequestedMoveDate: &dateStr},
		)
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(move3.ID, moves[0].ID)
	})

	suite.Run("returns moves filtered by ppm type", func() {
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			PPMType: models.StringPointer("PARTIAL"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(partialPPMMove.Locator, moves[0].Locator)

		// Search for FULL PPM moves
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			CloseoutLocation: models.StringPointer("fT bR"),
			NeedsPPMCloseout: models.BoolPointer(true),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(ppmShipment.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})

	suite.Run("returns moves filtered by closeout initiated date", func() {
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(createdPPM.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.NotEqual(createdPPM2.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})

	suite.Run("latest closeout initiated date is used for filter", func() {
		// Under test: ListOriginRequestsOrders
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate,
		})
		suite.Empty(moves)
		suite.FatalNoError(err)

		// Search for PPMs submitted on April 2nd
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			CloseoutInitiated: &closeoutInitiatedDate2,
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(createdPPM.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
	})

	suite.Run("task order queue does not return move with ONLY a destination address update request", func() {
		officeUser, _, session := setupTestData()
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		testUUID := uuid.UUID{}
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DestinationAddressID: &testUUID,
				},
			},
		}, nil)

		suite.NotNil(shipment)

		shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ShipmentAddressUpdate{
					NewAddressID: testUUID,
				},
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})
		suite.NotNil(shipmentAddressUpdate)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		// even though 2 moves were created, one by setupTestData(), only one will be returned from the call to List Orders since we filter out
		// the one with only a shipment address update to be routed to the destination requests queue
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with origin service items and destination address update request", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only destination shuttle service item
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		testUUID := uuid.UUID{}
		shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ShipmentAddressUpdate{
					NewAddressID: testUUID,
				},
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})
		suite.NotNil(shipmentAddressUpdate)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue does not return move with ONLY requested destination SIT service items", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		// build a move with destination service items
		move2 := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)
		shipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move2,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})

		destinationSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment2,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.Equal(models.MTOServiceItemStatusSubmitted, destinationSITServiceItem.Status)

		suite.FatalNoError(err)
		// even though 2 moves were created, only one will be returned from the call to List Orders since we filter out
		// the one with only destination service items
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with service item names", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		cratingServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDCRT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)
		suite.NotNil(cratingServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
		move = moves[0]
		suite.Len(move.MTOServiceItems, 2)

		var foundDOFSIT, foundDCRT bool
		for _, serviceItem := range move.MTOServiceItems {
			switch serviceItem.ReService.Code {
			case models.ReServiceCode("DOFSIT"):
				foundDOFSIT = true
			case models.ReServiceCode("DCRT"):
				foundDCRT = true
			}
		}
		suite.True(foundDOFSIT, "expected DOFSIT service item was not found")
		suite.True(foundDCRT, "expected DCRT service item was not found")
	})

	suite.Run("task order queue returns a move with origin requested SIT service items", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with both origin and destination requested SIT service items", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		destinationSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.Equal(models.MTOServiceItemStatusSubmitted, destinationSITServiceItem.Status)
		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with origin shuttle service item", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin shuttle service item
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		internationalOriginShuttleServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSHUT,
				},
			},
		}, nil)
		suite.NotNil(internationalOriginShuttleServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue does not return a move with ONLY destination shuttle service item", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only destination shuttle service item
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		domesticDestinationShuttleServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
		}, nil)
		suite.NotNil(domesticDestinationShuttleServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(0, moveCount)
		suite.Equal(0, len(moves))
	})

	suite.Run("task order queue returns a move with both origin and destination shuttle service item", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only destination shuttle service item
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		domesticOriginShuttleServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOSHUT,
				},
			},
		}, nil)
		suite.NotNil(domesticOriginShuttleServiceItem)

		internationalDestinationShuttleServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDSHUT,
				},
			},
		}, nil)
		suite.NotNil(internationalDestinationShuttleServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with excess weight flagged for review", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		now := time.Now()
		// build a move with a ExcessWeightQualifiedAt value
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:                  models.MoveStatusAPPROVALSREQUESTED,
					Show:                    models.BoolPointer(true),
					ExcessWeightQualifiedAt: &now,
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with unaccompanied baggage excess weight flagged for review", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		now := time.Now()
		// build a move with a ExcessUnaccompaniedBaggageWeightQualifiedAt value
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
					ExcessUnaccompaniedBaggageWeightQualifiedAt: &now,
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with a pending SIT extension and origin service items to review", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status: models.SITExtensionStatusPending,
				},
			},
		}, nil)
		suite.NotNil(sitExtension)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue returns a move with pending amended orders", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Document{},
				Type:  &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})

	suite.Run("task order queue does NOT return a move with a pending SIT extension and only destination service items to review", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with destination service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		destinationSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)
		suite.NotNil(destinationSITServiceItem)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status: models.SITExtensionStatusPending,
				},
			},
		}, nil)
		suite.NotNil(sitExtension)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.FatalNoError(err)
		suite.Equal(0, moveCount)
		suite.Equal(0, len(moves))
	})

	suite.Run("task order queue returns a move with a pending SIT extension and BOTH origin and destination service items to review", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// build a move with only origin service items
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		suite.NotNil(shipment)
		originSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		suite.NotNil(originSITServiceItem)

		destinationSITServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status: models.SITExtensionStatusPending,
				},
			},
		}, nil)
		suite.NotNil(sitExtension)

		moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{})

		suite.Equal(models.MTOServiceItemStatusSubmitted, destinationSITServiceItem.Status)
		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Equal(1, len(moves))
	})
}
func (suite *OrderServiceSuite) TestListOrderWithAssignedUserSingle() {
	// Under test: ListOriginRequestsOrders
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
	_, updateError := assignedOfficeUserUpdater.UpdateAssignedOfficeUser(appCtx, createdMove.ID, &scUser, models.QueueTypeCounseling)

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
	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

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
	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

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
	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

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
	waf := entitlements.NewWeightAllotmentFetcher()
	suite.Run("does not return moves where the service member affiliation is Marines for non-USMC office user", func() {

		orderFetcher := NewOrderFetcher(waf)
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
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListOrdersWithEmptyFields() {
	expectedOrder := factory.BuildOrder(suite.DB(), nil, nil)
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)
	moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{PerPage: models.Int64Pointer(1), Page: models.Int64Pointer(1)})

	suite.FatalNoError(err)
	suite.Nil(moves)
}

func (suite *OrderServiceSuite) TestListOrdersWithPagination() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)
	params := services.ListOrderParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(1)}
	moves, count, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

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
	waf := entitlements.NewWeightAllotmentFetcher()
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
					Status:  models.MoveStatusServiceCounselingCompleted,
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

	orderFetcher := NewOrderFetcher(waf)

	suite.Run("Sort by locator code", func() {
		expectedMove1, expectedMove2, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		params = services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove2.Locator, moves[0].Locator)
		suite.Equal(expectedMove1.Locator, moves[1].Locator)
	})

	suite.Run("Sort by move status", func() {
		expectedMove1, expectedMove2, session := setupTestData()

		moveNeedsSC := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusNeedsServiceCounseling,
					Locator: "A1SC01",
					Show:    models.BoolPointer(true),
				},
			},
		}, nil)

		moveSubmitted := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusSUBMITTED,
					Locator: "B2SUB2",
					Show:    models.BoolPointer(true),
				},
			},
		}, nil)

		moveApprovalsRequested := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVALSREQUESTED,
					Locator: "C3APP3",
					Show:    models.BoolPointer(true),
				},
			},
		}, nil)

		moveApproved := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "D4APR4",
					Show:    models.BoolPointer(true),
				},
			},
		}, nil)

		expectedStatusesAsc := []models.MoveStatus{
			models.MoveStatusServiceCounselingCompleted,
			models.MoveStatusSUBMITTED,
			models.MoveStatusNeedsServiceCounseling,
			models.MoveStatusAPPROVALSREQUESTED,
			models.MoveStatusAPPROVED,
			models.MoveStatusAPPROVED,
		}

		expectedLocatorsAsc := []string{
			expectedMove2.Locator,          // SERVICE COUNSELING COMPLETED
			moveSubmitted.Locator,          // SUBMITTED (NEW MOVE)
			moveNeedsSC.Locator,            // NEEDS SERVICE COUNSELING
			moveApprovalsRequested.Locator, // APPROVALS REQUESTED
			expectedMove1.Locator,          // APPROVED
			moveApproved.Locator,           // APPROVED
		}

		params := services.ListOrderParams{
			Sort:  models.StringPointer("status"),
			Order: models.StringPointer("asc"),
		}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(6, len(moves))

		for i, move := range moves {
			suite.Equal(expectedStatusesAsc[i], move.Status, fmt.Sprintf("Unexpected status at index %d (asc)", i))
			suite.Equal(expectedLocatorsAsc[i], move.Locator, fmt.Sprintf("Unexpected locator at index %d (asc)", i))
		}

		expectedStatusesDesc := []models.MoveStatus{
			models.MoveStatusAPPROVED,
			models.MoveStatusAPPROVED,
			models.MoveStatusAPPROVALSREQUESTED,
			models.MoveStatusNeedsServiceCounseling,
			models.MoveStatusSUBMITTED,
			models.MoveStatusServiceCounselingCompleted,
		}

		expectedLocatorsDesc := []string{
			expectedMove1.Locator,          // APPROVED
			moveApproved.Locator,           // APPROVED
			moveApprovalsRequested.Locator, // APPROVALS REQUESTED
			moveNeedsSC.Locator,            // NEEDS SERVICE COUNSELING
			moveSubmitted.Locator,          // SUBMITTED
			expectedMove2.Locator,          // SERVICE COUNSELING COMPLETED
		}

		params.Order = models.StringPointer("desc")
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(6, len(moves))

		for i, move := range moves {
			suite.Equal(expectedStatusesDesc[i], move.Status, fmt.Sprintf("Unexpected status at index %d (asc)", i))
			suite.Equal(expectedLocatorsDesc[i], move.Locator, fmt.Sprintf("Unexpected locator at index %d (asc)", i))
		}
	})

	suite.Run("Sort by service member affiliations", func() {
		expectedMove1, expectedMove2, session := setupTestData()
		params := services.ListOrderParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)

		params = services.ListOrderParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(*expectedMove2.Orders.ServiceMember.Affiliation, *moves[0].Orders.ServiceMember.Affiliation)
		suite.Equal(*expectedMove1.Orders.ServiceMember.Affiliation, *moves[1].Orders.ServiceMember.Affiliation)
	})

	suite.Run("Sort by request move date", func() {
		_, _, session := setupTestData()

		params := services.ListOrderParams{Sort: models.StringPointer("requestedMoveDate"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(2, len(moves[0].MTOShipments)) // the move with two shipments has the earlier date
		suite.Equal(1, len(moves[1].MTOShipments))
		// NOTE: You have to use Jan 02, 2006 as the example for date/time formatting in Go
		suite.Equal(requestedMoveDate1.Format("2006/01/02"), moves[1].MTOShipments[0].RequestedPickupDate.Format("2006/01/02"))

		params = services.ListOrderParams{Sort: models.StringPointer("requestedMoveDate"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(1, len(moves[0].MTOShipments)) // the move with one shipment should be first
		suite.Equal(2, len(moves[1].MTOShipments))
		suite.Equal(requestedMoveDate1.Format("2006/01/02"), moves[0].MTOShipments[0].RequestedPickupDate.Format("2006/01/02"))
	})

	suite.Run("Sort by request move date including pickup, delivery, and PPM expected departure", func() {
		_, _, session := setupTestData()

		expectedDepartureDate := time.Date(testdatagen.GHCTestYear, 01, 01, 0, 0, 0, 0, time.UTC)
		requestedDeliveryDate := time.Date(testdatagen.GHCTestYear, 01, 02, 0, 0, 0, 0, time.UTC)

		// PPM (expected departure date only)
		movePPM := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "PPM001",
					Status:  models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		shipmentPPM := factory.BuildMTOShipmentWithMove(&movePPM, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ShipmentID:            shipmentPPM.ID,
					ExpectedDepartureDate: expectedDepartureDate,
				},
			},
		}, nil)

		// NTSr (delivery date only)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "NTS001",
					Status:  models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType:          models.MTOShipmentTypeHHGOutOfNTS,
					RequestedPickupDate:   nil,
					RequestedDeliveryDate: &requestedDeliveryDate,
				},
			},
		}, nil)

		// sort by requestedMoveDate asc and validate order
		params := services.ListOrderParams{Sort: models.StringPointer("requestedMoveDate"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.True(len(moves) >= 4)

		suite.Equal("PPM001", moves[0].Locator) // jan 1
		suite.Equal("NTS001", moves[1].Locator) // jan 2
		suite.Equal("TTZ123", moves[2].Locator) // jan 3 (pickup)
		suite.Equal("AA1234", moves[3].Locator) // feb 20

		params.Order = models.StringPointer("desc")
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)

		suite.Equal("AA1234", moves[0].Locator)
		suite.Equal("TTZ123", moves[1].Locator)
		suite.Equal("NTS001", moves[2].Locator)
		suite.Equal("PPM001", moves[3].Locator)
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
		factory.BuildMTOShipmentWithMove(&move2, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApprovalsRequested,
				},
			},
		}, nil)
		move3 := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentWithMove(&move3, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApprovalsRequested,
				},
			},
		}, nil)

		params := services.ListOrderParams{Sort: models.StringPointer("appearedInTooAt"), Order: models.StringPointer("asc")}

		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal(moves[0].ID, move1.ID)
		suite.Equal(moves[1].ID, move2.ID)
		suite.Equal(moves[2].ID, move3.ID)
	})

	// ADDS EXTRA MOVE
	suite.Run("Sort by service member last name", func() {
		_, _, session := setupTestData()

		// Last name sort is the only one that needs 3 moves for a complete test, so add that here after all tests that require 2 moves
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{ // Leo Zephyer
					LastName: &serviceMemberLastName,
				},
			},
		}, nil)
		params := services.ListOrderParams{Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Spacemen, Lea", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Zephyer, Leo", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)

		params = services.ListOrderParams{Sort: models.StringPointer("customerName"), Order: models.StringPointer("desc")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal("Zephyer, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Lea", *moves[2].Orders.ServiceMember.LastName+", "+*moves[2].Orders.ServiceMember.FirstName)
	})

	// ADDS EXTRA MOVES
	suite.Run("Listed orders are alphabetical by move code within a non-unique sort column", func() {
		_, _, session := setupTestData()

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "BB1234",
				},
			},
		}, nil)

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusServiceCounselingCompleted,
					Locator: "AA5678",
				},
			},
		}, nil)

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVED,
					Locator: "UU1234",
				},
			},
		}, nil)

		// Check at multiple page sizes becuase without a secondary sort the order within statuses is inconsistent at different page sizes
		params := services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc"), PerPage: models.Int64Pointer(1)}
		moves, count, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(5, count)

		suite.Equal("AA5678", moves[0].Locator)

		params = services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc"), PerPage: models.Int64Pointer(3)}
		moves, count, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(3, len(moves))
		suite.Equal(5, count)

		suite.Equal("AA5678", moves[0].Locator)
		suite.Equal("TTZ123", moves[1].Locator)
		suite.Equal("AA1234", moves[2].Locator)

		// Sorting by a column with non-unique values
		params = services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc")}
		moves, count, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.NoError(err)
		suite.Equal(5, len(moves))
		suite.Equal(5, count)

		suite.Equal(models.MoveStatusServiceCounselingCompleted, moves[0].Status)
		suite.Equal(models.MoveStatusServiceCounselingCompleted, moves[1].Status)
		suite.Equal(models.MoveStatusAPPROVED, moves[2].Status)
		suite.Equal(models.MoveStatusAPPROVED, moves[3].Status)
		suite.Equal(models.MoveStatusAPPROVED, moves[4].Status)

		suite.Equal("AA5678", moves[0].Locator)
		suite.Equal("TTZ123", moves[1].Locator)
		suite.Equal("AA1234", moves[2].Locator)
		suite.Equal("BB1234", moves[3].Locator)
		suite.Equal("UU1234", moves[4].Locator)
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
	waf := entitlements.NewWeightAllotmentFetcher()
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
	orderFetcher := NewOrderFetcher(waf)

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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutInitiated"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppm1.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppm2.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout initiated date (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("closeoutLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by closeout location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("destinationDutyLocation"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentA.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentB.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by destination duty location (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmType"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(ppmShipmentFull.Shipment.MoveTaskOrder.Locator, moves[0].Locator)
		suite.Equal(ppmShipmentPartial.Shipment.MoveTaskOrder.Locator, moves[1].Locator)

		// Sort by PPM type (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
		moves, _, err := orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
			NeedsPPMCloseout: models.BoolPointer(true),
			Sort:             models.StringPointer("ppmStatus"),
			Order:            models.StringPointer("asc"),
		})

		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(ppmShipmentNeedsCloseout.Status, moves[0].MTOShipments[0].PPMShipment.Status)

		// Sort by PPM type (descending)
		moves, _, err = orderFetcher.ListOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeServicesCounselor, &services.ListOrderParams{
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
	waf := entitlements.NewWeightAllotmentFetcher()
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
					PostalCode:     "89828",
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
		orderFetcher := NewOrderFetcher(waf)
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
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVALSREQUESTED,
			},
		},
	}, nil)
	waf := entitlements.NewWeightAllotmentFetcher()
	// Make a TOO user and the postal code to GBLOC link.
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           tooOfficeUser.User.Roles,
		OfficeUserID:    tooOfficeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	orderFetcher := NewOrderFetcher(waf)
	moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, &services.ListOrderParams{})

	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPM() {
	postalCode := "50309"
	partialPPMType := models.MovePPMTypePARTIAL
	waf := entitlements.NewWeightAllotmentFetcher()
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

	// GBLOC for the below doesn't really matter, it just means the query for the moves passes the inner join in ListOriginRequestsOrders
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), ppmShipment.PickupAddress.PostalCode, tooOfficeUser.TransportationOffice.Gbloc)

	orderFetcher := NewOrderFetcher(waf)
	moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListOrdersWithViewAsGBLOCParam() {
	var hqOfficeUser models.OfficeUser
	var hqOfficeUserAGFM models.OfficeUser
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)

	suite.Run("Sort by locator code", func() {
		expectedMove1, expectedMove2, expectedShipment3, hqSession, hqSessionAGFM := setupTestData()

		// Request as an HQ user with their default GBLOC, KKFA
		params := services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		// Expect the same results with a ViewAsGBLOC that equals the user's default GBLOC, KKFA
		params = services.ListOrderParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc"), ViewAsGBLOC: models.StringPointer("KKFA")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal(expectedMove1.Locator, moves[0].Locator)
		suite.Equal(expectedMove2.Locator, moves[1].Locator)

		// Expect the AGFM move when using the ViewAsGBLOC param set to AGFM
		params = services.ListOrderParams{ViewAsGBLOC: models.StringPointer("AGFM")}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&hqSession), hqOfficeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(expectedShipment3.ID, moves[0].MTOShipments[0].ID)

		// Expect the same results without a ViewAsGBLOC for a user whose default GBLOC is AGFM
		params = services.ListOrderParams{}
		moves, _, err = orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&hqSessionAGFM), hqOfficeUserAGFM.ID, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal(expectedShipment3.ID, moves[0].MTOShipments[0].ID)
	})
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithDeletedShipment() {
	postalCode := "50309"
	deletedAt := time.Now()
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)
	moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, &services.ListOrderParams{Status: []string{string(models.MoveStatusSUBMITTED)}})
	suite.FatalNoError(err)
	suite.Equal(0, moveCount)
	suite.Len(moves, 0)
}

func (suite *OrderServiceSuite) TestListOrdersForTOOWithPPMWithOneDeletedShipmentButOtherExists() {
	postalCode := "50309"
	deletedAt := time.Now()
	waf := entitlements.NewWeightAllotmentFetcher()
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

	orderFetcher := NewOrderFetcher(waf)
	moves, moveCount, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), tooOfficeUser.ID, &services.ListOrderParams{})
	suite.FatalNoError(err)
	suite.Equal(1, moveCount)
	suite.Len(moves, 1)
}

func (suite *OrderServiceSuite) TestListAllOrderLocations() {
	waf := entitlements.NewWeightAllotmentFetcher()
	suite.Run("returns a list of all order locations in the current users queue", func() {
		orderFetcher := NewOrderFetcher(waf)
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
	waf := entitlements.NewWeightAllotmentFetcher()
	requestedMoveDate1 := time.Date(testdatagen.GHCTestYear, 05, 20, 0, 0, 0, 0, time.UTC)
	requestedMoveDate2 := time.Date(testdatagen.GHCTestYear, 07, 03, 0, 0, 0, 0, time.UTC)

	setupData := func() {
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
	}

	orderFetcher := NewOrderFetcher(waf)

	suite.Run("list moves by customer name - full name (last, first)", func() {
		setupData()
		// Search "Spacemen, Margaret"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Spacemen, Margaret"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - full name (first last)", func() {
		setupData()
		// Search "Margaret Spacemen"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Margaret Spacemen"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial last (multiple)", func() {
		setupData()
		// Search "space"
		params := services.ListOrderParams{CustomerName: models.StringPointer("space"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Spacemen, Margaret", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial last (single)", func() {
		setupData()
		// Search "Light"
		params := services.ListOrderParams{CustomerName: models.StringPointer("Light"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(1, len(moves))
		suite.Equal("Starlight, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial first", func() {
		setupData()
		// Search "leo"
		params := services.ListOrderParams{CustomerName: models.StringPointer("leo"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Leo", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Starlight, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - partial matching within first or last", func() {
		setupData()
		// Search "ar"
		params := services.ListOrderParams{CustomerName: models.StringPointer("ar"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(2, len(moves))
		suite.Equal("Spacemen, Margaret", *moves[0].Orders.ServiceMember.LastName+", "+*moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Starlight, Leo", *moves[1].Orders.ServiceMember.LastName+", "+*moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("list moves by customer name - empty", func() {
		setupData()
		// Search "johnny"
		params := services.ListOrderParams{CustomerName: models.StringPointer("johnny"), Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)
		suite.Equal(0, len(moves))
	})
}

func (suite *OrderServiceSuite) TestListAllOrderLocationsWithViewAsGBLOCParam() {
	waf := entitlements.NewWeightAllotmentFetcher()
	suite.Run("returns a list of all order locations in the current users queue", func() {
		orderFetcher := NewOrderFetcher(waf)
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

		// Create two default moves with shipment, should be in KKFA and have the status SUBMITTED
		KKFAMove1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		KKFAMove2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		// Create third move with the same origin duty location as one of the above
		KKFAMove3 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ID: KKFAMove2.Orders.OriginDutyLocation.ID,
				},
				Type:     &factory.DutyLocations.OriginDutyLocation,
				LinkOnly: true,
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

		// Confirm the factory created moves have the desired GBLOCS, 3x KKFA,
		suite.Equal("KKFA", *KKFAMove1.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove2.Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAMove3.Orders.OriginDutyLocationGBLOC)

		// Fetch and check secondary GBLOC
		KKFA := "KKFA"
		params := services.ListOrderParams{
			ViewAsGBLOC: &KKFA,
		}
		KKFAmoves, err := orderFetcher.ListAllOrderLocations(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)

		suite.FatalNoError(err)
		// This value should be updated to 3 if ListAllOrderLocations is updated to return distinct locations
		suite.Equal(3, len(KKFAmoves))

		suite.Equal("KKFA", *KKFAmoves[0].Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAmoves[1].Orders.OriginDutyLocationGBLOC)
		suite.Equal("KKFA", *KKFAmoves[2].Orders.OriginDutyLocationGBLOC)

		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove1.Orders.OriginDutyLocation.Name), "Should contain first KKFA move's origin duty location")
		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove2.Orders.OriginDutyLocation.Name), "Should contain second KKFA move's origin duty location")
		suite.Condition(movesContainOriginDutyLocation(KKFAmoves, KKFAMove3.Orders.OriginDutyLocation.Name), "Should contain third KKFA move's origin duty location")
	})
}

func (suite *OrderServiceSuite) TestOriginDutyLocationFilter() {
	var session auth.Session
	waf := entitlements.NewWeightAllotmentFetcher()
	var expectedMove models.Move
	var officeUser models.OfficeUser
	orderFetcher := NewOrderFetcher(waf)
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

	suite.Run("Returns orders matching full originDutyLocation name filter", func() {
		officeUser, expectedMove, session = setupTestData()
		locationName := expectedMove.Orders.OriginDutyLocation.Name
		expectedMoves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{OriginDutyLocation: strings.Split(locationName, " ")})
		suite.NoError(err)
		suite.Equal(1, len(expectedMoves))
		suite.Equal(locationName, string(expectedMoves[0].Orders.OriginDutyLocation.Name))
	})

	suite.Run("Returns orders matching partial originDutyLocation name filter", func() {
		officeUser, expectedMove, session = setupTestData()
		locationName := expectedMove.Orders.OriginDutyLocation.Name
		//Split the location name and retrieve a substring (first string) for the search param
		partialParamSearch := strings.Split(locationName, " ")[0]
		expectedMoves, _, err := orderFetcher.ListOriginRequestsOrders(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.ListOrderParams{OriginDutyLocation: strings.Split(partialParamSearch, " ")})
		suite.NoError(err)
		suite.Equal(1, len(expectedMoves))
		suite.Equal(locationName, string(expectedMoves[0].Orders.OriginDutyLocation.Name))
	})
}

func (suite *OrderServiceSuite) TestListDestinationRequestsOrders() {
	army := models.AffiliationARMY
	airForce := models.AffiliationAIRFORCE
	spaceForce := models.AffiliationSPACEFORCE
	usmc := models.AffiliationMARINES

	setupTestData := func(officeUserGBLOC string) (models.OfficeUser, auth.Session) {

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: officeUserGBLOC,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "99501", officeUser.TransportationOffice.Gbloc)

		fetcher := &mocks.OfficeUserGblocFetcher{}
		fetcher.On("FetchGblocForOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID,
		).Return(officeUserGBLOC, nil)

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		return officeUser, session
	}

	buildMoveKKFA := func(moveCode string, lastName string) (models.Move, models.MTOShipment) {
		postalCode := "90210"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "90210", "KKFA")

		// setting up two moves, each with requested destination SIT service items
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{PostalCode: postalCode},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:  models.MoveStatusAPPROVALSREQUESTED,
					Show:    models.BoolPointer(true),
					Locator: moveCode,
				},
			},
			{
				Model: models.ServiceMember{
					LastName: &lastName,
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		return move, shipment
	}

	buildMoveAGFM := func() (models.Move, models.MTOShipment) {
		postalCode := "43077"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), postalCode, "AGFM")

		// setting up two moves, each with requested destination SIT service items
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{PostalCode: postalCode},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			}}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		return move, shipment
	}

	buildMoveZone2AK := func(branch models.ServiceMemberAffiliation) (models.Move, models.MTOShipment) {
		// Create a USAF move in Alaska Zone II
		// this is a use a us_post_region_cities_id within AK Zone II
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					City:       "Anchorage",
					State:      "AK",
					PostalCode: "99501",
				},
			},
		}, nil)

		// setting up two moves, each with requested destination SIT service items
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &branch,
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		return move, shipment
	}

	buildMoveZone4AK := func(branch models.ServiceMemberAffiliation) (models.Move, models.MTOShipment) {
		// Create a USAF move in Alaska Zone II
		// this will use a us_post_region_cities_id within AK Zone II

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					City:       "Anchorage",
					State:      "AK",
					PostalCode: "99501",
				},
			},
		}, nil)
		// setting up two moves, each with requested destination SIT service items
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
					Show:   models.BoolPointer(true),
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &branch,
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		return move, shipment
	}

	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

	suite.Run("returns move in destination queue with all locked information", func() {
		officeUser, session := setupTestData("MBFL")
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		now := time.Now()

		// setting up two moves, each with requested destination SIT service items
		// a move associated with an air force customer containing AK Zone II shipment
		move, shipment := buildMoveZone2AK(airForce)
		move.LockedByOfficeUserID = &tooUser.ID
		move.LockExpiresAt = &now
		suite.MustSave(&move)

		// destination service item in SUBMITTED status so it shows in queue
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		params := services.ListOrderParams{Status: []string{string(models.MoveStatusAPPROVALSREQUESTED)}}
		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		// we should get both moves back because one is in Zone II & the other is within the postal code GBLOC
		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)
		lockedMove := moves[0]
		suite.NotNil(lockedMove.LockedByOfficeUserID)
		suite.NotNil(lockedMove.LockExpiresAt)
	})

	suite.Run("returns moves for KKFA GBLOC when destination address is in KKFA GBLOC, and uses secondary sort column", func() {
		officeUser, session := setupTestData("KKFA")
		// setting up two moves, each with requested destination SIT service items
		move, shipment := buildMoveKKFA("CC1234", "Spaceman")

		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		move2, shipment2 := buildMoveKKFA("BB1234", "Spaceman")

		// destination shuttle
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
			{
				Model:    shipment2,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		move3, shipment3 := buildMoveKKFA("AA6789", "Landman")
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeMS,
				},
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		move4, shipment4 := buildMoveKKFA("AA1234", "Spaceman")
		// build the destination SIT service items and update their status to SUBMITTED
		oneMonthLater := time.Now().AddDate(0, 1, 0)
		factory.BuildDestSITServiceItems(suite.DB(), move4, shipment4, &oneMonthLater, nil)

		// build the SIT extension update
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    move4,
				LinkOnly: true,
			},
			{
				Model:    shipment4,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status:            models.SITExtensionStatusPending,
					ContractorRemarks: models.StringPointer("gimme some more plz"),
				},
			},
		}, nil)

		// Sort by a column with non-unique values
		params := services.ListOrderParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc")}
		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		suite.FatalNoError(err)
		suite.Equal(4, moveCount)
		suite.Len(moves, 4)

		// Verify primary sort
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moves[0].Status)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moves[1].Status)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moves[2].Status)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moves[3].Status)

		// Verify secondary sort
		suite.Equal("AA1234", moves[0].Locator)
		suite.Equal("AA6789", moves[1].Locator)
		suite.Equal("BB1234", moves[2].Locator)
		suite.Equal("CC1234", moves[3].Locator)

		// Sort by a column with non-unique values
		params = services.ListOrderParams{Sort: models.StringPointer("customerName"), Order: models.StringPointer("asc")}
		moves, moveCount, err = orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		suite.FatalNoError(err)
		suite.Equal(4, moveCount)
		suite.Len(moves, 4)

		// Verify primary sort
		suite.Equal("Landman", *moves[0].Orders.ServiceMember.LastName)
		suite.Equal("Spaceman", *moves[1].Orders.ServiceMember.LastName)
		suite.Equal("Spaceman", *moves[2].Orders.ServiceMember.LastName)
		suite.Equal("Spaceman", *moves[3].Orders.ServiceMember.LastName)

		// Verify secondary sort
		suite.Equal("AA6789", moves[0].Locator)
		suite.Equal("AA1234", moves[1].Locator)
		suite.Equal("BB1234", moves[2].Locator)
		suite.Equal("CC1234", moves[3].Locator)
	})

	suite.Run("returns moves for MBFL GBLOC including USAF/SF in Alaska Zone II", func() {
		officeUser, session := setupTestData("MBFL")

		// setting up two moves, each with requested destination SIT service items
		// a move associated with an air force customer containing AK Zone II shipment
		move, shipment := buildMoveZone2AK(airForce)

		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		// Create a move outside Alaska Zone II (Zone IV in this case)
		move2, shipment2 := buildMoveZone4AK(spaceForce)

		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
			{
				Model:    shipment2,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		params := services.ListOrderParams{Status: []string{string(models.MoveStatusAPPROVALSREQUESTED)}}
		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		// we should get both moves back because one is in Zone II & the other is within the postal code GBLOC
		suite.FatalNoError(err)
		suite.Equal(2, moveCount)
		suite.Len(moves, 2)
	})

	suite.Run("returns moves for JEAT GBLOC excluding USAF/SF in Alaska Zone II", func() {
		officeUser, session := setupTestData("JEAT")

		// Create a move in Zone II, but not an air force or space force service member
		move, shipment := buildMoveZone4AK(army)

		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{},
		)

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)
	})

	suite.Run("returns moves for USMC GBLOC when moves belong to USMC servicemembers", func() {
		officeUser, session := setupTestData("USMC")

		// setting up two moves, each with requested destination SIT service items
		// both will be USMC moves, one in Zone II AK and the other not
		move, shipment := buildMoveZone2AK(usmc)

		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		// this one won't be in Zone II
		move2, shipment2 := buildMoveZone4AK(usmc)

		// destination shuttle
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
			{
				Model:    shipment2,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		move3, shipment3 := buildMoveZone4AK(usmc)
		// we need to create a service item and attach it to the move/shipment
		// else the query will exclude the move since it doesn't use LEFT JOINs
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeMS,
				},
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &services.ListOrderParams{},
		)

		// we should get three moves back since they're USMC moves and zone doesn't matter
		suite.FatalNoError(err)
		suite.Equal(3, moveCount)
		suite.Len(moves, 3)
	})

	suite.Run("returns moves for secondary GBLOC (KKFA) and primary GBLOC (AGFM)", func() {
		// build office user with primary GBLOC of AGFM
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Fort Punxsutawney",
					Gbloc: "AGFM",
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		// Add a secondary GBLOC to the above office user, this should default to KKFA
		secondaryTransportationOfficeAssignment := factory.BuildAlternateTransportationOfficeAssignment(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					ID: officeUser.ID,
				},
				LinkOnly: true,
			},
		}, nil)
		suite.Equal("AGFM", officeUser.TransportationOffice.Gbloc)
		suite.Equal("KKFA", secondaryTransportationOfficeAssignment.TransportationOffice.Gbloc)
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// setting up four moves in KKFA, each with destination requests
		move, shipment := buildMoveKKFA("CC1234", "Spaceman")
		// destination service item in SUBMITTED status
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		move2, shipment2 := buildMoveKKFA("BB1234", "Spaceman")
		// destination shuttle
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
			{
				Model:    shipment2,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		move3, shipment3 := buildMoveKKFA("AA6789", "Landman")
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeMS,
				},
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment3,
				LinkOnly: true,
			},
			{
				Model:    move3,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		move4, shipment4 := buildMoveKKFA("AA1234", "Spaceman")
		// build the destination SIT service items and update their status to SUBMITTED
		oneMonthLater := time.Now().AddDate(0, 1, 0)
		factory.BuildDestSITServiceItems(suite.DB(), move4, shipment4, &oneMonthLater, nil)
		// build the SIT extension update
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    move4,
				LinkOnly: true,
			},
			{
				Model:    shipment4,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status:            models.SITExtensionStatusPending,
					ContractorRemarks: models.StringPointer("gimme some more plz"),
				},
			},
		}, nil)

		// setting up one move in AGFM with destination requests
		AGFMmove1, AGFMshipment1 := buildMoveAGFM()
		// destination shuttle
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				Model:    AGFMmove1,
				LinkOnly: true,
			},
			{
				Model:    AGFMshipment1,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		// Fetch and check secondary GBLOC destination queue
		KKFA := "KKFA"
		params := services.ListOrderParams{
			ViewAsGBLOC: &KKFA,
		}

		moves, moveCount, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		suite.FatalNoError(err)
		suite.Equal(4, moveCount)
		suite.Len(moves, 4)

		// Fetch and check primary GBLOC destination queue
		AGFM := "AGFM"
		params = services.ListOrderParams{
			ViewAsGBLOC: &AGFM,
		}

		moves, moveCount, err = orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session), officeUser.ID, roles.RoleTypeTOO, &params,
		)

		suite.FatalNoError(err)
		suite.Equal(1, moveCount)
		suite.Len(moves, 1)
	})

	suite.Run("filters by new duty location name", func() {
		officeUser, session := setupTestData("KKFA")

		// First move with new duty location "Camp Alpha"
		campAlpha := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "Camp Alpha",
				},
			},
		}, nil)
		moveA := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model:    campAlpha,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation},
		}, nil)
		shipA := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveA,
				LinkOnly: true,
			},
		}, nil)
		// attaching dest service item so it shows up
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    moveA,
				LinkOnly: true,
			},
			{
				Model:    shipA,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		// First move with new duty location "Camp Alpha"
		campBravo := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "Camp Bravo",
				},
			},
		}, nil)
		moveB := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model:    campBravo,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation},
		}, nil)
		shipB := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveB,
				LinkOnly: true,
			},
		}, nil)
		// attaching dest service item so it shows up
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model:    moveA,
				LinkOnly: true,
			},
			{
				Model:    shipB,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		params := services.ListOrderParams{DestinationDutyLocation: swag.String("Camp Alpha")}
		moves, count, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session),
			officeUser.ID,
			roles.RoleTypeTOO,
			&params,
		)
		suite.FatalNoError(err)
		suite.Equal(1, count)
		suite.Len(moves, 1)
		suite.Equal("Camp Alpha", moves[0].Orders.NewDutyLocation.Name)
	})

	suite.Run("customerName sort on first name", func() {
		officeUser, session := setupTestData("KKFA")

		// we got two Smiths, Adam and Bob
		makeMove := func(first string, last string) models.Move {
			move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
				{
					Model: models.ServiceMember{
						FirstName: &first,
						LastName:  &last,
					},
				},
			}, nil)
			shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil)
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDFSIT,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    shipment,
					LinkOnly: true,
				},
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusSubmitted,
					},
				},
			}, nil)
			return move
		}

		makeMove("Adam", "Smith")
		makeMove("Bob", "Smith")

		params := services.ListOrderParams{Sort: swag.String("customerName"), Order: swag.String("asc")}
		moves, count, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session),
			officeUser.ID,
			roles.RoleTypeTOO,
			&params,
		)
		suite.FatalNoError(err)
		suite.Equal(2, count)
		suite.Len(moves, 2)

		// both last names Smith, but Adam should come before Bob
		suite.Equal("Adam", *moves[0].Orders.ServiceMember.FirstName)
		suite.Equal("Bob", *moves[1].Orders.ServiceMember.FirstName)
	})

	suite.Run("sorts by requested move date ascending", func() {
		officeUser, session := setupTestData("KKFA")

		postalCode := "90210"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), postalCode, "KKFA")

		now := time.Now().UTC()
		older := now.Add(-48 * time.Hour)
		newer := now.Add(48 * time.Hour)

		makeMove := func(locator string, reqDate time.Time) models.Move {
			move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
				{
					Model: models.Move{
						Status:  models.MoveStatusAPPROVALSREQUESTED,
						Locator: locator,
					},
				},
			}, nil)

			shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						RequestedPickupDate: &reqDate,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil)

			// attach a destination‑SIT service item so it appears in the queue
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{Model: models.ReService{Code: models.ReServiceCodeDDFSIT}},
				{Model: move, LinkOnly: true},
				{Model: shipment, LinkOnly: true},
				{Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusSubmitted}},
			}, nil)

			return move
		}

		makeMove("OLD1", older)
		makeMove("NEW1", newer)

		params := services.ListOrderParams{
			Sort:  swag.String("requestedMoveDate"),
			Order: swag.String("asc"),
		}
		moves, count, err := orderFetcher.ListDestinationRequestsOrders(
			suite.AppContextWithSessionForTest(&session),
			officeUser.ID,
			roles.RoleTypeTOO,
			&params,
		)
		suite.FatalNoError(err)
		suite.Equal(2, count)
		suite.Len(moves, 2)

		// oldest‐date move should come first
		suite.Equal("OLD1", moves[0].Locator)
		suite.Equal("NEW1", moves[1].Locator)
	})

}
