package supportapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services/invoice"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

// UpdatePaymentRequestStatusHandler updates payment requests status
type UpdatePaymentRequestStatusHandler struct {
	handlers.HandlerContext
	services.PaymentRequestStatusUpdater
	services.PaymentRequestFetcher
}

// Handle updates payment requests status
func (h UpdatePaymentRequestStatusHandler) Handle(params paymentrequestop.UpdatePaymentRequestStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(paymentRequestID)

	if err != nil {
		msg := fmt.Sprintf("Error finding Payment Request for status update with ID: %s", params.PaymentRequestID.String())
		logger.Error(msg, zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
	}

	status := existingPaymentRequest.Status
	mtoID := existingPaymentRequest.MoveTaskOrderID

	var reviewedDate time.Time
	var recGexDate time.Time
	var sentGexDate time.Time
	var paidAtDate time.Time

	if existingPaymentRequest.ReviewedAt != nil {
		reviewedDate = *existingPaymentRequest.ReviewedAt
	}
	if existingPaymentRequest.ReceivedByGexAt != nil {
		recGexDate = *existingPaymentRequest.ReceivedByGexAt
	}
	if existingPaymentRequest.SentToGexAt != nil {
		sentGexDate = *existingPaymentRequest.SentToGexAt
	}
	if existingPaymentRequest.PaidAt != nil {
		paidAtDate = *existingPaymentRequest.PaidAt
	}

	// Let's map the incoming status to our enumeration type
	switch params.Body.Status {
	case "PENDING":
		status = models.PaymentRequestStatusPending
	case "REVIEWED":
		status = models.PaymentRequestStatusReviewed
		reviewedDate = time.Now()
	case "SENT_TO_GEX":
		status = models.PaymentRequestStatusSentToGex
		sentGexDate = time.Now()
	case "RECEIVED_BY_GEX":
		status = models.PaymentRequestStatusReceivedByGex
		recGexDate = time.Now()
	case "PAID":
		status = models.PaymentRequestStatusPaid
		paidAtDate = time.Now()
	}

	// If we got a rejection reason let's use it
	rejectionReason := existingPaymentRequest.RejectionReason
	if params.Body.RejectionReason != nil {
		rejectionReason = params.Body.RejectionReason
	}

	paymentRequestForUpdate := models.PaymentRequest{
		ID:                   existingPaymentRequest.ID,
		MoveTaskOrder:        existingPaymentRequest.MoveTaskOrder,
		MoveTaskOrderID:      existingPaymentRequest.MoveTaskOrderID,
		IsFinal:              existingPaymentRequest.IsFinal,
		Status:               status,
		RejectionReason:      rejectionReason,
		RequestedAt:          existingPaymentRequest.RequestedAt,
		ReviewedAt:           &reviewedDate,
		SentToGexAt:          &sentGexDate,
		ReceivedByGexAt:      &recGexDate,
		PaidAt:               &paidAtDate,
		PaymentRequestNumber: existingPaymentRequest.PaymentRequestNumber,
		SequenceNumber:       existingPaymentRequest.SequenceNumber,
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&paymentRequestForUpdate, params.IfMatch)

	if err != nil {
		switch err.(type) {
		case services.NotFoundError:
			return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.ConflictError:
			return paymentrequestop.NewUpdatePaymentRequestStatusConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		default:
			logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
			return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.PaymentRequestUpdateEventKey,
		MtoID:           mtoID,
		UpdatedObjectID: updatedPaymentRequest.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.SupportUpdatePaymentRequestStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("supportapi.UpdatePaymentRequestStatusHandler could not generate the event")
	}

	returnPayload := payloads.PaymentRequest(updatedPaymentRequest)
	return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload)
}

// ListMTOPaymentRequestsHandler gets all payment requests for a given MTO
type ListMTOPaymentRequestsHandler struct {
	handlers.HandlerContext
}

// Handle getting payment requests for a given MTO
func (h ListMTOPaymentRequestsHandler) Handle(params paymentrequestop.ListMTOPaymentRequestsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoID, err := uuid.FromString(params.MoveTaskOrderID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing move task order id: %s", params.MoveTaskOrderID.String()), zap.Error(err))
		return paymentrequestop.NewListMTOPaymentRequestsInternalServerError()
	}

	var paymentRequests models.PaymentRequests

	query := h.DB().Where("move_id = ?", mtoID)

	err = query.All(&paymentRequests)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return paymentrequestop.NewListMTOPaymentRequestsInternalServerError()
	}

	payload := payloads.PaymentRequests(&paymentRequests)

	return paymentrequestop.NewListMTOPaymentRequestsOK().WithPayload(*payload)
}

