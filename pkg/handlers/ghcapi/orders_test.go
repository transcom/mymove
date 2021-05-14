package ghcapi

import (
	"net/http/httptest"
	"testing"
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
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
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

func (suite *HandlerSuite) TestUpdateOrderHandlerIntegration() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")

	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	body := &ghcmessages.UpdateOrderPayload{
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
		orderservice.NewOrderFetcher(suite.DB()),
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
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateOrderEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

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
		Tac:                 handlers.FmtString("E19A"),
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
		orderservice.NewOrderFetcher(suite.DB()),
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
	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

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
			Tac:                 handlers.FmtString("E19A"),
			Sac:                 handlers.FmtString("987654321"),
		},
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
		orderservice.NewOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderNotFound{}, response)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerPreconditionsFailed() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

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
		Tac:                 handlers.FmtString("E19A"),
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
		orderservice.NewOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerValidationError() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	err := move.Submit()
	if err != nil {
		suite.T().Fatal("Should transition.")
	}
	suite.MustSave(&move)

	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

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
		Sac:                 handlers.FmtString("987654321"),
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
		orderservice.NewOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
	invalidResponse := response.(*orderop.UpdateOrderUnprocessableEntity).Payload
	errorDetail := invalidResponse.Detail

	updatedOrder, _ := models.FetchOrder(suite.DB(), order.ID)

	suite.Contains(*errorDetail, "NewDutyStationID can not be blank.")
	suite.NotNil(updatedOrder.TAC)
}

func (suite *HandlerSuite) TestUpdateOrderHandlerWithoutTac() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	order.TAC = nil
	suite.MustSave(&order)

	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	newDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
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
		NewDutyStationID:    handlers.FmtUUID(newDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Sac:                 handlers.FmtString("987654321"),
	}

	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	suite.Run("When Move is still in draft status, TAC can be nil", func() {
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
			orderservice.NewOrderFetcher(suite.DB()),
		}
		response := handler.Handle(params)

		suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)
		payload := response.(*orderop.UpdateOrderOK).Payload

		updatedOrder, _ := models.FetchOrder(suite.DB(), order.ID)

		suite.EqualValues(body.OrdersNumber, updatedOrder.OrdersNumber)
		suite.Nil(updatedOrder.TAC)
		suite.Equal(move.Locator, payload.MoveCode)
	})

	suite.Run("When Move is no longer in draft status, TAC must be present", func() {
		// Submit the move to change its status
		err := move.Submit()
		if err != nil {
			suite.T().Fatal("Should transition.")
		}
		suite.MustSave(&move)
		updatedMove, _ := models.FetchMoveByMoveID(suite.DB(), move.ID)
		updatedOrder, _ := models.FetchOrder(suite.DB(), order.ID)

		suite.EqualValues(models.MoveStatusSUBMITTED, updatedMove.Status)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(updatedOrder.UpdatedAt),
			Body:        body,
		}

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := UpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(suite.DB()),
			orderservice.NewOrderFetcher(suite.DB()),
		}
		response := handler.Handle(params)

		suite.Assertions.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
		invalidResponse := response.(*orderop.UpdateOrderUnprocessableEntity).Payload
		errorDetail := invalidResponse.Detail

		suite.Contains(*errorDetail, "TransportationAccountingCode cannot be blank.")
	})

	suite.Run("TAC can only contain 4 alphanumeric characters", func() {
		existingOrder, _ := models.FetchOrder(suite.DB(), order.ID)

		invalidCases := []struct {
			desc string
			tac  string
		}{
			{"TestOneCharacter", "A"},
			{"TestTwoCharacters", "AB"},
			{"TestThreeCharacters", "ABC"},
			{"TestGreaterThanFourChars", "ABCD1"},
			{"TestNonAlphaNumChars", "AB-C"},
		}
		for _, invalidCase := range invalidCases {
			body.Tac = &invalidCase.tac
			params := orderop.UpdateOrderParams{
				HTTPRequest: request,
				OrderID:     strfmt.UUID(order.ID.String()),
				IfMatch:     etag.GenerateEtag(existingOrder.UpdatedAt),
				Body:        body,
			}
			context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
			handler := UpdateOrderHandler{
				context,
				orderservice.NewOrderUpdater(suite.DB()),
				orderservice.NewOrderFetcher(suite.DB()),
			}
			response := handler.Handle(params)

			suite.Assertions.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
			invalidResponse := response.(*orderop.UpdateOrderUnprocessableEntity).Payload
			errorDetail := invalidResponse.Detail

			suite.Contains(*errorDetail, "TAC must be exactly 4 alphanumeric characters.")
		}
	})
}

