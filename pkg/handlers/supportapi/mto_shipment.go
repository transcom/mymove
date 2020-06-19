package supportapi

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"

	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMTOShipmentStatusHandlerFunc updates the status of a MTO Shipment
type UpdateMTOShipmentStatusHandlerFunc struct {
	handlers.HandlerContext
	services.Fetcher
	services.MTOShipmentStatusUpdater
}

// Handle updates the status of a MTO Shipment
func (h UpdateMTOShipmentStatusHandlerFunc) Handle(params mtoshipmentops.UpdateMTOShipmentStatusParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())
	status := models.MTOShipmentStatus(params.Body.Status)
	rejectionReason := params.Body.RejectionReason
	eTag := params.IfMatch

	shipment, err := h.UpdateMTOShipmentStatus(shipmentID, status, rejectionReason, eTag)

	if err != nil {
		logger.Error("UpdateMTOShipmentStatus error: ", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusUnprocessableEntity().WithPayload(
				payloads.ValidationError("The input provided did not pass validation.", h.GetTraceID(), e.ValidationErrors))
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case mtoshipment.ConflictStatusError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		default:
			return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	payload := payloads.MTOShipment(shipment)
	return mtoshipmentops.NewUpdateMTOShipmentStatusOK().WithPayload(payload)
}
