package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"

	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForMoveTaskOrder(mto models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {

	payload := &ghcmessages.MoveTaskOrder{
		ID:     *handlers.FmtUUID(mto.ID),
		MoveID: *handlers.FmtUUID(mto.MoveID),
	}
	if mto.ActualWeight != nil {
		payload.ActualWeight = *handlers.FmtInt64(int64(*mto.ActualWeight))
	}
	return payload
}

// UpdateMoveTaskOrderActualWeightHandler updates the actual weight for a move task order
type UpdateMoveTaskOrderActualWeightHandler struct {
	handlers.HandlerContext
}

// Handle updating the actual weight for a move task order
func (h UpdateMoveTaskOrderActualWeightHandler) Handle(params movetaskordercodeop.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID, _ := uuid.FromString(params.MoveTaskOrderID)
	mto, err := models.FetchMoveTaskOrder(h.DB(), moveTaskOrderID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.PatchActualWeight
	actualWeight := unit.Pound(payload.ActualWeight)
	mto.ActualWeight = &actualWeight

	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)

	return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightOK().WithPayload(moveTaskOrderPayload)
}
