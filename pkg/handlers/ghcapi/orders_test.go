package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrderHandlerIntegration() {
	officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	params := orderop.GetOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetOrdersHandler{
		context,
		orderservice.NewOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	orderOK := response.(*orderop.GetOrderOK)
	ordersPayload := orderOK.Payload

	suite.Assertions.IsType(&orderop.GetOrderOK{}, response)
	suite.Equal(order.ID.String(), ordersPayload.ID.String())
	suite.Equal(move.Locator, ordersPayload.MoveCode)
	suite.Equal(order.ServiceMemberID.String(), ordersPayload.Customer.ID.String())
	suite.Equal(order.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
	suite.NotNil(order.NewDutyStation)
	payloadEntitlement := ordersPayload.Entitlement
	suite.Equal((*order.EntitlementID).String(), payloadEntitlement.ID.String())
	orderEntitlement := order.Entitlement
	suite.NotNil(orderEntitlement)
	suite.EqualValues(orderEntitlement.ProGearWeight, payloadEntitlement.ProGearWeight)
	suite.EqualValues(orderEntitlement.ProGearWeightSpouse, payloadEntitlement.ProGearWeightSpouse)
	suite.EqualValues(orderEntitlement.RequiredMedicalEquipmentWeight, payloadEntitlement.RequiredMedicalEquipmentWeight)
	suite.EqualValues(orderEntitlement.OrganizationalClothingAndIndividualEquipment, payloadEntitlement.OrganizationalClothingAndIndividualEquipment)
	suite.Equal(order.OriginDutyStation.ID.String(), ordersPayload.OriginDutyStation.ID.String())
	suite.NotZero(order.OriginDutyStation)
	suite.NotZero(ordersPayload.DateIssued)
}

func (suite *HandlerSuite) TestWeightAllowances() {
	suite.Run("With E-1 rank and no dependents", func() {
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Order: models.Order{
				ID:            uuid.Must(uuid.NewV4()),
				HasDependents: *swag.Bool(false),
			},
			Entitlement: models.Entitlement{
				ID:                   uuid.Must(uuid.NewV4()),
				DependentsAuthorized: swag.Bool(false),
				ProGearWeight:        2000,
				ProGearWeightSpouse:  500,
			},
		})
		request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
		params := orderop.GetOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
		}
		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("FetchOrder", order.ID).Return(&order, nil)

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetOrdersHandler{
			context,
			&orderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload
		payloadEntitlement := orderPayload.Entitlement
		orderEntitlement := order.Entitlement
		expectedAllowance := int64(orderEntitlement.WeightAllotment().TotalWeightSelf)

		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
		suite.Equal(expectedAllowance, payloadEntitlement.TotalWeight)
		suite.Equal(int64(*orderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	})

	suite.Run("With E-1 rank and dependents", func() {
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Order: models.Order{
				ID:            uuid.Must(uuid.NewV4()),
				HasDependents: *swag.Bool(true),
			},
		})

		request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
		params := orderop.GetOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
		}

		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("FetchOrder", order.ID).Return(&order, nil)

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetOrdersHandler{
			context,
			&orderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload
		payloadEntitlement := orderPayload.Entitlement
		orderEntitlement := order.Entitlement
		expectedAllowance := int64(orderEntitlement.WeightAllotment().TotalWeightSelfPlusDependents)

		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
		suite.Equal(expectedAllowance, payloadEntitlement.TotalWeight)
		suite.Equal(int64(*orderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	})
}

type updateOrderHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateOrderPayload
}

func (suite *HandlerSuite) makeUpdateOrderHandlerSubtestData() (subtestData *updateOrderHandlerSubtestData) {
	subtestData = &updateOrderHandlerSubtestData{}

	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())
	subtestData.order = subtestData.move.Orders

	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	subtestData.body = &ghcmessages.UpdateOrderPayload{
		DepartmentIndicator: &deptIndicator,
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    &ordersTypeDetail,
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("E19A"),
		Sac:                 handlers.FmtString("987654321"),
	}

	return subtestData
}

