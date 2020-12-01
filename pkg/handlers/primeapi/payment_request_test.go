package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/assert"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreatePaymentRequestHandler() {
	moveTaskOrderID, _ := uuid.FromString("96e21765-3e29-4acf-89a2-1317a9f7f0da")
	paymentRequestID, _ := uuid.FromString("70c0c9c1-cf3f-4195-b15c-d185dc5cd0bf")

	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	suite.T().Run("successful create payment request", func(t *testing.T) {
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   paymentRequestID,
			MoveTaskOrderID:      moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		serviceItemID1, _ := uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code:     models.ReServiceCodeDLH,
				Priority: 1,
			},
			MTOServiceItem: models.MTOServiceItem{
				ID: serviceItemID1,
			},
		})
		serviceItemID2, _ := uuid.FromString("119f0a05-34d7-4d86-9745-009c0707b4c2")
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code:     models.ReServiceCodeFSC,
				Priority: 99,
			},
			MTOServiceItem: models.MTOServiceItem{
				ID: serviceItemID2,
			},
		})

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.MatchedBy(func(paymentRequest *models.PaymentRequest) bool {
				// Making sure the service items are ordered by priority regardless of the order in which they come in through the payment request parameters
				return paymentRequest.PaymentServiceItems[0].MTOServiceItemID == serviceItemID1
			})).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID2),
					},
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
				},
				PointOfContact: "user@prime.com",
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
		suite.Equal(returnedPaymentRequest.PaymentRequestNumber, typedResponse.Payload.PaymentRequestNumber)
	})

	suite.T().Run("create payment request without adding service item params passed into payload", func(t *testing.T) {
		serviceItemID1, _ := uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   paymentRequestID,
			MoveTaskOrderID:      moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					ID: serviceItemID1,
				},
			},
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)
		typedResponse := response.(*paymentrequestop.CreatePaymentRequestCreated)

		paymentServiceItemParams := typedResponse.Payload.PaymentServiceItems[0].PaymentServiceItemParams

		suite.Equal(len(paymentServiceItemParams), 0)
		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
	})

	suite.T().Run("fail to create payment request adding service item params passed into payload", func(t *testing.T) {
		serviceItemID1, _ := uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   paymentRequestID,
			MoveTaskOrderID:      moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					ID: serviceItemID1,
				},
			},
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
						Params: []*primemessages.ServiceItemParamsItems0{
							{
								Key:   "weight",
								Value: "5678",
							},
						},
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create payment request -- nil body", func(t *testing.T) {

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

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

	suite.T().Run("failed create payment request -- creator failed with error", func(t *testing.T) {

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, errors.New("creator failed")).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				PointOfContact:  "user@prime.com",
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestInternalServerError{}, response)

		errResponse := response.(*paymentrequestop.CreatePaymentRequestInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong") // check body (body was written before panic)

	})

	suite.T().Run("failed create payment request -- invalid MTO ID format", func(t *testing.T) {

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		badFormatID := strfmt.UUID("hb7b134a-7c44-45f2-9114-bb0831cc5db3")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: &badFormatID,
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create payment request -- invalid service item ID format", func(t *testing.T) {

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				PointOfContact:  "user@prime.com",
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: badFormatID,
					},
				},
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create payment request - validation errors", func(t *testing.T) {
		verrs := &validate.Errors{
			Errors: map[string][]string{
				"violation": {"invalid value"},
			},
		}
		err := services.NewInvalidCreateInputError(verrs, "can't create payment request for MTO ID 1234")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

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
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
	})

	suite.T().Run("failed create payment request due to conflict in model", func(t *testing.T) {

		ordersID, _ := uuid.FromString("2b8b141a-7c44-45f2-9114-bb0831cc5db3")
		err := services.NewConflictError(ordersID, "incomplete orders")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

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
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestConflict{}, response)
	})

	suite.T().Run("failed create payment request due to bad data", func(t *testing.T) {

		err := services.NewBadDataError("sent some bad data, foo!")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

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
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)
	})

	suite.T().Run("successful create payment request payload audit", func(t *testing.T) {

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		serviceItemID1, _ := uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
		serviceItemID2, _ := uuid.FromString("119f0a05-34d7-4d86-9745-009c0707b4c2")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				PointOfContact:  "user@prime.com",
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(serviceItemID2),
					},
				},
			},
		}

		session := auth.Session{}

		// Capture creation attempt in audit log
		zapFields, err := audit.Capture(&params.Body, nil, suite.TestLogger(), &session, params.HTTPRequest)

		var eventType string
		for _, field := range zapFields {
			if field.Key == "event_type" {
				eventType = field.String
			}
		}
		suite.Nil(err, "No error for audit.Capture call")
		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "event_type", zapFields[0].Key)
			assert.Equal(t, "audit_post_payment_requests", eventType)
		}
	})
}
