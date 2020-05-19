package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
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
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError(handlers.ValidationErrMessage, "The information you provided is invalid", h.GetTraceID(), e.ValidationErrors)
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusPreconditionFailed().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ConflictError:
			payload := &supportmessages.ClientError{
				Title:    handlers.FmtString("This MTO service item is not in a state for the status to be changed"),
				Detail:   handlers.FmtString("Make sure the MTO service item's status has not already been changed or it has a rejectionReason if you are trying to reject it"),
				Instance: handlers.FmtUUID(h.GetTraceID()),
			}
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusConflict().WithPayload(payload)
		default:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusInternalServerError()
		}
	}

	payload := payloads.MTOServiceItem(mtoServiceItem)
	return mtoserviceitemops.NewUpdateMTOServiceItemStatusOK().WithPayload(payload)
}
