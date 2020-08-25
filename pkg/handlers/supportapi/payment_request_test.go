package supportapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/gen/supportmessages"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	paymentRequestID := paymentRequest.ID

	suite.T().Run("successful status update of payment request", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		eTag := etag.GenerateEtag(paymentRequest.UpdatedAt)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil, ETag: eTag},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
			IfMatch:          eTag,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
	})

	suite.T().Run("successful status update of prime-available payment request", func(t *testing.T) {
		availableMove := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		availablePaymentRequestID := availablePaymentRequest.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		eTag := etag.GenerateEtag(availablePaymentRequest.UpdatedAt)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil, ETag: eTag},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
			IfMatch:          eTag,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		// Set up the observer so we can check the log output later:
		logCore, logs := observer.New(zap.InfoLevel)

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), zap.New(logCore)),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
		}
		traceID, err := uuid.NewV4()
		suite.NoError(err, "Error creating a new trace ID.")
		handler.SetTraceID(traceID)

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

		notification := &models.WebhookNotification{}
		err = suite.DB().Where("object_id = $1 AND trace_id = $2", availablePaymentRequestID.String(), traceID.String()).First(notification)
		suite.NoError(err)
		suite.Equal(string(notification.Status), "PENDING")

		suite.Greater(logs.Len(), 0)
		if logs.Len() > 0 {
			entry := logs.All()[0]
			suite.Contains(entry.Message, "Notification created and stored.")
		}
	})

	suite.T().Run("unsuccessful status update of payment request (500)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.UpdatePaymentRequestStatusInternalServerError{}, response)

		errResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})

	suite.T().Run("unsuccessful status update of payment request, not found (404)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, services.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
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

		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
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
}

func (suite *HandlerSuite) TestListMTOPaymentRequestHandler() {
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	mto := testdatagen.MakeDefaultMove(suite.DB())
	suite.T().Run("successful get an MTO with payment requests", func(t *testing.T) {
		mtoID := paymentRequest.MoveTaskOrderID
		req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/payment-requests", mtoID), nil)

		params := paymentrequestop.ListMTOPaymentRequestsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: strfmt.UUID(mtoID.String()),
		}

		handler := ListMTOPaymentRequestsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ListMTOPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		suite.IsType(paymentrequestop.NewListMTOPaymentRequestsOK(), response)
		suite.Equal(len(paymentRequestsPayload), 1)
	})

	suite.T().Run("successful get an MTO with no payment requests", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/payment-requests", mto.ID), nil)

		params := paymentrequestop.ListMTOPaymentRequestsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: strfmt.UUID(mto.ID.String()),
		}

		handler := ListMTOPaymentRequestsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ListMTOPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		suite.IsType(paymentrequestop.NewListMTOPaymentRequestsOK(), response)
		suite.Equal(len(paymentRequestsPayload), 0)
	})
}
