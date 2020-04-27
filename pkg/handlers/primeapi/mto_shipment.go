package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {

	logger := h.LoggerFromRequest(params.HTTPRequest)

	payload := params.Body
	moveTaskOrderID := params.MoveTaskOrderID
	eTag := params.IfMatch

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload, moveTaskOrderID)

	//create a fn that loops and uses the payload to create mtoservice items
	mtoServiceItemsList, verrs := payloads.MTOServiceItemList(payload)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating mto service item: ", zap.Error(verrs))

		return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity()
	}

	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest()
	}

	mtoShipment.MTOServiceItems = mtoServiceItemsList

	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(mtoShipment, eTag)

	// return any errors
	if err != nil {
		logger.Error("Error creating mto service item: ", zap.Error(err))
		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	if err != nil {
		return nil
	}
	returnPayload := payloads.MTOShipmentFromCreate(mtoShipment)
	return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoShipment := payloads.MTOShipmentModel(params.Body)
	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))
	mtoShipment, err := h.mtoShipmentUpdater.UpdateMTOShipment(mtoShipment, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound()
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError()
		}
	}
	mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload)
}
