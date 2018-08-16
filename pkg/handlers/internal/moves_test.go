package internal

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *utils.HandlerSuite) TestCreateMoveHandlerAllValues() {
	// Given: a set of orders, user and servicemember
	orders := testdatagen.MakeDefaultOrder(suite.db)

	req := httptest.NewRequest("POST", "/orders/orderid/moves", nil)
	req = suite.authenticateRequest(req, orders.ServiceMember)

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
	handler := CreateMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.CreateMoveCreated{}, response)
	okResponse := response.(*moveop.CreateMoveCreated)

	suite.Assertions.Equal(orders.ID.String(), okResponse.Payload.OrdersID.String())
}

func (suite *utils.HandlerSuite) TestPatchMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

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
	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	// And: Returned query to include our added move
	suite.Assertions.Equal(&newType, okResponse.Payload.SelectedMoveType)
}

func (suite *utils.HandlerSuite) TestPatchMoveHandlerWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.db)
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.db)

	// And: the context contains a different user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, anotherUser)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseForbidden(response)
}

func (suite *utils.HandlerSuite) TestPatchMoveHandlerNoMove() {
	// Given: a logged in user and no Move
	user := testdatagen.MakeDefaultServiceMember(suite.db)

	moveUUID := uuid.Must(uuid.NewV4())

	// And: the context contains a logged in user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, user)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseNotFound(response)
}

func (suite *utils.HandlerSuite) TestPatchMoveHandlerNoType() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

	patchPayload := internalmessages.PatchMovePayload{}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	suite.Assertions.Equal(move.ID.String(), okResponse.Payload.ID.String())
}

func (suite *utils.HandlerSuite) TestShowMoveHandler() {

	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&moveop.ShowMoveOK{}, showResponse)
	okResponse := showResponse.(*moveop.ShowMoveOK)

	// And: Returned query to include our added move
	suite.Assertions.Equal(move.OrdersID.String(), okResponse.Payload.OrdersID.String())

}

func (suite *utils.HandlerSuite) TestShowMoveWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.db)
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.db)

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, anotherUser)

	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showMoveParams)
	// Then: expect a forbidden response
	suite.checkResponseForbidden(showResponse)

}

func (suite *utils.HandlerSuite) TestSubmitPPMMoveForApprovalHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakeDefaultPPM(suite.db)
	move := ppm.Move

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.SubmitMoveForApprovalParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is submitted
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetNotificationSender(notifications.NewStubNotificationSender(suite.logger))
	handler := SubmitMoveHandler(context)
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