func (suite *HandlerSuite) TestUpdateOrderHandler() {
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := UpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(suite.DB()),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.UpdateOrderOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), ordersPayload.OriginDutyStation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.OrdersTypeDetail, ordersPayload.OrderTypeDetail)
		suite.Equal(body.OrdersNumber, ordersPayload.OrderNumber)
		suite.Equal(body.DepartmentIndicator, ordersPayload.DepartmentIndicator)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac, ordersPayload.Sac)
	})

	suite.Run("Returns a 403 when the user does not have TXO role", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderForbidden{}, response)
	})

	// We need to confirm whether a user who only has the TIO role should indeed
	// be authorized to update orders. If not, we also need to prevent them from
	// clicking the Edit Orders button in the frontend.
	suite.Run("Allows a TIO to update orders", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		move := subtestData.move
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsTOO", order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)
		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderOK{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateOrderEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateOrderPayload{}

	requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateOrderAsTOO", order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		updater,
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)

	orderOK := response.(*orderop.UpdateOrderOK)
	ordersPayload := orderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.IsType(&orderop.UpdateOrderOK{}, response)
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

type counselingUpdateOrderHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.CounselingUpdateOrderPayload
}

func (suite *HandlerSuite) makeCounselingUpdateOrderHandlerSubtestData() (subtestData *counselingUpdateOrderHandlerSubtestData) {
	subtestData = &counselingUpdateOrderHandlerSubtestData{}

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	subtestData.move = testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
	subtestData.order = subtestData.move.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	subtestData.body = &ghcmessages.CounselingUpdateOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
	}

	return subtestData
}

func (suite *HandlerSuite) TestCounselingUpdateOrderHandler() {
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := CounselingUpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(suite.DB()),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.CounselingUpdateOrderOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.CounselingUpdateOrderOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), ordersPayload.OriginDutyStation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(body.OrdersType, ordersPayload.OrderType)
	})

	suite.Run("Returns a 403 when the user does not have Counselor role", func() {
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderUnprocessableEntity{}, response)
	})
}

type updateAllowanceHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateAllowancePayload
}

func (suite *HandlerSuite) makeUpdateAllowanceHandlerSubtestData() (subtestData *updateAllowanceHandlerSubtestData) {
	subtestData = &updateAllowanceHandlerSubtestData{}

	subtestData.move = testdatagen.MakeServiceCounselingCompletedMove(suite.DB())
	subtestData.order = subtestData.move.Orders
	newAuthorizedWeight := int64(10000)
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.BranchAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	subtestData.body = &ghcmessages.UpdateAllowancePayload{
		Agency:               affiliation,
		AuthorizedWeight:     &newAuthorizedWeight,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
	}
	return subtestData
}

func (suite *HandlerSuite) TestUpdateAllowanceHandler() {

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := UpdateAllowanceHandler{
			context,
			orderservice.NewOrderUpdater(suite.DB()),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.UpdateAllowanceOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
	})

	suite.Run("Returns a 403 when the user does not have TOO role", func() {
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowancePreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceUnprocessableEntity{}, response)
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateAllowanceEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateAllowancePayload{}

	requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	params := orderop.UpdateAllowanceParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateAllowanceAsTOO", order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateAllowanceHandler{
		context,
		updater,
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)

	orderOK := response.(*orderop.UpdateAllowanceOK)
	ordersPayload := orderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.IsType(&orderop.UpdateAllowanceOK{}, response)
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

func (suite *HandlerSuite) TestCounselingUpdateAllowanceHandler() {
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.BranchAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	body := &ghcmessages.CounselingUpdateAllowancePayload{
		Agency:               affiliation,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := CounselingUpdateAllowanceHandler{
			context,
			orderservice.NewOrderUpdater(suite.DB()),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.CounselingUpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.CounselingUpdateAllowanceOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
	})

	suite.Run("Returns a 403 when the user does not have Counselor role", func() {
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowancePreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceUnprocessableEntity{}, response)
	})
}
