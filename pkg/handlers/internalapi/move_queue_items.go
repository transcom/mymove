package internalapi

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/auth"
	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForMoveQueueItem(MoveQueueItem models.MoveQueueItem, StorageInTransits internalmessages.StorageInTransits) *internalmessages.MoveQueueItem {
	MoveQueueItemPayload := internalmessages.MoveQueueItem{
		ID:                         handlers.FmtUUID(MoveQueueItem.ID),
		CreatedAt:                  handlers.FmtDateTime(MoveQueueItem.CreatedAt),
		Edipi:                      swag.String(MoveQueueItem.Edipi),
		Rank:                       MoveQueueItem.Rank,
		CustomerName:               swag.String(MoveQueueItem.CustomerName),
		Locator:                    swag.String(MoveQueueItem.Locator),
		GblNumber:                  handlers.FmtStringPtr(MoveQueueItem.GBLNumber),
		Status:                     swag.String(MoveQueueItem.Status),
		PpmStatus:                  handlers.FmtStringPtr(MoveQueueItem.PpmStatus),
		HhgStatus:                  handlers.FmtStringPtr(MoveQueueItem.HhgStatus),
		OrdersType:                 swag.String(MoveQueueItem.OrdersType),
		MoveDate:                   handlers.FmtDatePtr(MoveQueueItem.MoveDate),
		SubmittedDate:              handlers.FmtDateTimePtr(MoveQueueItem.SubmittedDate),
		LastModifiedDate:           handlers.FmtDateTime(MoveQueueItem.LastModifiedDate),
		OriginDutyStationName:      swag.String(MoveQueueItem.OriginDutyStationName),
		DestinationDutyStationName: swag.String(MoveQueueItem.DestinationDutyStationName),
		StorageInTransits:          StorageInTransits,
	}
	return &MoveQueueItemPayload
}

// ShowQueueHandler returns a list of all MoveQueueItems in the moves queue
type ShowQueueHandler struct {
	handlers.HandlerContext
	storageInTransitsIndexer services.StorageInTransitsIndexer
}

// Handle retrieves a list of all MoveQueueItems in the system in the moves queue
func (h ShowQueueHandler) Handle(params queueop.ShowQueueParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return queueop.NewShowQueueForbidden()
	}

	lifecycleState := params.QueueType

	MoveQueueItems, err := models.GetMoveQueueItems(h.DB(), lifecycleState)
	if err != nil {
		h.Logger().Error("Loading Queue", zap.String("State", lifecycleState), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	MoveQueueItemPayloads := make([]*internalmessages.MoveQueueItem, len(MoveQueueItems))
	for i, MoveQueueItem := range MoveQueueItems {
		var storageInTransits []models.StorageInTransit
		if MoveQueueItem.ShipmentID != uuid.Nil {
			storageInTransits, err = h.storageInTransitsIndexer.IndexStorageInTransits(MoveQueueItem.ShipmentID, session)

			if err != nil {
				h.Logger().Error(fmt.Sprintf("SITs Retrieval failed for shipment: %s", MoveQueueItem.ShipmentID), zap.Error(err))
				return handlers.ResponseForError(h.Logger(), err)
			}
		}

		storageInTransitsList := make(internalmessages.StorageInTransits, len(storageInTransits))

		for i, storageInTransit := range storageInTransits {
			storageInTransitsList[i] = payloadForStorageInTransitModel(&storageInTransit)
		}

		MoveQueueItemPayload := payloadForMoveQueueItem(MoveQueueItem, storageInTransitsList)
		MoveQueueItemPayloads[i] = MoveQueueItemPayload

	}
	return queueop.NewShowQueueOK().WithPayload(MoveQueueItemPayloads)
}
