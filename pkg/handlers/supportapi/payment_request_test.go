//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used set up environment variables
//RA: Given the functions causing the lint errors are used to set environment variables for testing purposes, it does not present a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package supportapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/etag"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_request"
	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	suite.Run("successful status update of payment request", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequest.ID), nil)
		eTag := etag.GenerateEtag(paymentRequest.UpdatedAt)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil, ETag: eTag},
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
			IfMatch:          eTag,
		}
		queryBuilder := query.NewQueryBuilder()

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
		}

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
	})

	suite.Run("successful status update of prime-available payment request", func() {
		availableMove := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		availablePaymentRequestID := availablePaymentRequest.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		eTag := etag.GenerateEtag(availablePaymentRequest.UpdatedAt)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil, ETag: eTag},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
			IfMatch:          eTag,
		}
		queryBuilder := query.NewQueryBuilder()

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
		}

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		suite.HasWebhookNotification(availablePaymentRequestID, traceID)
	})

	suite.Run("unsuccessful status update of payment request (500)", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.UpdatePaymentRequestStatusInternalServerError{}, response)

		errResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})

	suite.Run("unsuccessful status update of payment request, not found (404)", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusNotFound(), response)

	})

	suite.Run("unsuccessful status update of payment request, precondition failed (412)", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.PreconditionFailedError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed(), response)

	})
	suite.Run("unsuccessful status update of payment request, conflict error (409)", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.ConflictError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerConfig:               handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusConflict(), response)
	})
}

func (suite *HandlerSuite) TestListMTOPaymentRequestHandler() {
	suite.Run("successful get an MTO with payment requests", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		mtoID := paymentRequest.MoveTaskOrderID
		req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/payment-requests", mtoID), nil)

		params := paymentrequestop.ListMTOPaymentRequestsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: strfmt.UUID(mtoID.String()),
		}

		handler := ListMTOPaymentRequestsHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ListMTOPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		suite.IsType(paymentrequestop.NewListMTOPaymentRequestsOK(), response)
		suite.Equal(1, len(paymentRequestsPayload))
	})

	suite.Run("successful get an MTO with no payment requests", func() {
		mto := testdatagen.MakeDefaultMove(suite.DB())
		req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/payment-requests", mto.ID), nil)

		params := paymentrequestop.ListMTOPaymentRequestsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: strfmt.UUID(mto.ID.String()),
		}

		handler := ListMTOPaymentRequestsHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ListMTOPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		suite.IsType(paymentrequestop.NewListMTOPaymentRequestsOK(), response)
		suite.Equal(0, len(paymentRequestsPayload))
	})
}

func (suite *HandlerSuite) TestGetPaymentRequestEDIHandler() {
	setupTestData := func() models.PaymentRequest {
		currentTimeStr := time.Now().Format("20060102")
		basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTimeStr,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTimeStr,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "4242",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2424",
			},
		}
		paymentServiceItem := testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
		)

		// Add a price to the service item.
		priceCents := unit.Cents(250000)
		paymentServiceItem.PriceCents = &priceCents
		suite.MustSave(&paymentServiceItem)

		return paymentServiceItem.PaymentRequest
	}

	setupHandler := func() GetPaymentRequestEDIHandler {
		icnSequencer := sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)
		return GetPaymentRequestEDIHandler{
			HandlerConfig:                     handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock()),
		}
	}

	urlFormat := "/payment-requests/%s/edi"

	suite.Run("successful get of EDI for payment request", func() {
		paymentRequest := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequest.ID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		response := setupHandler().Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIOK(), response)

		ediResponse := response.(*paymentrequestop.GetPaymentRequestEDIOK)
		ediPayload := ediResponse.Payload

		suite.Equal(ediPayload.ID, strfmt.UUID(paymentRequest.ID.String()))

		// Check to make sure EDI is there and starts with expected segment.
		edi := ediPayload.Edi
		if suite.NotEmpty(edi) {
			suite.Regexp("^ISA*", edi)
		}
	})

	suite.Run("failure due to incorrectly formatted payment request ID", func() {
		invalidID := "12345"
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, invalidID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(invalidID),
		}

		response := setupHandler().Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDINotFound(), response)
	})

	suite.Run("failure due to a validation error", func() {
		paymentRequest := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequest.ID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, apperror.NewInvalidInputError(paymentRequest.ID, nil, validate.NewErrors(), ""))

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerConfig:                     handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIUnprocessableEntity(), response)
	})

	suite.Run("failure due to a conflict error", func() {
		paymentRequest := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequest.ID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, apperror.NewConflictError(paymentRequest.ID, "conflict error"))

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerConfig:                     handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIConflict(), response)
	})

	suite.Run("failure due to payment request ID not found", func() {
		notFoundID := uuid.Must(uuid.NewV4())
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, notFoundID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(notFoundID.String()),
		}

		response := setupHandler().Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDINotFound(), response)
	})

	suite.Run("failure when generating EDI", func() {
		paymentRequest := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequest.ID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		errStr := "some error"
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New(errStr)).Once()

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerConfig:                     handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIInternalServerError(), response)
		errResponse := response.(*paymentrequestop.GetPaymentRequestEDIInternalServerError)
		suite.Contains(*errResponse.Payload.Detail, errStr)
	})
}

