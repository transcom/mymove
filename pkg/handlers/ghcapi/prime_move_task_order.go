package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

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

// UpdatePostCounselingInfoHandler updates the actual weight for a move task order
type UpdatePostCounselingInfoHandler struct {
	handlers.HandlerContext
	moveTaskOrderPostCounselingInfoUpdater services.MoveTaskOrderPostCounselingInfoUpdater
}

// Handle updating the actual weight for a move task order
func (h UpdatePostCounselingInfoHandler) Handle(params movetaskordercodeop.UpdatePostCounselingInfoParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderPostCounselingInfoUpdater.UpdatePostCounselingInfo(moveTaskOrderID, params.Body.ScheduledMoveDate, *params.Body.SecondaryPickupAddress, *params.Body.SecondaryDeliveryAddress, params.Body.PpmIsIncluded)
	if err != nil {
		logger.Error("ghciapi.UUpdatePostCounselingInfoHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskordercodeop.NewUpdatePostCounselingInfoNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskordercodeop.NewUpdatePostCounselingInfoBadRequest()
		default:
			return movetaskordercodeop.NewUpdatePostCounselingInfoInternalServerError()
		}
	}

	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return movetaskordercodeop.NewUpdatePostCounselingInfoOK().WithPayload(moveTaskOrderPayload)
}

type UpdateDestinationAddressHandler struct {
	handlers.HandlerContext
	moveTaskOrderDestinationAddressUpdater services.MoveTaskOrderDestinationAddressUpdater
}

func (h UpdateDestinationAddressHandler) Handle(params movetaskordercodeop.UpdateDestinationAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderDestinationAddressUpdater.UpdateDestinationAddress(moveTaskOrderID, *params.Body.DestinationAddress)
	if err != nil {
		logger.Error("ghciapi.UpdateDestinationAddressHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskordercodeop.NewUpdateDestinationAddressNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskordercodeop.NewUpdateDestinationAddressBadRequest()
		default:
			return movetaskordercodeop.NewUpdateDestinationAddressInternalServerError()
		}
	}

	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return movetaskordercodeop.NewUpdateDestinationAddressOK().WithPayload(moveTaskOrderPayload)
}