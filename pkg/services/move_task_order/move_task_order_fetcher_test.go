package movetaskorder_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {

	setupTestData := func() (models.Move, models.MTOShipment) {

		expectedMTO := factory.BuildMove(suite.DB(), nil, nil)

		// Make a couple of shipments for the move; one prime, one external
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DeletedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)

		return expectedMTO, primeShipment
	}

	mtoFetcher := NewMoveTaskOrderFetcher()

	suite.Run("Success with Prime-available move by ID, fetch all non-deleted shipments", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get two shipments back since we didn't set searchParams to exclude external ones.
		suite.Len(actualMTO.MTOShipments, 2)
	})

	suite.Run("Success with fetching move with a related service item", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		address := factory.BuildAddress(suite.DB(), nil, nil)
		sitEntryDate := time.Now()
		customerContact := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{})
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &sitEntryDate,
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact},
				},
			},
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT, // DDFSIT - Domestic destination 1st day SIT
				},
			},
		}, nil)
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
		}, nil)

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		if suite.Len(actualMTO.MTOServiceItems, 2) {
			serviceItem1 := actualMTO.MTOServiceItems[0]
			suite.Equal(models.ReServiceCodeDDFSIT, serviceItem1.ReService.Code)
			suite.Equal(address.StreetAddress1, serviceItem1.SITDestinationFinalAddress.StreetAddress1)
			suite.Equal(address.State, serviceItem1.SITDestinationFinalAddress.State)
			suite.Equal(address.City, serviceItem1.SITDestinationFinalAddress.City)
			suite.Equal(1, len(serviceItem1.CustomerContacts))
		}
	})

	suite.Run("Success with Prime-available move by Locator, no deleted or external shipments", func() {
		expectedMTO, primeShipment := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			Locator:                  expectedMTO.Locator,
			ExcludeExternalShipments: true,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get one shipment back since we requested no external shipments.
		if suite.Len(actualMTO.MTOShipments, 1) {
			suite.Equal(expectedMTO.ID.String(), actualMTO.ID.String())
			suite.Equal(primeShipment.ID.String(), actualMTO.MTOShipments[0].ID.String())
		}
	})

	suite.Run("Success with move that has only deleted shipments", func() {
		mtoWithAllShipmentsDeleted := factory.BuildMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mtoWithAllShipmentsDeleted,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DeletedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: mtoWithAllShipmentsDeleted.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Equal(mtoWithAllShipmentsDeleted.ID, actualMTO.ID)
		suite.Len(actualMTO.MTOShipments, 0)
	})

	suite.Run("Failure - nil searchParams", func() {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), nil)
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should not be nil")
	})

	suite.Run("Failure - searchParams with no ID/locator set", func() {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &services.MoveTaskOrderFetcherParams{})
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should have either a move ID or locator set")
	})

	suite.Run("Failure - Not Found with Bad ID", func() {
		badID, _ := uuid.NewV4()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: badID,
		}

		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.Error(err)
	})
}

// Checks that there are expectedMatchCount matches between the moves and move ID list
// Allows caller to check that expected moves are in the list
// Also returns a map from uuids → moves
func doIDsMatch(moves []models.Move, moveIDs []*uuid.UUID, expectedMatchCount int) (bool, map[uuid.UUID]models.Move) {
	matched := 0
	mapIDs := make(map[uuid.UUID]models.Move)
	for _, expectedID := range moveIDs {
		for _, actualMove := range moves {
			if expectedID != nil && actualMove.ID == *expectedID {
				mapIDs[*expectedID] = actualMove
				expectedID = nil // if found, clear so we don't match again
				matched++        // keep count so we match all
			}
		}
	}
	return expectedMatchCount == matched, mapIDs
}

