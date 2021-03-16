package ghcapi

import (
	"net/http/httptest"
	"time"

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
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	order := moveTaskOrder.Orders
	request := httptest.NewRequest("GET", "/move-orders/{orderID}", nil)
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
	suite.Equal(moveTaskOrder.Locator, ordersPayload.MoveCode)
	suite.Equal(order.ServiceMemberID.String(), ordersPayload.Customer.ID.String())
	suite.Equal(order.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
	suite.NotNil(order.NewDutyStation)
	payloadEntitlement := ordersPayload.Entitlement
	suite.Equal((*order.EntitlementID).String(), payloadEntitlement.ID.String())
	orderEntitlement := order.Entitlement
	suite.NotNil(orderEntitlement)
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
			},
		})
		request := httptest.NewRequest("GET", "/move-orders/{orderID}", nil)
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

		request := httptest.NewRequest("GET", "/move-orders/{orderID}", nil)
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

func (suite *HandlerSuite) TestUpdateOrderHandlerIntegration() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	order := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	request := httptest.NewRequest("PATCH", "/move-orders/{orderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")

	newAuthorizedWeight := int64(10000)
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	affiliation := ghcmessages.BranchAIRFORCE
	grade := ghcmessages.GradeO5
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	body := &ghcmessages.UpdateOrderPayload{
		AuthorizedWeight:     &newAuthorizedWeight,
		Agency:               affiliation,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		IssueDate:            handlers.FmtDatePtr(&issueDate),
		ReportByDate:         handlers.FmtDatePtr(&reportByDate),
		OrdersType:           "RETIREMENT",
		OrdersTypeDetail:     &ordersTypeDetail,
		DepartmentIndicator:  &deptIndicator,
		OrdersNumber:         handlers.FmtString("ORDER100"),
		NewDutyStationID:     handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID:  handlers.FmtUUID(originDutyStation.ID),
		Tac:                  handlers.FmtString("012345678"),
		Sac:                  handlers.FmtString("987654321"),
	}

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
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
	suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
	suite.Equal(body.Grade, ordersPayload.Grade)
	suite.Equal(body.Agency, ordersPayload.Agency)
	suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
}

// Test that a move order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateOrderEventTrigger() {
	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())
	order := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{orderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    &ordersTypeDetail,
		DepartmentIndicator: &deptIndicator,
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("012345678"),
		Sac:                 handlers.FmtString("987654321"),
	}

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	// Set up handler:
	handler := UpdateOrderHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB
	suite.IsNotErrResponse(response)
	orderOK := response.(*orderop.UpdateOrderOK)
	ordersPayload := orderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")

	suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerNotFound() {
	request := httptest.NewRequest("PATCH", "/move-orders/{orderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     "8d013ebb-9561-467b-ae6d-853d2bceadde",
		IfMatch:     "",
		Body: &ghcmessages.UpdateOrderPayload{
			IssueDate:           handlers.FmtDatePtr(&issueDate),
			ReportByDate:        handlers.FmtDatePtr(&reportByDate),
			OrdersType:          "RETIREMENT",
			OrdersTypeDetail:    &ordersTypeDetail,
			DepartmentIndicator: &deptIndicator,
			OrdersNumber:        handlers.FmtString("ORDER100"),
			NewDutyStationID:    handlers.FmtUUID(uuid.Nil),
			OriginDutyStationID: handlers.FmtUUID(uuid.Nil),
			Tac:                 handlers.FmtString("012345678"),
			Sac:                 handlers.FmtString("987654321"),
		},
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderNotFound{}, response)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerPreconditionsFailed() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	order := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{orderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    &ordersTypeDetail,
		DepartmentIndicator: &deptIndicator,
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("012345678"),
		Sac:                 handlers.FmtString("987654321"),
	}

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt.Add(time.Second * 30)),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerBadRequest() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	order := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{orderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    &ordersTypeDetail,
		DepartmentIndicator: &deptIndicator,
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(uuid.Nil), // An unknown duty station will result in a invalid input error
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("012345678"),
		Sac:                 handlers.FmtString("987654321"),
	}

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt.Add(time.Second * 30)),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
}
