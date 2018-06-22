package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestApproveMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)
	// Given: and office User
	officeUser, _ := testdatagen.MakeOfficeUser(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveMoveOK{}, response)
	okResponse := response.(*officeop.ApproveMoveOK)

	// And: Returned query to have an approved status
	suite.Assertions.Equal(internalmessages.MoveStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestCancelMoveHandler() {
	// Given: a set of orders, a move, and office user
	orders, err := testdatagen.MakeOrder(suite.db)
	suite.Nil(err)
	var selectedType = internalmessages.SelectedMoveTypePPM
	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	officeUser, err := testdatagen.MakeOfficeUser(suite.db)
	suite.Nil(err)

	// Move is submitted
	err = move.Submit()
	suite.Nil(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// And: Orders are submitted and saved on move
	err = orders.Submit()
	suite.Nil(err)
	suite.Equal(models.OrderStatusSUBMITTED, orders.Status, "expected Submitted")
	suite.mustSave(&orders)
	move.Orders = orders
	suite.mustSave(move)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := officeop.CancelMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And params include the cancel reason
	reason := "Orders revoked."
	params.Reason = &reason

	// And: a move is canceled
	handler := CancelMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.CancelMoveOK{}, response)
	okResponse := response.(*officeop.CancelMoveOK)

	// And: Returned query to have an canceled status
	suite.Equal(internalmessages.MoveStatusCANCELED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApprovePPMHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm, _ := testdatagen.MakePPM(suite.db)
	officeUser, _ := testdatagen.MakeOfficeUser(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/personally_procured_moves/some_id/approve", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetSesService(suite.sesService)
	handler := ApprovePPMHandler(context)
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApprovePPMOK{}, response)
	okResponse := response.(*officeop.ApprovePPMOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.PPMStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveReimbursementHandler() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement, _ := testdatagen.MakeRequestedReimbursement(suite.db)
	officeUser, _ := testdatagen.MakeOfficeUser(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveReimbursementOK{}, response)
	okResponse := response.(*officeop.ApproveReimbursementOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.ReimbursementStatusAPPROVED, *okResponse.Payload.Status)
}