func (suite *MoveTaskOrderServiceSuite) TestListAllMoveTaskOrdersFetcher() {
	// Set up a hidden move so we can check if it's in the output:
	now := time.Now()
	show := false
	setupTestData := func() (models.Move, models.Move, models.MTOShipment) {
		hiddenMTO := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &show,
				},
			},
		}, nil)

		mto := factory.BuildMove(suite.DB(), nil, nil)

		// Make a couple of shipments for the default move; one prime, one external
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		return hiddenMTO, mto, primeShipment
	}

	mtoFetcher := NewMoveTaskOrderFetcher()

	suite.Run("all move task orders", func() {
		hiddenMTO, mto, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: false,
			IncludeHidden:      true,
			Since:              nil,
		}

		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		move := moves[0]

		suite.Equal(2, len(moves))
		// These are the move ids we want to find in the return list
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &mto.ID}
		suite.True(doIDsMatch(moves, moveIDs, 2))
		suite.NotNil(move.Orders)
		suite.NotNil(move.Orders.OriginDutyLocation)
		suite.NotNil(move.Orders.NewDutyLocation)
	})

	suite.Run("default search - excludes hidden move task orders", func() {
		hiddenMTO, expectedMove, _ := setupTestData()
		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moves {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		// These are the move ids we have but we expect to find just one
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &expectedMove.ID}
		result, mapMoves := doIDsMatch(moves, moveIDs, 1)
		suite.True(result)                                   // exactly one ID should match
		suite.Contains(mapMoves, expectedMove.ID)            // it should be the expectedMove
		suite.Len(mapMoves[expectedMove.ID].MTOShipments, 2) // That move should have 2 shipments

	})

	suite.Run("returns shipments that respect the external searchParams flag", func() {
		hiddenMTO, expectedMove, primeShipment := setupTestData()

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			ExcludeExternalShipments: true,
		}

		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// We should only get back the one move that's not hidden.
		// These are the move ids we have but we expect to find just one
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &expectedMove.ID}
		result, mapMoves := doIDsMatch(moves, moveIDs, 1)
		suite.True(result)                                   // exactly one ID should match
		suite.Len(mapMoves[expectedMove.ID].MTOShipments, 1) // That move should have one shipment
		suite.Equal(primeShipment.ID.String(), mapMoves[expectedMove.ID].MTOShipments[0].ID.String())
	})

	suite.Run("all move task orders that are available to prime and using since", func() {
		// Under test: ListAllMoveOrders
		// Mocked:     None
		// Set up:     Create an old and new move, first check that both are returned.
		//             Then make one move "old" by moving the updatedAt time back
		//             Request all move orders again
		// Expected outcome: Only the new move should be returned
		now = time.Now()

		// Create the two moves
		newMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldMTO := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Check that they're both returned
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: true,
			// IncludeHidden should be false by default
			Since: nil,
		}
		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(moves, 2)

		// Put 1 Move updatedAt in the past
		suite.NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
			now.Add(-2*time.Hour), oldMTO.ID).Exec())

		// Make search params search for moves newer than timestamp
		searchParams.Since = &now
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// Expect one move, the newer one
		suite.Equal(1, len(mtosWithSince))
		moveIDs := []*uuid.UUID{&oldMTO.ID, &newMove.ID}
		result, mapMoves2 := doIDsMatch(mtosWithSince, moveIDs, 1) // only one should match
		suite.True(result)
		suite.Contains(mapMoves2, newMove.ID)                                                                            // and it should be the new one
		suite.NotContains(mapMoves2, oldMTO.ID, "Returned moves should not contain the old move %s", oldMTO.ID.String()) // it's not the hiddenMTO
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListPrimeMoveTaskOrdersFetcher() {
	now := time.Now()
	// Set up a hidden move so we can check if it's in the output:
	hiddenMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Show: models.BoolPointer(false),
			},
		},
	}, nil)
	// Make a default, not Prime-available move:
	nonPrimeMove := factory.BuildMove(suite.DB(), nil, nil)
	// Make some Prime moves:
	primeMove1 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	primeMove2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	primeMove3 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	factory.BuildMTOShipmentWithMove(&primeMove3, suite.DB(), nil, nil)

	// Move primeMove1 and primeMove3 into the past so we can exclude them:
	suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=$1 WHERE id IN ($2, $3);",
		now.Add(-10*time.Second), primeMove1.ID, primeMove3.ID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id=$2;",
		now.Add(-10*time.Second), primeMove1.OrdersID).Exec())

	fetcher := NewMoveTaskOrderFetcher()
	searchParams := services.MoveTaskOrderFetcherParams{}

	// Run the fetcher without `since` to get all Prime moves:
	primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(primeMoves, 3)

	moveIDs := []uuid.UUID{primeMoves[0].ID, primeMoves[1].ID, primeMoves[2].ID}
	suite.NotContains(moveIDs, hiddenMove.ID)
	suite.NotContains(moveIDs, nonPrimeMove.ID)
	suite.Contains(moveIDs, primeMove1.ID)
	suite.Contains(moveIDs, primeMove2.ID)
	suite.Contains(moveIDs, primeMove3.ID)

	// Run the fetcher with `since` to get primeMove2 and primeMove3 (because of the shipment)
	since := now.Add(-5 * time.Second)
	searchParams.Since = &since
	sinceMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(sinceMoves, 2)

	sinceMoveIDs := []uuid.UUID{sinceMoves[0].ID, sinceMoves[1].ID}
	suite.Contains(sinceMoveIDs, primeMove2.ID)
	suite.Contains(sinceMoveIDs, primeMove3.ID)
}
