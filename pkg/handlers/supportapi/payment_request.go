package supportapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"

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
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(appCtx, paymentRequestID)

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
	case supportmessages.PaymentRequestStatusPENDING:
		status = models.PaymentRequestStatusPending
	case supportmessages.PaymentRequestStatusREVIEWED:
		status = models.PaymentRequestStatusReviewed
		reviewedDate = time.Now()
	case supportmessages.PaymentRequestStatusSENTTOGEX:
		status = models.PaymentRequestStatusSentToGex
		sentGexDate = time.Now()
	case supportmessages.PaymentRequestStatusRECEIVEDBYGEX:
		status = models.PaymentRequestStatusReceivedByGex
		recGexDate = time.Now()
	case supportmessages.PaymentRequestStatusPAID:
		status = models.PaymentRequestStatusPaid
		paidAtDate = time.Now()
	case supportmessages.PaymentRequestStatusEDIERROR:
		status = models.PaymentRequestStatusEDIError
	case supportmessages.PaymentRequestStatusDEPRECATED:
		status = models.PaymentRequestStatusDeprecated
	}

	// If we got a rejection reason let's use it
	rejectionReason := existingPaymentRequest.RejectionReason
	if params.Body.RejectionReason != nil {
		rejectionReason = params.Body.RejectionReason
	}

	paymentRequestForUpdate := models.PaymentRequest{
		ID:                              existingPaymentRequest.ID,
		MoveTaskOrder:                   existingPaymentRequest.MoveTaskOrder,
		MoveTaskOrderID:                 existingPaymentRequest.MoveTaskOrderID,
		IsFinal:                         existingPaymentRequest.IsFinal,
		Status:                          status,
		RejectionReason:                 rejectionReason,
		RequestedAt:                     existingPaymentRequest.RequestedAt,
		ReviewedAt:                      &reviewedDate,
		SentToGexAt:                     &sentGexDate,
		ReceivedByGexAt:                 &recGexDate,
		PaidAt:                          &paidAtDate,
		PaymentRequestNumber:            existingPaymentRequest.PaymentRequestNumber,
		RecalculationOfPaymentRequestID: existingPaymentRequest.RecalculationOfPaymentRequestID,
		SequenceNumber:                  existingPaymentRequest.SequenceNumber,
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(appCtx, &paymentRequestForUpdate, params.IfMatch)

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
		DBConnection:    appCtx.DB(),
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
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	paymentRequestID := uuid.FromStringOrNil(params.PaymentRequestID.String())
	paymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(appCtx, paymentRequestID)
	if err != nil {
		msg := fmt.Sprintf("Error finding Payment Request for EDI generation with ID: %s", params.PaymentRequestID.String())
		logger.Error(msg, zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestEDINotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
	}

	var payload supportmessages.PaymentRequestEDI
	payload.ID = *handlers.FmtUUID(paymentRequestID)

	edi858c, err := h.GHCPaymentRequestInvoiceGenerator.Generate(appCtx, paymentRequest, false)
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

		// InvalidInputError -> Unprocessable Entity response
		case services.InvalidInputError:
			return paymentrequestop.NewGetPaymentRequestEDIUnprocessableEntity().
				WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))

		// ConflictError -> Conflict Error response
		case services.ConflictError:
			return paymentrequestop.NewGetPaymentRequestEDIConflict().
				WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))

		// QueryError -> Internal Server error
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				// Note we do not expose this detail in the payload
				logger.Error("Error retrieving an EDI for the payment request", zap.Error(e.Unwrap()))
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
	appCtx := appcontext.NewAppContext(h.DB(), logger)

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
		reviewedPaymentRequestProcessor, err := paymentrequest.InitNewPaymentRequestReviewedProcessor(appCtx, true, h.ICNSequencer(), h.GexSender())
		if err != nil {
			msg := "failed to initialize InitNewPaymentRequestReviewedProcessor"
			logger.Error(msg, zap.Error(err))
			return paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
		reviewedPaymentRequestProcessor.ProcessReviewedPaymentRequest(appCtx)
	} else {
		if paymentRequestID != uuid.Nil {
			pr, err := h.PaymentRequestFetcher.FetchPaymentRequest(appCtx, paymentRequestID)
			if err != nil {
				msg := fmt.Sprintf("Error finding Payment Request with ID: %s", params.Body.PaymentRequestID.String())
				logger.Error(msg, zap.Error(err))
				return paymentrequestop.NewProcessReviewedPaymentRequestsNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
			}
			paymentRequests = append(paymentRequests, pr)
		} else {
			reviewedPaymentRequests, err := h.PaymentRequestReviewedFetcher.FetchReviewedPaymentRequest(appCtx)
			if err != nil {
				msg := "function ProcessReviewedPaymentRequest failed call to FetchReviewedPaymentRequest"
				logger.Error(msg, zap.Error(err))
				return paymentrequestop.NewProcessReviewedPaymentRequestsInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
			}
			paymentRequests = append(paymentRequests, reviewedPaymentRequests...)
		}

		// Update each payment request to have the given status
		for _, pr := range paymentRequests {
			switch paymentRequestStatus {
			case supportmessages.PaymentRequestStatusPENDING:
				pr.Status = models.PaymentRequestStatusPending
			case supportmessages.PaymentRequestStatusREVIEWED:
				reviewedAt := time.Now()
				pr.Status = models.PaymentRequestStatusReviewed
				pr.ReviewedAt = &reviewedAt
			case supportmessages.PaymentRequestStatusSENTTOGEX:
				sentToGex := time.Now()
				pr.Status = models.PaymentRequestStatusSentToGex
				pr.SentToGexAt = &sentToGex
			case supportmessages.PaymentRequestStatusRECEIVEDBYGEX:
				recByGex := time.Now()
				pr.Status = models.PaymentRequestStatusReceivedByGex
				pr.ReceivedByGexAt = &recByGex
			case supportmessages.PaymentRequestStatusPAID:
				paidAt := time.Now()
				pr.Status = models.PaymentRequestStatusPaid
				pr.PaidAt = &paidAt
			case supportmessages.PaymentRequestStatusEDIERROR:
				pr.Status = models.PaymentRequestStatusEDIError
			case supportmessages.PaymentRequestStatusDEPRECATED:
				pr.Status = models.PaymentRequestStatusDeprecated
			case "":
				sentToGex := time.Now()
				pr.Status = models.PaymentRequestStatusSentToGex
				pr.SentToGexAt = &sentToGex
			default:
				return paymentrequestop.NewProcessReviewedPaymentRequestsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "bad request, an invalid status type was used", h.GetTraceID()))
			}

			newPr := pr
			var nilEtag string
			updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(appCtx, &newPr, nilEtag)

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
		path997 := v.GetString(cli.GEXSFTP997PickupDirectory)
		path824 := v.GetString(cli.GEXSFTP824PickupDirectory)

		sshClient, err := cli.InitGEXSSH(v, logger)
		if err != nil {
			logger.Fatal("couldn't initialize SSH client", zap.Error(err))
		}
		defer func() {
			if closeErr := sshClient.Close(); closeErr != nil {
				logger.Fatal("could not close SFTP client", zap.Error(closeErr))
			}
		}()

		sftpClient, err := cli.InitGEXSFTP(sshClient, logger)
		if err != nil {
			logger.Fatal("couldn't initialize SFTP client", zap.Error(err))
		}
		defer func() {
			if closeErr := sftpClient.Close(); closeErr != nil {
				logger.Fatal("could not close SFTP client", zap.Error(closeErr))
			}
		}()

		wrappedSFTPClient := invoice.NewSFTPClientWrapper(sftpClient)
		syncadaSFTPSession := invoice.NewSyncadaSFTPReaderSession(wrappedSFTPClient, *deleteFromSyncada)

		_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(appCtx, path997, time.Time{}, invoice.NewEDI997Processor())
		if err != nil {
			logger.Error("Error reading 997 responses", zap.Error(err))
		} else {
			logger.Info("Successfully processed 997 responses")
		}
		_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(appCtx, path824, time.Time{}, invoice.NewEDI824Processor())
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