const testDateFormat = "060102"

func (suite *HandlerSuite) createPaymentRequest(num int) models.PaymentRequests {
	var prs models.PaymentRequests
	for i := 0; i < num; i++ {
		currentTimeStr := time.Now().Format(testDateFormat)
		basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTimeStr,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTimeStr,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "4242",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2424",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "24245",
			},
		}

		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: mto,
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		})
		prs = append(prs, paymentRequest)
		requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
		actualPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 22, 0, 0, 0, 0, time.UTC)

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
				ScheduledPickupDate: &scheduledPickupDate,
				ActualPickupDate:    &actualPickupDate,
			},
		})

		assertions := testdatagen.Assertions{
			Move:           mto,
			MTOShipment:    mtoShipment,
			PaymentRequest: paymentRequest,
		}

		// dlh
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			assertions,
		)
	}
	return prs
}

func (suite *HandlerSuite) TestProcessReviewedPaymentRequestsHandler() {

	err := os.Setenv("GEX_SFTP_PORT", "1234")
	suite.FatalNoError(err)
	err = os.Setenv("GEX_SFTP_USER_ID", "FAKE_USER_ID")
	suite.FatalNoError(err)
	err = os.Setenv("GEX_SFTP_IP_ADDRESS", "127.0.0.1")
	suite.FatalNoError(err)
	err = os.Setenv("GEX_SFTP_PASSWORD", "FAKE PASSWORD")
	suite.FatalNoError(err)
	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	err = os.Setenv("GEX_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")
	suite.FatalNoError(err)

	setupTestData := func() (models.PaymentRequest, models.PaymentRequests) {
		reviewedPRs := suite.createPaymentRequest(4)

		sentToGEXTime := time.Now()
		return models.PaymentRequest{
			ID:                              reviewedPRs[0].ID,
			MoveTaskOrder:                   reviewedPRs[0].MoveTaskOrder,
			MoveTaskOrderID:                 reviewedPRs[0].MoveTaskOrderID,
			IsFinal:                         reviewedPRs[0].IsFinal,
			Status:                          models.PaymentRequestStatusSentToGex,
			RejectionReason:                 reviewedPRs[0].RejectionReason,
			SentToGexAt:                     &sentToGEXTime,
			PaymentRequestNumber:            reviewedPRs[0].PaymentRequestNumber,
			SequenceNumber:                  reviewedPRs[0].SequenceNumber,
			RecalculationOfPaymentRequestID: reviewedPRs[0].RecalculationOfPaymentRequestID,
		}, reviewedPRs
	}

	setupHandler := func(paymentRequestForUpdate models.PaymentRequest, reviewedPRs models.PaymentRequests) ProcessReviewedPaymentRequestsHandler {
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(reviewedPRs, nil)

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(reviewedPRs[0], nil)
		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:         paymentRequestFetcher,
			PaymentRequestStatusUpdater:   paymentRequestStatusUpdater,
			PaymentRequestReviewedFetcher: paymentRequestReviewedFetcher,
		}

		handler.SetICNSequencer(sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName))

		return handler
	}

	urlFormat := "/payment-requests/process-reviewed"

	suite.Run("successful update of reviewed payment requests with send to syncada true", func() {
		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)

		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:     &sendToSyncada,
				ReadFromSyncada:   &readFromSyncada,
				DeleteFromSyncada: &deleteFromSyncada,
				Status:            "SENT_TO_GEX",
			},
		}

		response := setupHandler(setupTestData()).Handle(params)

		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.Run("successful update of reviewed payment requests with send to syncada false", func() {
		handler := setupHandler(setupTestData())
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)

		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:     &sendToSyncada,
				ReadFromSyncada:   &readFromSyncada,
				DeleteFromSyncada: &deleteFromSyncada,
				Status:            "SENT_TO_GEX",
			},
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ProcessReviewedPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
		// Ensure that the previously reviewed status have been updated to match the status flag
		for _, pr := range paymentRequestsPayload {
			suite.Equal(params.Body.Status, pr.Status)
		}
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.Run("successful update of reviewed payment requests with send to syncada false, when no status flag is set", func() {
		handler := setupHandler(setupTestData())
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)

		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body:        &supportmessages.ProcessReviewedPaymentRequests{ReadFromSyncada: &readFromSyncada, SendToSyncada: &sendToSyncada, DeleteFromSyncada: &deleteFromSyncada},
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ProcessReviewedPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
		// Ensure that the previously reviewed status have been updated to match the default status of SENT_TO_GEX
		for _, pr := range paymentRequestsPayload {
			suite.Equal(supportmessages.PaymentRequestStatusSENTTOGEX, pr.Status)
		}
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.Run("successful update of a given reviewed payment request", func() {
		handler := setupHandler(setupTestData())
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())

		paymentRequestID := reviewedPaymentRequests[0].ID
		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)

		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:     &sendToSyncada,
				ReadFromSyncada:   &readFromSyncada,
				DeleteFromSyncada: &deleteFromSyncada,
				PaymentRequestID:  strfmt.UUID(paymentRequestID.String()),
			},
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ProcessReviewedPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
		// Ensure that the previously reviewed status have been updated to match the default status of SENT_TO_GEX
		suite.Equal(supportmessages.PaymentRequestStatusSENTTOGEX, paymentRequestsPayload[0].Status)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.Run("fail if required send to syncada flag is not set", func() {
		handler := setupHandler(setupTestData())
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body:        &supportmessages.ProcessReviewedPaymentRequests{},
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest(), response)
	})

	suite.Run("fail if paymentRequestId is supplied but not found", func() {
		paymentRequestForUpdate, reviewedPRs := setupTestData()
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(reviewedPRs, nil)

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		var nilPr models.PaymentRequest
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(nilPr, errors.New("could not fetch payment request"))

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:         paymentRequestFetcher,
			PaymentRequestStatusUpdater:   paymentRequestStatusUpdater,
			PaymentRequestReviewedFetcher: paymentRequestReviewedFetcher,
		}

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)
		prID := reviewedPRs[0].ID
		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:     &sendToSyncada,
				ReadFromSyncada:   &readFromSyncada,
				DeleteFromSyncada: &deleteFromSyncada,
				PaymentRequestID:  strfmt.UUID(prID.String()),
			},
		}

		response := handler.Handle(params)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsNotFound(), response)
	})

	suite.Run("fail if reviewed payment request are not retrieved", func() {
		paymentRequestForUpdate, reviewedPRs := setupTestData()
		var nilReviewedPrs models.PaymentRequests
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nilReviewedPrs, errors.New("Reviewed Payment Requests notretrieved"))

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(reviewedPRs[0], nil)

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:         paymentRequestFetcher,
			PaymentRequestStatusUpdater:   paymentRequestStatusUpdater,
			PaymentRequestReviewedFetcher: paymentRequestReviewedFetcher,
		}

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprint(urlFormat), nil)
		sendToSyncada := false
		readFromSyncada := false
		deleteFromSyncada := false

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:     &sendToSyncada,
				ReadFromSyncada:   &readFromSyncada,
				DeleteFromSyncada: &deleteFromSyncada,
			},
		}

		response := handler.Handle(params)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError(), response)
	})
}

