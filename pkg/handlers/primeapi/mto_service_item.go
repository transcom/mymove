package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMTOServiceItemHandler is the handler to update MTO shipments
type CreateMTOServiceItemHandler struct {
	handlers.HandlerContext
	mtoServiceItemCreator services.MTOServiceItemCreator
}

// Handle handler that updates a mto shipment
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemops.CreateMTOServiceItemParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	params.Body.SetMoveTaskOrderID(params.MoveTaskOrderID)
	params.Body.SetMtoShipmentID(params.MtoShipmentID)
	mtoServiceItem := payloads.MTOServiceItemModel(params.Body)

	mtoServiceItem, verrs, err := h.mtoServiceItemCreator.CreateMTOServiceItem(mtoServiceItem)
	if verrs != nil && verrs.HasAny() {
		return mtoserviceitemops.NewCreateMTOServiceItemBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(verrs.Error())})
	}

	if err != nil {
		logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return mtoserviceitemops.NewCreateMTOServiceItemBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return mtoserviceitemops.NewCreateMTOServiceItemPreconditionFailed().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}
	mtoServiceItemPayload := payloads.MTOServiceItem(mtoServiceItem)
	return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemPayload)
}
