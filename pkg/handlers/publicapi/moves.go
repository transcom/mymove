package publicapi

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForMoveModel(move *models.Move) *apimessages.Move {
	if move == nil {
		return nil
	}

	var SelectedMoveType = apimessages.SelectedMoveTypeHHG
	if move.SelectedMoveType != nil {
		SelectedMoveType = apimessages.SelectedMoveType(*move.SelectedMoveType)
	}

	cancelReason := ""
	if move.CancelReason != nil {
		cancelReason = *move.CancelReason
	}
	return &apimessages.Move{
		SelectedMoveType: &SelectedMoveType,
		OrdersID:         handlers.FmtUUID(move.OrdersID),
		Status:           apimessages.MoveStatus(move.Status),
		Locator:          swag.String(move.Locator),
		CancelReason:     swag.String(cancelReason),
	}
}
