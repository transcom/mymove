package supportapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	supportMocks "github.com/transcom/mymove/pkg/services/support/mocks"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"
)

func (suite *HandlerSuite) TestListMTOsHandler() {
	// unavailable MTO
	factory.BuildMove(suite.DB(), nil, nil)

	moveTaskOrder := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    moveTaskOrder,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    moveTaskOrder,
			LinkOnly: true,
		},
	}, nil)

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.ListMTOsParams{HTTPRequest: request}
	handlerConfig := suite.HandlerConfig()

	handler := ListMTOsHandler{
		HandlerConfig:        handlerConfig,
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

	suite.Run("successfully hide fake moves", func() {
		handler := HideNonFakeMoveTaskOrdersHandlerFunc{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderHider(),
		}
		var moves models.Moves

		mto1 := factory.BuildMove(suite.DB(), nil, nil)
		mto2 := factory.BuildMove(suite.DB(), nil, nil)
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
		move := factory.BuildMove(suite.DB(), nil, nil)
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
			suite.HandlerConfig(),
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
		mto := factory.BuildMove(suite.DB(), nil, nil)
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
			suite.HandlerConfig(),
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
	move := factory.BuildSubmittedMove(suite.DB(), nil, nil)
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/available-to-prime", nil)
	params := movetaskorderops.MakeMoveTaskOrderAvailableParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
		IfMatch:         etag.GenerateEtag(move.UpdatedAt),
	}
	handlerConfig := suite.HandlerConfig()
	queryBuilder := query.NewQueryBuilder()
	moveRouter, err := moverouter.NewMoveRouter()
	suite.FatalNoError(err)
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticPackPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticDestinationPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewFuelSurchargePricer(), suite.HandlerConfig().FeatureFlagFetcher())

	ppmEstimator := &mocks.PPMEstimator{}
	// make the request
	handler := MakeMoveTaskOrderAvailableHandlerFunc{handlerConfig,
		movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator),
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
	move := factory.BuildMove(suite.DB(), nil, nil)
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := movetaskorderops.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
	}

	handlerConfig := suite.HandlerConfig()
	handler := GetMoveTaskOrderHandlerFunc{handlerConfig,
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
		destinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		originDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		customer := factory.BuildServiceMember(suite.DB(), nil, nil)
		contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		document := factory.BuildDocument(suite.DB(), nil, nil)

		// Create the mto payload we will be requesting to create
		issueDate := models.TimePointer(time.Now())
		reportByDate := models.TimePointer(time.Now().AddDate(0, 0, -1))
		ordersTypedetail := supportmessages.OrdersTypeDetailHHGPERMITTED
		deptIndicator := supportmessages.DeptIndicatorAIRANDSPACEFORCE

		grade := (supportmessages.Rank)("E_6")
		mtoPayload := &supportmessages.MoveTaskOrder{
			PpmType:      "FULL",
			ContractorID: handlers.FmtUUID(contractor.ID),
			Order: &supportmessages.Order{
				Rank:                      &grade, // Convert support API "Rank" into our internal tracking of "Grade"
				OrderNumber:               models.StringPointer("4554"),
				DestinationDutyLocationID: handlers.FmtUUID(destinationDutyLocation.ID),
				OriginDutyLocationID:      handlers.FmtUUID(originDutyLocation.ID),
				Entitlement: &supportmessages.Entitlement{
					DependentsAuthorized: models.BoolPointer(true),
					TotalDependents:      5,
					NonTemporaryStorage:  models.BoolPointer(false),
				},
				IssueDate:           handlers.FmtDatePtr(issueDate),
				ReportByDate:        handlers.FmtDatePtr(reportByDate),
				OrdersType:          supportmessages.NewOrdersType("PERMANENT_CHANGE_OF_STATION"),
				OrdersTypeDetail:    &ordersTypedetail,
				UploadedOrdersID:    handlers.FmtUUID(document.ID),
				Status:              supportmessages.NewOrdersStatus(supportmessages.OrdersStatusDRAFT),
				Tac:                 models.StringPointer("E19A"),
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
			suite.HandlerConfig(),
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
		moveRouter, err := moverouter.NewMoveRouter()
		suite.FatalNoError(err)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)

		setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
			mockCreator := &mocks.SignedCertificationCreator{}

			mockCreator.On(
				"CreateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
			).Return(returnValue...)

			return mockCreator
		}

		setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
			mockUpdater := &mocks.SignedCertificationUpdater{}

			mockUpdater.On(
				"UpdateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
				mock.AnythingOfType("string"),
			).Return(returnValue...)

			return mockUpdater
		}

		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticPackPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewDomesticDestinationPricer(suite.HandlerConfig().FeatureFlagFetcher()), ghcrateengine.NewFuelSurchargePricer(), suite.HandlerConfig().FeatureFlagFetcher())

		ppmEstimator := &mocks.PPMEstimator{}
		// Submit the request to approve the MTO
		approvalHandler := MakeMoveTaskOrderAvailableHandlerFunc{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator),
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
			LastName:  models.StringPointer("Griffin"),
			Agency:    models.StringPointer("Marines"),
			DodID:     models.StringPointer("1209457894"),
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

		handler := CreateMoveTaskOrderHandler{
			suite.HandlerConfig(),
			internalmovetaskorder.NewInternalMoveTaskOrderCreator(),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// VERIFY RESULTS
		suite.IsType(&movetaskorderops.CreateMoveTaskOrderNotFound{}, response)
	})
}
