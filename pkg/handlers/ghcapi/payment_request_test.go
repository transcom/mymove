package ghcapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	"net/http/httptest"
	"time"
)

func (suite *HandlerSuite) TestListPaymentRequestsHandler() {

	paymentRequestID1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	paymentRequestID2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	paymentRequestID3, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	IDs := []uuid.UUID{
		paymentRequestID1,
		paymentRequestID2,
		paymentRequestID3,
	}

	var paymentRequests models.PaymentRequests

	for _, id := range IDs {
		paymentRequest := models.PaymentRequest{
			ID:        id,
			IsFinal:   false,
			Status:    models.PaymentRequestStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		paymentRequests = append(paymentRequests, paymentRequest)
	}

	suite.T().Run("successful fetch of payment requests", func(t *testing.T) {
		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestList").Return(&paymentRequests, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_requests"), nil)

		params := paymentrequestop.ListPaymentRequestsParams{
			HTTPRequest: req,
		}

		handler := ListPaymentRequestsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestListFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.ListPaymentRequestsOK{}, response)
		okResponse := response.(*paymentrequestop.ListPaymentRequestsOK)
		suite.Equal(len(IDs), len(okResponse.Payload))
		suite.Equal(paymentRequestID1.String(), okResponse.Payload[0].ID.String())
	})

	suite.T().Run("failed fetch of payment requests", func(t *testing.T) {
		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestList").Return(nil, errors.New("test failed to create with err returned")).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_requests"), nil)

		params := paymentrequestop.ListPaymentRequestsParams{
			HTTPRequest: req,
		}

		handler := ListPaymentRequestsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestListFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.ListPaymentRequestsInternalServerError{}, response)
	})
}