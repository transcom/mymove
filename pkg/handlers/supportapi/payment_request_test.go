package supportapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/etag"
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
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
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
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
	mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
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
