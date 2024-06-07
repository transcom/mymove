// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestGetMoveTaskOrderHandlerIntegration() {
	moveTaskOrder := factory.BuildMove(suite.DB(), nil, nil)
	factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeMS)
	factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeCS)

	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := movetaskorderops.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetMoveTaskOrderHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveTaskOrderResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	moveTaskOrderPayload := moveTaskOrderResponse.Payload

	// Validate outgoing payload
	suite.NoError(moveTaskOrderPayload.Validate(strfmt.Default))

	suite.Assertions.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)
	suite.Equal(strfmt.UUID(moveTaskOrder.ID.String()), moveTaskOrderPayload.ID)
	suite.Nil(moveTaskOrderPayload.AvailableToPrimeAt)
	// TODO: Check that the *moveTaskOrderPayload.Status is not "canceled"
	// suite.False(*moveTaskOrderPayload.IsCanceled)
	suite.Equal(strfmt.UUID(moveTaskOrder.OrdersID.String()), moveTaskOrderPayload.OrderID)
	suite.NotNil(moveTaskOrderPayload.ReferenceID)
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationSuccess() {
	factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeMS)
	factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeCS)

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

	validStatuses := []struct {
		desc   string
		status models.MoveStatus
	}{
		{"Submitted", models.MoveStatusSUBMITTED},
		{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
	}
	for _, validStatus := range validStatuses {
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: validStatus.status,
				},
			},
		}, nil)

		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveID}/status", nil)
		requestUser := factory.BuildOfficeUser(nil, nil, nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		request = request.WithContext(trace.NewContext(request.Context(), traceID))

		serviceItemCodes := ghcmessages.MTOApprovalServiceItemCodes{
			ServiceCodeMS: true,
			ServiceCodeCS: true,
		}
		params := movetaskorderops.UpdateMoveTaskOrderStatusParams{
			HTTPRequest:      request,
			MoveTaskOrderID:  move.ID.String(),
			IfMatch:          etag.GenerateEtag(move.UpdatedAt),
			ServiceItemCodes: &serviceItemCodes,
		}
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal"))
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		// setup the handler
		handler := UpdateMoveTaskOrderStatusHandlerFunc{handlerConfig,
			movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil)),
		}

		// Validate incoming payload: no body to validate

		// make the request
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		moveResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		updatedMove := models.Move{}
		suite.DB().Find(&updatedMove, movePayload.ID)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)

		suite.Assertions.IsType(&movetaskorderops.UpdateMoveTaskOrderStatusOK{}, response)
		suite.Equal(strfmt.UUID(move.ID.String()), movePayload.ID)
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.HasWebhookNotification(move.ID, traceID) // this action always creates a notification for the Prime

		// also check MTO level service items are properly created
		var serviceItems models.MTOServiceItems
		suite.DB().Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
		suite.Len(serviceItems, 2, "Expected to find at most 2 service items")

		containsServiceCode := func(items models.MTOServiceItems, target models.ReServiceCode) bool {
			for _, si := range items {
				if si.ReService.Code == target {
					return true
				}
			}

			return false
		}

		suite.True(containsServiceCode(serviceItems, models.ReServiceCodeMS), fmt.Sprintf("Expected to find reServiceCode, %s, in array.", models.ReServiceCodeMS))
		suite.True(containsServiceCode(serviceItems, models.ReServiceCodeCS), fmt.Sprintf("Expected to find reServiceCode, %s, in array.", models.ReServiceCodeCS))
	}
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationWithStaleEtag() {
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		},
	}, nil)

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	requestUser := factory.BuildUser(nil, nil, nil)
	request = suite.AuthenticateUserRequest(request, requestUser)
	params := movetaskorderops.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
		IfMatch:         etag.GenerateEtag(time.Now()),
	}
	handlerConfig := suite.HandlerConfig()

	// Stale ETags are already unit tested in the move_task_order_updater_test,
	// so we can mock this here to speed up the test and avoid hitting the DB
	moveUpdater := &mocks.MoveTaskOrderUpdater{}
	moveUpdater.On("MakeAvailableToPrime",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
		false,
		false,
	).Return(nil, apperror.PreconditionFailedError{})

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{handlerConfig, moveUpdater}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.Assertions.IsType(&movetaskorderops.UpdateMoveTaskOrderStatusPreconditionFailed{}, response)
	payload := response.(*movetaskorderops.UpdateMoveTaskOrderStatusPreconditionFailed).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationWithIncompleteOrder() {
	orderWithoutDefaults := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusServiceCounselingCompleted,
			},
		},
		{
			Model:    orderWithoutDefaults,
			LinkOnly: true,
		},
	}, nil)

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	requestUser := factory.BuildUser(nil, nil, nil)
	request = suite.AuthenticateUserRequest(request, requestUser)
	params := movetaskorderops.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
		IfMatch:         etag.GenerateEtag(move.UpdatedAt),
	}
	handlerConfig := suite.HandlerConfig()
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
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

	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{handlerConfig,
		movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil)),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)

	suite.Assertions.IsType(&movetaskorderops.UpdateMoveTaskOrderStatusUnprocessableEntity{}, response)
	invalidResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusUnprocessableEntity).Payload

	// Validate outgoing payload
	suite.NoError(invalidResponse.Validate(strfmt.Default))

	errorDetail := invalidResponse.Detail

	suite.Contains(*errorDetail, "TransportationAccountingCode cannot be blank.")
	suite.Contains(*errorDetail, "OrdersNumber cannot be blank.")
	suite.Contains(*errorDetail, "DepartmentIndicator cannot be blank.")
	suite.Contains(*errorDetail, "OrdersTypeDetail cannot be blank.")
}

