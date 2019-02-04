package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dates"
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

func (suite *ModelSuite) TestCreateNewMoveValidLocatorString() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	selectedMoveType := SelectedMoveTypeHHG

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)

	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	// Verify valid items are in locator
	suite.Regexp("^[346789BCDFGHJKMPQRTVWXY]+$", move.Locator)
	// Verify invalid items are not in locator - this should produce "non-word" locators
	suite.NotRegexp("[0125AEIOULNSZ]", move.Locator)
}

func (suite *ModelSuite) TestFetchMove() {
	order1 := testdatagen.MakeDefaultOrder(suite.DB())
	order2 := testdatagen.MakeDefaultOrder(suite.DB())

	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	deliveryDate := dates.NextWorkday(*calendar, pickupDate)

	tdl := testdatagen.MakeDefaultTDL(suite.DB())

	market := "dHHG"
	sourceGBLOC := "BMLK"

	session := &auth.Session{
		UserID:          order1.ServiceMember.UserID,
		ServiceMemberID: order1.ServiceMemberID,
		ApplicationName: auth.MyApp,
	}
	selectedMoveType := SelectedMoveTypeHHG

	move, verrs, err := order1.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	suite.Equal(6, len(move.Locator))

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			ServiceMember:           order1.ServiceMember,
			Move:                    *move,
			MoveID:                  move.ID,
		},
	})

	// All correct
	fetchedMove, err := FetchMove(suite.DB(), session, move.ID)
	suite.Nil(err, "Expected to get moveResult back.")
	suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")
	suite.Equal(fetchedMove.Shipments[0].PickupAddressID, shipment.PickupAddressID)

	// Bad Move
	fetchedMove, err = FetchMove(suite.DB(), session, uuid.Must(uuid.NewV4()))
	suite.Equal(ErrFetchNotFound, err, "Expected to get FetchNotFound.")

	// Bad User
	session.UserID = order2.ServiceMember.UserID
	session.ServiceMemberID = order2.ServiceMemberID
	fetchedMove, err = FetchMove(suite.DB(), session, move.ID)
	suite.Equal(ErrFetchForbidden, err, "Expected to get a Forbidden back.")
}

func (suite *ModelSuite) TestMoveCancellationWithReason() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
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
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
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

	// Create hhg (shipment) on this move
	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	deliveryDate := dates.NextWorkday(*calendar, pickupDate)
	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	market := "dHHG"
	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: Shipment{
			MoveID:                  move.ID,
			Move:                    *move,
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
			Status:                  ShipmentStatusDRAFT,
		},
	})

	move.Shipments = append(move.Shipments, shipment)

	// Once submitted
	err = move.Submit()
	suite.MustSave(move)
	suite.DB().Reload(move)
	suite.Nil(err)
	suite.Equal(MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.Equal(PPMStatusSUBMITTED, move.PersonallyProcuredMoves[0].Status, "expected Submitted")
	suite.Equal(ShipmentStatusSUBMITTED, move.Shipments[0].Status, "expected Submitted")
	// Can cancel move
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Nil(move.CancelReason)
}

func (suite *ModelSuite) TestCancelMoveCancelsOrdersPPM() {
	// Given: A move with Orders, PPM and Move all in submitted state
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
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
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = ""

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = SaveMoveDependencies(suite.DB(), move)
	suite.True(verrs.HasAny(), "saving invalid statuses should yield an error")
}

func (suite *ModelSuite) TestSaveMoveDependenciesSuccess() {
	// Given: A move with Orders with acceptable status
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	verrs, err = SaveMoveDependencies(suite.DB(), move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.Nil(err)
}

func (suite *ModelSuite) TestSaveMoveDependenciesSetsGBLOCSuccess() {
	// Given: A shipment's move with orders in acceptable status

	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMember.DutyStationID = &dutyStation.ID
	serviceMember.DutyStation = dutyStation
	suite.MustSave(&serviceMember)

	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	deliveryDate := dates.NextWorkday(*calendar, pickupDate)
	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	market := "dHHG"
	sourceGBLOC := "BMLK"

	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED

	selectedMoveType := SelectedMoveTypeHHGPPM
	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			ServiceMember:           serviceMember,
			Move:                    *move,
			MoveID:                  move.ID,
		},
	})
	// Associate Shipment with the move it's on.
	move.Shipments = append(move.Shipments, shipment)
	move.Orders = orders

	// And: Move is in SUBMITTED state
	move.Status = MoveStatusSUBMITTED
	verrs, err = SaveMoveDependencies(suite.DB(), move)
	suite.False(verrs.HasAny(), "failed to save valid statuses")
	suite.Nil(err)
	suite.DB().Reload(&shipment)

	// Then: Shipment GBLOCs will be equal to:
	// destination GBLOC: orders' new duty station's transportation office's GBLOC
	suite.Equal(orders.NewDutyStation.TransportationOffice.Gbloc, *shipment.DestinationGBLOC)
	// source GBLOC: service member's current duty station's transportation office's GBLOC
	suite.Equal(serviceMember.DutyStation.TransportationOffice.Gbloc, *shipment.SourceGBLOC)
	// GBL number should be set
	suite.NotNil(shipment.GBLNumber)
}
