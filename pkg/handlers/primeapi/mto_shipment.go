package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/models"
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

	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest()
	}

	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(mtoShipment, eTag)

	//create a fn that loops and uses the payload to create mtoservice items

	createdServiceItem, verrs, err := h.MTOServiceItemCreator.CreateMTOServiceItem(&serviceItem)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating mto service item: ", zap.Error(verrs))
		payload := payloadForValidationError(handlers.ValidationErrMessage, "The information you provided is invalid.", h.GetTraceID(), verrs)

		return mtoserviceitemop.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload)
	}

	// return any errors
	if err != nil {
		logger.Error("Error creating mto service item: ", zap.Error(err))

		if strings.Contains(errors.Cause(err).Error(), models.ViolatesForeignKeyConstraint) {
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to create a mto service item.", h.GetTraceID())

			return mtoserviceitemop.NewCreateMTOServiceItemNotFound().WithPayload(payload)
		}

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