func (suite *HandlerSuite) TestRecalculatePaymentRequestHandler() {
	paymentRequestID := uuid.Must(uuid.NewV4())
	strfmtPaymentRequestID := strfmt.UUID(paymentRequestID.String())

	method := "POST"
	urlFormat := "/payment-requests/%s/recalculate"

	suite.Run("golden path", func() {
		samplePaymentRequest := models.PaymentRequest{
			ID:                              uuid.Must(uuid.NewV4()),
			MoveTaskOrderID:                 uuid.Must(uuid.NewV4()),
			Status:                          models.PaymentRequestStatusPending,
			PaymentRequestNumber:            "1111-2222-1",
			SequenceNumber:                  1,
			RecalculationOfPaymentRequestID: &paymentRequestID,
		}

		mockRecalculator := &mocks.PaymentRequestRecalculator{}
		mockRecalculator.On("RecalculatePaymentRequest",
			mock.AnythingOfType("*appcontext.appContext"),
			paymentRequestID,
		).Return(&samplePaymentRequest, nil).Once()
		handler := RecalculatePaymentRequestHandler{
			HandlerConfig:              handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentRequestRecalculator: mockRecalculator,
		}

		req := httptest.NewRequest(method, fmt.Sprintf(urlFormat, paymentRequestID), nil)
		params := paymentrequestop.RecalculatePaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}
		response := handler.Handle(params)

		mockRecalculator.AssertExpectations(suite.T())

		if suite.IsType(paymentrequestop.NewRecalculatePaymentRequestCreated(), response) {
			paymentRequestResponse := response.(*paymentrequestop.RecalculatePaymentRequestCreated)
			payload := paymentRequestResponse.Payload
			suite.Equal(samplePaymentRequest.ID.String(), payload.ID.String())
			suite.Equal(samplePaymentRequest.MoveTaskOrderID.String(), payload.MoveTaskOrderID.String())
			suite.Equal(samplePaymentRequest.Status.String(), string(payload.Status))
			suite.Equal(samplePaymentRequest.PaymentRequestNumber, payload.PaymentRequestNumber)
			// SequenceNumber is not on payload at all as it's an internal representation.
			if suite.NotNil(payload.RecalculationOfPaymentRequestID) {
				suite.Equal(samplePaymentRequest.RecalculationOfPaymentRequestID.String(), payload.RecalculationOfPaymentRequestID.String())
			}
		}
	})

	errorTestCases := []struct {
		testErr      error
		responseType interface{}
	}{
		{
			apperror.NewBadDataError("test"),
			paymentrequestop.NewRecalculatePaymentRequestBadRequest(),
		},
		{
			apperror.NewNotFoundError(paymentRequestID, "test"),
			paymentrequestop.NewRecalculatePaymentRequestNotFound(),
		},
		{
			apperror.NewConflictError(paymentRequestID, "test"),
			paymentrequestop.NewRecalculatePaymentRequestConflict(),
		},
		{
			apperror.NewPreconditionFailedError(paymentRequestID, errors.New("test")),
			paymentrequestop.NewRecalculatePaymentRequestPreconditionFailed(),
		},
		{
			apperror.NewInvalidInputError(paymentRequestID, errors.New("test"), validate.NewErrors(), "test"),
			paymentrequestop.NewRecalculatePaymentRequestUnprocessableEntity(),
		},
		{
			apperror.NewInvalidCreateInputError(validate.NewErrors(), "test"),
			paymentrequestop.NewRecalculatePaymentRequestUnprocessableEntity(),
		},
		{
			apperror.NewQueryError("TestObject", errors.New("test"), "test"),
			paymentrequestop.NewRecalculatePaymentRequestInternalServerError(),
		},
		{
			errors.New("test"),
			paymentrequestop.NewRecalculatePaymentRequestInternalServerError(),
		},
	}

	for _, testCase := range errorTestCases {
		testName := fmt.Sprintf("%T error from service should produce %T response type", testCase.testErr, testCase.responseType)
		suite.Run(testName, func() {
			mockRecalculator := &mocks.PaymentRequestRecalculator{}
			mockRecalculator.On("RecalculatePaymentRequest",
				mock.AnythingOfType("*appcontext.appContext"),
				paymentRequestID,
			).Return(nil, testCase.testErr)
			handler := RecalculatePaymentRequestHandler{
				HandlerConfig:              handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
				PaymentRequestRecalculator: mockRecalculator,
			}

			req := httptest.NewRequest(method, fmt.Sprintf(urlFormat, paymentRequestID), nil)
			params := paymentrequestop.RecalculatePaymentRequestParams{
				HTTPRequest:      req,
				PaymentRequestID: strfmtPaymentRequestID,
			}
			response := handler.Handle(params)

			mockRecalculator.AssertExpectations(suite.T())

			suite.IsType(testCase.responseType, response)
		})
	}
}
