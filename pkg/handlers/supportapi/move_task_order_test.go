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

		for i, m := range mtoIDsRequestsPayload.Moves {
			suite.Equal(moves[i].ID.String(), m.MoveTaskOrderID.String())
			suite.NotEqual("{}", *m.HideReason)
		}
	})

	suite.T().Run("unsuccessfully hide fake moves", func(t *testing.T) {
		var moves models.Moves
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		moves = append(moves, move)
		var hiddenMoves services.HiddenMoves
		for _, m := range moves {
			hm := services.HiddenMove{
				MTOID:  m.ID,
				Reason: "move is hidden",
			}
			hiddenMoves = append(hiddenMoves, hm)
		}
		mockHider := &mocks.MoveTaskOrderHider{}
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			mockHider,
		}

		mockHider.On("Hide").Return(hiddenMoves, errors.New("MTOs not retrieved"))

		response := handler.Handle(params)
		suite.IsType(movetaskorderops.NewHideNonFakeMoveTaskOrdersInternalServerError(), response)
	})

	suite.T().Run("Do not include mto in payload when it's missing a contractor id", func(t *testing.T) {
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		mto.ContractorID = nil
		moves = append(moves, mto)
		var hiddenMoves services.HiddenMoves
		for _, m := range moves {
			hm := services.HiddenMove{
				MTOID:  m.ID,
				Reason: "move is hidden",
			}
			hiddenMoves = append(hiddenMoves, hm)
		}

		mockHider := &mocks.MoveTaskOrderHider{}
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			mockHider,
		}
		mockHider.On("Hide").Return(hiddenMoves, nil)

		response := handler.Handle(params)
		moveTaskOrdersResponse := response.(*movetaskorderops.HideNonFakeMoveTaskOrdersOK)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		// Ensure that mto without a contractorID is NOT included in the payload
		for i, mto := range moveTaskOrdersPayload.Moves {
			suite.Equal(moves[i].ID.String(), mto.MoveTaskOrderID.String())
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
	suite.NotNil(responsePayload.Order.Customer.FirstName)

	suite.Equal(destinationDutyStation.Name, responsePayload.Order.DestinationDutyStation.Name)
	suite.Equal(originDutyStation.Name, responsePayload.Order.OriginDutyStation.Name)

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
	ordersTypedetail := supportmessages.OrdersTypeDetailHHGPERMITTED
	deptIndicator := supportmessages.DeptIndicatorAIRFORCE
	selectedMoveType := supportmessages.SelectedMoveTypeHHG

	mtoPayload := &supportmessages.MoveTaskOrder{
		PpmType:          "FULL",
		SelectedMoveType: &selectedMoveType,
		ContractorID:     handlers.FmtUUID(contractor.ID),
		Order: &supportmessages.Order{
			Rank:                     (supportmessages.Rank)("E_6"),
			OrderNumber:              swag.String("4554"),
			DestinationDutyStationID: handlers.FmtUUID(destinationDutyStation.ID),
			OriginDutyStationID:      handlers.FmtUUID(originDutyStation.ID),
			Entitlement: &supportmessages.Entitlement{
				DependentsAuthorized: swag.Bool(true),
				TotalDependents:      5,
				NonTemporaryStorage:  swag.Bool(false),
			},
			IssueDate:           handlers.FmtDatePtr(issueDate),
			ReportByDate:        handlers.FmtDatePtr(reportByDate),
			OrdersType:          "PERMANENT_CHANGE_OF_STATION",
			OrdersTypeDetail:    &ordersTypedetail,
			UploadedOrdersID:    handlers.FmtUUID(document.ID),
			Status:              (supportmessages.OrdersStatus)(models.OrderStatusDRAFT),
			Tac:                 swag.String("E19A"),
			DepartmentIndicator: &deptIndicator,
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
		// Set up:     We pass in a new moveTaskOrder and order associated with an existing customer,
		//             existing duty stations and existing uploaded orders document
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty station data are pulled in.
		//			   Status should be default value which is DRAFT

		// We only provide an existing customerID not the whole object.
		// We expect the handler to link the correct objects
		mtoPayload.Order.CustomerID = handlers.FmtUUID(dbCustomer.ID)

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

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the DB
		suite.Equal(dbCustomer.FirstName, responsePayload.Order.Customer.FirstName)
		// Check that status has defaulted to DRAFT
		suite.Equal(models.MoveStatusDRAFT, (models.MoveStatus)(responsePayload.Status))
		// Check that SelectedMoveType was set
		suite.Equal(string(models.SelectedMoveTypeHHG), string(*responsePayload.SelectedMoveType))
	})

	suite.T().Run("Successful integration test with createMoveTaskOrder", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We successfully create a new MTO, and then test that this move can be successfully approved.
		// Expected outcome:
		//             New MTO and orders are created. MTO can be approved and marked as available to Prime.

		// Let's copy the default mtoPayload so we don't affect the other tests:
		var integrationMTO supportmessages.MoveTaskOrder
		integrationMTO = *mtoPayload

		// We have to set the status for the orders to APPROVED and the move to SUBMITTED so that we can try to approve
		// this move later on. We can't approve a DRAFT move.
		integrationMTO.Status = supportmessages.MoveStatusSUBMITTED
		integrationMTO.Order.Status = supportmessages.OrdersStatusAPPROVED

		// We only provide an existing customerID not the whole object.
		// We expect the handler to link the correct objects
		integrationMTO.Order.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        &integrationMTO,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		createdMTO := moveTaskOrdersResponse.Payload

		// Check that status has been set to SUBMITTED
		suite.Equal(models.MoveStatusSUBMITTED, (models.MoveStatus)(createdMTO.Status))

		// Now we'll try to approve this MTO and verify that it was successfully made available to the Prime
		approvalRequest := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/available-to-prime", nil)
		approvalParams := move_task_order.MakeMoveTaskOrderAvailableParams{
			HTTPRequest:     approvalRequest,
			MoveTaskOrderID: createdMTO.ID.String(),
			IfMatch:         createdMTO.ETag,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

		// Submit the request to approve the MTO
		approvalHandler := MakeMoveTaskOrderAvailableHandlerFunc{context,
			movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
		}
		approvalResponse := approvalHandler.Handle(approvalParams)

		// VERIFY RESULTS
		suite.IsNotErrResponse(approvalResponse)
		suite.Assertions.IsType(&move_task_order.MakeMoveTaskOrderAvailableOK{}, approvalResponse)
		approvalOKResponse := approvalResponse.(*movetaskorderops.MakeMoveTaskOrderAvailableOK)
		approvedMTO := approvalOKResponse.Payload

		suite.Equal(approvedMTO.ID, strfmt.UUID(createdMTO.ID.String()))
		suite.NotNil(approvedMTO.AvailableToPrimeAt)
		suite.Equal(string(approvedMTO.Status), string(models.MoveStatusAPPROVED))
	})

	suite.T().Run("Successful createMoveTaskOrder 201 with canceled status", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder and order associated with an existing customer,
		//             existing duty stations and existing uploaded orders document.
		//             The status is canceled.
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty station data are pulled in.
		//             Status is canceled.

		// We only provide an existing customerID not the whole object.
		mtoPayload.Order.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		// Set the status to CANCELED
		mtoPayload.Status = supportmessages.MoveStatusCANCELED

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

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the DB
		suite.Equal(dbCustomer.FirstName, responsePayload.Order.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)
	})

	suite.T().Run("Successful createMoveTaskOrder 201 with customer creation", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, and new customer.
		//             The order is associated with existing duty stations.
		// Expected outcome:
		//             New MTO, orders and customer are created.

		// This time we provide customer details to create
		newCustomerFirstName := "Grace"
		mtoPayload.Order.Customer = &supportmessages.Customer{
			FirstName: &newCustomerFirstName,
			LastName:  swag.String("Griffin"),
			Agency:    swag.String("Marines"),
			DodID:     swag.String("1209457894"),
			Rank:      (supportmessages.Rank)("ACADEMY_CADET"),
		}
		mtoPayload.Order.CustomerID = nil

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

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
		// Check that customer name matches the passed in value
		suite.Equal(newCustomerFirstName, *responsePayload.Order.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)

	})
	suite.T().Run("Success createMoveTaskOrder discarded readOnly referenceID", func(t *testing.T) {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, and existing customer.
		//             We use a nonsense referenceID. ReferenceID is readOnly, this should not be applied.
		// Expected outcome:
		//             A new referenceID is generated.
		//             Default status is draft.

		// Running the same request should result in the same reference id
		mtoPayload.ReferenceID = "some terrible reference id"
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

		// Check that the referenceID DOES NOT match what was sent in
		suite.NotEqual(mtoPayload.ReferenceID, responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyStation, &originDutyStation)
	})

	suite.T().Run("Failed createMoveTaskOrder 422 UnprocessableEntity due to no customer", func(t *testing.T) {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, but no customer info
		// Expected outcome:
		//             Failure due to no customer info, so unprocessableEntity

		mtoPayload.Order.Customer = nil
		mtoPayload.Order.CustomerID = nil

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
		// Set up:     We pass in a new moveTaskOrder, order, and existing customer.
		//             The order has a bad duty station ID.
		// Expected outcome:
		//             Failure of 404 Not Found since the dutystation is not found.

		// We only provide an existing customerID not the whole object.
		mtoPayload.Order.CustomerID = handlers.FmtUUID(dbCustomer.ID)

		// Using a randomID as a dutyStationID should cause a query error
		mtoPayload.Order.OriginDutyStationID = handlers.FmtUUID(uuid.Must(uuid.NewV4()))

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
