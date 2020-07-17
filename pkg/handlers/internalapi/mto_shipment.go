package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {

	logger := h.LoggerFromRequest(params.HTTPRequest)

	payload := params.Body
	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment request body cannot be empty.", h.GetTraceID()))
	}

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload)

	mtoServiceItemsList, verrs := payloads.MTOServiceItemList(payload)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating mto service item list: ", zap.Error(verrs))

		return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
			"The MTO service item list is invalid.", h.GetTraceID(), nil))
	}

	//mtoShipment.MTOServiceItems = *mtoServiceItemsList
	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(mtoShipment, mtoServiceItemsList)
	if err != nil {
		logger.Error("primeapi.CreateMTOShipmentHandler", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	returnPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
}