package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForQueueMoveModel(queueMove models.QueueMove) *internalmessages.QueueMove {
	queueMovePayload := internalmessages.QueueMove{
		ID:               fmtUUID(queueMove.ID),
		CreatedAt:        fmtDateTime(queueMove.CreatedAt),
		UpdatedAt:        fmtDateTime(queueMove.UpdatedAt),
		Edipi:            queueMove.Edipi,
		Rank:             queueMove.Rank,
		CustomerName:     queueMove.CustomerName,
		LocatorNumber:    queueMove.LocatorNumber,
		Status:           queueMove.Status,
		MoveType:         queueMove.MoveType,
		MoveDate:         fmtDate(queueMove.MoveDate),
		CustomerDeadline: fmtDate(queueMove.CustomerDeadline),
		LastModified:     queueMove.LastModified,
	}
	return &queueMovePayload
}

// ShowQueueHandler returns a list of all queueMoves in the new moves queue
type ShowQueueHandler HandlerContext

// Handle retrieves a list of all queueMoves in the system in the new moves queue
func (h ShowQueueHandler) Handle(params queueop.ShowQueueParams) middleware.Responder {
	var response middleware.Responder
	// TODO: Only authorized users should be able to view
	_, ok := authctx.GetUser(params.HTTPRequest.Context())
	if !ok {
		h.logger.Error("No user logged in, this should never happen.", zap.Error(err))
	}

	lifecycleState, err := params.queueType.String()
	if err != nil {
		response = queueop.NewShowQueueBadRequest()
		return response
	}

	queueMoves, err := models.GetQueueMoves(h.db, lifecycleState)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = queueop.ShowQueueBadRequest()
	} else {
		queueMovePayloads := make(internalmessages.ShowQueuePayload, len(queueMoves))
		for i, queueMove := range queueMoves {
			queueMovePayload := payloadForQueueMoveModel(queueMove)
			queueMovePayloads[i] = &queueMovePayload
		}
		response = queueop.ShowQueueOK().WithPayload(queueMovePayload)
	}
	return response
}
