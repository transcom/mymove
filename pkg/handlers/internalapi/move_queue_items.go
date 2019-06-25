package internalapi

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
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

type JSONDate time.Time

// Dates without timestamps need custom unmarshalling
func (j *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JSONDate(t)
	return nil
}

type QueueSitData struct {
	ID              uuid.UUID `json:"id"`
	Status          string    `json:"status"`
	ActualStartDate JSONDate  `json:"actual_start_date"`
	OutDate         JSONDate  `json:"out_date"`
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
		var sits []QueueSitData
		err := json.Unmarshal([]byte(MoveQueueItem.SitArray), &sits)

		if err != nil {
			h.Logger().Error("Unmarshalling SITs", zap.Error(err))
		}

		if len(sits) == 1 && sits[0].ID == uuid.Nil {
			sits = []QueueSitData{}
		}

		storageInTransitsList := make(internalmessages.StorageInTransits, len(sits))

		for i, storageInTransit := range sits {
			actualStartDate := time.Time(storageInTransit.ActualStartDate)
			outDate := time.Time(storageInTransit.OutDate)

			sitObject := models.StorageInTransit{
				ID:              storageInTransit.ID,
				Status:          models.StorageInTransitStatus(storageInTransit.Status),
				ActualStartDate: &actualStartDate,
				OutDate:         &outDate,
			}

			storageInTransitsList[i] = payloadForStorageInTransitModel(&sitObject)
		}

		MoveQueueItemPayload := payloadForMoveQueueItem(MoveQueueItem, storageInTransitsList)
		MoveQueueItemPayloads[i] = MoveQueueItemPayload

	}
	return queueop.NewShowQueueOK().WithPayload(MoveQueueItemPayloads)
}
