package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
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

func payloadForClientError(title string, detail string, instance uuid.UUID) *supportmessages.ClientError {
	return &supportmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

func payloadForValidationError(title string, detail string, instance uuid.UUID, validationErrors *validate.Errors) *supportmessages.ValidationError {
	return &supportmessages.ValidationError{
		InvalidFields: handlers.NewValidationErrorsResponse(validationErrors).Errors,
		ClientError:   *payloadForClientError(title, detail, instance),
	}
}

// Handle updates mto server item statuses
func (h UpdateMTOServiceItemStatusHandler) Handle(params mtoserviceitemops.UpdateMTOServiceItemStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtoServiceItemID := uuid.FromStringOrNil(params.MtoServiceItemID)
	status := models.MTOServiceItemStatus(*params.Body.Status)
	eTag := params.IfMatch
	reason := params.Body.Reason

	mtoServiceItem, err := h.UpdateMTOServiceItemStatus(mtoServiceItemID, status, reason, eTag)

	if err != nil {
		logger.Error("UpdateMTOServiceItemStatus error: ", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "UpdateMTOServiceItemStatus", h.GetTraceID(), e.ValidationErrors)
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusPreconditionFailed().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ConflictError:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusConflict().WithPayload(&supportmessages.Error{Message: handlers.FmtString("Can only update status from SUBMITTED to APPROVED or REJECTED")})
		default:
			return mtoserviceitemops.NewUpdateMTOServiceItemStatusInternalServerError()
		}
	}

	payload := payloads.MTOServiceItem(mtoServiceItem)
	return mtoserviceitemops.NewUpdateMTOServiceItemStatusOK().WithPayload(payload)
}
