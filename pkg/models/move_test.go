// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package models_test

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &m.Move{}

	expErrors := map[string][]string{
		"locator":   {"Locator can not be blank."},
		"orders_id": {"OrdersID can not be blank."},
		"status":    {"Status can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestCreateNewMoveValidLocatorString() {
	orders := factory.BuildOrder(suite.DB(), nil, nil)
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	office := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	moveOptions := m.MoveOptions{
		Show:               m.BoolPointer(true),
		CounselingOfficeID: &office.ID,
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	// Verify valid items are in locator
	suite.Regexp("^[346789BCDFGHJKMPQRTVWXY]+$", move.Locator)
	// Verify invalid items are not in locator - this should produce "non-word" locators
	suite.NotRegexp("[0125AEIOULNSZ]", move.Locator)
}

func (suite *ModelSuite) TestNeedsReweighFollowsBusinessLogic() {
	suite.Run("Prime estimated weights trigger reweigh", func() {
		maxWeight := 7000

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: m.Entitlement{
					DBAuthorizedWeight: &maxWeight,
				},
			},
		}, nil)

		now := time.Now()

		moveNeedsReweigh := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: m.Move{
					AvailableToPrimeAt: &now,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		testWeightsNeedsReweigh := []unit.Pound{3500, 3460}
		for i := range testWeightsNeedsReweigh {
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: m.MTOShipment{
						PrimeEstimatedWeight: &testWeightsNeedsReweigh[i],
					},
				},
				{
					Model:    moveNeedsReweigh,
					LinkOnly: true,
				},
			}, nil)
		}

		needsReweigh, err := moveNeedsReweigh.NeedsReweigh(suite.AppContextForTest())
		suite.NoError(err)
		suite.Equal(true, *needsReweigh)
	})

	suite.Run("Prime actual weights trigger reweigh", func() {
		maxWeight := 7000

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: m.Entitlement{
					DBAuthorizedWeight: &maxWeight,
				},
			},
		}, nil)

		now := time.Now()

		moveNeedsReweigh := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: m.Move{
					AvailableToPrimeAt: &now,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		testWeightsNeedsReweigh := []unit.Pound{3500, 3460}
		for i := range testWeightsNeedsReweigh {
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: m.MTOShipment{
						PrimeEstimatedWeight: &testWeightsNeedsReweigh[i],
					},
				},
				{
					Model:    moveNeedsReweigh,
					LinkOnly: true,
				},
			}, nil)
		}

		needsReweigh, err := moveNeedsReweigh.NeedsReweigh(suite.AppContextForTest())
		suite.NoError(err)
		suite.Equal(true, *needsReweigh)
	})

	suite.Run("Reweigh does not trigger when it should not trigger", func() {
		maxWeight := 7000

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: m.Entitlement{
					DBAuthorizedWeight: &maxWeight,
				},
			},
		}, nil)

		now := time.Now()

		moveNeedsReweigh := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: m.Move{
					AvailableToPrimeAt: &now,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		testWeightsDoesNotNeedReweigh := []unit.Pound{2500, 2000}
		for i := range testWeightsDoesNotNeedReweigh {
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: m.MTOShipment{
						PrimeEstimatedWeight: &testWeightsDoesNotNeedReweigh[i],
					},
				},
				{
					Model:    moveNeedsReweigh,
					LinkOnly: true,
				},
			}, nil)
		}

		needsReweigh, err := moveNeedsReweigh.NeedsReweigh(suite.AppContextForTest())
		suite.NoError(err)
		suite.Equal(false, *needsReweigh)
	})
}

func (suite *ModelSuite) TestGenerateReferenceID() {

	refID, err := m.GenerateReferenceID(suite.DB())
	suite.Run("reference id is properly created", func() {
		suite.NoError(err)
		suite.NotZero(refID)
		firstNum, _ := strconv.Atoi(strings.Split(refID, "-")[0])
		secondNum, _ := strconv.Atoi(strings.Split(refID, "-")[1])
		suite.Equal(reflect.TypeOf(refID).String(), "string")
		suite.Equal(firstNum >= 0 && firstNum <= 9999, true)
		suite.Equal(secondNum >= 0 && secondNum <= 9999, true)
		suite.Equal(string(refID[4]), "-")
	})
}

