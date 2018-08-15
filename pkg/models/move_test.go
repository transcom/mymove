package models_test

import (
	"time"

	"github.com/gobuffalo/uuid"

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
	order1 := testdatagen.MakeDefaultOrder(suite.db)
	order2 := testdatagen.MakeDefaultOrder(suite.db)

	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)

	tdl := testdatagen.MakeDefaultTDL(suite.db)

	market := "dHHG"
	sourceGBLOC := "BMLK"

	session := &auth.Session{
		UserID:          order1.ServiceMember.UserID,
		ServiceMemberID: order1.ServiceMemberID,
		ApplicationName: auth.MyApp,
	}
	var selectedType = internalmessages.SelectedMoveTypeHHG

	move, verrs, err := order1.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	suite.Equal(6, len(move.Locator))

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			ServiceMember:           &order1.ServiceMember,
			Move:                    move,
		},
	})

	// All correct
	fetchedMove, err := FetchMove(suite.db, session, move.ID)
	suite.Nil(err, "Expected to get moveResult back.")
	suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")
	suite.Equal(fetchedMove.Shipments[0].PickupAddressID, shipment.PickupAddressID)

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
	orders := testdatagen.MakeDefaultOrder(suite.db)
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
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// Can cancel move, and status changes as expected
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Equal(&reason, move.CancelReason, "expected 'SM's orders revoked'")

}

func (suite *ModelSuite) TestMoveStateMachine() {
	orders := testdatagen.MakeDefaultOrder(suite.db)
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.mustSave(&orders)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	reason := ""
	move.Orders = orders

	// Once submitted
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// Can cancel move
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Nil(move.CancelReason)
}

func (suite *ModelSuite) TestCancelMoveCancelsOrdersPPM() {
	// Given: A move with Orders, PPM and Move all in submitted state
	orders := testdatagen.MakeDefaultOrder(suite.db)
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

	// Associate PPM with the move it's on.
	move.PersonallyProcuredMoves = append(move.PersonallyProcuredMoves, *ppm)
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// When move is canceled, expect associated PPM and Order to be canceled
	reason := "Orders changed"
	err = move.Cancel(reason)
	suite.Nil(err)

	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Equal(PPMStatusCANCELED, move.PersonallyProcuredMoves[0].Status, "expected Canceled")
	suite.Equal(OrderStatusCANCELED, move.Orders.Status, "expected Canceled")
}

func (suite *ModelSuite) TestSaveMoveDependenciesFail() {
	// Given: A move with Orders with unacceptable status
	orders := testdatagen.MakeDefaultOrder(suite.db)
	orders.Status = ""

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = SaveMoveDependencies(suite.db, move)
	suite.True(verrs.HasAny(), "saving invalid statuses should yield an error")
}

func (suite *ModelSuite) TestSaveMoveDependenciesSuccess() {
	// Given: A move with Orders with acceptable status
	orders := testdatagen.MakeDefaultOrder(suite.db)
	orders.Status = OrderStatusSUBMITTED

	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = SaveMoveDependencies(suite.db, move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.Nil(err)
}

func (suite *ModelSuite) TestSaveMoveDependenciesSetsGBLOCSuccess() {
	// Given: A shipment's move with orders in acceptable status

	dutyStation := testdatagen.MakeDefaultDutyStation(suite.db)
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.db)
	serviceMember.DutyStationID = &dutyStation.ID
	serviceMember.DutyStation = dutyStation
	suite.mustSave(&serviceMember)

	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	market := "dHHG"
	sourceGBLOC := "BMLK"

	orders := testdatagen.MakeDefaultOrder(suite.db)
	orders.Status = OrderStatusSUBMITTED

	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			ServiceMember:           &serviceMember,
			Move:                    move,
		},
	})
	// Associate Shipment with the move it's on.
	move.Shipments = append(move.Shipments, shipment)
	move.Orders = orders

	// And: Move is in SUBMITTED state
	move.Status = MoveStatusSUBMITTED
	verrs, err = SaveMoveDependencies(suite.db, move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.Nil(err)
	suite.db.Reload(&shipment)

	// Then: Shipment GBLOCs will be equal to:
	// destination GBLOC: orders' new duty station's transportation office's GBLOC
	suite.Assertions.Equal(orders.NewDutyStation.TransportationOffice.Gbloc, *shipment.DestinationGBLOC)
	// source GBLOC: service member's current duty station's transportation office's GBLOC
	suite.Assertions.Equal(serviceMember.DutyStation.TransportationOffice.Gbloc, *shipment.SourceGBLOC)
}
