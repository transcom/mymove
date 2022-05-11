package supportapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	supportMocks "github.com/transcom/mymove/pkg/services/support/mocks"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"
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
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())

	handler := ListMTOsHandler{
		HandlerContext:       context,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(),
	}

	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	listMTOsResponse := response.(*movetaskorderops.ListMTOsOK)
	listMTOsPayload := listMTOsResponse.Payload

	suite.Equal(2, len(listMTOsPayload))
}

func (suite *HandlerSuite) TestHideNonFakeMoveTaskOrdersHandler() {
	request := httptest.NewRequest("PATCH", "/move-task-orders/hide", nil)
	params := movetaskorderops.HideNonFakeMoveTaskOrdersParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())

	suite.Run("successfully hide fake moves", func() {
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			context,
			movetaskorder.NewMoveTaskOrderHider(),
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

	suite.Run("unsuccessfully hide fake moves", func() {
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

		mockHider.On("Hide",
			mock.AnythingOfType("*appcontext.appContext"),
		).Return(hiddenMoves, errors.New("MTOs not retrieved"))

		response := handler.Handle(params)
		suite.IsType(movetaskorderops.NewHideNonFakeMoveTaskOrdersInternalServerError(), response)
	})

	suite.Run("Do not include mto in payload when it's missing a contractor id", func() {
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
		mockHider.On("Hide",
			mock.AnythingOfType("*appcontext.appContext"),
		).Return(hiddenMoves, nil)

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

func (suite *HandlerSuite) TestMakeMoveAvailableHandlerIntegrationSuccess() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
	})
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/available-to-prime", nil)
	params := movetaskorderops.MakeMoveTaskOrderAvailableParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
		IfMatch:         etag.GenerateEtag(move.UpdatedAt),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)

	// make the request
	handler := MakeMoveTaskOrderAvailableHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveResponse := response.(*movetaskorderops.MakeMoveTaskOrderAvailableOK)
	movePayload := moveResponse.Payload

	suite.Assertions.IsType(&movetaskorderops.MakeMoveTaskOrderAvailableOK{}, response)
	suite.Equal(movePayload.ID, strfmt.UUID(move.ID.String()))
	suite.NotNil(movePayload.AvailableToPrimeAt)
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := movetaskorderops.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMoveTaskOrderHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}
	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

	moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	movePayload := moveResponse.Payload
	suite.Equal(movePayload.ID, strfmt.UUID(move.ID.String()))
	suite.Nil(movePayload.AvailableToPrimeAt)
}

// moveTaskOrderPopulated function spot checks a few values in the Move, Orders, and Customer to
// ensure they are populated.
func (suite *HandlerSuite) moveTaskOrderPopulated(response *movetaskorderops.CreateMoveTaskOrderCreated,
	destinationDutyLocation *models.DutyLocation,
	originDutyLocation *models.DutyLocation) {

	responsePayload := response.Payload

	suite.NotNil(responsePayload.MoveCode)
	suite.NotNil(responsePayload.Order.Customer.FirstName)

	suite.Equal(destinationDutyLocation.Name, responsePayload.Order.DestinationDutyLocation.Name)
	suite.Equal(originDutyLocation.Name, responsePayload.Order.OriginDutyLocation.Name)

}

