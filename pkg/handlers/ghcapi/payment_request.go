package ghcapi

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
)

// GetPaymentRequestForMoveHandler gets payment requests associated with a move
type GetPaymentRequestForMoveHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestListFetcher
}

// Handle handles the HTTP handling for GetPaymentRequestForMoveHandler
func (h GetPaymentRequestForMoveHandler) Handle(
	params paymentrequestop.GetPaymentRequestsForMoveParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			locator := params.Locator

			paymentRequests, err := h.FetchPaymentRequestListByMove(
				appCtx,
				locator,
			)
			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error fetching Payment Request for locator: %s", locator), zap.Error(err))
				return paymentrequestop.NewGetPaymentRequestNotFound(), err
			}

			returnPayload, err := payloads.PaymentRequests(appCtx, paymentRequests, h.FileStorer())
			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error building payment requests payload for locator: %s", locator), zap.Error(err))
				return paymentrequestop.NewGetPaymentRequestsForMoveInternalServerError(), err
			}

			return paymentrequestop.NewGetPaymentRequestsForMoveOK().
					WithPayload(*returnPayload),
				nil
		})
}

// GetPaymentRequestHandler gets payment requests
type GetPaymentRequestHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestFetcher
}

// Handle gets payment requests
func (h GetPaymentRequestHandler) Handle(
	params paymentrequestop.GetPaymentRequestParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())
			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
				return paymentrequestop.NewGetPaymentRequestInternalServerError(), nil
			}

			paymentRequest, err := h.FetchPaymentRequest(appCtx, paymentRequestID)
			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error fetching Payment Request with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
				return paymentrequestop.NewGetPaymentRequestNotFound(), nil
			}

			if reflect.DeepEqual(paymentRequest, models.PaymentRequest{}) {
				paymentRequestUUID, _ := uuid.FromString(params.PaymentRequestID.String())
				notFoundErr := apperror.NewNotFoundError(
					paymentRequestUUID,
					"Could not find a Payment Request with ID",
				)
				appCtx.Logger().Info(notFoundErr.Error())
				return paymentrequestop.NewGetPaymentRequestNotFound(), notFoundErr
			}

			returnPayload, err := payloads.PaymentRequest(appCtx, &paymentRequest, h.FileStorer())
			if err != nil {
				return paymentrequestop.NewGetPaymentRequestInternalServerError(), err
			}

			response := paymentrequestop.NewGetPaymentRequestOK().WithPayload(returnPayload)

			return response, nil
		})
}

// UpdatePaymentRequestStatusHandler updates payment requests status
type UpdatePaymentRequestStatusHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestStatusUpdater
	services.PaymentRequestFetcher
}

