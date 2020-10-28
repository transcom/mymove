package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/gofrs/uuid"

	moveorder "github.com/transcom/mymove/pkg/services/move_order"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/etag"
	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestListMoveOrdersHandler() {

	moveOrderID1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	moveOrderID2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	moveOrderID3, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	IDs := []uuid.UUID{
		moveOrderID1,
		moveOrderID2,
		moveOrderID3,
	}

	var moveOrders []models.Order

	for _, id := range IDs {
		moveOrder := models.Order{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		moveOrders = append(moveOrders, moveOrder)
	}

	suite.T().Run("When office user is TOO and fetch is successful", func(t *testing.T) {
		moveOrderFetcher := &mocks.MoveOrderFetcher{}
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{
			Stub: true,
		})

		moveOrderFetcher.On("ListMoveOrders", officeUser.ID).Return(moveOrders, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/move_orders"), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := moveorderop.ListMoveOrdersParams{
			HTTPRequest: req,
		}

		handler := ListMoveOrdersHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			moveOrderFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&moveorderop.ListMoveOrdersOK{}, response)
		okResponse := response.(*moveorderop.ListMoveOrdersOK)
		suite.Equal(len(IDs), len(okResponse.Payload))
		suite.Equal(moveOrderID1.String(), okResponse.Payload[0].ID.String())
		suite.Equal(moveOrderID2.String(), okResponse.Payload[1].ID.String())
		suite.Equal(moveOrderID3.String(), okResponse.Payload[2].ID.String())
	})

	suite.T().Run("When office user is not TOO, response should be 403", func(t *testing.T) {
		moveOrderFetcher := &mocks.MoveOrderFetcher{}
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			Stub: true,
		})
		moveOrderFetcher.AssertNumberOfCalls(t, "ListMoveOrders", 0)

		req := httptest.NewRequest("GET", fmt.Sprintf("/move_orders"), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := moveorderop.ListMoveOrdersParams{
			HTTPRequest: req,
		}

		handler := ListMoveOrdersHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			moveOrderFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&moveorderop.ListMoveOrdersForbidden{}, response)
	})
}

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
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().TotalWeightSelf), payloadEntitlement.TotalWeight)
	suite.Equal(int64(*moveOrderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	suite.Equal(moveOrder.OriginDutyStation.ID.String(), moveOrdersPayload.OriginDutyStation.ID.String())
	suite.NotZero(moveOrder.OriginDutyStation)
	suite.NotZero(moveOrdersPayload.DateIssued)
}

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerIntegration() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	moveOrder := moveTaskOrder.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")

	body := &ghcmessages.UpdateMoveOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    "INSTRUCTION_20_WEEKS",
		DepartmentIndicator: "COAST_GUARD",
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("012345678"),
		Sac:                 handlers.FmtString("987654321"),
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

	body := &ghcmessages.UpdateMoveOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    "INSTRUCTION_20_WEEKS",
		DepartmentIndicator: "COAST_GUARD",
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
	moveOrderOK := response.(*moveorderop.UpdateMoveOrderOK) // Type assertion: Take the response, change the type to *moveorderop.UpdateMoveOrderOK, and set the variable moveOrderOK to the new type
	moveOrdersPayload := moveOrderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderOK{}, response)
	suite.Equal(moveOrdersPayload.ID, strfmt.UUID(moveOrder.ID.String()))
	suite.HasWebhookNotification(moveOrder.ID, traceID)
	// suite.HasWebhookNotification(updatedOrder.ID, traceID)
}

func (suite *HandlerSuite) TestUpdateMoveOrderHandlerNotFound() {
	request := httptest.NewRequest("PATCH", "/move-orders/{moveOrderID}", nil)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")

	params := moveorderop.UpdateMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: "8d013ebb-9561-467b-ae6d-853d2bceadde",
		IfMatch:     "",
		Body: &ghcmessages.UpdateMoveOrderPayload{
			IssueDate:           handlers.FmtDatePtr(&issueDate),
			ReportByDate:        handlers.FmtDatePtr(&reportByDate),
			OrdersType:          "RETIREMENT",
			OrdersTypeDetail:    "INSTRUCTION_20_WEEKS",
			DepartmentIndicator: "COAST_GUARD",
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

	body := &ghcmessages.UpdateMoveOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    "INSTRUCTION_20_WEEKS",
		DepartmentIndicator: "COAST_GUARD",
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

	body := &ghcmessages.UpdateMoveOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          "RETIREMENT",
		OrdersTypeDetail:    "INSTRUCTION_20_WEEKS",
		DepartmentIndicator: "COAST_GUARD",
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

	suite.Assertions.IsType(&moveorderop.UpdateMoveOrderBadRequest{}, response)
}
