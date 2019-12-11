package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

func PayloadForMoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	payload := &ghcmessages.MoveTaskOrder{
		ID:                strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:         strfmt.Date(moveTaskOrder.CreatedAt),
		IsAvaiableToPrime: moveTaskOrder.IsAvailableToPrime,
		IsCanceled:        moveTaskOrder.IsCancelled,
		MoveOrdersID:      strfmt.UUID(moveTaskOrder.MoveOrderID.String()),
		ReferenceID:       moveTaskOrder.ReferenceID,
		UpdatedAt:         strfmt.Date(moveTaskOrder.UpdatedAt),
	}
	return payload
}
