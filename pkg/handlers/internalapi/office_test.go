package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestApproveMoveHandler() {
	// Given: a set of complete orders, a move, office user and servicemember user
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersNumber:        handlers.FmtString("1234"),
				OrdersTypeDetail:    &hhgPermitted,
				TAC:                 handlers.FmtString("1234"),
				SAC:                 handlers.FmtString("sac"),
				DepartmentIndicator: handlers.FmtString("17 Navy and Marine Corps"),
			},
		},
	}, nil)
	// Given: an office User
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	moveRouter := moverouter.NewMoveRouter()

	// Move is submitted and saved
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
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
	move := factory.BuildMove(suite.DB(), nil, nil)
	// Given: an office User
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	moveRouter := moverouter.NewMoveRouter()

	// Move is submitted and saved
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
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
	suite.Assertions.IsType(&officeop.ApproveMoveBadRequest{}, response)
}

func (suite *HandlerSuite) TestApproveMoveHandlerForbidden() {
	// Given: a set of orders, a move, office user and servicemember user
	move := factory.BuildMove(suite.DB(), nil, nil)
	// Given: an non-office User
	user := factory.BuildServiceMember(suite.DB(), nil, nil)
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
	suite.Run("Successfully cancels move", func() {
		// Given: a set of orders, a move, and office user
		// Orders has service member with transportation office and phone nums
		moveRouter := moverouter.NewMoveRouter()

		// Given: a set of orders, a move, user and servicemember
		move := factory.BuildMove(suite.DB(), nil, nil)

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := officeop.CancelMoveParams{
			HTTPRequest: req,
			MoveID:      strfmt.UUID(move.ID.String()),
		}

		// And: a move is canceled
		handlerConfig := suite.HandlerConfig()
		handler := CancelMoveHandler{handlerConfig, moveRouter}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&officeop.CancelMoveOK{}, response)
		okResponse := response.(*officeop.CancelMoveOK)

		// And: Returned query to have an canceled status
		suite.Equal(internalmessages.MoveStatusCANCELED, okResponse.Payload.Status)
	})

	suite.Run("Fails to cancel someone elses move", func() {
		// Given: a set of orders, a move, and office user
		// Orders has service member with transportation office and phone nums
		moveRouter := moverouter.NewMoveRouter()

		// Given: a set of orders, a move, user and servicemember
		move := factory.BuildMove(suite.DB(), nil, nil)
		other_user := factory.BuildServiceMember(suite.DB(), nil, nil)

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
		req = suite.AuthenticateRequest(req, other_user)

		params := officeop.CancelMoveParams{
			HTTPRequest: req,
			MoveID:      strfmt.UUID(move.ID.String()),
		}

		// And: a move is canceled
		handlerConfig := suite.HandlerConfig()
		handler := CancelMoveHandler{handlerConfig, moveRouter}
		response := handler.Handle(params)

		// Then: expect a 403 status code
		suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	})

	suite.Run("Fails to cancel submitted move", func() {
		// Given: a set of orders, a move, and office user
		// Orders has service member with transportation office and phone nums
		moveRouter := moverouter.NewMoveRouter()

		// Given: a set of orders, a move, user and servicemember
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
				},
			},
		}, nil)

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/cancel", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := officeop.CancelMoveParams{
			HTTPRequest: req,
			MoveID:      strfmt.UUID(move.ID.String()),
		}

		// And: a move is canceled
		handlerConfig := suite.HandlerConfig()
		handler := CancelMoveHandler{handlerConfig, moveRouter}
		response := handler.Handle(params)

		// Then: expect a error status code
		suite.Assertions.IsType(&officeop.ApproveMoveConflict{}, response)
	})
}

// TODO: Determine whether we need to complete remove reimbursements handler from Office handlers
func (suite *HandlerSuite) TestApproveReimbursementHandler() {
	// Given: a set of orders, a move, user and servicemember
	reimbursement := testdatagen.MakeDefaultRequestedReimbursement(suite.DB())
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

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
	user := factory.BuildServiceMember(suite.DB(), nil, nil)

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
