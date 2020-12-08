package supportapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services/mocks"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	supportMocks "github.com/transcom/mymove/pkg/services/support/mocks"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMTOsHandler() {
	// unavailable MTO
	testdatagen.MakeDefaultMove(suite.DB())

	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
	})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.ListMTOsParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	handler := ListMTOsHandler{
		HandlerContext:       context,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	listMTOsResponse := response.(*movetaskorderops.ListMTOsOK)
	listMTOsPayload := listMTOsResponse.Payload

	suite.Equal(2, len(listMTOsPayload))

}

func (suite *HandlerSuite) TestHideNonFakeMoveTaskOrdersHandler() {
	request := httptest.NewRequest("GET", "/move-task-orders/hide", nil)
	params := move_task_order.HideNonFakeMoveTaskOrdersParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	moves := []models.Move{
		testdatagen.MakeAvailableMove(suite.DB()),
		testdatagen.MakeAvailableMove(suite.DB()),
	}
	mockHider := &mocks.MoveTaskOrderHider{}
	handler := HideNonFakeMoveTaskOrdersHandlerFunc{
		context,
		mockHider,
	}
	mockHider.On("Hide").Return(moves, nil)

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(move_task_order.NewHideNonFakeMoveTaskOrdersOK, response)
}

func (suite *HandlerSuite) TestMakeMoveTaskOrderAvailableHandlerIntegrationSuccess() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/available-to-prime", nil)
	params := move_task_order.MakeMoveTaskOrderAvailableParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
		IfMatch:         etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

	// make the request
	handler := MakeMoveTaskOrderAvailableHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.MakeMoveTaskOrderAvailableOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.MakeMoveTaskOrderAvailableOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.NotNil(moveTaskOrdersPayload.AvailableToPrimeAt)
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := move_task_order.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveTaskOrderHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}
	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.GetMoveTaskOrderOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.Nil(moveTaskOrdersPayload.AvailableToPrimeAt)
}

func (suite *HandlerSuite) TestCreateMoveTaskOrderRequestHandler() {

	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	dbCustomer := testdatagen.MakeDefaultServiceMember(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())
	document := testdatagen.MakeDefaultDocument(suite.DB())
	issueDate := swag.Time(time.Now())
	reportByDate := swag.Time(time.Now().AddDate(0, 0, -1))
	referenceID, _ := models.GenerateReferenceID(suite.DB())

	mtoWithoutCustomer := models.Move{
		ReferenceID:        &referenceID,
		AvailableToPrimeAt: swag.Time(time.Now()),
		PPMType:            swag.String("FULL"),
		ContractorID:       &contractor.ID,
		Orders: models.Order{
			Grade:               swag.String("E_6"),
			OrdersNumber:        swag.String("4554"),
			NewDutyStationID:    destinationDutyStation.ID,
			OriginDutyStationID: &originDutyStation.ID,
			Entitlement: &models.Entitlement{
				DependentsAuthorized: swag.Bool(true),
				TotalDependents:      swag.Int(5),
				NonTemporaryStorage:  swag.Bool(false),
			},
			Status:           models.OrderStatusDRAFT,
			IssueDate:        *issueDate,
			ReportByDate:     *reportByDate,
			OrdersType:       "PERMANENT_CHANGE_OF_STATION",
			UploadedOrders:   document,
			UploadedOrdersID: document.ID,
		},
	}

	request := httptest.NewRequest("POST", "/move-task-orders", nil)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	suite.T().Run("successful create movetaskorder request 201", func(t *testing.T) {

		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(originDutyStation.ID.String())
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)

		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload
		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal(*mtoWithoutCustomer.ReferenceID, moveTaskOrdersPayload.ReferenceID)
		suite.NotNil(moveTaskOrdersPayload.Locator)
		suite.NotNil(moveTaskOrdersPayload.AvailableToPrimeAt)
		suite.Equal((models.MoveStatus)(moveTaskOrdersPayload.Status), models.MoveStatusDRAFT)
	})

	suite.T().Run("successful cancel movetaskorder request 201", func(t *testing.T) {
		// Regenerate the ReferenceID because it needs to be unique
		referenceID, _ := models.GenerateReferenceID(suite.DB())
		mtoWithoutCustomer.ReferenceID = &referenceID

		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(originDutyStation.ID.String())
		// Set IsCanceled to true to set the Move's status to CANCELED
		mtoPayload.IsCanceled = swag.Bool(true)
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)

		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload
		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal((models.MoveStatus)(moveTaskOrdersPayload.Status), models.MoveStatusCANCELED)
	})

	suite.T().Run("move status stays as DRAFT when IsCanceled is false", func(t *testing.T) {
		// Regenerate the ReferenceID because it needs to be unique
		referenceID, _ := models.GenerateReferenceID(suite.DB())
		mtoWithoutCustomer.ReferenceID = &referenceID

		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(originDutyStation.ID.String())
		// Set IsCanceled to true to set the Move's status to CANCELED
		mtoPayload.IsCanceled = swag.Bool(false)
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)

		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload
		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal((models.MoveStatus)(moveTaskOrdersPayload.Status), models.MoveStatusDRAFT)
	})

	suite.T().Run("successful create movetaskorder request -- with customer creation", func(t *testing.T) {

		newCustomer := models.ServiceMember{
			FirstName: swag.String("Grace"),
			LastName:  swag.String("Griffin"),
		}
		// Regenerate the ReferenceID because it needs to be unique
		referenceID, _ := models.GenerateReferenceID(suite.DB())
		mtoWithoutCustomer.ReferenceID = &referenceID

		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		customerPayload := payloads.Customer(&newCustomer)
		mtoPayload.MoveOrder.Customer = customerPayload
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(originDutyStation.ID.String())

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal(*mtoWithoutCustomer.ReferenceID, moveTaskOrdersPayload.ReferenceID)
		suite.NotNil(moveTaskOrdersPayload.Locator)
		suite.NotNil(moveTaskOrdersPayload.AvailableToPrimeAt)
	})
	suite.T().Run("failed create movetaskorder request 400 -- repeat ReferenceID", func(t *testing.T) {

		// Running the same request should result in the same reference id
		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(originDutyStation.ID.String())

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderBadRequest{}, response)
	})
	suite.T().Run("failed create movetaskorder request 422 -- unprocessable entity", func(t *testing.T) {

		// Running the same request should result in the same reference id
		// If customerID is provided create MTO without creating a new customer
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create movetaskorder request 404 -- not found", func(t *testing.T) {
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())

		mockCreator := supportMocks.InternalMoveTaskOrderCreator{}
		handler := CreateMoveTaskOrderHandler{context,
			&mockCreator,
		}

		notFoundError := services.NotFoundError{}

		mockCreator.On("InternalCreateMoveTaskOrder",
			mock.Anything,
			mock.Anything,
		).Return(nil, notFoundError)

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderNotFound{}, response)
	})

	suite.T().Run("failed create movetaskorder request 400 -- Invalid Request", func(t *testing.T) {
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())
		mtoPayload.MoveOrder.DestinationDutyStationID = strfmt.UUID(destinationDutyStation.ID.String())
		// using a customerID as a dutyStationID should cause a query error
		mtoPayload.MoveOrder.OriginDutyStationID = strfmt.UUID(dbCustomer.ID.String())
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// make the request
		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderBadRequest{}, response)
	})
}
