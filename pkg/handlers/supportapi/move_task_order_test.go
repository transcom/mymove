package supportapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services/mocks"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	supportMocks "github.com/transcom/mymove/pkg/services/support/mocks"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"

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
	request := httptest.NewRequest("PATCH", "/move-task-orders/hide", nil)
	params := move_task_order.HideNonFakeMoveTaskOrdersParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	suite.T().Run("successfully hide fake moves", func(t *testing.T) {
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			movetaskorder.NewMoveTaskOrderHider(suite.DB()),
		}
		var moves models.Moves

		mto1 := testdatagen.MakeDefaultMove(suite.DB())
		mto2 := testdatagen.MakeDefaultMove(suite.DB())
		moves = append(moves, mto1, mto2)

		response := handler.Handle(params)
		mtoIDsRequestsResponse := response.(*movetaskorderops.HideNonFakeMoveTaskOrdersOK)
		mtoIDsRequestsPayload := mtoIDsRequestsResponse.Payload
		suite.IsNotErrResponse(response)
		suite.IsType(movetaskorderops.NewHideNonFakeMoveTaskOrdersOK(), response)

		for i, mtoID := range mtoIDsRequestsPayload {
			suite.Equal(supportmessages.MoveTaskOrderID(strfmt.UUID(moves[i].ID.String())), mtoID)
		}
	})

	suite.T().Run("unsuccessfully hide fake moves", func(t *testing.T) {
		var moves models.Moves
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		moves = append(moves, move)
		mockHider := &mocks.MoveTaskOrderHider{}
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			mockHider,
		}

		mockHider.On("Hide").Return(moves, errors.New("MTOs not retrieved"))

		response := handler.Handle(params)
		suite.IsType(movetaskorderops.NewHideNonFakeMoveTaskOrdersInternalServerError(), response)
	})

	suite.T().Run("Do not include mto in payload when it's missing a contractor id", func(t *testing.T) {
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		mto.ContractorID = nil
		moves = append(moves, mto)

		mockHider := &mocks.MoveTaskOrderHider{}
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			mockHider,
		}
		mockHider.On("Hide").Return(moves, nil)

		response := handler.Handle(params)
		moveTaskOrdersResponse := response.(*movetaskorderops.HideNonFakeMoveTaskOrdersOK)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		// Ensure that mto without a contractorID is NOT included in the payload
		for i, mtoID := range moveTaskOrdersPayload {
			suite.Equal(supportmessages.MoveTaskOrderID(strfmt.UUID(moves[i].ID.String())), mtoID)
		}
		suite.IsType(movetaskorderops.NewHideNonFakeMoveTaskOrdersOK(), response)
	})
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

// moveTaskOrderPopulated function spot checks a few values in the Move, Orders, and Customer to
// ensure they are populated.
func (suite *HandlerSuite) moveTaskOrderPopulated(response *movetaskorderops.CreateMoveTaskOrderCreated,
	destinationDutyStation *models.DutyStation,
	originDutyStation *models.DutyStation) {

	responsePayload := response.Payload

	suite.NotNil(responsePayload.MoveCode)
	suite.NotNil(responsePayload.AvailableToPrimeAt)
	suite.NotNil(responsePayload.MoveOrder.Customer.FirstName)

	suite.Equal(destinationDutyStation.Name, responsePayload.MoveOrder.DestinationDutyStation.Name)
	suite.Equal(originDutyStation.Name, responsePayload.MoveOrder.OriginDutyStation.Name)

}

