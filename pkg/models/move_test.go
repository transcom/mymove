//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package models_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"locator":   {"Locator can not be blank."},
		"orders_id": {"OrdersID can not be blank."},
		"status":    {"Status can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestCreateNewMoveValidLocatorString() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())
	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	// Verify valid items are in locator
	suite.Regexp("^[346789BCDFGHJKMPQRTVWXY]+$", move.Locator)
	// Verify invalid items are not in locator - this should produce "non-word" locators
	suite.NotRegexp("[0125AEIOULNSZ]", move.Locator)
}

func (suite *ModelSuite) TestGenerateReferenceID() {

	refID, err := GenerateReferenceID(suite.DB())
	suite.T().Run("reference id is properly created", func(t *testing.T) {
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
	order1 := testdatagen.MakeDefaultOrder(suite.DB())
	order2 := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

	session := &auth.Session{
		UserID:          order1.ServiceMember.UserID,
		ServiceMemberID: order1.ServiceMemberID,
		ApplicationName: auth.MilApp,
	}
	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := order1.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	suite.Equal(6, len(move.Locator))

	// All correct
	fetchedMove, err := FetchMove(suite.DB(), session, move.ID)
	suite.Nil(err, "Expected to get moveResult back.")
	suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")

	// We're asserting that if for any reason
	// a move gets into the remove "COMPLETED" state
	// it does not fail being queried
	move.Status = "COMPLETED"
	suite.DB().Save(move)

	actualMove, err := FetchMove(suite.DB(), session, move.ID)

	suite.NoError(err, "Failed fetching completed move")
	suite.Equal("COMPLETED", string(actualMove.Status))

	move.Status = MoveStatusDRAFT
	suite.DB().Save(move) // teardown/reset back to draft

	// Bad Move
	_, err = FetchMove(suite.DB(), session, uuid.Must(uuid.NewV4()))
	suite.Equal(ErrFetchNotFound, err, "Expected to get FetchNotFound.")

	// Bad User
	session.UserID = order2.ServiceMember.UserID
	session.ServiceMemberID = order2.ServiceMemberID
	_, err = FetchMove(suite.DB(), session, move.ID)
	suite.Equal(ErrFetchForbidden, err, "Expected to get a Forbidden back.")

	suite.T().Run("Hidden move is not returned", func(t *testing.T) {
		// Create a hidden move
		hiddenMove := testdatagen.MakeHiddenHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

		// Attempt to fetch this move. We should receive an error.
		_, err := FetchMove(suite.DB(), session, hiddenMove.ID)
		suite.Equal(ErrFetchNotFound, err, "Expected to get FetchNotFound.")
	})

	suite.T().Run("deleted shipments are excluded from the results", func(t *testing.T) {
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		mto := mtoShipment.MoveTaskOrder
		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: MTOShipment{
				ShipmentType: MTOShipmentTypeHHG,
				Status:       MTOShipmentStatusSubmitted,
				DeletedAt:    TimePointer(time.Now()),
			},
			Move: mto,
		})

		session := &auth.Session{
			UserID:          mto.Orders.ServiceMember.UserID,
			ServiceMemberID: mto.Orders.ServiceMemberID,
			ApplicationName: auth.MilApp,
		}

		actualMove, err := FetchMove(suite.DB(), session, mto.ID)

		suite.NoError(err)
		suite.Len(actualMove.MTOShipments, 1)

		suite.Equal(mtoShipment.ID, actualMove.MTOShipments[0].ID)
	})
}

func (suite *ModelSuite) TestSaveMoveDependenciesFail() {
	// Given: A move with Orders with unacceptable status
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = ""
	testdatagen.MakeDefaultContractor(suite.DB())
	selectedMoveType := SelectedMoveTypeHHGPPM

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)

	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, _ = SaveMoveDependencies(suite.DB(), move)
	suite.True(verrs.HasAny(), "saving invalid statuses should yield an error")
}

func (suite *ModelSuite) TestSaveMoveDependenciesSuccess() {
	// Given: A move with Orders with acceptable status
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED
	testdatagen.MakeDefaultContractor(suite.DB())
	selectedMoveType := SelectedMoveTypeHHGPPM

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = SaveMoveDependencies(suite.DB(), move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.NoError(err)
}

func (suite *ModelSuite) TestFetchMoveByOrderID() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")
	invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: Order{
			ID: orderID,
		},
	})
	testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: Move{
			ID:       moveID,
			OrdersID: orderID,
			Orders:   order,
		},
	})

	tests := []struct {
		lookupID  uuid.UUID
		resultID  uuid.UUID
		resultErr bool
	}{
		{lookupID: orderID, resultID: moveID, resultErr: false},
		{lookupID: invalidID, resultID: invalidID, resultErr: true},
	}

	for _, ts := range tests {
		move, err := FetchMoveByOrderID(suite.DB(), ts.lookupID)
		if ts.resultErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
		suite.Equal(move.ID, ts.resultID, "Wrong moveID: %s", ts.lookupID)
	}
}
