package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestApproveMoveHandler() {
	// Given: a set of complete orders, a move, office user and servicemember user
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	assertions := testdatagen.Assertions{
		Order: models.Order{
			OrdersNumber:        handlers.FmtString("1234"),
			OrdersTypeDetail:    &hhgPermitted,
			TAC:                 handlers.FmtString("1234"),
			DepartmentIndicator: handlers.FmtString("17 - United States Marines"),
		},
	}
	move := testdatagen.MakeMove(suite.DB(), assertions)
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Move is submitted and saved
	err := move.Submit()
	suite.Nil(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.MustSave(&move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveMoveOK{}, response)
	okResponse := response.(*officeop.ApproveMoveOK)

	// And: Returned query to have an approved status
	suite.Assertions.Equal(internalmessages.MoveStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveMoveHandlerIncompleteOrders() {
	// Given: a set of incomplete orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.DB())
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Move is submitted and saved
	err := move.Submit()
	suite.Nil(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.MustSave(&move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: move handler is hit
	handler := ApproveMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&officeop.ApprovePPMBadRequest{}, response)
}

func (suite *HandlerSuite) TestApproveMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.DB())
	// Given: an non-office User
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.AuthenticateRequest(req, user)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: response is Forbidden
	suite.Assertions.IsType(&officeop.ApproveMoveForbidden{}, response)
}
func (suite *HandlerSuite) TestCancelMoveHandler() {
	// Given: a set of orders, a move, and office user
	// Orders has service member with transportation office and phone nums
	orders := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := models.SelectedMoveTypePPM
	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	suite.Nil(err)

	// Move is submitted
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// And: Orders are submitted and saved on move
	err = orders.Submit()
	suite.Nil(err)
	suite.Equal(models.OrderStatusSUBMITTED, orders.Status, "expected Submitted")
	suite.MustSave(&orders)
	move.Orders = orders
	suite.MustSave(move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	// And params include the cancel reason
	reason := "Orders revoked."
	reasonPayload := &internalmessages.CancelMove{
		CancelReason: &reason,
	}
	params := officeop.CancelMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		CancelMove:  reasonPayload,
	}

	// And: a move is canceled
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetNotificationSender(suite.TestNotificationSender())
	handler := CancelMoveHandler{context}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.CancelMoveOK{}, response)
	okResponse := response.(*officeop.CancelMoveOK)

	// And: Returned query to have an canceled status
	suite.Equal(internalmessages.MoveStatusCANCELED, okResponse.Payload.Status)
}
func (suite *HandlerSuite) TestCancelMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.DB())
	// Given: an non-office User
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
	req = suite.AuthenticateRequest(req, user)

	// And params include the cancel reason
	reason := "Orders revoked."
	reasonPayload := &internalmessages.CancelMove{
		CancelReason: &reason,
	}
	params := officeop.CancelMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		CancelMove:  reasonPayload,
	}
	// And: a move is canceled
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetNotificationSender(suite.TestNotificationSender())
	handler := CancelMoveHandler{context}
	response := handler.Handle(params)

	// Then: response is Forbidden
	suite.Assertions.IsType(&officeop.CancelMoveForbidden{}, response)
}
func (suite *HandlerSuite) TestApprovePPMHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusSUBMITTED,
		},
	})

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/personally_procured_moves/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetNotificationSender(suite.TestNotificationSender())
	handler := ApprovePPMHandler{context}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApprovePPMOK{}, response)
	okResponse := response.(*officeop.ApprovePPMOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.PPMStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApprovePPMHandlerForbidden() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/personally_procured_moves/some_id/approve", nil)
	req = suite.AuthenticateRequest(req, user)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetNotificationSender(suite.TestNotificationSender())
	handler := ApprovePPMHandler{context}
	response := handler.Handle(params)

	// Then: expect a Forbidden status code
	suite.Assertions.IsType(&officeop.ApprovePPMForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveReimbursementHandler() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement, _ := testdatagen.MakeRequestedReimbursement(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveReimbursementOK{}, response)
	okResponse := response.(*officeop.ApproveReimbursementOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.ReimbursementStatusAPPROVED, *okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveReimbursementHandlerForbidden() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement, _ := testdatagen.MakeRequestedReimbursement(suite.DB())
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.AuthenticateRequest(req, user)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect Forbidden response
	suite.Assertions.IsType(&officeop.ApproveReimbursementForbidden{}, response)
}