func (suite *ModelSuite) TestFetchMove() {

	setupTestData := func() (*auth.Session, m.Order) {

		order := factory.BuildOrder(suite.DB(), nil, nil)
		factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		session := &auth.Session{
			UserID:          order.ServiceMember.UserID,
			ServiceMemberID: order.ServiceMemberID,
			ApplicationName: auth.MilApp,
		}
		return session, order

	}

	suite.Run("Fetch a move", func() {
		// Under test:       FetchMove fetches a move associated with a specific order
		// Mocked:           None
		// Set up:           Create an HHG move, then fetch it, then move to status completed, fetch again
		// Expected outcome: Move found, in both cases
		session, order := setupTestData()
		office := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		// Create HHG Move
		moveOptions := m.MoveOptions{
			Show:               m.BoolPointer(true),
			CounselingOfficeID: &office.ID,
		}
		move, verrs, err := order.CreateNewMove(suite.DB(), moveOptions)
		suite.NoError(err)
		suite.Zero(verrs.Count())
		suite.Equal(6, len(move.Locator))

		// Fetch move
		fetchedMove, err := m.FetchMove(suite.DB(), session, move.ID)
		suite.Nil(err, "Expected to get moveResult back.")
		suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")

		// We're asserting that if for any reason
		// a move gets into the remove "COMPLETED" state
		// it does not fail being queried
		move.Status = "COMPLETED"
		suite.DB().Save(move)

		// Fetch move again
		actualMove, err := m.FetchMove(suite.DB(), session, move.ID)
		suite.NoError(err, "Failed fetching completed move")
		suite.Equal("COMPLETED", string(actualMove.Status))

	})

	suite.Run("Fetch a move not found", func() {
		// Under test:       FetchMove
		// Mocked:           None
		// Set up:           Fetch a non-existent move
		// Expected outcome: Move not found, ErrFetchNotFound error

		session, _ := setupTestData()

		// Bad Move
		_, err := m.FetchMove(suite.DB(), session, uuid.Must(uuid.NewV4()))
		suite.Equal(m.ErrFetchNotFound, err, "Expected to get FetchNotFound.")
	})

	suite.Run("Fetch a move bad user", func() {
		// Under test:       FetchMove
		// Mocked:           None
		// Set up:           Create a user and orders, no move. Create a second user and a move.
		//                   Fetch the second user's move, but with the first user logged in.
		// Expected outcome: Move not found, ErrFetchForbidden
		session, _ := setupTestData()
		office := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		// Create a second sm and a move only on that sm
		order2 := factory.BuildOrder(suite.DB(), nil, nil)
		moveOptions := m.MoveOptions{
			Show:               m.BoolPointer(true),
			CounselingOfficeID: &office.ID,
		}
		move2, verrs, err := order2.CreateNewMove(suite.DB(), moveOptions)
		suite.NoError(err)
		suite.Zero(verrs.Count())
		suite.Equal(6, len(move2.Locator))

		// A fetch on the second moveID, with the first user logged in, should fail
		_, err = m.FetchMove(suite.DB(), session, move2.ID)

		suite.Equal(m.ErrFetchForbidden, err, "Expected to get a Forbidden back.")
	})

	suite.Run("Hidden move is not returned", func() {
		// Under test:       FetchMove
		// Mocked:           None
		// Set up:           Create an sm with orders, then create a hidden move
		// Expected outcome: Move not found, ErrFetchNotFound error
		session, order := setupTestData()
		// Create a hidden move
		hiddenMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: m.Move{
					Show: m.BoolPointer(false),
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		// Attempt to fetch this move. We should receive an error.
		_, err := m.FetchMove(suite.DB(), session, hiddenMove.ID)
		suite.Equal(m.ErrFetchNotFound, err, "Expected to get FetchNotFound.")
	})

	suite.Run("deleted shipments are excluded from the results", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mto := mtoShipment.MoveTaskOrder
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: m.MTOShipment{
					ShipmentType: m.MTOShipmentTypeHHG,
					Status:       m.MTOShipmentStatusSubmitted,
					DeletedAt:    m.TimePointer(time.Now()),
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)

		session := &auth.Session{
			UserID:          mto.Orders.ServiceMember.UserID,
			ServiceMemberID: mto.Orders.ServiceMemberID,
			ApplicationName: auth.MilApp,
		}

		actualMove, err := m.FetchMove(suite.DB(), session, mto.ID)

		suite.NoError(err)
		suite.Len(actualMove.MTOShipments, 1)

		suite.Equal(mtoShipment.ID, actualMove.MTOShipments[0].ID)
	})
}

