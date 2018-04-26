package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"

	authctx "github.com/transcom/mymove/pkg/auth/context"
)

func payloadForMoveQueueItem(MoveQueueItem models.MoveQueueItem) *internalmessages.MoveQueueItem {
	MoveQueueItemPayload := internalmessages.MoveQueueItem{
		ID:               fmtUUID(MoveQueueItem.ID),
		CreatedAt:        fmtDateTime(MoveQueueItem.CreatedAt),
		Edipi:            MoveQueueItem.Edipi,
		Rank:             MoveQueueItem.Rank,
		CustomerName:     MoveQueueItem.CustomerName,
		Locator:          MoveQueueItem.Locator,
		Status:           MoveQueueItem.Status,
		OrdersType:       MoveQueueItem.OrdersType,
		MoveDate:         fmtDate(MoveQueueItem.MoveDate),
		CustomerDeadline: fmtDate(MoveQueueItem.CustomerDeadline),
		LastModifiedDate: fmtDateTime(MoveQueueItem.LastModifiedDate),
		LastModifiedName: MoveQueueItem.LastModifiedName,
	}
	return &MoveQueueItemPayload
}

// ShowQueueHandler returns a list of all MoveQueueItems in the moves queue
type ShowQueueHandler HandlerContext

// Handle retrieves a list of all MoveQueueItems in the system in the moves queue
func (h ShowQueueHandler) Handle(params queueop.ShowQueueParams) middleware.Responder {
	var response middleware.Responder
	// TODO: Only authorized users should be able to view
	_, ok := authctx.GetUser(params.HTTPRequest.Context())
	if !ok {
		h.logger.Error("No user logged in, this should never happen.")
	}

	lifecycleState := params.QueueType

	MoveQueueItems, err := models.GetMoveQueueItems(h.db, lifecycleState)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = queueop.ShowQueueBadRequest()
	} else {
		MoveQueueItemPayloads := make(internalmessages.MoveQueueItem, len(MoveQueueItems))
		for i, MoveQueueItem := range MoveQueueItems {
			MoveQueueItemPayload := payloadForMoveQueueItem(MoveQueueItem)
			MoveQueueItemPayloads[i] = &MoveQueueItemPayload
		}
		response = queueop.ShowQueueOK().WithPayload(MoveQueueItemPayloads)
	}
	return response
}