func (suite *HandlerSuite) TestUpdateAllowanceHandlerIntegration() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders

	newAuthorizedWeight := int64(10000)
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.BranchAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	requestUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)
	handler := UpdateAllowanceHandler{
		context,
		orderservice.NewOrderUpdater(suite.DB()),
		orderservice.NewOrderFetcher(suite.DB()),
	}

	suite.T().Run("successfully updates order allowance", func(t *testing.T) {
		body := &ghcmessages.UpdateAllowancePayload{
			Agency:               affiliation,
			AuthorizedWeight:     &newAuthorizedWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		allowanceOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := allowanceOK.Payload

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

	suite.T().Run("successfully updates order allowance without ocie, rme, pro-gear, and pro-gear spouse fields", func(t *testing.T) {
		newMove := testdatagen.MakeDefaultMove(suite.DB())
		newOrder := newMove.Orders

		body := &ghcmessages.UpdateAllowancePayload{
			Agency:               affiliation,
			AuthorizedWeight:     &newAuthorizedWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
		}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(newOrder.ID.String()),
			IfMatch:     etag.GenerateEtag(newOrder.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		allowanceOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := allowanceOK.Payload

		suite.Assertions.IsType(&orderop.UpdateAllowanceOK{}, response)
		suite.Equal(newOrder.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)

		// should be defaults
		suite.EqualValues(newOrder.Entitlement.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.EqualValues(newOrder.Entitlement.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.EqualValues(newOrder.Entitlement.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.EqualValues(newOrder.Entitlement.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
	})

	suite.T().Run("successfully updates order allowance without updating authorized weight field as Service Counselor role", func(t *testing.T) {
		newMove := testdatagen.MakeDefaultMove(suite.DB())
		newOrder := newMove.Orders
		authWeight := swag.Int64(1234)

		scRoleUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})
		req := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
		req = suite.AuthenticateOfficeRequest(req, scRoleUser)

		body := &ghcmessages.UpdateAllowancePayload{
			Agency:               affiliation,
			AuthorizedWeight:     authWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
		}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: req,
			OrderID:     strfmt.UUID(newOrder.ID.String()),
			IfMatch:     etag.GenerateEtag(newOrder.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		allowanceOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := allowanceOK.Payload

		suite.Assertions.IsType(&orderop.UpdateAllowanceOK{}, response)
		suite.Equal(newOrder.ID.String(), ordersPayload.ID.String())
		suite.NotEqual(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
		// equals the original value
		suite.EqualValues(*newOrder.Entitlement.DBAuthorizedWeight, *ordersPayload.Entitlement.AuthorizedWeight)
	})

	suite.T().Run("returns 400 bad request", func(t *testing.T) {
		body := &ghcmessages.UpdateAllowancePayload{}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(""),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&orderop.UpdateAllowanceBadRequest{}, response)
	})

	suite.T().Run("returns 403 forbidden", func(t *testing.T) {
		stubbedUser := testdatagen.MakeStubbedUser(suite.DB())
		newRequest := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
		newRequest = suite.AuthenticateUserRequest(newRequest, stubbedUser)
		body := &ghcmessages.UpdateAllowancePayload{}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: newRequest,
			OrderID:     strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&orderop.UpdateAllowanceForbidden{}, response)
	})

	suite.T().Run("returns 404 not found", func(t *testing.T) {
		body := &ghcmessages.UpdateAllowancePayload{}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&orderop.UpdateAllowanceNotFound{}, response)
	})

	suite.T().Run("returns 412 pre-condition failed", func(t *testing.T) {
		body := &ghcmessages.UpdateAllowancePayload{
			Agency:               affiliation,
			AuthorizedWeight:     &newAuthorizedWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
		}

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&orderop.UpdateAllowancePreconditionFailed{}, response)
	})
}
