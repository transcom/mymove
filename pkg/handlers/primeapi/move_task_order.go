package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/unit"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

func (h ListMoveTaskOrdersHandler) Handle(params movetaskorderops.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders
	err := h.DB().All(&mtos)

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

type UpdateMoveTaskOrderEstimatedWeightHandler struct {
	handlers.HandlerContext
	moveTaskOrderPrimeEstimatedWeightUpdater services.MoveTaskOrderPrimeEstimatedWeightUpdater
}

func (h UpdateMoveTaskOrderEstimatedWeightHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderEstimatedWeightParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	primeEstimatedWeight := unit.Pound(params.Body.PrimeEstimatedWeight)
	mto, err := h.moveTaskOrderPrimeEstimatedWeightUpdater.UpdatePrimeEstimatedWeight(moveTaskOrderID, primeEstimatedWeight, time.Now())
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightBadRequest()
		default:
			return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightOK().WithPayload(moveTaskOrderPayload)
}
