package supportapi

import (
	"github.com/transcom/mymove/pkg/services/support"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationSuccess() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
		IfMatch:         etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.Equal(*moveTaskOrdersPayload.IsAvailableToPrime, true)
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
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
	suite.Equal(*moveTaskOrdersPayload.IsAvailableToPrime, false)

}

func (suite *HandlerSuite) TestCreateMoveTaskOrderRequestHandler() {

	destinationDutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{})
	originDutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{})
	dbCustomer := testdatagen.MakeCustomer(suite.DB(), testdatagen.Assertions{})
	contractor := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{})

	mtoWithoutCustomer := models.MoveTaskOrder{
		ReferenceID:        "4857363",
		IsAvailableToPrime: true,
		PPMType:            swag.String("FULL"),
		ContractorID:       contractor.ID,
		MoveOrder: models.MoveOrder{
			Grade:                    swag.String("E_6"),
			OrderNumber:              swag.String("4554"),
			DestinationDutyStationID: &destinationDutyStation.ID,
			OriginDutyStationID:      &originDutyStation.ID,
			Entitlement: &models.Entitlement{
				DependentsAuthorized: swag.Bool(true),
				TotalDependents:      swag.Int(5),
				NonTemporaryStorage:  swag.Bool(false),
			},
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
			support.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)

		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload
		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal(mtoWithoutCustomer.ReferenceID, moveTaskOrdersPayload.ReferenceID)
		suite.Equal(true, *moveTaskOrdersPayload.IsAvailableToPrime)
	})

	suite.T().Run("successful create movetaskorder request -- with customer creation", func(t *testing.T) {

		newCustomer := models.Customer{
			FirstName: swag.String("Grace"),
			LastName:  swag.String("Griffin"),
		}
		mtoWithoutCustomer.ReferenceID = "346523"

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
			support.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		suite.Assertions.IsType(&move_task_order.CreateMoveTaskOrderCreated{}, response)
		suite.Equal(mtoWithoutCustomer.ReferenceID, moveTaskOrdersPayload.ReferenceID)
		suite.Equal(true, *moveTaskOrdersPayload.IsAvailableToPrime)
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
			support.NewInternalMoveTaskOrderCreator(context.DB()),
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
			support.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create movetaskorder request 404 -- not found", func(t *testing.T) {
		mtoPayload := payloads.MoveTaskOrder(&mtoWithoutCustomer)
		mtoPayload.MoveOrder.CustomerID = strfmt.UUID(dbCustomer.ID.String())

		mockFetcher := mocks.InternalMoveTaskOrderCreator{}
		handler := CreateMoveTaskOrderHandler{context,
			&mockFetcher,
		}

		notFoundError := services.NotFoundError{}

		mockFetcher.On("FetchCustomer",
			mock.Anything,
			mock.Anything,
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
			support.NewInternalMoveTaskOrderCreator(context.DB()),
		}
		response := handler.Handle(params)

		suite.IsType(&movetaskorderops.CreateMoveTaskOrderBadRequest{}, response)
	})
}
