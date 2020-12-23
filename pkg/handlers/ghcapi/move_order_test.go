package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveorder "github.com/transcom/mymove/pkg/services/move_order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveOrderHandlerIntegration() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	request := httptest.NewRequest("GET", "/move-orders/{moveOrderID}", nil)
	params := moveorderop.GetMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveOrdersHandler{
		context,
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveOrderOK := response.(*moveorderop.GetMoveOrderOK)
	moveOrdersPayload := moveOrderOK.Payload

	suite.Assertions.IsType(&moveorderop.GetMoveOrderOK{}, response)
	suite.Equal(moveOrder.ID.String(), moveOrdersPayload.ID.String())
	suite.Equal(moveOrder.ServiceMemberID.String(), moveOrdersPayload.CustomerID.String())
	suite.Equal(moveOrder.NewDutyStationID.String(), moveOrdersPayload.DestinationDutyStation.ID.String())
	suite.NotNil(moveOrder.NewDutyStation)
	payloadEntitlement := moveOrdersPayload.Entitlement
	suite.Equal((*moveOrder.EntitlementID).String(), payloadEntitlement.ID.String())
	moveOrderEntitlement := moveOrder.Entitlement
	suite.NotNil(moveOrderEntitlement)
	suite.Equal(moveOrder.OriginDutyStation.ID.String(), moveOrdersPayload.OriginDutyStation.ID.String())
	suite.NotZero(moveOrder.OriginDutyStation)
	suite.NotZero(moveOrdersPayload.DateIssued)
}

func (suite *HandlerSuite) TestWeightAllowances() {
	suite.Run("With E-1 rank and no dependents", func() {
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Order: models.Order{
				HasDependents: *swag.Bool(false),
			},
			Entitlement: models.Entitlement{
				ID:                   uuid.Must(uuid.NewV4()),
				DependentsAuthorized: swag.Bool(false),
			},
		})
		request := httptest.NewRequest("GET", "/move-orders/{moveOrderID}", nil)
		params := moveorderop.GetMoveOrderParams{
			HTTPRequest: request,
			MoveOrderID: strfmt.UUID(order.ID.String()),
		}
		moveOrderFetcher := mocks.MoveOrderFetcher{}
		moveOrderFetcher.On("FetchMoveOrder", order.ID).Return(&order, nil)

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetMoveOrdersHandler{
			context,
			&moveOrderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*moveorderop.GetMoveOrderOK)
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
				HasDependents: *swag.Bool(true),
			},
		})

		request := httptest.NewRequest("GET", "/move-orders/{orderID}", nil)
		params := moveorderop.GetMoveOrderParams{
			HTTPRequest: request,
			MoveOrderID: strfmt.UUID(order.ID.String()),
		}
		moveOrderFetcher := mocks.MoveOrderFetcher{}
		moveOrderFetcher.On("FetchMoveOrder", order.ID).Return(&order, nil)

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetMoveOrdersHandler{
			context,
			&moveOrderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*moveorderop.GetMoveOrderOK)
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

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerIntegration() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")

	newAuthorizedWeight := int64(10000)
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	affiliation := ghcmessages.BranchAIRFORCE
	grade := ghcmessages.GradeO5
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	body := &ghcmessages.UpdateMoveOrderPayload{
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

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
		IfMatch:     etag.GenerateEtag(moveOrder.UpdatedAt),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(context.DB())
	handler := UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(suite.DB(), queryBuilder),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveOrderOK := response.(*moveorderop.UpdateMoveOrderOK)
	moveOrdersPayload := moveOrderOK.Payload

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderOK{}, response)
	suite.Equal(moveOrder.ID.String(), moveOrdersPayload.ID.String())
	suite.Equal(body.NewDutyStationID.String(), moveOrdersPayload.DestinationDutyStation.ID.String())
	suite.Equal(body.OriginDutyStationID.String(), moveOrdersPayload.OriginDutyStation.ID.String())
	suite.Equal(*body.IssueDate, moveOrdersPayload.DateIssued)
	suite.Equal(*body.ReportByDate, moveOrdersPayload.ReportByDate)
	suite.Equal(body.OrdersType, moveOrdersPayload.OrderType)
	suite.Equal(body.OrdersTypeDetail, moveOrdersPayload.OrderTypeDetail)
	suite.Equal(body.OrdersNumber, moveOrdersPayload.OrderNumber)
	suite.Equal(body.DepartmentIndicator, moveOrdersPayload.DepartmentIndicator)
	suite.Equal(body.Tac, moveOrdersPayload.Tac)
	suite.Equal(body.Sac, moveOrdersPayload.Sac)
	suite.Equal(body.AuthorizedWeight, moveOrdersPayload.Entitlement.AuthorizedWeight)
	suite.Equal(body.Grade, moveOrdersPayload.Grade)
	suite.Equal(body.Agency, moveOrdersPayload.Agency)
	suite.Equal(body.DependentsAuthorized, moveOrdersPayload.Entitlement.DependentsAuthorized)
}

// Test that a move order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateMoveOrderEventTrigger() {
	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateMoveOrderPayload{
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

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
		IfMatch:     etag.GenerateEtag(moveOrder.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(context.DB())
	// Set up handler:
	handler := UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(suite.DB(), queryBuilder),
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB
	suite.IsNotErrResponse(response)
	moveOrderOK := response.(*moveorderop.UpdateMoveOrderOK)
	moveOrdersPayload := moveOrderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderOK{}, response)
	suite.Equal(moveOrdersPayload.ID, strfmt.UUID(moveOrder.ID.String()))
	suite.HasWebhookNotification(moveOrder.ID, traceID)
}

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerNotFound() {
	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: "8d013ebb-9561-467b-ae6d-853d2bceadde",
		IfMatch:     "",
		Body: &ghcmessages.UpdateMoveOrderPayload{
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
	queryBuilder := query.NewQueryBuilder(context.DB())
	handler := UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(suite.DB(), queryBuilder),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderNotFound{}, response)
}

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerPreconditionsFailed() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateMoveOrderPayload{
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

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
		IfMatch:     etag.GenerateEtag(moveOrder.UpdatedAt.Add(time.Second * 30)),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(context.DB())
	handler := UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(suite.DB(), queryBuilder),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderPreconditionFailed{}, response)
}

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerBadRequest() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicator("COAST_GUARD")
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")

	body := &ghcmessages.UpdateMoveOrderPayload{
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

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
		IfMatch:     etag.GenerateEtag(moveOrder.UpdatedAt.Add(time.Second * 30)),
		Body:        body,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(context.DB())
	handler := UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(suite.DB(), queryBuilder),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderPreconditionFailed{}, response)
}
