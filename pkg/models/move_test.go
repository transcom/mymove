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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

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
}

func (suite *ModelSuite) TestMoveCancellationWithReason() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)
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
	reason := "SM's orders revoked"

	// Check to ensure move shows SUBMITTED before Cancel()
	err = move.Submit()
	suite.NoError(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// Can cancel move, and status changes as expected
	err = move.Cancel(reason)
	suite.NoError(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Equal(&reason, move.CancelReason, "expected 'SM's orders revoked'")

}

func (suite *ModelSuite) TestMoveStateMachine() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)
	testdatagen.MakeDefaultContractor(suite.DB())

	selectedMoveType := SelectedMoveTypeHHGPPM

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	reason := ""
	move.Orders = orders

	// Create PPM on this move
	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: PersonallyProcuredMove{
			Move:      *move,
			MoveID:    move.ID,
			Status:    PPMStatusDRAFT,
			Advance:   &advance,
			AdvanceID: &advance.ID,
		},
	})
	move.PersonallyProcuredMoves = append(move.PersonallyProcuredMoves, ppm)

	// Once submitted
	err = move.Submit()
	suite.MustSave(move)
	suite.DB().Reload(move)
	suite.NoError(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.Equal(PPMStatusSUBMITTED, move.PersonallyProcuredMoves[0].Status, "expected Submitted")
	// Can cancel move
	err = move.Cancel(reason)
	suite.NoError(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Nil(move.CancelReason)
}

func (suite *ModelSuite) TestCancelMoveCancelsOrdersPPM() {
	// Given: A move with Orders, PPM and Move all in submitted state
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)
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

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	// Associate PPM with the move it's on.
	move.PersonallyProcuredMoves = append(move.PersonallyProcuredMoves, *ppm)
	err = move.Submit()
	suite.NoError(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// When move is canceled, expect associated PPM and Order to be canceled
	reason := "Orders changed"
	err = move.Cancel(reason)
	suite.NoError(err)

	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Equal(PPMStatusCANCELED, move.PersonallyProcuredMoves[0].Status, "expected Canceled")
	suite.Equal(OrderStatusCANCELED, move.Orders.Status, "expected Canceled")
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

func (suite *ModelSuite) TestMoveApproval() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

	suite.Run("from valid statuses", func() {
		validStatuses := []struct {
			desc   string
			status MoveStatus
		}{
			{"Submitted", MoveStatusSUBMITTED},
			{"Approvals Requested", MoveStatusAPPROVALSREQUESTED},
			{"Service Counseling Completed", MoveStatusServiceCounselingCompleted},
			{"Approved", MoveStatusAPPROVED},
		}
		for _, validStatus := range validStatuses {
			move.Status = validStatus.status

			err := move.Approve()

			suite.NoError(err)
			suite.Equal(MoveStatusAPPROVED, move.Status)
		}
	})

	suite.Run("from invalid statuses", func() {
		invalidStatuses := []struct {
			desc   string
			status MoveStatus
		}{
			{"Draft", MoveStatusDRAFT},
			{"Canceled", MoveStatusCANCELED},
			{"Needs Service Counseling", MoveStatusNeedsServiceCounseling},
		}
		for _, invalidStatus := range invalidStatuses {
			move.Status = invalidStatus.status

			err := move.Approve()

			suite.Error(err)
			suite.Contains(err.Error(), "A move can only be approved if it's in one of these states")
			suite.Contains(err.Error(), fmt.Sprintf("However, its current status is: %s", invalidStatus.status))
		}
	})
}
