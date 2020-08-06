package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

//
// CREATE
//

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil || (!session.IsMilApp() && session.ServiceMemberID == uuid.Nil) {
		return mtoshipmentops.NewCreateMTOShipmentUnauthorized()
	}

	payload := params.Body
	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment request body cannot be empty.", h.GetTraceID()))
	}

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
	// TODO: remove this status change once the UpdateMTOShipment api is implemented and can update to Submitted
	mtoShipment.Status = models.MTOShipmentStatusSubmitted
	serviceItemsList := make(models.MTOServiceItems, 0)
	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(mtoShipment, serviceItemsList)

	if err != nil {
		logger.Error("internalapi.CreateMTOShipmentHandler", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("internalapi.CreateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}
	returnPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
}

//
// UPDATE
//

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle updates the mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return mtoshipmentops.NewUpdateMTOShipmentUnauthorized()
	}

	if !session.IsMilApp() && session.ServiceMemberID == uuid.Nil {
		return mtoshipmentops.NewUpdateMTOShipmentForbidden()
	}

	copyParams := params
	mtoShipmentID := copyParams.MtoShipmentID

	payload := params.Body
	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment request body cannot be empty.", h.GetTraceID()))
	}

	// TODO: incorporate draft status, only push update thru mto updater if status not draft
	mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
	mtoShipment.ID = uuid.FromStringOrNil(mtoShipmentID.String())

	mtoShipment, err := h.mtoShipmentUpdater.UpdateMTOShipment(mtoShipment, params.IfMatch)

	if err != nil {
		logger.Error("internalapi.UpdateMTOShipmentHandler", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("internalapi.UpdateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	returnPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload)
}
