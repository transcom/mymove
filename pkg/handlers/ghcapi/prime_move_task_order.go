package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// UpdateMoveTaskOrderActualWeightHandler updates the actual weight for a move task order
type UpdateMoveTaskOrderActualWeightHandler struct {
	handlers.HandlerContext
}

// Handle updating the actual weight for a move task order
func (h UpdateMoveTaskOrderActualWeightHandler) Handle(params movetaskordercodeop.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
	moveTaskOrderID, _ := uuid.FromString(params.MoveTaskOrderID)
	// fetch the move task order
	fmt.Println(moveTaskOrderID)

	payload := params.PatchActualWeight
	actualWeight := payload.ActualWeight
	fmt.Println(actualWeight)
	mto := &ghcmessages.MoveTaskOrder{
		ActualWeight: actualWeight,
	}
	return movetaskordercodeop.NewUpdateMoveTaskOrderActualWeightOK().WithPayload(mto)
}
