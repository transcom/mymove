package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	"testing"
)

func (suite *HandlerSuite) TestCreatePaymentRequestHandler() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	moveTaskOrderID := moveTaskOrder.ID
	serviceItemID, _ := uuid.NewV4()
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	rejectionReason := "Missing documentation"

	paymentRequest := models.PaymentRequest{
		ID:              paymentRequestID,
		MoveTaskOrderID: moveTaskOrderID,
		ServiceItemIDs:  []uuid.UUID{serviceItemID},
		RejectionReason: rejectionReason,
	}

	suite.T().Run("successful create payment request", func(t *testing.T) {
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
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequestPayload{
				IsFinal:               &paymentRequest.IsFinal,
				MoveTaskOrderID:       *handlers.FmtUUID(paymentRequest.MoveTaskOrderID),
				ProofOfServicePackage: nil,
				ServiceItemIDs:        []strfmt.UUID{*handlers.FmtUUID(paymentRequest.ServiceItemIDs[0])},
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
	})

	suite.T().Run("failed create payment request", func(t *testing.T) {
		badPaymentRequest := models.PaymentRequest{
			ID:              paymentRequestID,
			MoveTaskOrderID: uuid.UUID{},
			ServiceItemIDs:  []uuid.UUID{serviceItemID},
			RejectionReason: rejectionReason,
		}
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			&paymentRequest,
			mock.Anything).Return(&badPaymentRequest, nil, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)

	})

}