func (suite *HandlerSuite) TestUpdateMTOStatusServiceCounselingCompletedHandler() {
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

	setupTestData := func() UpdateMTOStatusServiceCounselingCompletedHandlerFunc {
		handlerConfig := suite.HandlerConfig()
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := UpdateMTOStatusServiceCounselingCompletedHandlerFunc{
			handlerConfig,
			movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil)),
		}
		return handler
	}

	suite.Run("Successful move status update to Service Counseling Completed - Integration", func() {
		handler := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
		}, nil)

		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: move.ID.String(),
			IfMatch:         etag.GenerateEtag(move.UpdatedAt),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		moveTaskOrderResponse := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedOK)
		moveTaskOrderPayload := moveTaskOrderResponse.Payload

		// Validate outgoing payload
		suite.NoError(moveTaskOrderPayload.Validate(strfmt.Default))

		suite.IsType(&movetaskorderops.UpdateMTOStatusServiceCounselingCompletedOK{}, response)
		suite.Equal(strfmt.UUID(move.ID.String()), moveTaskOrderPayload.ID)
		suite.NotNil(moveTaskOrderPayload.ServiceCounselingCompletedAt)
		suite.EqualValues(models.MoveStatusServiceCounselingCompleted, moveTaskOrderPayload.Status)
	})

	suite.Run("Unsuccessful move status update to Service Counseling Completed, forbidden - Integration", func() {
		handler := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		forbiddenUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		forbiddenRequest := suite.AuthenticateOfficeRequest(request, forbiddenUser)

		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest: forbiddenRequest,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.UpdateMTOStatusServiceCounselingCompletedForbidden{}, response)
		payload := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedForbidden).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful move status update to Service Counseling Completed, not found - Integration", func() {
		handler := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)
		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: uuid.Must(uuid.NewV4()).String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.UpdateMTOStatusServiceCounselingCompletedNotFound{}, response)
		payload := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move status update to Service Counseling Completed, eTag does not match - Integration", func() {
		handler := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)
		move := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: move.ID.String(),
			IfMatch:         etag.GenerateEtag(time.Now()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.UpdateMTOStatusServiceCounselingCompletedPreconditionFailed{}, response)
		payload := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful move status update to Service Counseling Completed, state conflict - Integration", func() {
		handler := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)
		draftMove := factory.BuildMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    draftMove,
				LinkOnly: true,
			},
		}, nil)

		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: draftMove.ID.String(),
			IfMatch:         etag.GenerateEtag(draftMove.UpdatedAt),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.UpdateMTOStatusServiceCounselingCompletedConflict{}, response)
		payload := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful move status update to Service Counseling Completed, misc mocked errors - Integration", func() {
		handlerOrig := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
		requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)
		testCases := []struct {
			mockError       error
			handlerResponse middleware.Responder
		}{
			{apperror.InvalidInputError{}, &movetaskorderops.UpdateMTOStatusServiceCounselingCompletedUnprocessableEntity{}},
			{apperror.QueryError{}, &movetaskorderops.UpdateMTOStatusServiceCounselingCompletedInternalServerError{}},
			{errors.New("generic error"), &movetaskorderops.UpdateMTOStatusServiceCounselingCompletedInternalServerError{}},
		}

		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusNeedsServiceCounseling)
		eTag := etag.GenerateEtag(move.UpdatedAt)
		params := movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: move.ID.String(),
			IfMatch:         eTag,
		}

		for _, testCase := range testCases {
			mockUpdater := mocks.MoveTaskOrderUpdater{}
			mockUpdater.On("UpdateStatusServiceCounselingCompleted",
				mock.AnythingOfType("*appcontext.appContext"),
				move.ID,
				eTag,
			).Return(&models.Move{}, testCase.mockError)

			handler := UpdateMTOStatusServiceCounselingCompletedHandlerFunc{
				handlerOrig.HandlerConfig,
				&mockUpdater,
			}

			// Validate incoming payload: no body to validate

			response := handler.Handle(params)
			suite.IsNotErrResponse(response)
			suite.IsType(testCase.handlerResponse, response)

			// Validate outgoing payload
			switch response := response.(type) {
			case *movetaskorderops.UpdateMTOStatusServiceCounselingCompletedUnprocessableEntity:
				suite.NoError(response.Payload.Validate(strfmt.Default))
			case *movetaskorderops.UpdateMTOStatusServiceCounselingCompletedInternalServerError:
				suite.Nil(response.Payload)
			default:
				suite.Fail(fmt.Sprintf("unexpected response type of %T", response))
			}
		}
	})
}

