package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListPaymentRequestHandler() {

	// Happy path. API end point right now just returns mock data. TODO: Add more test conditions when actual stuff is added.
	suite.T().Run("The list endpoint returns the mocked objects", func(t *testing.T) {
		paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

		returnedPaymentRequests := []models.PaymentRequest{
			{
				ID:        paymentRequestID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		paymentRequestLister := &mocks.PaymentRequestLister{}
		paymentRequestLister.On("ListPaymentRequests").Return(&returnedPaymentRequests, nil, nil).Once()

		handler := ListPaymentRequestsHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestLister: paymentRequestLister,
		}

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment_requests"), nil)
		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := payment_requests.ListPaymentRequestsParams{
			HTTPRequest: req,
		}
		response := handler.Handle(params)

		suite.IsType(&payment_requests.ListPaymentRequestsOK{}, response)
	})

}