func (suite *HandlerSuite) TestCreateMoveTaskOrderRequestHandler() {

	// Create the objects that are already in the db
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	dbCustomer := testdatagen.MakeDefaultServiceMember(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())
	document := testdatagen.MakeDefaultDocument(suite.DB())

	// Create the mto payload we will be requesting to create
	issueDate := swag.Time(time.Now())
	reportByDate := swag.Time(time.Now().AddDate(0, 0, -1))
	referenceID, _ := models.GenerateReferenceID(suite.DB())

	mtoPayload := &supportmessages.MoveTaskOrder{
		ReferenceID:        referenceID,
		AvailableToPrimeAt: handlers.FmtDateTime(time.Now()),
		PpmType:            "FULL",
		ContractorID:       handlers.FmtUUID(contractor.ID),
		IsCanceled:         swag.Bool(false),
		MoveOrder: &supportmessages.MoveOrder{
			Rank:                     (supportmessages.Rank)("E_6"),
			OrderNumber:              swag.String("4554"),
			DestinationDutyStationID: handlers.FmtUUID(destinationDutyStation.ID),
			OriginDutyStationID:      handlers.FmtUUID(originDutyStation.ID),
			Entitlement: &supportmessages.Entitlement{
				DependentsAuthorized: swag.Bool(true),
				TotalDependents:      5,
				NonTemporaryStorage:  swag.Bool(false),
			},
			IssueDate:        handlers.FmtDatePtr(issueDate),
			ReportByDate:     handlers.FmtDatePtr(reportByDate),
			OrdersType:       "PERMANENT_CHANGE_OF_STATION",
			UploadedOrdersID: handlers.FmtUUID(document.ID),
			Status:           (supportmessages.OrdersStatus)(models.OrderStatusDRAFT),
			Tac:              swag.String("47475"),
		},
	}

	// Create the handler object
	request := httptest.NewRequest("POST", "/move-task-orders", nil)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := CreateMoveTaskOrderHandler{context,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
	}

	suite.T().Run("Successful createMoveTaskOrder 201", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder and moveOrder associated with an existing customer,
		//             existing duty stations and existing uploaded orders document
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty station data are pulled in.
		//			   isCanceled should be default value which is draft

		// We only provide an existing customerID not the whole object.
		// We expect the handler to link the correct objects
		mtoPayload.MoveOrder.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID matches what was sent in
		suite.Equal(mtoPayload.ReferenceID, responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the DB
		suite.Equal(dbCustomer.FirstName, responsePayload.MoveOrder.Customer.FirstName)
		// Check that status has defaulted to DRAFT
		suite.Equal(models.MoveStatusDRAFT, (models.MoveStatus)(responsePayload.Status))
	})

	suite.T().Run("Successful createMoveTaskOrder 201 with canceled status", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder and moveOrder associated with an existing customer,
		//             existing duty stations and existing uploaded orders document.
		//             The status is canceled.
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty station data are pulled in.
		//             Status is canceled.

		// Regenerate the ReferenceID because it needs to be unique
		referenceID, _ := models.GenerateReferenceID(suite.DB())
		mtoPayload.ReferenceID = referenceID

		// We only provide an existing customerID not the whole object.
		mtoPayload.MoveOrder.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		// Set IsCanceled to true to set the Move's status to CANCELED
		mtoPayload.IsCanceled = swag.Bool(true)

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID matches what was sent in
		suite.Equal(mtoPayload.ReferenceID, responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the DB
		suite.Equal(dbCustomer.FirstName, responsePayload.MoveOrder.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)
	})

	suite.T().Run("Successful createMoveTaskOrder 201 with customer creation", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, moveOrder, and new customer.
		//             The move order is associated with existing duty stations.
		// Expected outcome:
		//             New MTO, orders and customer are created.

		// This time we provide customer details to create
		newCustomerFirstName := "Grace"
		mtoPayload.MoveOrder.Customer = &supportmessages.Customer{
			FirstName: &newCustomerFirstName,
			LastName:  swag.String("Griffin"),
			Agency:    swag.String("Marines"),
			DodID:     swag.String("1209457894"),
			Rank:      (supportmessages.Rank)("ACADEMY_CADET"),
		}
		mtoPayload.MoveOrder.CustomerID = nil

		// Regenerate the ReferenceID because it needs to be unique
		referenceID, _ := models.GenerateReferenceID(suite.DB())
		mtoPayload.ReferenceID = referenceID

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID matches what was sent in
		suite.Equal(mtoPayload.ReferenceID, responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the passed in value
		suite.Equal(newCustomerFirstName, *responsePayload.MoveOrder.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)

	})
	suite.T().Run("Failed createMoveTaskOrder 400 BadRequest due to repeat ReferenceID", func(t *testing.T) {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, moveOrder, and existing customer.
		//             We do not generate a new referenceID that has been used already
		// Expected outcome:
		//             Failure due to bad referenceID, so unprocessableEntity
		//             Default status is draft.

		// Running the same request should result in the same reference id
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderBadRequest{}, response)
	})

	suite.T().Run("Failed createMoveTaskOrder 422 UnprocessableEntity due to no customer", func(t *testing.T) {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, moveOrder, but no customer info
		// Expected outcome:
		//             Failure due to no customer info, so unprocessableEntity

		mtoPayload.MoveOrder.Customer = nil
		mtoPayload.MoveOrder.CustomerID = nil

		// Running the same request should result in the same reference id
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderUnprocessableEntity{}, response)
	})

	suite.T().Run("Failed createMoveTaskOrder 404 NotFound", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle
		// Mocked:     MoveTaskOrderCreator.CreateMoveTaskOrder
		// Set up:     We call the handler but force the mocked service object to return a notFoundError
		// Expected outcome:
		//             NotFound Response

		mockCreator := supportMocks.InternalMoveTaskOrderCreator{}
		handler.moveTaskOrderCreator = &mockCreator

		// Set expectation that a call to InternalCreateMoveTaskOrder will return a notFoundError
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

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderNotFound{}, response)
	})

	suite.T().Run("Failed createMoveTaskOrder 404 NotFound Bad Dutystation", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, moveOrder, and existing customer.
		//             The move order has a bad duty station ID.
		// Expected outcome:
		//             Failure of 404 Not Found since the dutystation is not found.

		// We only provide an existing customerID not the whole object.
		mtoPayload.MoveOrder.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		// Using a randomID as a dutyStationID should cause a query error
		mtoPayload.MoveOrder.OriginDutyStationID = handlers.FmtUUID(uuid.Must(uuid.NewV4()))

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		handler := CreateMoveTaskOrderHandler{context,
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderNotFound{}, response)
	})
}