// Handle updates payment requests status
func (h UpdatePaymentRequestStatusHandler) Handle(
	params paymentrequestop.UpdatePaymentRequestStatusParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())
			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
				return paymentrequestop.NewGetPaymentRequestInternalServerError(), err
			}

			// Let's fetch the existing payment request using the PaymentRequestFetcher service object
			existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(
				appCtx,
				paymentRequestID,
			)
			if err != nil {
				paymentRequestUUID, _ := uuid.FromString(params.PaymentRequestID.String())
				notFoundErr := apperror.NewNotFoundError(
					paymentRequestUUID,
					"Could not find a Payment Request for status update with ID",
				)
				appCtx.Logger().Info(notFoundErr.Error())
				return paymentrequestop.NewGetPaymentRequestNotFound(), notFoundErr
			}

			now := time.Now()
			existingPaymentRequest.Status = models.PaymentRequestStatus(params.Body.Status)

			if existingPaymentRequest.Status != models.PaymentRequestStatusReviewed &&
				existingPaymentRequest.Status != models.PaymentRequestStatusReviewedAllRejected {
				errMessage := fmt.Sprintf(
					"Incoming payment request status should be REVIEWED or REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED instead it was: %s",
					existingPaymentRequest.Status.String(),
				)
				unprocessableErr := apperror.NewUnprocessableEntityError(errMessage)
				payload := payloadForValidationError(
					"Unable to complete request",
					errMessage,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)
				return paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity().
					WithPayload(payload), unprocessableErr
			}

			existingPaymentRequest.ReviewedAt = &now

			// If we got a rejection reason let's use it
			if params.Body.RejectionReason != nil {
				existingPaymentRequest.RejectionReason = params.Body.RejectionReason
			}

			// Capture update attempt in audit log
			_, err = audit.Capture(appCtx, &existingPaymentRequest, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().
					Error("Auditing service error for payment request update.", zap.Error(err))
				return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError(), err
			}

			// And now let's save our updated model object using the PaymentRequestUpdater service object.
			updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(
				appCtx,
				&existingPaymentRequest,
				params.IfMatch,
			)
			if err != nil {
				switch err.(type) {
				case apperror.NotFoundError:
					return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.PreconditionFailedError:
					return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity().WithPayload(payload), err
				default:
					appCtx.Logger().Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
					return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError(), err
				}
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.PaymentRequestUpdateEventKey,
				MtoID:           updatedPaymentRequest.MoveTaskOrderID,
				UpdatedObjectID: updatedPaymentRequest.ID,
				EndpointKey:     event.GhcUpdatePaymentRequestStatusEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().
					Error("ghcapi.UpdatePaymentRequestStatusHandler could not generate the event")
			}

			returnPayload, err := payloads.PaymentRequest(appCtx, updatedPaymentRequest, h.FileStorer())
			if err != nil {
				return paymentrequestop.NewGetPaymentRequestInternalServerError(), err
			}

			return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload), nil
		})
}

// ShipmentsSITBalanceHandler is the handler type for getShipmentsPaymentSITBalance
type ShipmentsSITBalanceHandler struct {
	handlers.HandlerConfig
	services.ShipmentsPaymentSITBalance
}

// Handle handles the getShipmentsPaymentSITBalance request
func (h ShipmentsSITBalanceHandler) Handle(
	params paymentrequestop.GetShipmentsPaymentSITBalanceParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			paymentRequestID := uuid.FromStringOrNil(params.PaymentRequestID.String())

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetShipmentsPaymentSITBalance error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return paymentrequestop.NewGetShipmentsPaymentSITBalanceNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return paymentrequestop.NewGetShipmentsPaymentSITBalanceForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return paymentrequestop.NewGetShipmentsPaymentSITBalanceInternalServerError(), err
				default:
					return paymentrequestop.NewGetShipmentsPaymentSITBalanceInternalServerError(), err
				}
			}

			shipmentSITBalances, err := h.ListShipmentPaymentSITBalance(appCtx, paymentRequestID)
			if err != nil {
				return handleError(err)
			}

			payload := payloads.ShipmentsPaymentSITBalance(shipmentSITBalances)

			return paymentrequestop.NewGetShipmentsPaymentSITBalanceOK().WithPayload(payload), nil
		})
}

type PaymentRequestBulkDownloadHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestBulkDownloadCreator
}

func (h PaymentRequestBulkDownloadHandler) Handle(params paymentrequestop.BulkDownloadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			logger := appCtx.Logger()

			paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
			if err != nil {
				errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

				errPayload := &ghcmessages.Error{Message: &errInstance}

				appCtx.Logger().Error(err.Error())
				return paymentrequestop.NewBulkDownloadBadRequest().WithPayload(errPayload), err
			}

			paymentRequestPacket, err := h.PaymentRequestBulkDownloadCreator.CreatePaymentRequestBulkDownload(appCtx, paymentRequestID)
			if err != nil {
				logger.Error("Error creating Payment Request Downloads Packet", zap.Error(err))
				errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
				errPayload := &ghcmessages.Error{Message: &errInstance}
				return paymentrequestop.NewBulkDownloadInternalServerError().
					WithPayload(errPayload), err
			}

			payload := io.NopCloser(paymentRequestPacket)
			filename := fmt.Sprintf("inline; filename=\"PaymentRequestBulkPacket-%s.pdf\"", time.Now().Format("01-02-2006_15-04-05"))

			return paymentrequestop.NewBulkDownloadOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}
