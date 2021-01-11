package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/models"

	//"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	//"go.uber.org/zap"
)

// UpdateMTOServiceItemStatusHandler patches shipments
type UpdateMTOServiceItemStatusHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemUpdater
}

// Handle updates mto server item statuses
func (h UpdateMTOServiceItemStatusHandler) Handle(params mtoserviceitemops.UpdateMTOServiceItemStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtoServiceItemID := uuid.FromStringOrNil(params.MtoServiceItemID)
	status := models.MTOServiceItemStatus(params.Body.Status)
	eTag := params.IfMatch
	reason := params.Body.RejectionReason

	mtoServiceItem, err := h.UpdateMTOServiceItemStatus(mtoServiceItemID, status, reason, eTag)

	if err != nil {
		logger.Error("UpdateMTOServiceItemStatus error: ", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			payload := payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID())
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusNotFound().WithPayload(payload)
		case services.InvalidInputError:
			payload := payloads.ValidationError("The information you provided is invalid", h.GetTraceID(), e.ValidationErrors)
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		default:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	payload := payloads.MTOServiceItem(mtoServiceItem)
	return mtoserviceitemops.NewUpdateMTOServiceItemStatusOK().WithPayload(payload)
}
