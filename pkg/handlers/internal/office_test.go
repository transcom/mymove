package internal

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestApproveMoveHandler() {
	// Given: a set of complete orders, a move, office user and servicemember user
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	assertions := testdatagen.Assertions{
		Order: models.Order{
			OrdersNumber:        utils.FmtString("1234"),
			OrdersTypeDetail:    &hhgPermitted,
			TAC:                 utils.FmtString("1234"),
			DepartmentIndicator: utils.FmtString("17 - United States Marines"),
		},
	}
	move := testdatagen.MakeMove(suite.parent.Db, assertions)
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	// Move is submitted and saved
	err := move.Submit()
	suite.parent.Nil(err)
	suite.parent.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.parent.MustSave(&move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&officeop.ApproveMoveOK{}, response)
	okResponse := response.(*officeop.ApproveMoveOK)

	// And: Returned query to have an approved status
	suite.parent.Assertions.Equal(internalmessages.MoveStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveMoveHandlerIncompleteOrders() {
	// Given: a set of incomplete orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.parent.Db)
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	// Move is submitted and saved
	err := move.Submit()
	suite.parent.Nil(err)
	suite.parent.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.parent.MustSave(&move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: move handler is hit
	handler := ApproveMoveHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.parent.Assertions.IsType(&officeop.ApprovePPMBadRequest{}, response)
}

func (suite *HandlerSuite) TestApproveMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.parent.Db)
	// Given: an non-office User
	user := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.parent.AuthenticateRequest(req, user)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: response is Forbidden
	suite.parent.Assertions.IsType(&officeop.ApproveMoveForbidden{}, response)
}
func (suite *HandlerSuite) TestCancelMoveHandler() {
	// Given: a set of orders, a move, and office user
	// Orders has service member with transportation office and phone nums
	orders := testdatagen.MakeDefaultOrder(suite.parent.Db)

	var selectedType = internalmessages.SelectedMoveTypePPM
	move, verrs, err := orders.CreateNewMove(suite.parent.Db, &selectedType)
	suite.parent.Nil(err)
	suite.parent.False(verrs.HasAny(), "failed to validate move")
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)
	suite.parent.Nil(err)

	// Move is submitted
	err = move.Submit()
	suite.parent.Nil(err)
	suite.parent.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// And: Orders are submitted and saved on move
	err = orders.Submit()
	suite.parent.Nil(err)
	suite.parent.Equal(models.OrderStatusSUBMITTED, orders.Status, "expected Submitted")
	suite.parent.MustSave(&orders)
	move.Orders = orders
	suite.parent.MustSave(move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

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
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetNotificationSender(suite.parent.NotificationSender)
	handler := CancelMoveHandler(context)
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&officeop.CancelMoveOK{}, response)
	okResponse := response.(*officeop.CancelMoveOK)

	// And: Returned query to have an canceled status
	suite.parent.Equal(internalmessages.MoveStatusCANCELED, okResponse.Payload.Status)
}
func (suite *HandlerSuite) TestCancelMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.parent.Db)
	// Given: an non-office User
	user := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
	req = suite.parent.AuthenticateRequest(req, user)

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
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetNotificationSender(suite.parent.NotificationSender)
	handler := CancelMoveHandler(context)
	response := handler.Handle(params)

	// Then: response is Forbidden
	suite.parent.Assertions.IsType(&officeop.CancelMoveForbidden{}, response)
}
func (suite *HandlerSuite) TestApprovePPMHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/personally_procured_moves/some_id/approve", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetNotificationSender(suite.parent.NotificationSender)
	handler := ApprovePPMHandler(context)
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&officeop.ApprovePPMOK{}, response)
	okResponse := response.(*officeop.ApprovePPMOK)

	// And: Returned query to have an approved status
	suite.parent.Equal(internalmessages.PPMStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApprovePPMHandlerForbidden() {
	// Given: a set of orders, a move, user and servicemember
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)
	user := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/personally_procured_moves/some_id/approve", nil)
	req = suite.parent.AuthenticateRequest(req, user)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetNotificationSender(suite.parent.NotificationSender)
	handler := ApprovePPMHandler(context)
	response := handler.Handle(params)

	// Then: expect a Forbidden status code
	suite.parent.Assertions.IsType(&officeop.ApprovePPMForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveReimbursementHandler() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement, _ := testdatagen.MakeRequestedReimbursement(suite.parent.Db)
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&officeop.ApproveReimbursementOK{}, response)
	okResponse := response.(*officeop.ApproveReimbursementOK)

	// And: Returned query to have an approved status
	suite.parent.Equal(internalmessages.ReimbursementStatusAPPROVED, *okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveReimbursementHandlerForbidden() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement, _ := testdatagen.MakeRequestedReimbursement(suite.parent.Db)
	user := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.parent.AuthenticateRequest(req, user)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect Forbidden response
	suite.parent.Assertions.IsType(&officeop.ApproveReimbursementForbidden{}, response)
}
