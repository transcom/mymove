package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gobuffalo/validate"

	"github.com/stretchr/testify/mock"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/testdatagen"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestCreatePaymentRequestHandler() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	moveTaskOrderID := moveTaskOrder.ID
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: moveTaskOrderID,
		IsFinal:         false,
	}

	suite.T().Run("successful create payment request", func(t *testing.T) {
		returnedPaymentRequest := models.PaymentRequest{
			ID:              paymentRequestID,
			MoveTaskOrderID: moveTaskOrderID,
			MoveTaskOrder:   moveTaskOrder,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, validate.NewErrors(), nil).Once()

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
				IsFinal:         &paymentRequest.IsFinal,
				MoveTaskOrderID: *handlers.FmtUUID(paymentRequest.MoveTaskOrderID),
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
		typedResponse := response.(*paymentrequestop.CreatePaymentRequestCreated)
		suite.Equal(returnedPaymentRequest.ID.String(), typedResponse.Payload.ID.String())
		if suite.NotNil(typedResponse.Payload.IsFinal) {
			suite.Equal(returnedPaymentRequest.IsFinal, *typedResponse.Payload.IsFinal)
		}
		suite.Equal(returnedPaymentRequest.MoveTaskOrderID.String(), typedResponse.Payload.MoveTaskOrderID.String())
	})

	suite.T().Run("failed create payment request", func(t *testing.T) {
		badPaymentRequest := models.PaymentRequest{
			ID:              paymentRequestID,
			MoveTaskOrderID: uuid.UUID{},
		}
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&badPaymentRequest, nil, nil).Once()

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

	suite.T().Run("failed create payment request -- err is returned", func(t *testing.T) {
		returnedPaymentRequest := models.PaymentRequest{
			ID:              paymentRequestID,
			MoveTaskOrderID: moveTaskOrderID,
			MoveTaskOrder:   moveTaskOrder,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, validate.NewErrors(), errors.New("test failed to create with err returned")).Once()

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
				IsFinal:         &paymentRequest.IsFinal,
				MoveTaskOrderID: *handlers.FmtUUID(paymentRequest.MoveTaskOrderID),
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestInternalServerError{}, response)

	})

	suite.T().Run("failed create payment request -- non-existent MTO ID", func(t *testing.T) {
		returnedPaymentRequest := models.PaymentRequest{}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, validate.NewErrors(), nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, requestUser)

		badMTOID, _ := uuid.FromString("23000000-0000-0000-0000-000000000000")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequestPayload{
				IsFinal:         &paymentRequest.IsFinal,
				MoveTaskOrderID: *handlers.FmtUUID(badMTOID),
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestNotFound{}, response)
	})

}
