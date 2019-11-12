package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	entitlementscodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"
)

// UpdateMoveTaskOrderActualWeightHandler updates the actual weight for a move task order
type UpdateMoveTaskOrderActualWeightHandler struct {
	handlers.HandlerContext
	moveTaskOrderActualWeightUpdater services.MoveTaskOrderActualWeightUpdater
}

// Handle updating the actual weight for a move task order
func (h UpdateMoveTaskOrderActualWeightHandler) Handle(params movetaskordercodeop.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderActualWeightUpdater.UpdateMoveTaskOrderActualWeight(moveTaskOrderID, params.Body.ActualWeight)
	if err != nil {
		logger.Error("ghciapi.UpdateMoveTaskOrderActualWeightHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightBadRequest()
		default:
			return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightInternalServerError()
		}
	}

	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightOK().WithPayload(moveTaskOrderPayload)
}

// GetPrimeEntitlementsHandler fetches the entitlements for a move task order
type GetPrimeEntitlementsHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle getting the entitlements for a move task order
func (h GetPrimeEntitlementsHandler) Handle(params entitlementscodeop.GetPrimeEntitlementsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("ghciapi.GetPrimeEntitlementsHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return entitlementscodeop.NewGetPrimeEntitlementsNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return entitlementscodeop.NewGetPrimeEntitlementsBadRequest()
		default:
			return entitlementscodeop.NewGetPrimeEntitlementsInternalServerError()
		}
	}
	entitlements := payloadForEntitlements(&mto.Entitlements)

	return entitlementscodeop.NewGetPrimeEntitlementsOK().WithPayload(entitlements)
}
