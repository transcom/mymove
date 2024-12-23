package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/move"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestFetchPaymentRequestHandler() {
	expectedServiceItemName := "Domestic linehaul"
	expectedShipmentType := models.MTOShipmentTypeHHG

	setupTestData := func() (models.PaymentServiceItemParam, models.OfficeUser) {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		// This should create all the other associated records we need.
		paymentServiceItemParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ServiceItemParamKey{
					Key:  models.ServiceItemParamNameRequestedPickupDate,
					Type: models.ServiceItemParamTypeDate,
				},
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTIO,
		})
		return paymentServiceItemParam, officeUser
	}

	suite.Run("successful fetch of payment request", func() {
		paymentServiceItemParam, officeUser := setupTestData()
		paymentRequest := paymentServiceItemParam.PaymentServiceItem.PaymentRequest

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s", paymentRequest.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := GetPaymentRequestHandler{
			suite.HandlerConfig(),
			paymentrequest.NewPaymentRequestFetcher(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestOK{}, response)
		okResponse := response.(*paymentrequestop.GetPaymentRequestOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		paymentServiceItemParamPayload := payload.ServiceItems[0].PaymentServiceItemParams[0]

		suite.Equal(paymentRequest.ID.String(), payload.ID.String())
		suite.Equal(expectedServiceItemName, payload.ServiceItems[0].MtoServiceItemName)
		suite.EqualValues(expectedShipmentType, payload.ServiceItems[0].MtoShipmentType)

		suite.Equal(1, len(payload.ServiceItems))
		suite.Equal(paymentServiceItemParam.PaymentServiceItemID.String(), payload.ServiceItems[0].ID.String())
		suite.Equal(1, len(payload.ServiceItems[0].PaymentServiceItemParams))
		suite.Equal(paymentServiceItemParam.ID.String(), paymentServiceItemParamPayload.ID.String())
		suite.EqualValues(models.ServiceItemParamNameRequestedPickupDate, paymentServiceItemParamPayload.Key)
		suite.Equal(paymentServiceItemParam.Value, paymentServiceItemParamPayload.Value)
	})

	suite.Run("payment request not found", func() {
		paymentServiceItemParam, officeUser := setupTestData()
		paymentRequest := paymentServiceItemParam.PaymentServiceItem.PaymentRequest
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(models.PaymentRequest{}, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s", paymentRequest.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := GetPaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
		payload := response.(*paymentrequestop.GetPaymentRequestNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestGetPaymentRequestsForMoveHandler() {
	expectedServiceItemName := "Domestic linehaul"
	expectedShipmentType := models.MTOShipmentTypeHHG
	var moveLocator string

	setupTestData := func() (models.PaymentServiceItemParam, models.OfficeUser) {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		moveLocator = move.Locator

		// This should create all the other associated records we need.
		paymentServiceItemParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ServiceItemParamKey{
					Key:  models.ServiceItemParamNameRequestedPickupDate,
					Type: models.ServiceItemParamTypeDate,
				},
			},
		}, nil)

		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTIO,
		})
		return paymentServiceItemParam, officeUser
	}

	suite.Run("Successful list fetch", func() {
		paymentServiceItemParam, officeUser := setupTestData()
		paymentRequests := models.PaymentRequests{paymentServiceItemParam.PaymentServiceItem.PaymentRequest}

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", moveLocator), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     moveLocator,
		}
		handlerConfig := suite.HandlerConfig()
		handler := GetPaymentRequestForMoveHandler{
			HandlerConfig:             handlerConfig,
			PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.GetPaymentRequestsForMoveOK{}, response)
		paymentRequestsResponse := response.(*paymentrequestop.GetPaymentRequestsForMoveOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		// Validate outgoing payload
		suite.NoError(paymentRequestsPayload.Validate(strfmt.Default))

		paymentServiceItemParamPayload := paymentRequestsPayload[0].ServiceItems[0].PaymentServiceItemParams[0]

		suite.Equal(1, len(paymentRequestsPayload))
		suite.Equal(paymentRequests[0].ID.String(), paymentRequestsPayload[0].ID.String())
		suite.Equal(expectedServiceItemName, paymentRequestsPayload[0].ServiceItems[0].MtoServiceItemName)
		suite.EqualValues(expectedShipmentType, paymentRequestsPayload[0].ServiceItems[0].MtoShipmentType)

		suite.Equal(1, len(paymentRequestsPayload[0].ServiceItems))
		suite.Equal(paymentServiceItemParam.PaymentServiceItemID.String(), paymentRequestsPayload[0].ServiceItems[0].ID.String())
		suite.Equal(1, len(paymentRequestsPayload[0].ServiceItems[0].PaymentServiceItemParams))
		suite.Equal(paymentServiceItemParam.ID.String(), paymentServiceItemParamPayload.ID.String())
		suite.EqualValues(models.ServiceItemParamNameRequestedPickupDate, paymentServiceItemParamPayload.Key)
		suite.Equal(paymentServiceItemParam.Value, paymentServiceItemParamPayload.Value)
	})

	suite.Run("Failed list fetch - Not found error ", func() {
		_, officeUser := setupTestData()
		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestListByMove", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything).Return(nil, errors.New("not found")).Once()

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", "ABC123"), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     "ABC123",
		}
		handlerConfig := suite.HandlerConfig()
		handler := GetPaymentRequestForMoveHandler{
			HandlerConfig:             handlerConfig,
			PaymentRequestListFetcher: paymentRequestListFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
		payload := response.(*paymentrequestop.GetPaymentRequestNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")

	setupTestData := func() models.OfficeUser {

		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: true,
				},
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		}, []roles.RoleType{roles.RoleTypeTIO})

		return officeUser
	}
	paymentRequest := models.PaymentRequest{
		ID:        paymentRequestID,
		IsFinal:   false,
		Status:    models.PaymentRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	suite.Run("successful status update of payment request", func() {
		officeUser := setupTestData()
		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil, ETag: etag.GenerateEtag(pendingPaymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			PaymentRequestStatusUpdater:   statusUpdater,
			PaymentRequestFetcher:         paymentRequestFetcher,
			MoveAssignedOfficeUserUpdater: move.AssignedOfficeUserUpdater{},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(models.PaymentRequestStatusReviewed.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.Run("successful status update of rejected payment request", func() {
		officeUser := setupTestData()
		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED", RejectionReason: nil, ETag: etag.GenerateEtag(paymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			PaymentRequestStatusUpdater:   statusUpdater,
			PaymentRequestFetcher:         paymentRequestFetcher,
			MoveAssignedOfficeUserUpdater: move.AssignedOfficeUserUpdater{},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.Logger().Error("")
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(models.PaymentRequestStatusReviewedAllRejected.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.Run("prevent handler from updating payment request status to unapproved statuses", func() {
		officeUser := setupTestData()
		nonApprovedPRStatuses := [...]ghcmessages.PaymentRequestStatus{
			ghcmessages.PaymentRequestStatusSENTTOGEX,
			ghcmessages.PaymentRequestStatusTPPSRECEIVED,
			ghcmessages.PaymentRequestStatusPAID,
			ghcmessages.PaymentRequestStatusEDIERROR,
			ghcmessages.PaymentRequestStatusPENDING,
			ghcmessages.PaymentRequestStatusDEPRECATED,
		}

		for _, nonApprovedPRStatus := range nonApprovedPRStatuses {
			pendingPaymentRequest := factory.BuildPaymentRequest(nil, nil, nil)

			paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
			paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
				pendingPaymentRequest.ID).Return(pendingPaymentRequest, nil).Once()

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
			req = suite.AuthenticateOfficeRequest(req, officeUser)
			params := paymentrequestop.UpdatePaymentRequestStatusParams{
				HTTPRequest:      req,
				Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: nonApprovedPRStatus, RejectionReason: nil, ETag: etag.GenerateEtag(paymentRequest.UpdatedAt)},
				PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
			}

			handler := UpdatePaymentRequestStatusHandler{
				HandlerConfig:               suite.HandlerConfig(),
				PaymentRequestStatusUpdater: statusUpdater,
				PaymentRequestFetcher:       paymentRequestFetcher,
			}

			// Validate incoming payload
			suite.NoError(params.Body.Validate(strfmt.Default))

			response := handler.Handle(params)
			suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity(), response)
			payload := response.(*paymentrequestop.UpdatePaymentRequestStatusUnprocessableEntity).Payload

			// Validate outgoing payload
			suite.NoError(payload.Validate(strfmt.Default))
		}
	})

	suite.Run("successful status update of prime-available payment request", func() {
		officeUser := setupTestData()
		availableMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		availablePaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		availablePaymentRequestID := availablePaymentRequest.ID

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(&availablePaymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(availablePaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               suite.HandlerConfig(),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(availablePaymentRequestID, traceID)
	})

	suite.Run("unsuccessful status update of payment request (500)", func() {
		officeUser := setupTestData()
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               suite.HandlerConfig(),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("unsuccessful status update of payment request, not found (404)", func() {
		officeUser := setupTestData()
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               suite.HandlerConfig(),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusNotFound(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("unsuccessful status update of payment request, precondition failed (412)", func() {
		officeUser := setupTestData()
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.PreconditionFailedError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               suite.HandlerConfig(),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("unsuccessful status update of payment request, validation errors (422)", func() {
		officeUser := setupTestData()
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.NewInvalidInputError(paymentRequestID, nil, nil, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               suite.HandlerConfig(),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity(), response)
		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestShipmentsSITBalanceHandler() {

	setupTestData := func() models.OfficeUser {
		officeUserTIO := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTIO})
		return officeUserTIO
	}

	suite.Run("successful response of the shipments SIT Balance handler", func() {
		officeUserTIO := setupTestData()

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		thirtyDaySITExtensionRequest := 30
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status:        models.SITExtensionStatusApproved,
					RequestedDays: thirtyDaySITExtensionRequest,
					ApprovedDays:  &thirtyDaySITExtensionRequest,
				},
			},
		}, nil)

		reviewedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					SequenceNumber: 2,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)

		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
		originDepartureDate := originEntryDate.AddDate(0, 0, 90)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &originEntryDate,
					SITDepartureDate: &originDepartureDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		destinationEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		ddasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &destinationEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationPaymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "60",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s/shipments-payment-sit-balance", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTIO)

		params := paymentrequestop.GetShipmentsPaymentSITBalanceParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := ShipmentsSITBalanceHandler{
			HandlerConfig:              suite.HandlerConfig(),
			ShipmentsPaymentSITBalance: paymentrequest.NewPaymentRequestShipmentsSITBalance(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetShipmentsPaymentSITBalanceOK{}, response)
		payload := response.(*paymentrequestop.GetShipmentsPaymentSITBalanceOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.NotNil(payload)
		suite.Len(payload, 1)
		shipmentSITBalance := payload[0]

		// Destination SIT had a SIT entry date of 90 days before today
		// Because Destination SIT did not receive a departure date, we go based off today
		// Meaning Destination SIT has spent 90 days in SIT so far
		// Origin SIT received a SIT departure date 90 days after its entry date
		// Meaning 90 + 90 = 180 days spent in SIT, meaning this test went well over its entitlement
		// of 120
		suite.Equal(shipment.ID.String(), shipmentSITBalance.ShipmentID.String())
		suite.Equal(int64(120), shipmentSITBalance.TotalSITDaysAuthorized)
		suite.Equal(int64(60), shipmentSITBalance.PendingSITDaysInvoiced)
		// Since there is no departure date on one of the SITs, +1 is added to the count to count the last day
		suite.Equal(int64(-62), shipmentSITBalance.TotalSITDaysRemaining) // Well over entitlement
		suite.Equal(destinationPaymentEndDate.Format("2006-01-02"), shipmentSITBalance.PendingBilledEndDate.String())
		suite.Equal(int64(30), *shipmentSITBalance.PreviouslyBilledDays)
	})

	suite.Run("returns 404 not found when payment request does not exist", func() {
		officeUserTIO := setupTestData()
		paymentRequestID := uuid.Must(uuid.NewV4())

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s/shipments-payment-sit-balance", paymentRequestID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTIO)

		params := paymentrequestop.GetShipmentsPaymentSITBalanceParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := ShipmentsSITBalanceHandler{
			HandlerConfig:              suite.HandlerConfig(),
			ShipmentsPaymentSITBalance: paymentrequest.NewPaymentRequestShipmentsSITBalance(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetShipmentsPaymentSITBalanceNotFound{}, response)
		payload := response.(*paymentrequestop.GetShipmentsPaymentSITBalanceNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
