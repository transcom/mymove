package models_test

import (
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"orders_id": {"OrdersID can not be blank."},
		"status":    {"Status can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestFetchMove() {
	order1, _ := testdatagen.MakeOrder(suite.db)
	order2, _ := testdatagen.MakeOrder(suite.db)

	session := &auth.Session{
		UserID:          order1.ServiceMember.UserID,
		ServiceMemberID: order1.ServiceMemberID,
		ApplicationName: auth.MyApp,
	}
	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := order1.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	suite.Equal(6, len(move.Locator))

	// All correct
	fetchedMove, err := FetchMove(suite.db, session, move.ID)
	suite.Nil(err, "Expected to get moveResult back.")
	suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")

	// Bad Move
	fetchedMove, err = FetchMove(suite.db, session, uuid.Must(uuid.NewV4()))
	suite.Equal(ErrFetchNotFound, err, "Expected to get FetchNotFound.")

	// Bad User
	session.UserID = order2.ServiceMember.UserID
	session.ServiceMemberID = order2.ServiceMemberID
	fetchedMove, err = FetchMove(suite.db, session, move.ID)
	suite.Equal(ErrFetchForbidden, err, "Expected to get a Forbidden back.")
}

func (suite *ModelSuite) TestMoveCancellationWithReason() {
	orders, err := testdatagen.MakeOrder(suite.db)
	suite.Nil(err)
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.mustSave(&orders)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders
	reason := "SM's orders revoked"

	// Check to ensure move shows SUBMITTED before Cancel()
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.GetStatus(), "expected Submitted")

	// Can cancel move, and status changes as expected
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.GetStatus(), "expected Canceled")
	suite.Equal(&reason, move.CancelReason, "expected 'SM's orders revoked'")

}

func (suite *ModelSuite) TestMoveStateMachine() {
	orders, err := testdatagen.MakeOrder(suite.db)
	suite.Nil(err)
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.mustSave(&orders)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	reason := ""
	move.Orders = orders

	// Can't cancel a move with DRAFT status
	err = move.Cancel(reason)
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

	// Once submitted
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.GetStatus(), "expected Submitted")

	// Can cancel move
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.GetStatus(), "expected Canceled")
	suite.Nil(move.CancelReason)
}

func (suite *ModelSuite) TestCancelMoveCancelsOrdersPPM() {
	// Given: A move with Orders, PPM and Move all in submitted state
	orders, err := testdatagen.MakeOrder(suite.db)
	suite.Nil(err)
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.mustSave(&orders)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.db, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	ppm.Status = PPMStatusSUBMITTED // NEVER do this outside of a test.

	// Associate PPM with the move it's on.
	move.PersonallyProcuredMoves = append(move.PersonallyProcuredMoves, *ppm)
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.GetStatus(), "expected Submitted")

	// When move is canceled, expect associated PPM and Order to be canceled
	reason := "Orders changed"
	err = move.Cancel(reason)
	suite.Nil(err)

	suite.Equal(MoveStatusCANCELED, move.GetStatus(), "expected Canceled")
	suite.Equal(PPMStatusCANCELED, move.PersonallyProcuredMoves[0].Status, "expected Canceled")
	suite.Equal(OrderStatusCANCELED, move.Orders.Status, "expected Canceled")
}
