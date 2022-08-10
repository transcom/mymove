package primeapi

import (
	"errors"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/assert"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/services/audit"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	dlhTestServiceArea = "004"
	dlhTestWeight      = unit.Pound(4000)
)

type createPaymentRequestHandlerSubtestData struct {
	moveTaskOrderID  uuid.UUID
	paymentRequestID uuid.UUID
	serviceItemID1   uuid.UUID
	serviceItemID2   uuid.UUID
	serviceItemID3   uuid.UUID
	requestUser      models.User
}

func (suite *HandlerSuite) makeCreatePaymentRequestHandlerSubtestData() (subtestData *createPaymentRequestHandlerSubtestData) {
	subtestData = &createPaymentRequestHandlerSubtestData{}
	subtestData.moveTaskOrderID, _ = uuid.FromString("96e21765-3e29-4acf-89a2-1317a9f7f0da")
	subtestData.paymentRequestID, _ = uuid.FromString("70c0c9c1-cf3f-4195-b15c-d185dc5cd0bf")

	subtestData.requestUser = testdatagen.MakeStubbedUser(suite.DB())

	subtestData.serviceItemID1, _ = uuid.FromString("1b7b134a-7c44-45f2-9114-bb0831cc5db3")
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code:     models.ReServiceCodeDLH,
			Priority: 1,
		},
		MTOServiceItem: models.MTOServiceItem{
			ID: subtestData.serviceItemID1,
		},
	})
	subtestData.serviceItemID2, _ = uuid.FromString("119f0a05-34d7-4d86-9745-009c0707b4c2")
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code:     models.ReServiceCodeFSC,
			Priority: 99,
		},
		MTOServiceItem: models.MTOServiceItem{
			ID: subtestData.serviceItemID2,
		},
	})

	subtestData.serviceItemID3, _ = uuid.FromString("d01a9002-7ce5-4c07-9187-c00de15293ed")
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code:     models.ReServiceCodeDOASIT,
			Priority: 99,
		},
		MTOServiceItem: models.MTOServiceItem{
			ID: subtestData.serviceItemID3,
		},
	})
	return subtestData
}

func (suite *HandlerSuite) TestCreatePaymentRequestHandler() {
	suite.Run("successful create payment request", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   subtestData.paymentRequestID,
			MoveTaskOrderID:      subtestData.moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(paymentRequest *models.PaymentRequest) bool {
				// Making sure the service items are ordered by priority regardless of the order in which they come in through the payment request parameters
				return paymentRequest.PaymentServiceItems[0].MTOServiceItemID == subtestData.serviceItemID1
			})).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID2),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID3),
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

	suite.Run("create payment request without adding service item params passed into payload", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   subtestData.paymentRequestID,
			MoveTaskOrderID:      subtestData.moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					ID: subtestData.serviceItemID1,
				},
			},
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
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

	suite.Run("successfully create payment request with service item params passed into payload", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()
		returnedPaymentRequest := models.PaymentRequest{
			ID:                   subtestData.paymentRequestID,
			MoveTaskOrderID:      subtestData.moveTaskOrderID,
			PaymentRequestNumber: "1234-5678-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					ID: subtestData.serviceItemID3,
				},
			},
		}

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&returnedPaymentRequest, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID3),
						Params: []*primemessages.ServiceItemParamsItems0{
							{
								Key:   string(models.ServiceItemParamNameSITPaymentRequestStart),
								Value: "2021-08-05",
							},
						},
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
	})

	suite.Run("failed create payment request -- nil body", func() {
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)
	})

	suite.Run("failed create payment request -- creator failed with error", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, errors.New("creator failed")).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				PointOfContact:  "user@prime.com",
			},
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreatePaymentRequestInternalServerError{}, response)

		errResponse := response.(*paymentrequestop.CreatePaymentRequestInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong") // check body (body was written before panic)

	})

	suite.Run("failed create payment request -- invalid MTO ID format", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

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

	suite.Run("failed create payment request -- invalid service item ID format", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(&models.PaymentRequest{}, nil).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")
		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
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

	suite.Run("failed create payment request - validation errors", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		verrs := &validate.Errors{
			Errors: map[string][]string{
				"violation": {"invalid value"},
			},
		}
		err := apperror.NewInvalidCreateInputError(verrs, "can't create payment request for MTO ID 1234")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
	})

	suite.Run("failed create payment request due to conflict in model", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		ordersID, _ := uuid.FromString("2b8b141a-7c44-45f2-9114-bb0831cc5db3")
		err := apperror.NewConflictError(ordersID, "incomplete orders")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestConflict{}, response)
	})

	suite.Run("failed create payment request due to bad data", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		err := apperror.NewBadDataError("sent some bad data, foo!")
		paymentRequestCreator := &mocks.PaymentRequestCreator{}
		paymentRequestCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest")).Return(nil, err).Once()

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID2),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestBadRequest{}, response)
	})

	suite.Run("successful create payment request payload audit", func() {
		subtestData := suite.makeCreatePaymentRequestHandlerSubtestData()

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, subtestData.requestUser)

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(subtestData.moveTaskOrderID),
				PointOfContact:  "user@prime.com",
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID1),
					},
					{
						ID: *handlers.FmtUUID(subtestData.serviceItemID2),
					},
				},
			},
		}

		session := auth.Session{}
		appCtx := suite.AppContextWithSessionForTest(&session)

		// Capture creation attempt in audit log
		zapFields, err := audit.Capture(appCtx, &params.Body, nil, params.HTTPRequest)

		var eventType string
		for _, field := range zapFields {
			if field.Key == "event_type" {
				eventType = field.String
			}
		}
		suite.Nil(err, "No error for audit.Capture call")
		if assert.NotEmpty(suite.T(), zapFields) {
			assert.Equal(suite.T(), "event_type", zapFields[0].Key)
			assert.Equal(suite.T(), "audit_post_payment_requests", eventType)
		}
	})
}