func (suite *HandlerSuite) TestUpdateMoveTIORemarksHandler() {
	setupTestData := func() (models.Move, UpdateMoveTIORemarksHandlerFunc, models.User) {
		moveTaskOrder := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
		}, nil)
		requestUser := factory.BuildUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
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

		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := UpdateMoveTIORemarksHandlerFunc{
			handlerConfig,
			movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil)),
		}
		return moveTaskOrder, handler, requestUser
	}

	remarks := "Reweigh requested"
	suite.Run("Successfully update the Move's TIORemarks field", func() {
		move, handler, requestUser := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/tio-remarks", nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := movetaskorderops.UpdateMoveTIORemarksParams{
			HTTPRequest:     request,
			MoveTaskOrderID: move.ID.String(),
			Body:            &ghcmessages.Move{TioRemarks: &remarks},
			IfMatch:         etag.GenerateEtag(move.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		moveTaskOrderResponse := response.(*movetaskorderops.UpdateMoveTIORemarksOK)
		moveTaskOrderPayload := moveTaskOrderResponse.Payload

		// Validate outgoing payload
		suite.NoError(moveTaskOrderPayload.Validate(strfmt.Default))

		suite.Assertions.IsType(&movetaskorderops.UpdateMoveTIORemarksOK{}, response)
		updatedMove := models.Move{}
		suite.DB().Find(&updatedMove, moveTaskOrderPayload.ID)
		suite.Equal(moveTaskOrderPayload.TioRemarks, updatedMove.TIORemarks)
	})

	suite.Run("Unsuccessful move TIO Remarks, eTag does not match", func() {
		move, handler, requestUser := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/tio-remarks", nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := movetaskorderops.UpdateMoveTIORemarksParams{
			HTTPRequest:     request,
			MoveTaskOrderID: move.ID.String(),
			Body:            &ghcmessages.Move{TioRemarks: &remarks},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&movetaskorderops.UpdateMoveTIORemarksPreconditionFailed{}, response)
		payload := response.(*movetaskorderops.UpdateMoveTIORemarksPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful move TIO Remarks update, not found", func() {
		_, handler, requestUser := setupTestData()
		request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/tio-remarks", nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := movetaskorderops.UpdateMoveTIORemarksParams{
			HTTPRequest:     request,
			MoveTaskOrderID: uuid.Must(uuid.NewV4()).String(),
			Body:            &ghcmessages.Move{TioRemarks: &remarks},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&movetaskorderops.UpdateMoveTIORemarksNotFound{}, response)
		payload := response.(*movetaskorderops.UpdateMoveTIORemarksNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
