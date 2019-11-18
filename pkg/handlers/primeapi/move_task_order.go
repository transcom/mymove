package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

// ListMoveTaskOrdersHandler fetches all the move task orders
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

// Handle fetching all the move task orders with an option to only get those since a certain date
func (h ListMoveTaskOrdersHandler) Handle(params movetaskorderops.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Q()
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
	}

	payload := make(primemessages.MoveTaskOrders, len(mtos))

	for i, m := range mtos {
		payload[i] = payloads.MoveTaskOrder(m)
	}

	return movetaskorderops.NewListMoveTaskOrdersOK().WithPayload(payload)
}

// UpdateMoveTaskOrderEstimatedWeightHandler updates the move task order's estimated weight
type UpdateMoveTaskOrderEstimatedWeightHandler struct {
	handlers.HandlerContext
	moveTaskOrderPrimeEstimatedWeightUpdater services.MoveTaskOrderPrimeEstimatedWeightUpdater
}

// Handle updating the move task order's estimated weight
func (h UpdateMoveTaskOrderEstimatedWeightHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderEstimatedWeightParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	primeEstimatedWeight := unit.Pound(params.Body.PrimeEstimatedWeight)
	mto, err := h.moveTaskOrderPrimeEstimatedWeightUpdater.UpdatePrimeEstimatedWeight(moveTaskOrderID, primeEstimatedWeight, time.Now())
	if err != nil {
		logger.Error("primeapi.UpdateMoveTaskOrderEstimatedWeightHandler error", zap.Error(err))
		switch e := err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightNotFound()
		case movetaskorderservice.ErrInvalidInput:
			payload := &primemessages.ValidationError{
				InvalidFields: e.InvalidFields(),
				ClientError: primemessages.ClientError{
					Title:    handlers.FmtString(handlers.ValidationErrMessage),
					Detail:   handlers.FmtString(e.Error()),
					Instance: handlers.FmtUUID(h.GetTraceID()),
				},
			}
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightUnprocessableEntity().WithPayload(payload)
		default:
			return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightOK().WithPayload(moveTaskOrderPayload)
}

// UpdateMoveTaskOrderActualWeightHandler updates the actual weight for a move task order
type UpdateMoveTaskOrderActualWeightHandler struct {
	handlers.HandlerContext
	moveTaskOrderActualWeightUpdater services.MoveTaskOrderActualWeightUpdater
}

// Handle updating the actual weight for a move task order
func (h UpdateMoveTaskOrderActualWeightHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderActualWeightUpdater.UpdateMoveTaskOrderActualWeight(moveTaskOrderID, params.Body.ActualWeight)
	if err != nil {
		logger.Error("primeapi.UpdateMoveTaskOrderActualWeightHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderActualWeightNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskorderops.NewUpdateMoveTaskOrderActualWeightBadRequest()
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderActualWeightInternalServerError()
		}
	}

	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderActualWeightOK().WithPayload(moveTaskOrderPayload)
}

// GetPrimeEntitlementsHandler fetches the entitlements for a move task order
type GetPrimeEntitlementsHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle getting the entitlements for a move task order
func (h GetPrimeEntitlementsHandler) Handle(params movetaskorderops.GetPrimeEntitlementsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("primeapi.GetPrimeEntitlementsHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewGetPrimeEntitlementsNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskorderops.NewGetPrimeEntitlementsBadRequest()
		default:
			return movetaskorderops.NewGetPrimeEntitlementsInternalServerError()
		}
	}
	entitlements := payloads.Entitlements(&mto.Entitlements)

	return movetaskorderops.NewGetPrimeEntitlementsOK().WithPayload(entitlements)
}

type GetMoveTaskOrderCustomerHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

func (h GetMoveTaskOrderCustomerHandler) Handle(params movetaskorderops.GetMoveTaskOrderCustomerParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("ghciap.GetMoveTaskOrderCustomerHandler error", zap.Error(err))
		switch e := err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewGetMoveTaskOrderCustomerNotFound()
		case movetaskorderservice.ErrInvalidInput:
			payload := &primemessages.ValidationError{
				InvalidFields: e.InvalidFields(),
				ClientError: primemessages.ClientError{
					Title:    handlers.FmtString(handlers.ValidationErrMessage),
					Detail:   handlers.FmtString(e.Error()),
					Instance: handlers.FmtUUID(h.GetTraceID()),
				},
			}
			return movetaskorderops.NewGetMoveTaskOrderCustomerUnprocessableEntity().WithPayload(payload)
		default:
			return movetaskorderops.NewGetMoveTaskOrderCustomerInternalServerError()
		}
	}
	customer := payloads.CustomerWithMTO(mto)
	return movetaskorderops.NewGetMoveTaskOrderCustomerOK().WithPayload(customer)
}
