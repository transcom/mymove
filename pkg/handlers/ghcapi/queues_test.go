package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveorder "github.com/transcom/mymove/pkg/services/move_order"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveQueuesHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Station GBLOC is LKNQ
	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create a move with an origin duty station outside of office user GBLOC
	excludedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "AGFM",
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	order := hhgMove.Orders
	result := payload.QueueMoves[0]

	suite.Len(payload.QueueMoves, 1)
	suite.Equal(order.ServiceMember.ID.String(), result.Customer.ID.String())
	suite.Equal(*order.DepartmentIndicator, string(result.DepartmentIndicator))
	suite.Equal(order.OriginDutyStation.TransportationOffice.Gbloc, string(result.OriginGBLOC))
	suite.Equal(order.NewDutyStation.ID.String(), result.DestinationDutyStation.ID.String())
	suite.Equal(hhgMove.Locator, result.Locator)
	suite.Equal(int64(1), result.ShipmentsCount)

}

func (suite *HandlerSuite) TestGetMoveQueuesFilter() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Create an order with AIR_FORCE department_indicator (default)
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create an order with ARMY department_indicator
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("ARMY"),
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
		Branch:      models.StringPointer("ARMY"),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	suite.Equal(1, len(payload.QueueMoves))
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerStatuses() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Station GBLOC is LKNQ
	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create a shipment on hhgMove that has Rejected status
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusRejected,
		},
	})

	// Create an order with an origin duty station outside of office user GBLOC
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Name:  "Fort Punxsutawney",
			Gbloc: "AGFM",
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
	})

	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	payload := response.(*queues.GetMovesQueueOK).Payload
	result := payload.QueueMoves[0]

	suite.Equal(ghcmessages.QueueMoveStatus("New move"), result.Status)

	// let's test for the Move approved status
	hhgMove.Status = models.MoveStatusAPPROVED
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.QueueMoveStatus("Move approved"), result.Status)

	// Now let's test Approvals requested
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusSubmitted,
		},
		Move: hhgMove,
	})

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.QueueMoveStatus("Approvals requested"), result.Status)

}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedRole() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedUser() {
	serviceUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceUser.User.Roles = append(serviceUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeCustomer,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateRequest(request, serviceUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerEmptyResults() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	// Create an order with an origin duty station outside of office user GBLOC
	hhgMoveType := models.SelectedMoveTypeHHG
	excludedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "AGFM",
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMovesQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	suite.Len(payload.QueueMoves, 0)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandler() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Station GBLOC is LKNQ
	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))

	actualPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		PaymentRequest: models.PaymentRequest{
			CreatedAt: prevCreatedAt,
		},
	})

	// Create an order with an origin duty station outside of office user GBLOC
	excludedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "AGFM",
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedPaymentRequest.MoveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		paymentrequest.NewPaymentRequestListFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	suite.Len(payload.QueuePaymentRequests, 1)

	paymentRequest := *payload.QueuePaymentRequests[0]

	suite.Equal(actualPaymentRequest.ID.String(), paymentRequest.ID.String())
	suite.Equal(actualPaymentRequest.MoveTaskOrderID.String(), paymentRequest.MoveID.String())
	suite.Equal(hhgMove.Orders.ServiceMemberID.String(), paymentRequest.Customer.ID.String())
	suite.Equal(actualPaymentRequest.Status.String(), string(paymentRequest.Status))

	createdAt := actualPaymentRequest.CreatedAt
	age := time.Since(createdAt).Hours() / 24.0

	suite.Equal(fmt.Sprintf("%.2f", age), fmt.Sprintf("%.2f", paymentRequest.Age))
	suite.Equal(createdAt.Format("2006-01-02T15:04:05.000Z07:00"), paymentRequest.SubmittedAt.String()) // swagger formats to milliseconds
	suite.Equal(hhgMove.Locator, paymentRequest.Locator)

	suite.Equal(*hhgMove.Orders.DepartmentIndicator, string(paymentRequest.DepartmentIndicator))
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerUnauthorizedRole() {
	officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		paymentrequest.NewPaymentRequestListFetcher(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerServerError() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	paymentRequestListFetcher := mocks.PaymentRequestListFetcher{}

	paymentRequestListFetcher.On("FetchPaymentRequestList", officeUser.ID).Return(nil, errors.New("database query error"))

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		&paymentRequestListFetcher,
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueInternalServerError{}, response)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerEmptyResults() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	paymentRequestListFetcher := mocks.PaymentRequestListFetcher{}

	paymentRequestListFetcher.On("FetchPaymentRequestList", officeUser.ID).Return(&models.PaymentRequests{}, nil)

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(query.NewQueryBuilder(context.DB())),
		&paymentRequestListFetcher,
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	suite.Len(payload.QueuePaymentRequests, 0)
	suite.Equal(int64(0), payload.TotalCount)
}
