package ghcapi

import (
	"errors"
	"fmt"
	"testing"

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

	suite.T().Run("successful fetch of payment request", func(t *testing.T) {
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_request"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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
	suite.T().Run("payment request not found", func(t *testing.T) {
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(models.PaymentRequest{}, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_request"), nil)

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

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")

	paymentRequest := models.PaymentRequest{
		ID:        paymentRequestID,
		IsFinal:   false,
		Status:    models.PaymentRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.T().Run("successful status update of payment request", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&paymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

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