// GetPaymentRequestEDIHandler returns the EDI for a given payment request
type GetPaymentRequestEDIHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
	services.GHCPaymentRequestInvoiceGenerator
}

// Handle getting the EDI for a given payment request
func (h GetPaymentRequestEDIHandler) Handle(params paymentrequestop.GetPaymentRequestEDIParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	paymentRequestID := uuid.FromStringOrNil(params.PaymentRequestID.String())
	paymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(paymentRequestID)
	if err != nil {
		msg := fmt.Sprintf("Error finding Payment Request for EDI generation with ID: %s", params.PaymentRequestID.String())
		logger.Error(msg, zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestEDINotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
	}

	var payload supportmessages.PaymentRequestEDI
	payload.ID = *handlers.FmtUUID(paymentRequestID)

	edi858c, err := h.GHCPaymentRequestInvoiceGenerator.Generate(paymentRequest, false)
	if err == nil {
		payload.Edi, err = edi858c.EDIString(logger)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error generating EDI string for payment request ID: %s: %s", paymentRequestID, err))
		switch e := err.(type) {

		// NotFoundError -> Not Found response
		case services.NotFoundError:
			return paymentrequestop.NewGetPaymentRequestEDINotFound().
				WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))

		// InvalidInputError -> Unprocessable Entity reponse
		case services.InvalidInputError:
			return paymentrequestop.NewGetPaymentRequestEDIUnprocessableEntity().
				WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))

		// ConflictError -> Conflict Error reponse
		case services.ConflictError:
			return paymentrequestop.NewGetPaymentRequestEDIConflict().
				WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))

		// QueryError -> Internal Server error
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				// Note we do not expose this detail in the payload
				logger.Error("Error retrieving an EDI for thepayment request", zap.Error(e.Unwrap()))
			}
			return paymentrequestop.NewGetPaymentRequestEDIInternalServerError().
				WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return paymentrequestop.NewGetPaymentRequestEDIInternalServerError().
				WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	return paymentrequestop.NewGetPaymentRequestEDIOK().WithPayload(&payload)
}

// ProcessReviewedPaymentRequestsHandler returns the EDI for a given payment request
type ProcessReviewedPaymentRequestsHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
	services.PaymentRequestReviewedFetcher
	services.PaymentRequestStatusUpdater
	// Unable to get logger to pass in for the instantiation of
	// paymentrequest.InitNewPaymentRequestReviewedProcessor(h.DB(), logger, true),
	// This limitation has come up a few times
	// - https://dp3.atlassian.net/browse/MB-2352 (story to address issue)
	// - https://ustcdp3.slack.com/archives/CP6F568DC/p1592508325118600
	// - https://github.com/transcom/mymove/blob/c42adf61735be8ee8e5e83f41a656206f1e59b9d/pkg/handlers/primeapi/api.go
	// As a temporary workaround paymentrequest.InitNewPaymentRequestReviewedProcessor
	// is called directly in the handler
}