func (suite *HandlerSuite) TestCreateMoveTaskOrderRequestHandler() {
	setupTestData := func() (models.DutyLocation, models.DutyLocation, models.ServiceMember, *supportmessages.MoveTaskOrder) {
		// Create the objects that are already in the db
		destinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		originDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		customer := testdatagen.MakeDefaultServiceMember(suite.DB())
		contractor := testdatagen.MakeDefaultContractor(suite.DB())
		document := testdatagen.MakeDefaultDocument(suite.DB())

		// Create the mto payload we will be requesting to create
		issueDate := swag.Time(time.Now())
		reportByDate := swag.Time(time.Now().AddDate(0, 0, -1))
		ordersTypedetail := supportmessages.OrdersTypeDetailHHGPERMITTED
		deptIndicator := supportmessages.DeptIndicatorAIRFORCE
		selectedMoveType := supportmessages.SelectedMoveTypeHHG

		rank := (supportmessages.Rank)("E_6")
		mtoPayload := &supportmessages.MoveTaskOrder{
			PpmType:          "FULL",
			SelectedMoveType: &selectedMoveType,
			ContractorID:     handlers.FmtUUID(contractor.ID),
			Order: &supportmessages.Order{
				Rank:                      &rank,
				OrderNumber:               swag.String("4554"),
				DestinationDutyLocationID: handlers.FmtUUID(destinationDutyLocation.ID),
				OriginDutyLocationID:      handlers.FmtUUID(originDutyLocation.ID),
				Entitlement: &supportmessages.Entitlement{
					DependentsAuthorized: swag.Bool(true),
					TotalDependents:      5,
					NonTemporaryStorage:  swag.Bool(false),
				},
				IssueDate:           handlers.FmtDatePtr(issueDate),
				ReportByDate:        handlers.FmtDatePtr(reportByDate),
				OrdersType:          supportmessages.NewOrdersType("PERMANENT_CHANGE_OF_STATION"),
				OrdersTypeDetail:    &ordersTypedetail,
				UploadedOrdersID:    handlers.FmtUUID(document.ID),
				Status:              supportmessages.NewOrdersStatus(supportmessages.OrdersStatusDRAFT),
				Tac:                 swag.String("E19A"),
				DepartmentIndicator: &deptIndicator,
			},
		}
		// We only provide an existing customerID not the whole object.
		// We expect the handler to link the correct objects
		mtoPayload.Order.CustomerID = handlers.FmtUUID(customer.ID)
		return destinationDutyLocation, originDutyLocation, customer, mtoPayload
	}

	setupHandler := func() CreateMoveTaskOrderHandler {
		return CreateMoveTaskOrderHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(),
		}
	}

	request := httptest.NewRequest("POST", "/move-task-orders", nil)

	suite.Run("Successful createMoveTaskOrder 201", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder and order associated with an existing customer,
		//             existing duty locations and existing uploaded orders document
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty location data are pulled in.
		//			   Status should be default value which is DRAFT
		destinationDutyLocation, originDutyLocation, customer, mtoPayload := setupTestData()

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyLocation, &originDutyLocation)
		// Check that customer name matches the DB
		suite.Equal(customer.FirstName, responsePayload.Order.Customer.FirstName)
		// Check that status has defaulted to DRAFT
		suite.Equal(models.MoveStatusDRAFT, (models.MoveStatus)(responsePayload.Status))
		// Check that SelectedMoveType was set
		suite.Equal(string(models.SelectedMoveTypeHHG), string(*responsePayload.SelectedMoveType))
	})

	suite.Run("Successful integration test with createMoveTaskOrder", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We successfully create a new MTO, and then test that this move can be successfully approved.
		// Expected outcome:
		//             New MTO and orders are created. MTO can be approved and marked as available to Prime.
		_, _, _, integrationMTO := setupTestData()

		// We have to set the status for the orders to APPROVED and the move to SUBMITTED so that we can try to approve
		// this move later on. We can't approve a DRAFT move.
		integrationMTO.Status = supportmessages.MoveStatusSUBMITTED
		integrationMTO.Order.Status = supportmessages.NewOrdersStatus(supportmessages.OrdersStatusAPPROVED)

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        integrationMTO,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		createdMTO := moveTaskOrdersResponse.Payload

		// Check that status has been set to SUBMITTED
		suite.Equal(models.MoveStatusSUBMITTED, (models.MoveStatus)(createdMTO.Status))

		// Now we'll try to approve this MTO and verify that it was successfully made available to the Prime
		approvalRequest := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/available-to-prime", nil)
		approvalParams := movetaskorderops.MakeMoveTaskOrderAvailableParams{
			HTTPRequest:     approvalRequest,
			MoveTaskOrderID: createdMTO.ID.String(),
			IfMatch:         createdMTO.ETag,
		}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)

		// Submit the request to approve the MTO
		approvalHandler := MakeMoveTaskOrderAvailableHandlerFunc{handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter),
		}
		approvalResponse := approvalHandler.Handle(approvalParams)

		// VERIFY RESULTS
		suite.IsNotErrResponse(approvalResponse)
		suite.Assertions.IsType(&movetaskorderops.MakeMoveTaskOrderAvailableOK{}, approvalResponse)
		approvalOKResponse := approvalResponse.(*movetaskorderops.MakeMoveTaskOrderAvailableOK)
		approvedMTO := approvalOKResponse.Payload

		suite.Equal(approvedMTO.ID, strfmt.UUID(createdMTO.ID.String()))
		suite.NotNil(approvedMTO.AvailableToPrimeAt)
		suite.Equal(string(approvedMTO.Status), string(models.MoveStatusAPPROVED))
	})

	suite.Run("Successful createMoveTaskOrder 201 with canceled status", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder and order associated with an existing customer,
		//             existing duty locations and existing uploaded orders document.
		//             The status is canceled.
		// Expected outcome:
		//             New MTO and orders are created. Customer data and duty location data are pulled in.
		//             Status is canceled.
		destinationDutyLocation, originDutyLocation, customer, mtoPayload := setupTestData()

		// Set the status to CANCELED
		mtoPayload.Status = supportmessages.MoveStatusCANCELED

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyLocation, &originDutyLocation)
		// Check that customer name matches the DB
		suite.Equal(customer.FirstName, responsePayload.Order.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)
	})

	suite.Run("Successful createMoveTaskOrder 201 with customer creation", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, and new customer.
		//             The order is associated with existing duty locations.
		// Expected outcome:
		//             New MTO, orders and customer are created.
		destinationDutyLocation, originDutyLocation, _, mtoPayload := setupTestData()

		// Set the status to CANCELED
		mtoPayload.Status = supportmessages.MoveStatusCANCELED

		// This time we provide customer details to create
		newCustomerFirstName := "Grace"
		mtoPayload.Order.Customer = &supportmessages.Customer{
			FirstName: &newCustomerFirstName,
			LastName:  swag.String("Griffin"),
			Agency:    swag.String("Marines"),
			DodID:     swag.String("1209457894"),
			Rank:      supportmessages.NewRank("ACADEMY_CADET"),
		}
		mtoPayload.Order.CustomerID = nil

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID was populated
		suite.NotEmpty(responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyLocation, &originDutyLocation)
		// Check that customer name matches the passed in value
		suite.Equal(newCustomerFirstName, *responsePayload.Order.Customer.FirstName)
		// Check that status has been set to CANCELED
		suite.Equal((models.MoveStatus)(responsePayload.Status), models.MoveStatusCANCELED)

	})
	suite.Run("Success createMoveTaskOrder discarded readOnly referenceID", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, and existing customer.
		//             We use a nonsense referenceID. ReferenceID is readOnly, this should not be applied.
		// Expected outcome:
		//             A new referenceID is generated.
		//             Default status is draft.
		destinationDutyLocation, originDutyLocation, _, mtoPayload := setupTestData()

		// Running the same request should result in the same reference id
		mtoPayload.ReferenceID = "some terrible reference id"
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderCreated{}, response)
		moveTaskOrdersResponse := response.(*movetaskorderops.CreateMoveTaskOrderCreated)
		responsePayload := moveTaskOrdersResponse.Payload

		// Check that the referenceID DOES NOT match what was sent in
		suite.NotEqual(mtoPayload.ReferenceID, responsePayload.ReferenceID)
		// Check that moveTaskOrder was populated, including nested objects
		suite.moveTaskOrderPopulated(moveTaskOrdersResponse, &destinationDutyLocation, &originDutyLocation)
	})

	suite.Run("Failed createMoveTaskOrder 422 UnprocessableEntity due to no customer", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, but no customer info
		// Expected outcome:
		//             Failure due to no customer info, so unprocessableEntity
		_, _, _, mtoPayload := setupTestData()

		mtoPayload.Order.Customer = nil
		mtoPayload.Order.CustomerID = nil

		// Running the same request should result in the same reference id
		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderUnprocessableEntity{}, response)
	})

	suite.Run("Failed createMoveTaskOrder 404 NotFound", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle
		// Mocked:     MoveTaskOrderCreator.CreateMoveTaskOrder
		// Set up:     We call the handler but force the mocked service object to return a notFoundError
		// Expected outcome:
		//             NotFound Response
		_, _, _, mtoPayload := setupTestData()

		mockCreator := supportMocks.InternalMoveTaskOrderCreator{}
		handler := setupHandler()
		handler.moveTaskOrderCreator = &mockCreator

		// Set expectation that a call to InternalCreateMoveTaskOrder will return a notFoundError
		notFoundError := apperror.NotFoundError{}
		mockCreator.On("InternalCreateMoveTaskOrder",
			mock.AnythingOfType("*appcontext.appContext"),
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

	suite.Run("Failed createMoveTaskOrder 404 NotFound Bad DutyLocation", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMoveTaskOrderHandler.Handle and MoveTaskOrderCreator.CreateMoveTaskOrder
		// Mocked:     None
		// Set up:     We pass in a new moveTaskOrder, order, and existing customer.
		//             The order has a bad duty location ID.
		// Expected outcome:
		//             Failure of 404 Not Found since the dutyLocation is not found.
		_, _, _, mtoPayload := setupTestData()

		// Using a randomID as a dutyLocationID should cause a query error
		mtoPayload.Order.OriginDutyLocationID = handlers.FmtUUID(uuid.Must(uuid.NewV4()))

		params := movetaskorderops.CreateMoveTaskOrderParams{
			HTTPRequest: request,
			Body:        mtoPayload,
		}

		handler := CreateMoveTaskOrderHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderNotFound{}, response)
	})
}
