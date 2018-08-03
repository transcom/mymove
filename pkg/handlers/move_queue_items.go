package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForMoveQueueItem(MoveQueueItem models.MoveQueueItem) *internalmessages.MoveQueueItem {
	MoveQueueItemPayload := internalmessages.MoveQueueItem{
		ID:               fmtUUID(MoveQueueItem.ID),
		CreatedAt:        fmtDateTime(MoveQueueItem.CreatedAt),
		Edipi:            swag.String(MoveQueueItem.Edipi),
		Rank:             MoveQueueItem.Rank,
		CustomerName:     swag.String(MoveQueueItem.CustomerName),
		Locator:          swag.String(MoveQueueItem.Locator),
		Status:           swag.String(MoveQueueItem.Status),
		PpmStatus:        MoveQueueItem.PpmStatus,
		OrdersType:       swag.String(MoveQueueItem.OrdersType),
		MoveDate:         fmtDatePtr(MoveQueueItem.MoveDate),
		CustomerDeadline: fmtDate(MoveQueueItem.CustomerDeadline),
		LastModifiedDate: fmtDateTime(MoveQueueItem.LastModifiedDate),
		LastModifiedName: swag.String(MoveQueueItem.LastModifiedName),
	}
	return &MoveQueueItemPayload
}

// ShowQueueHandler returns a list of all MoveQueueItems in the moves queue
type ShowQueueHandler HandlerContext

// Handle retrieves a list of all MoveQueueItems in the system in the moves queue
func (h ShowQueueHandler) Handle(params queueop.ShowQueueParams) middleware.Responder {
	// TODO: Check user is authorized to see office queues
	lifecycleState := params.QueueType

	MoveQueueItems, err := models.GetMoveQueueItems(h.db, lifecycleState)
	if err != nil {
		h.logger.Error("Loading Queue", zap.String("State", lifecycleState), zap.Error(err))
		return responseForError(h.logger, err)
	}

	MoveQueueItemPayloads := make([]*internalmessages.MoveQueueItem, len(MoveQueueItems))
	for i, MoveQueueItem := range MoveQueueItems {
		MoveQueueItemPayload := payloadForMoveQueueItem(MoveQueueItem)
		MoveQueueItemPayloads[i] = MoveQueueItemPayload
	}
	return queueop.NewShowQueueOK().WithPayload(MoveQueueItemPayloads)
}
