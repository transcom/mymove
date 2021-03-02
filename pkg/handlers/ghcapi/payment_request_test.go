package ghcapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"

	"net/http/httptest"
	"time"
)

func (suite *HandlerSuite) TestFetchPaymentRequestHandler() {

	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")

	paymentRequest := models.PaymentRequest{
		ID:        paymentRequestID,
		IsFinal:   false,
		Status:    models.PaymentRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	officeUserUUID, _ := uuid.NewV4()
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true, OfficeUser: models.OfficeUser{ID: officeUserUUID}})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	suite.T().Run("successful fetch of payment request", func(t *testing.T) {
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_request"), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestOK{}, response)
		okResponse := response.(*paymentrequestop.GetPaymentRequestOK)
		suite.Equal(paymentRequestID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("failed fetch for payment request - forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTOO,
		})
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_request"), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTOO)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestForbidden{}, response)
	})

	suite.T().Run("payment request not found", func(t *testing.T) {
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(models.PaymentRequest{}, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_request"), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestGetPaymentRequestsForMoveHandler() {
	expectedServiceItemName := "Test Service"
	expectedShipmentType := models.MTOShipmentTypeHHG

	move := testdatagen.MakeAvailableMove(suite.DB())
	// This should create all the other associated records we need.
	paymentServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
		Move: move,
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameRequestedPickupDate,
		},
	})
	paymentRequest := paymentServiceItemParam.PaymentServiceItem.PaymentRequest
	paymentRequests := models.PaymentRequests{paymentRequest}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	suite.T().Run("Successful list fetch", func(t *testing.T) {
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", move.Locator), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     move.Locator,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(suite.DB()),
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestsForMoveOK{}, response)
		paymentRequestsResponse := response.(*paymentrequestop.GetPaymentRequestsForMoveOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
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

	suite.T().Run("Failed list fetch - Not found error ", func(t *testing.T) {
		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestListByMove", mock.Anything,
			mock.Anything).Return(nil, errors.New("not found")).Once()

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", "ABC123"), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     "ABC123",
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentRequestListFetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
	})

	suite.T().Run("Failed list fetch - Forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestListByMove", mock.Anything,
			mock.Anything).Return(&paymentRequests, nil).Once()

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", "ABC123"), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUserTOO)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     "ABC123",
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentRequestListFetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestsForMoveForbidden{}, response)
	})
}

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	officeUserUUID, _ := uuid.NewV4()
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true, OfficeUser: models.OfficeUser{ID: officeUserUUID}})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	paymentRequest := models.PaymentRequest{
		ID:        paymentRequestID,
		IsFinal:   false,
		Status:    models.PaymentRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder(suite.DB()))

	suite.T().Run("successful status update of payment request", func(t *testing.T) {
		pendingPaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil, ETag: etag.GenerateEtag(pendingPaymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: statusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload
		suite.Equal(models.PaymentRequestStatusReviewed.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.T().Run("successful status update of rejected payment request", func(t *testing.T) {
		pendingPaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED", RejectionReason: nil, ETag: etag.GenerateEtag(paymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: statusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)
		suite.TestLogger().Error("")
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload
		suite.Equal(models.PaymentRequestStatusReviewedAllRejected.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.T().Run("failed status update of payment request - forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTOO,
		})

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&paymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTOO)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusForbidden(), response)

	})

	suite.T().Run("successful status update of prime-available payment request", func(t *testing.T) {
		availableMove := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		availablePaymentRequestID := availablePaymentRequest.ID

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&availablePaymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(availablePaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		handler.SetTraceID(traceID)

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		suite.HasWebhookNotification(availablePaymentRequestID, traceID)
	})

	suite.T().Run("unsuccessful status update of payment request (500)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, not found (404)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, services.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusNotFound(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, precondition failed (412)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, services.PreconditionFailedError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, validation errors (422)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, services.NewInvalidInputError(paymentRequestID, nil, nil, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity(), response)
	})
}