func (suite *ModelSuite) TestSaveMoveDependenciesFail() {
	// Given: A move with Orders with unacceptable status
	orders := factory.BuildOrder(suite.DB(), nil, nil)
	orders.Status = ""
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	office := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	moveOptions := m.MoveOptions{
		Show:               m.BoolPointer(true),
		CounselingOfficeID: &office.ID,
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)

	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, _ = m.SaveMoveDependencies(suite.DB(), move)

	suite.True(verrs.HasAny(), "saving invalid statuses should yield an error")
}

func (suite *ModelSuite) TestSaveMoveDependenciesSuccess() {
	// Given: A move with Orders with acceptable status
	orders := factory.BuildOrder(suite.DB(), nil, nil)
	orders.Status = m.OrderStatusSUBMITTED
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	office := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	moveOptions := m.MoveOptions{
		Show:               m.BoolPointer(true),
		CounselingOfficeID: &office.ID,
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = m.SaveMoveDependencies(suite.DB(), move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.NoError(err)
}

func (suite *ModelSuite) TestFetchMoveByOrderID() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")
	invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: m.Move{
				ID: moveID,
			},
		},
		{
			Model: m.Order{
				ID: orderID,
			},
		},
	}, nil)

	tests := []struct {
		lookupID  uuid.UUID
		resultID  uuid.UUID
		resultErr bool
	}{
		{lookupID: orderID, resultID: moveID, resultErr: false},
		{lookupID: invalidID, resultID: invalidID, resultErr: true},
	}

	for _, ts := range tests {
		move, err := m.FetchMoveByOrderID(suite.DB(), ts.lookupID)
		if ts.resultErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
		suite.Equal(move.ID, ts.resultID, "Wrong moveID: %s", ts.lookupID)
	}
}

func (suite *ModelSuite) FetchMovesByOrderID() {
	// Given an order with multiple moves return all moves belonging to that order.
	orderID := uuid.Must(uuid.NewV4())

	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")
	moveID2, _ := uuid.FromString("e76b5dae-ae00-4147-b818-07eff29fca98")

	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: m.Move{
				ID: moveID,
			},
		},
		{
			Model: m.Order{
				ID: orderID,
			},
		},
	}, nil)
	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: m.Move{
				ID: moveID2,
			},
		},
		{
			Model: m.Order{
				ID: orderID,
			},
		},
	}, nil)

	tests := []struct {
		lookupID  uuid.UUID
		resultErr bool
	}{
		{lookupID: orderID, resultErr: false},
	}

	moves, err := m.FetchMovesByOrderID(suite.DB(), tests[0].lookupID)
	if err != nil {
		suite.Error(err)
	}

	suite.Greater(len(moves), 1)
}

func (suite *ModelSuite) TestMoveIsPPMOnly() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	isPPMOnly := move.IsPPMOnly()
	suite.False(isPPMOnly, "A move with no shipments will return false for isPPMOnly.")

	factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
		{
			Model: m.MTOShipment{
				ShipmentType: m.MTOShipmentTypePPM,
			},
		},
	}, nil)
	isPPMOnly = move.IsPPMOnly()
	suite.True(isPPMOnly, "A move with only PPM shipments will return true for isPPMOnly")

	factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
		{
			Model: m.MTOShipment{
				ShipmentType: m.MTOShipmentTypeHHG,
			},
		},
	}, nil)
	isPPMOnly = move.IsPPMOnly()
	suite.False(isPPMOnly, "A move with one PPM shipment and one HHG shipment will return false for isPPMOnly.")
}

func (suite *ModelSuite) TestMoveHasPPM() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	hasPPM := move.HasPPM()
	suite.False(hasPPM, "A move with no shipments will return false for hasPPM.")

	factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
		{
			Model: m.MTOShipment{
				ShipmentType: m.MTOShipmentTypePPM,
			},
		},
	}, nil)
	hasPPM = move.HasPPM()
	suite.True(hasPPM, "A move with only PPM shipments will return true for hasPPM")

	factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
		{
			Model: m.MTOShipment{
				ShipmentType: m.MTOShipmentTypeHHG,
			},
		},
	}, nil)
	hasPPM = move.HasPPM()
	suite.True(hasPPM, "A move with one PPM shipment and one HHG shipment will return true for hasPPM.")

	move2 := factory.BuildMove(suite.DB(), nil, nil)

	factory.BuildMTOShipmentWithMove(&move2, suite.DB(), []factory.Customization{
		{
			Model: m.MTOShipment{
				ShipmentType: m.MTOShipmentTypeHHG,
			},
		},
	}, nil)
	hasPPM = move2.HasPPM()
	suite.False(hasPPM, "A move with one HHG shipment will return false for hasPPM.")
}
