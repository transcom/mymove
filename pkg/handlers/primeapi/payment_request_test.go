package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
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
	moveTaskOrderID, _ := uuid.FromString("96e21765-3e29-4acf-89a2-1317a9f7f0da")
	paymentRequestID, _ := uuid.FromString("70c0c9c1-cf3f-4195-b15c-d185dc5cd0bf")

	requestUser := testdatagen.MakeDefaultUser(suite.DB())

	suite.T().Run("successful create payment request", func(t *testing.T) {
		returnedPaymentRequest := models.PaymentRequest{
			ID:              paymentRequestID,
			MoveTaskOrderID: moveTaskOrderID,
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
		req = suite.AuthenticateUserRequest(req, requestUser)

		serviceItemID1, _ := uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
		serviceItemID2, _ := uuid.FromString("119f0a05-34d7-4d86-9745-009c0707b4c2")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequestPayload{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: *handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
						Params: []*primemessages.ServiceItemParamsItems0{
							{
								Key:   "weight",
								Value: "1234",
							},
							{
								Key:   "pickup",
								Value: "2019-12-16",
							},
						},
					},
					{
						ID: *handlers.FmtUUID(serviceItemID2),
						Params: []*primemessages.ServiceItemParamsItems0{
							{
								Key:   "weight",
								Value: "5678",
							},
						},
					},
				},
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

	suite.T().Run("failed create payment request -- nil body", func(t *testing.T) {
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)
	})

	suite.T().Run("failed create payment request -- creator failed", func(t *testing.T) {
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, validate.NewErrors(), errors.New("creator failed")).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequestPayload{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: *handlers.FmtUUID(moveTaskOrderID),
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestInternalServerError{}, response)
	})

	suite.T().Run("failed create payment request -- invalid MTO ID format", func(t *testing.T) {
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, validate.NewErrors(), nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		badFormatID := strfmt.UUID("hb7b134a-7c44-45f2-9114-bb0831cc5db3")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequestPayload{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: badFormatID,
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)
	})
}
