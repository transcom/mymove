package internalapi

import (
	"fmt"
	"github.com/go-openapi/swag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMoveHandlerAllValues() {
	// Given: a set of orders, user and servicemember
	orders := testdatagen.MakeDefaultOrder(suite.TestDB())

	req := httptest.NewRequest("POST", "/orders/orderid/moves", nil)
	req = suite.AuthenticateRequest(req, orders.ServiceMember)

	// When: a new Move is posted
	var selectedType = internalmessages.SelectedMoveTypePPM
	newMovePayload := &internalmessages.CreateMovePayload{
		SelectedMoveType: &selectedType,
	}
	params := moveop.CreateMoveParams{
		OrdersID:          strfmt.UUID(orders.ID.String()),
		CreateMovePayload: newMovePayload,
		HTTPRequest:       req,
	}
	// Then: we expect a move to have been created based on orders
	handler := CreateMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.CreateMoveCreated{}, response)
	okResponse := response.(*moveop.CreateMoveCreated)

	suite.Assertions.Equal(orders.ID.String(), okResponse.Payload.OrdersID.String())
}

func (suite *HandlerSuite) TestPatchMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}
	// And: a move is patched
	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	// And: Returned query to include our added move
	suite.Assertions.Equal(&newType, okResponse.Payload.SelectedMoveType)
}

func (suite *HandlerSuite) TestPatchMoveHandlerWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.TestDB())
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.TestDB())

	// And: the context contains a different user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseForbidden(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoMove() {
	// Given: a logged in user and no Move
	user := testdatagen.MakeDefaultServiceMember(suite.TestDB())

	moveUUID := uuid.Must(uuid.NewV4())

	// And: the context contains a logged in user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, user)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoType() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	patchPayload := internalmessages.PatchMovePayload{}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	suite.Assertions.Equal(move.ID.String(), okResponse.Payload.ID.String())
}

func (suite *HandlerSuite) TestShowMoveHandler() {

	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&moveop.ShowMoveOK{}, showResponse)
	okResponse := showResponse.(*moveop.ShowMoveOK)

	// And: Returned query to include our added move
	suite.Assertions.Equal(move.OrdersID.String(), okResponse.Payload.OrdersID.String())

}

func (suite *HandlerSuite) TestShowMoveWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.TestDB())
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.TestDB())

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	showResponse := showHandler.Handle(showMoveParams)
	// Then: expect a forbidden response
	suite.CheckResponseForbidden(showResponse)

}

func (suite *HandlerSuite) TestSubmitPPMMoveForApprovalHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakeDefaultPPM(suite.TestDB())
	move := ppm.Move

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.SubmitMoveForApprovalParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is submitted
	context := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	context.SetNotificationSender(notifications.NewStubNotificationSender(suite.TestLogger()))
	handler := SubmitMoveHandler{context}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
	okResponse := response.(*moveop.SubmitMoveForApprovalOK)

	// And: Returned query to have an approved status
	suite.Assertions.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)
	// And: Expect move's PPM's advance to have "Requested" status
	suite.Assertions.Equal(
		internalmessages.ReimbursementStatusREQUESTED,
		*okResponse.Payload.PersonallyProcuredMoves[0].Advance.Status)
}

func (suite *HandlerSuite) TestSubmitHHGMoveForApprovalHandler() {
	// Given: a set of orders, a move, user and servicemember
	shipment := testdatagen.MakeDefaultShipment(suite.TestDB())

	// There is not a way to set a field to nil using testdatagen.Assertions
	shipment.BookDate = nil
	suite.MustSave(&shipment)
	suite.Nil(shipment.BookDate)

	move := shipment.Move

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.SubmitMoveForApprovalParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: a move is submitted
	context := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	context.SetNotificationSender(notifications.NewStubNotificationSender(suite.TestLogger()))
	handler := SubmitMoveHandler{context}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
	okResponse := response.(*moveop.SubmitMoveForApprovalOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)
	suite.Equal(internalmessages.ShipmentStatusSUBMITTED, okResponse.Payload.Shipments[0].Status)
	suite.NotNil(okResponse.Payload.Shipments[0].BookDate)
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryHandler() {
	dutyStationAddress := testdatagen.MakeAddress(suite.TestDB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        swag.String("United States"),
		},
	})

	dutyStation := testdatagen.MakeDutyStation(suite.TestDB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Sam Houston",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   dutyStationAddress.ID,
			Address:     dutyStationAddress,
		},
	})

	rank := models.ServiceMemberRankE4
	serviceMember := testdatagen.MakeServiceMember(suite.TestDB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Rank:          &rank,
			DutyStationID: &dutyStation.ID,
			DutyStation:   dutyStation,
		},
	})

	newDutyStationAddress := testdatagen.MakeAddress(suite.TestDB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "n/a",
			City:           "San Antonio",
			State:          "TX",
			PostalCode:     "78234",
			Country:        swag.String("United States"),
		},
	})

	newDutyStation := testdatagen.MakeDutyStation(suite.TestDB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Gordon",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   newDutyStationAddress.ID,
			Address:     newDutyStationAddress,
		},
	})

	move := testdatagen.MakeMove(suite.TestDB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID:  serviceMember.ID,
			ServiceMember:    serviceMember,
			ReportByDate:     time.Date(2018, 10, 20, 0, 0, 0, 0, time.UTC),
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
		},
	})

	path := fmt.Sprintf("/moves/%s/move_dates", move.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	moveID := strfmt.UUID(move.ID.String())
	moveDate := strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC))
	params := moveop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      moveID,
		MoveDate:    moveDate,
	}

	context := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1125))

	showHandler := ShowMoveDatesSummaryHandler{context}
	response := showHandler.Handle(params)

	suite.IsType(&moveop.ShowMoveDatesSummaryOK{}, response)
	okResponse := response.(*moveop.ShowMoveDatesSummaryOK)

	id := move.ID.String() + ":" + moveDate.String()
	suite.Equal(id, *okResponse.Payload.ID)
	suite.Equal(moveID, *okResponse.Payload.MoveID)
	suite.Equal(moveDate, *okResponse.Payload.MoveDate)

	pack := []strfmt.Date{
		strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pack, okResponse.Payload.Pack)

	pickup := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pickup, okResponse.Payload.Pickup)

	transit := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 3, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 4, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 9, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(transit, okResponse.Payload.Transit)

	delivery := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 16, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(delivery, okResponse.Payload.Delivery)

	report := []strfmt.Date{
		strfmt.Date(move.Orders.ReportByDate),
	}
	suite.Equal(report, okResponse.Payload.Report)
}
