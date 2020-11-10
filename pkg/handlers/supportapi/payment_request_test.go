package supportapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_request"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(suite.DB()),
		}

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
	})

	suite.T().Run("successful status update of prime-available payment request", func(t *testing.T) {
		availableMove := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		availablePaymentRequestID := availablePaymentRequest.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		eTag := etag.GenerateEtag(availablePaymentRequest.UpdatedAt)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil, ETag: eTag},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
			IfMatch:          eTag,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(suite.DB()),
		}
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		handler.SetTraceID(traceID)

		response := handler.Handle(params)

		paymentRequestStatusResponse := response.(*paymentrequestop.UpdatePaymentRequestStatusOK)
		paymentRequestStatusPayload := paymentRequestStatusResponse.Payload

		suite.Equal(paymentRequestStatusPayload.Status, params.Body.Status)
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		suite.HasWebhookNotification(availablePaymentRequestID, traceID)
	})

	suite.T().Run("unsuccessful status update of payment request (500)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
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

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
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
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	mto := testdatagen.MakeDefaultMove(suite.DB())
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

func (suite *HandlerSuite) TestGetPaymentRequestEDIHandler() {
	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   time.Now().Format("20060102"),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "2424",
		},
	}
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDLH,
		basicPaymentServiceItemParams,
	)

	// Add a price to the service item.
	priceCents := unit.Cents(250000)
	paymentServiceItem.PriceCents = &priceCents
	suite.MustSave(&paymentServiceItem)

	paymentRequestID := paymentServiceItem.PaymentRequestID
	strfmtPaymentRequestID := strfmt.UUID(paymentRequestID.String())

	handler := GetPaymentRequestEDIHandler{
		HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(suite.DB()),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB()),
	}

	urlFormat := "/payment-requests/%s/edi"

	suite.T().Run("successful get of EDI for payment request", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequestID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIOK(), response)

		ediResponse := response.(*paymentrequestop.GetPaymentRequestEDIOK)
		ediPayload := ediResponse.Payload

		suite.Equal(ediPayload.ID, strfmtPaymentRequestID)

		// Check to make sure EDI is there and starts with expected segment.
		edi := ediPayload.Edi
		if suite.NotEmpty(edi) {
			suite.Regexp("^ISA*", edi)
		}
	})

	suite.T().Run("failure due to incorrectly formatted payment request ID", func(t *testing.T) {
		invalidID := "12345"
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, invalidID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(invalidID),
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDINotFound(), response)
	})

	suite.T().Run("failure due to a validation error", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequestID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		mockGenerator.On("Generate", mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, services.NewInvalidInputError(paymentRequestID, nil, validate.NewErrors(), ""))

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(suite.DB()),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIUnprocessableEntity(), response)
	})

	suite.T().Run("failure due to payment request ID not found", func(t *testing.T) {
		notFoundID := uuid.Must(uuid.NewV4())
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, notFoundID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(notFoundID.String()),
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDINotFound(), response)
	})

	suite.T().Run("failure when generating EDI", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequestID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		errStr := "some error"
		mockGenerator.On("Generate", mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New(errStr)).Once()

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(suite.DB()),
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
	for i := 0; i < 4; i++ {
		currentTime := time.Now()
		basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTime.Format(testDateFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
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
	reviewedPRs := suite.createPaymentRequest(4)

	sentToGEXTime := time.Now()
	paymentRequestForUpdate := models.PaymentRequest{
		ID:                   reviewedPRs[0].ID,
		MoveTaskOrder:        reviewedPRs[0].MoveTaskOrder,
		MoveTaskOrderID:      reviewedPRs[0].MoveTaskOrderID,
		IsFinal:              reviewedPRs[0].IsFinal,
		Status:               models.PaymentRequestStatusSentToGex,
		RejectionReason:      reviewedPRs[0].RejectionReason,
		SentToGexAt:          &sentToGEXTime,
		PaymentRequestNumber: reviewedPRs[0].PaymentRequestNumber,
		SequenceNumber:       reviewedPRs[0].SequenceNumber,
	}
	paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
	paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.Anything, mock.Anything).Return(reviewedPRs, nil)

	paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
	paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

	paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
	paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(reviewedPRs[0], nil)

	paymentRequestProcessor := &mocks.PaymentRequestReviewedProcessor{}
	paymentRequestProcessor.On("ProcessReviewedPaymentRequest").Return(nil)

	handler := ProcessReviewedPaymentRequestsHandler{
		HandlerContext:                  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		PaymentRequestFetcher:           paymentRequestFetcher,
		PaymentRequestStatusUpdater:     paymentRequestStatusUpdater,
		PaymentRequestReviewedFetcher:   paymentRequestReviewedFetcher,
		PaymentRequestReviewedProcessor: paymentRequestProcessor,
	}

	urlFormat := "/payment-requests/process-reviewed"

	suite.T().Run("successful update of reviewed payment requests with send to syncada true", func(t *testing.T) {
		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)

		sendToSyncada := true
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada: &sendToSyncada,
				Status:        "SENT_TO_GEX",
			},
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.T().Run("successful update of reviewed payment requests with send to syncada false", func(t *testing.T) {
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest()
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)

		sendToSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada: &sendToSyncada,
				Status:        "SENT_TO_GEX",
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

	suite.T().Run("successful update of reviewed payment requests with send to syncada false, when no status flag is set", func(t *testing.T) {
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest()
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)

		sendToSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body:        &supportmessages.ProcessReviewedPaymentRequests{SendToSyncada: &sendToSyncada},
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

	suite.T().Run("successful update of a given reviewed payment request", func(t *testing.T) {
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest()

		paymentRequestID := reviewedPaymentRequests[0].ID
		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)

		sendToSyncada := false
		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:    &sendToSyncada,
				PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
			},
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ProcessReviewedPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
		// Ensure that the previously reviewed status have been updated to match the default status of SENT_TO_GEX
		suite.Equal(supportmessages.PaymentRequestStatusSENTTOGEX, paymentRequestsPayload[0].Status)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.T().Run("fail if required send to syncada flag is not set", func(t *testing.T) {
		// Ensure that there are reviewed payment requests
		reviewedPaymentRequests, _ := handler.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest()
		suite.Equal(4, len(reviewedPaymentRequests))

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body:        &supportmessages.ProcessReviewedPaymentRequests{},
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest(), response)
	})

	suite.T().Run("fail if paymentRequestId is supplied but not found", func(t *testing.T) {
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.Anything, mock.Anything).Return(reviewedPRs, nil)

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		var nilPr models.PaymentRequest
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(nilPr, errors.New("could not fetch payment request"))

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerContext:                  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestFetcher:           paymentRequestFetcher,
			PaymentRequestStatusUpdater:     paymentRequestStatusUpdater,
			PaymentRequestReviewedFetcher:   paymentRequestReviewedFetcher,
			PaymentRequestReviewedProcessor: paymentrequest.InitNewPaymentRequestReviewedProcessor(suite.DB(), suite.TestLogger(), true),
		}

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)
		prID := reviewedPRs[0].ID
		sendToSyncada := false

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada:    &sendToSyncada,
				PaymentRequestID: strfmt.UUID(prID.String()),
			},
		}

		response := handler.Handle(params)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsNotFound(), response)
	})

	suite.T().Run("fail if reviewed payment request are not retrieved", func(t *testing.T) {
		var nilReviewedPrs models.PaymentRequests
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.Anything, mock.Anything).Return(nilReviewedPrs, errors.New("Reviewed Payment Requests notretrieved"))

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.Anything).Return(reviewedPRs[0], nil)

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerContext:                  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			PaymentRequestFetcher:           paymentRequestFetcher,
			PaymentRequestStatusUpdater:     paymentRequestStatusUpdater,
			PaymentRequestReviewedFetcher:   paymentRequestReviewedFetcher,
			PaymentRequestReviewedProcessor: paymentrequest.InitNewPaymentRequestReviewedProcessor(suite.DB(), suite.TestLogger(), true),
		}

		// Call the handler to update all reviewed payment request to a "Sent_To_Gex" status (default status when no flag is set)
		req := httptest.NewRequest("PATCH", fmt.Sprintf(urlFormat), nil)
		sendToSyncada := false

		params := paymentrequestop.ProcessReviewedPaymentRequestsParams{
			HTTPRequest: req,
			Body: &supportmessages.ProcessReviewedPaymentRequests{
				SendToSyncada: &sendToSyncada,
			},
		}

		response := handler.Handle(params)
		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError(), response)
	})
}