func (suite *HandlerSuite) setupDomesticLinehaulData() (models.Move, models.MTOServiceItems) {
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "7 Q St",
			City:           "Birmingham",
			State:          "AL",
			PostalCode:     "90210",
		},
	})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "148 S East St",
			City:           "Miami",
			State:          "FL",
			PostalCode:     "94535",
		},
	})
	testEstWeight := dlhTestWeight
	testActualWeight := testEstWeight

	contractYear, serviceArea, _, _ := testdatagen.SetupServiceAreaRateArea(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			Escalation:           1.0197,
			EscalationCompounded: 1.04071,
		},
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea: dlhTestServiceArea,
		},
		ReRateArea: models.ReRateArea{
			Name: "Alabama",
		},
		ReZip3: models.ReZip3{
			Zip3:          pickupAddress.PostalCode[0:3],
			BasePointCity: pickupAddress.City,
			State:         pickupAddress.State,
		},
	})

	baseLinehaulPrice := testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          false,
		},
	})

	_ = testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          true,
			PriceMillicents:       baseLinehaulPrice.PriceMillicents - 2500, // minus $0.025
		},
	})

	csService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeCS,
		},
	})
	csTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      csService.ID,
		PriceCents:     unit.Cents(22399),
	}
	suite.MustSave(&csTaskOrderFee)

	msService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeMS,
		},
	})
	msTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      msService.ID,
		PriceCents:     unit.Cents(25513),
	}
	suite.MustSave(&msTaskOrderFee)

	availableToPrimeAt := time.Date(testdatagen.GHCTestYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	moveTaskOrder, mtoServiceItems := testdatagen.MakeFullDLHMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &availableToPrimeAt,
		},
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &testEstWeight,
			PrimeActualWeight:    &testActualWeight,
			PickupAddressID:      &pickupAddress.ID,
			PickupAddress:        &pickupAddress,
			DestinationAddressID: &destinationAddress.ID,
			DestinationAddress:   &destinationAddress,
		},
	})

	publicationDate := moveTaskOrder.MTOShipments[0].ActualPickupDate.AddDate(0, 0, -3) // 3 days earlier
	ghcDieselFuelPrice := models.GHCDieselFuelPrice{
		PublicationDate:       publicationDate,
		FuelPriceInMillicents: unit.Millicents(277600),
	}
	suite.MustSave(&ghcDieselFuelPrice)

	return moveTaskOrder, mtoServiceItems
}

