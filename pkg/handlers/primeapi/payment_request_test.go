package primeapi

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/testdatagen"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"

	"testing"
)

func (suite *HandlerSuite) TestCreatePaymentRequestHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	rejectionReason := "Missing documentation"
	paymentRequest := models.PaymentRequest{
		ID: paymentRequestID,
		MoveTaskOrderID: moveTaskOrderID,
		ServiceItemIDs: []uuid.UUID{serviceItemID},
		RejectionReason: rejectionReason,
	}
	suite.T().Run("create payment request", func(t *testing.T) {
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			&paymentRequest,
			mock.Anything).Return(&paymentRequest, nil, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest:     req,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
	})
}
