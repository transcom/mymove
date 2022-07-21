package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
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
			SAC:                 handlers.FmtString("sac"),
			DepartmentIndicator: handlers.FmtString("17 Navy and Marine Corps"),
		},
	}
	move := testdatagen.MakeMove(suite.DB(), assertions)
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	moveRouter := moverouter.NewMoveRouter()

	// Move is submitted and saved
	err := moveRouter.Submit(suite.AppContextForTest(), &move)
	suite.NoError(err)
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
	handler := ApproveMoveHandler{
		suite.HandlerConfig(),
		moveRouter,
	}
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
	moveRouter := moverouter.NewMoveRouter()

	// Move is submitted and saved
	err := moveRouter.Submit(suite.AppContextForTest(), &move)
	suite.NoError(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
	suite.MustSave(&move)

	move.Orders.OrdersNumber = nil
	move.Orders.OrdersTypeDetail = nil
	move.Orders.DepartmentIndicator = nil
	suite.MustSave(&move.Orders)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: move handler is hit
	handler := ApproveMoveHandler{
		suite.HandlerConfig(),
		moveRouter,
	}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&officeop.ApprovePPMBadRequest{}, response)
}

func (suite *HandlerSuite) TestApproveMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := testdatagen.MakeDefaultMove(suite.DB())
	// Given: an non-office User
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	moveRouter := moverouter.NewMoveRouter()

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.AuthenticateRequest(req, user)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler{
		suite.HandlerConfig(),
		moveRouter,
	}
	response := handler.Handle(params)

	// Then: response is Forbidden
	suite.Assertions.IsType(&officeop.ApproveMoveForbidden{}, response)
}

func (suite *HandlerSuite) TestCancelMoveHandler() {
	// Given: a set of orders, a move, and office user
	// Orders has service member with transportation office and phone nums
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())
	moveRouter := moverouter.NewMoveRouter()

	selectedMoveType := models.SelectedMoveTypePPM
	moveOptions := models.MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	officeUser := testdatagen.MakeStubbedOfficeUser(suite.DB())
	suite.NoError(err)

	// Move is submitted
	err = moveRouter.Submit(suite.AppContextForTest(), move)
	suite.NoError(err)
	suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")

	// And: Orders are submitted and saved on move
	err = orders.Submit()
	suite.NoError(err)
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
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetNotificationSender(suite.TestNotificationSender())
	handler := CancelMoveHandler{handlerConfig, moveRouter}
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

	moveRouter := moverouter.NewMoveRouter()

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
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetNotificationSender(suite.TestNotificationSender())
	handler := CancelMoveHandler{handlerConfig, moveRouter}
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
	approveDate := strfmt.DateTime(time.Now())

	newApprovePersonallyProcuredMovePayload := internalmessages.ApprovePersonallyProcuredMovePayload{
		ApproveDate: &approveDate,
	}
	params := officeop.ApprovePPMParams{
		HTTPRequest:                          req,
		PersonallyProcuredMoveID:             strfmt.UUID(ppm.ID.String()),
		ApprovePersonallyProcuredMovePayload: &newApprovePersonallyProcuredMovePayload,
	}

	// And: a ppm is approved
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetNotificationSender(suite.TestNotificationSender())
	handler := ApprovePPMHandler{handlerConfig}
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
	approveDate := strfmt.DateTime(time.Now())

	newApprovePersonallyProcuredMovePayload := internalmessages.ApprovePersonallyProcuredMovePayload{
		ApproveDate: &approveDate,
	}
	params := officeop.ApprovePPMParams{
		HTTPRequest:                          req,
		PersonallyProcuredMoveID:             strfmt.UUID(ppm.ID.String()),
		ApprovePersonallyProcuredMovePayload: &newApprovePersonallyProcuredMovePayload,
	}

	// And: a ppm is approved
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetNotificationSender(suite.TestNotificationSender())
	handler := ApprovePPMHandler{handlerConfig}
	response := handler.Handle(params)

	// Then: expect a Forbidden status code
	suite.Assertions.IsType(&officeop.ApprovePPMForbidden{}, response)
}

// TODO: Determine whether we need to complete remove reimbursements handler from Office handlers
func (suite *HandlerSuite) TestApproveReimbursementHandler() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement := testdatagen.MakeDefaultRequestedReimbursement(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveReimbursementOK{}, response)
	okResponse := response.(*officeop.ApproveReimbursementOK)

	// And: Returned query to have an approved status
	suite.Equal(internalmessages.ReimbursementStatusAPPROVED, *okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveReimbursementHandlerForbidden() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement := testdatagen.MakeDefaultRequestedReimbursement(suite.DB())
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/reimbursement/some_id/approve", nil)
	req = suite.AuthenticateRequest(req, user)
	params := officeop.ApproveReimbursementParams{
		HTTPRequest:     req,
		ReimbursementID: strfmt.UUID(reimbursement.ID.String()),
	}

	// And: a reimbursement is approved
	handler := ApproveReimbursementHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	// Then: expect Forbidden response
	suite.Assertions.IsType(&officeop.ApproveReimbursementForbidden{}, response)
}