// Handle getting the EDI for a given payment request
func (h ProcessReviewedPaymentRequestsHandler) Handle(params paymentrequestop.ProcessReviewedPaymentRequestsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	paymentRequestID := uuid.FromStringOrNil(params.Body.PaymentRequestID.String())
	sendToSyncada := params.Body.SendToSyncada
	readFromSyncada := params.Body.ReadFromSyncada
	deleteFromSyncada := params.Body.DeleteFromSyncada
	paymentRequestStatus := params.Body.Status
	var paymentRequests models.PaymentRequests
	var updatedPaymentRequests models.PaymentRequests

	if sendToSyncada == nil {
		return paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "bad request, sendToSyncada flag required", h.GetTraceID()))
	}
	if readFromSyncada == nil {
		return paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "bad request, readFromSyncada flag required", h.GetTraceID()))
	}
	if deleteFromSyncada == nil {
		return paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "bad request, deleteFromSyncada flag required", h.GetTraceID()))
	}

	if *sendToSyncada {
		reviewedPaymentRequestProcessor, err := paymentrequest.InitNewPaymentRequestReviewedProcessor(h.DB(), logger, true, h.ICNSequencer())
		if err != nil {
			msg := fmt.Sprintf("failed to initialize InitNewPaymentRequestReviewedProcessor")
			logger.Error(msg, zap.Error(err))
			return paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
		err = reviewedPaymentRequestProcessor.ProcessReviewedPaymentRequest()
		if err != nil {
			msg := fmt.Sprintf("Error processing reviewed payment requests")
			logger.Error(msg, zap.Error(err))
			return paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	} else {
		if paymentRequestID != uuid.Nil {
			pr, err := h.PaymentRequestFetcher.FetchPaymentRequest(paymentRequestID)
			if err != nil {
				msg := fmt.Sprintf("Error finding Payment Request with ID: %s", params.Body.PaymentRequestID.String())
				logger.Error(msg, zap.Error(err))
				return paymentrequestop.NewProcessReviewedPaymentRequestsNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
			}
			paymentRequests = append(paymentRequests, pr)
		} else {
			reviewedPaymentRequests, err := h.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest()
			if err != nil {
				msg := fmt.Sprintf("function ProcessReviewedPaymentRequest failed call to FetchReviewedPaymentRequest")
				logger.Error(msg, zap.Error(err))
				return paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
			}
			for _, pr := range reviewedPaymentRequests {
				paymentRequests = append(paymentRequests, pr)
			}
		}

		// Update each payment request to have the given status
		for _, pr := range paymentRequests {
			switch paymentRequestStatus {
			case "PENDING":
				pr.Status = models.PaymentRequestStatusPending
			case "REVIEWED":
				reviewedAt := time.Now()
				pr.Status = models.PaymentRequestStatusReviewed
				pr.ReviewedAt = &reviewedAt
			case "SENT_TO_GEX":
				sentToGex := time.Now()
				pr.Status = models.PaymentRequestStatusSentToGex
				pr.SentToGexAt = &sentToGex
			case "RECEIVED_BY_GEX":
				recByGex := time.Now()
				pr.Status = models.PaymentRequestStatusReceivedByGex
				pr.ReceivedByGexAt = &recByGex
			case "PAID":
				paidAt := time.Now()
				pr.Status = models.PaymentRequestStatusPaid
				pr.PaidAt = &paidAt
			case "":
				sentToGex := time.Now()
				pr.Status = models.PaymentRequestStatusSentToGex
				pr.SentToGexAt = &sentToGex
			default:
				return paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "bad request, an invalid status type was used", h.GetTraceID()))
			}

			newPr := pr
			var nilEtag string
			updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&newPr, nilEtag)

			if err != nil {
				switch err.(type) {
				case services.NotFoundError:
					return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
				case services.PreconditionFailedError:
					return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
				default:
					logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
					return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
				}
			}
			updatedPaymentRequests = append(updatedPaymentRequests, *updatedPaymentRequest)
		}
	}

	if *readFromSyncada {
		// Set up viper to read environment variables
		v := viper.New()
		v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		v.AutomaticEnv()

		sshClient, err := cli.InitSyncadaSSH(v, logger)
		if err != nil {
			logger.Fatal("couldn't initialize SSH client", zap.Error(err))
		}
		defer func() {
			if closeErr := sshClient.Close(); closeErr != nil {
				logger.Fatal("could not close SFTP client", zap.Error(closeErr))
			}
		}()

		sftpClient, err := cli.InitSyncadaSFTP(sshClient, logger)
		if err != nil {
			logger.Fatal("couldn't initialize SFTP client", zap.Error(err))
		}
		defer func() {
			if closeErr := sftpClient.Close(); closeErr != nil {
				logger.Fatal("could not close SFTP client", zap.Error(closeErr))
			}
		}()

		wrappedSFTPClient := invoice.NewSFTPClientWrapper(sftpClient)
		syncadaSFTPSession := invoice.NewSyncadaSFTPReaderSession(wrappedSFTPClient, h.DB(), logger, *deleteFromSyncada)

		// TODO GEX will put different response types in different directories, but
		// Syncada puts everything in the same directory. When we have access to GEX in staging
		// we will have to change this to use separate paths for different response types.
		path := "/" + v.GetString(cli.SyncadaSFTPUserIDFlag) + v.GetString(cli.SyncadaSFTPOutboundDirectory)

		_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(path, time.Time{}, invoice.NewEDI997Processor(h.DB(), logger))
		if err != nil {
			logger.Error("Error reading 997 responses", zap.Error(err))
		} else {
			logger.Info("Successfully processed 997 responses")
		}
		_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(path, time.Time{}, invoice.NewEDI824Processor(h.DB(), logger))
		if err != nil {
			logger.Error("Error reading 824 responses", zap.Error(err))
		} else {
			logger.Info("Successfully processed 824 responses")
		}
	} else {
		logger.Info("Skipping reading from Syncada")
	}
	payload := payloads.PaymentRequests(&updatedPaymentRequests)

	return paymentrequestop.NewProcessReviewedPaymentRequestsOK().WithPayload(*payload)
}
