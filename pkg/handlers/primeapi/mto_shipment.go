package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	mtoshipmentservice "github.com/transcom/mymove/pkg/services/mto_shipment"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
)

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoShipment := payloads.MTOShipmentModel(params.Body)
	unmodifiedSince := time.Time(params.IfUnmodifiedSince)

	mtoShipment, err := h.mtoShipmentUpdater.UpdateMTOShipment(mtoShipment, unmodifiedSince)
	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
		switch err.(type) {
		case mtoshipmentservice.ErrNotFound:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound()
		case mtoshipmentservice.ErrInvalidInput:
			return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipmentservice.ErrPreconditionFailed:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed()
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError()
		}
	}
	mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload)
}