func (suite *HandlerSuite) TestCreatePaymentRequestHandlerNewPaymentRequestCreator() {
	const defaultZipDistance = 1234

	suite.Run("successfully create payment request with real PaymentRequestCreator", func() {

		move, mtoServiceItems := suite.setupDomesticLinehaulData()
		moveTaskOrderID := move.ID

		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		planner := &routemocks.Planner{}
		planner.On("Zip5TransitDistanceLineHaul",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(defaultZipDistance, nil)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(defaultZipDistance, nil)

		paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
			planner,
			ghcrateengine.NewServiceItemPricer(),
		)

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(mtoServiceItems[0].ID),
					},
					{
						ID: *handlers.FmtUUID(mtoServiceItems[1].ID),
					},
					{
						ID: *handlers.FmtUUID(mtoServiceItems[2].ID),
					},
					{
						ID: *handlers.FmtUUID(mtoServiceItems[3].ID),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestCreated{}, response)
		typedResponse := response.(*paymentrequestop.CreatePaymentRequestCreated)
		suite.NotEmpty(typedResponse.Payload.ID.String(), "valid payload ID string")
		suite.NotEmpty(typedResponse.Payload.MoveTaskOrderID.String(), "valid MTO ID")
		suite.NotEmpty(typedResponse.Payload.PaymentRequestNumber, "valid Payment Request Number")
	})
}

func (suite *HandlerSuite) TestCreatePaymentRequestHandlerInvalidMTOReferenceID() {
	const defaultZipDistance = 48

	suite.Run("fail to create payment request with real PaymentRequestCreator and empty MTO Reference ID", func() {

		move, mtoServiceItems := suite.setupDomesticLinehaulData()
		moveTaskOrderID := move.ID

		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		planner := &routemocks.Planner{}
		planner.On("Zip5TransitDistanceLineHaul",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(defaultZipDistance, nil)
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(defaultZipDistance, nil)

		paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
			planner,
			ghcrateengine.NewServiceItemPricer(),
		)

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(mtoServiceItems[0].ID),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}

		// Set Reference ID to an empty string
		*move.ReferenceID = ""
		suite.MustSave(&move)

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
		typedResponse := response.(*paymentrequestop.CreatePaymentRequestUnprocessableEntity)
		suite.Contains(*typedResponse.Payload.Detail, "has missing ReferenceID")
	})

	suite.Run("fail to create payment request with real PaymentRequestCreator and nil MTO Reference ID", func() {

		move, mtoServiceItems := suite.setupDomesticLinehaulData()
		moveTaskOrderID := move.ID

		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		req := httptest.NewRequest("POST", "/payment_requests", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		planner := &routemocks.Planner{}
		planner.On("Zip5TransitDistanceLineHaul",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(defaultZipDistance, nil)
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(defaultZipDistance, nil)

		paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
			planner,
			ghcrateengine.NewServiceItemPricer(),
		)

		handler := CreatePaymentRequestHandler{
			suite.HandlerConfig(),
			paymentRequestCreator,
		}

		params := paymentrequestop.CreatePaymentRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreatePaymentRequest{
				IsFinal:         swag.Bool(false),
				MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
				ServiceItems: []*primemessages.ServiceItem{
					{
						ID: *handlers.FmtUUID(mtoServiceItems[0].ID),
					},
				},
				PointOfContact: "user@prime.com",
			},
		}

		// Set Reference ID to a nil string
		move.ReferenceID = nil
		suite.MustSave(&move)

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.CreatePaymentRequestUnprocessableEntity{}, response)
		typedResponse := response.(*paymentrequestop.CreatePaymentRequestUnprocessableEntity)
		suite.Contains(*typedResponse.Payload.Detail, "has missing ReferenceID")
	})
}
