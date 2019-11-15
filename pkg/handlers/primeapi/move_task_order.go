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

type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

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
		logger.Error("ghciap.UpdateMoveTaskOrderEstimatedWeightHandler error", zap.Error(err))
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
