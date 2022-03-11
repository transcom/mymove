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
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/trace"

	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"

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
		queryBuilder := query.NewQueryBuilder()

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
			PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
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
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.T().Run("unsuccessful status update of payment request (500)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusNotFound(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, precondition failed (412)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.PreconditionFailedError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed(), response)

	})
	suite.T().Run("unsuccessful status update of payment request, conflict error (409)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nil, apperror.ConflictError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(paymentRequest, nil).Once()

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &supportmessages.UpdatePaymentRequestStatus{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusConflict(), response)

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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
		}

		response := handler.Handle(params)

		paymentRequestsResponse := response.(*paymentrequestop.ListMTOPaymentRequestsOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload

		suite.IsType(paymentrequestop.NewListMTOPaymentRequestsOK(), response)
		suite.Equal(len(paymentRequestsPayload), 0)
	})
}

func (suite *HandlerSuite) TestGetPaymentRequestEDIHandler() {
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

	paymentRequestID := paymentServiceItem.PaymentRequestID
	strfmtPaymentRequestID := strfmt.UUID(paymentRequestID.String())

	icnSequencer := sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)
	handler := GetPaymentRequestEDIHandler{
		HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock()),
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
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, apperror.NewInvalidInputError(paymentRequestID, nil, validate.NewErrors(), ""))

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIUnprocessableEntity(), response)
	})

	suite.T().Run("failure due to a conflict error", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf(urlFormat, paymentRequestID), nil)

		params := paymentrequestop.GetPaymentRequestEDIParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}

		mockGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, apperror.NewConflictError(paymentRequestID, "conflict error"))

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
			GHCPaymentRequestInvoiceGenerator: mockGenerator,
		}

		response := mockGeneratorHandler.Handle(params)

		suite.IsType(paymentrequestop.NewGetPaymentRequestEDIConflict(), response)
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
		mockGenerator.On("Generate", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New(errStr)).Once()

		mockGeneratorHandler := GetPaymentRequestEDIHandler{
			HandlerContext:                    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	reviewedPRs := suite.createPaymentRequest(4)

	sentToGEXTime := time.Now()
	paymentRequestForUpdate := models.PaymentRequest{
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
	}
	paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
	paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(reviewedPRs, nil)

	paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
	paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

	paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
	paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(reviewedPRs[0], nil)

	handler := ProcessReviewedPaymentRequestsHandler{
		HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
		PaymentRequestFetcher:         paymentRequestFetcher,
		PaymentRequestStatusUpdater:   paymentRequestStatusUpdater,
		PaymentRequestReviewedFetcher: paymentRequestReviewedFetcher,
	}

	handler.SetICNSequencer(sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName))

	urlFormat := "/payment-requests/process-reviewed"

	suite.T().Run("successful update of reviewed payment requests with send to syncada true", func(t *testing.T) {
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

		suite.IsType(paymentrequestop.NewProcessReviewedPaymentRequestsOK(), response)
	})

	suite.T().Run("successful update of reviewed payment requests with send to syncada false", func(t *testing.T) {
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

	suite.T().Run("successful update of reviewed payment requests with send to syncada false, when no status flag is set", func(t *testing.T) {
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

	suite.T().Run("successful update of a given reviewed payment request", func(t *testing.T) {
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

	suite.T().Run("fail if required send to syncada flag is not set", func(t *testing.T) {
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

	suite.T().Run("fail if paymentRequestId is supplied but not found", func(t *testing.T) {
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(reviewedPRs, nil)

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		var nilPr models.PaymentRequest
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(nilPr, errors.New("could not fetch payment request"))

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.T().Run("fail if reviewed payment request are not retrieved", func(t *testing.T) {
		var nilReviewedPrs models.PaymentRequests
		paymentRequestReviewedFetcher := &mocks.PaymentRequestReviewedFetcher{}
		paymentRequestReviewedFetcher.On("FetchReviewedPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(nilReviewedPrs, errors.New("Reviewed Payment Requests notretrieved"))

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(&paymentRequestForUpdate, nil)

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(reviewedPRs[0], nil)

		handler := ProcessReviewedPaymentRequestsHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.T().Run("golden path", func(t *testing.T) {
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
			HandlerContext:             handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestRecalculator: mockRecalculator,
		}

		req := httptest.NewRequest(method, fmt.Sprintf(urlFormat, paymentRequestID), nil)
		params := paymentrequestop.RecalculatePaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmtPaymentRequestID,
		}
		response := handler.Handle(params)

		mockRecalculator.AssertExpectations(t)

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
		suite.T().Run(testName, func(t *testing.T) {
			mockRecalculator := &mocks.PaymentRequestRecalculator{}
			mockRecalculator.On("RecalculatePaymentRequest",
				mock.AnythingOfType("*appcontext.appContext"),
				paymentRequestID,
			).Return(nil, testCase.testErr)
			handler := RecalculatePaymentRequestHandler{
				HandlerContext:             handlers.NewHandlerContext(suite.DB(), suite.Logger()),
				PaymentRequestRecalculator: mockRecalculator,
			}

			req := httptest.NewRequest(method, fmt.Sprintf(urlFormat, paymentRequestID), nil)
			params := paymentrequestop.RecalculatePaymentRequestParams{
				HTTPRequest:      req,
				PaymentRequestID: strfmtPaymentRequestID,
			}
			response := handler.Handle(params)

			mockRecalculator.AssertExpectations(t)

			suite.IsType(testCase.responseType, response)
		})
	}
}