// RecalculatePaymentRequestHandler recalculates a payment request
type RecalculatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestRecalculator
}

// Handle getting the EDI for a given payment request
func (h RecalculatePaymentRequestHandler) Handle(params paymentrequestop.RecalculatePaymentRequestParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	paymentRequestID := uuid.FromStringOrNil(params.PaymentRequestID.String())

	startNewTx := true // start new Tx in service
	newPaymentRequest, err := h.PaymentRequestRecalculator.RecalculatePaymentRequest(appCtx, paymentRequestID, startNewTx)

	if err != nil {
		switch e := err.(type) {
		case *services.BadDataError:
			return paymentrequestop.NewRecalculatePaymentRequestBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceID()))
		case services.NotFoundError:
			return paymentrequestop.NewRecalculatePaymentRequestNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.ConflictError:
			return paymentrequestop.NewRecalculatePaymentRequestConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return paymentrequestop.NewRecalculatePaymentRequestPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return paymentrequestop.NewRecalculatePaymentRequestUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		case services.InvalidCreateInputError:
			return paymentrequestop.NewRecalculatePaymentRequestUnprocessableEntity().WithPayload(payloads.ValidationError(err.Error(), h.GetTraceID(), e.ValidationErrors))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				// Note we do not expose this detail in the payload
				logger.Error("Error recalculating payment request", zap.Error(e.Unwrap()))
			}
			return paymentrequestop.NewRecalculatePaymentRequestInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			logger.Error(fmt.Sprintf("Error recalculating payment request for ID: %s: %s", paymentRequestID, err))
			return paymentrequestop.NewRecalculatePaymentRequestInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	returnPayload := payloads.PaymentRequest(newPaymentRequest)

	return paymentrequestop.NewRecalculatePaymentRequestCreated().WithPayload(returnPayload)
}
