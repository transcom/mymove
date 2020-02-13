package internalapi

import (
	"sort"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForMoveQueueItem(MoveQueueItem models.MoveQueueItem) *internalmessages.MoveQueueItem {
	MoveQueueItemPayload := internalmessages.MoveQueueItem{
		ID:                         handlers.FmtUUID(MoveQueueItem.ID),
		CreatedAt:                  handlers.FmtDateTime(MoveQueueItem.CreatedAt),
		Edipi:                      swag.String(MoveQueueItem.Edipi),
		Rank:                       (*internalmessages.ServiceMemberRank)(MoveQueueItem.Rank),
		CustomerName:               swag.String(MoveQueueItem.CustomerName),
		Locator:                    swag.String(MoveQueueItem.Locator),
		Status:                     swag.String(MoveQueueItem.Status),
		PpmStatus:                  handlers.FmtStringPtr(MoveQueueItem.PpmStatus),
		OrdersType:                 swag.String(MoveQueueItem.OrdersType),
		MoveDate:                   handlers.FmtDatePtr(MoveQueueItem.MoveDate),
		SubmittedDate:              handlers.FmtDateTimePtr(MoveQueueItem.SubmittedDate),
		LastModifiedDate:           handlers.FmtDateTime(MoveQueueItem.LastModifiedDate),
		OriginDutyStationName:      swag.String(MoveQueueItem.OriginDutyStationName),
		DestinationDutyStationName: swag.String(MoveQueueItem.DestinationDutyStationName),
		PmSurveyConductedDate:      handlers.FmtDateTimePtr(MoveQueueItem.PmSurveyConductedDate),
		OriginGbloc:                handlers.FmtStringPtr(MoveQueueItem.OriginGBLOC),
		DestinationGbloc:           handlers.FmtStringPtr(MoveQueueItem.DestinationGBLOC),
		DeliveredDate:              handlers.FmtDateTimePtr(MoveQueueItem.DeliveredDate),
		InvoiceApprovedDate:        handlers.FmtDateTimePtr(MoveQueueItem.InvoiceApprovedDate),
		WeightAllotment:            payloadForWeightAllotmentModel(models.GetWeightAllotment(*MoveQueueItem.Rank)),
		BranchOfService:            handlers.FmtString(MoveQueueItem.BranchOfService),
		ActualMoveDate:             handlers.FmtDatePtr(MoveQueueItem.ActualMoveDate),
		OriginalMoveDate:           handlers.FmtDatePtr(MoveQueueItem.OriginalMoveDate),
	}
	return &MoveQueueItemPayload
}

// ShowQueueHandler returns a list of all MoveQueueItems in the moves queue
type ShowQueueHandler struct {
	handlers.HandlerContext
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
	Location        string    `json:"location"`
}

// Implementation of a type and methods in order to use sort.Interface directly.
// This allows us to call sortQueueItemsByLastModifiedDate in the ShowQueueHandler which will
// sort the slice by the LastModfiedDate. Doing it this way allows us to avoid having reflect called
// which should act to speed the sort up.
type MoveQueueItems []models.MoveQueueItem

func (mqi MoveQueueItems) Less(i, j int) bool {
	return mqi[i].LastModifiedDate.Before(mqi[j].LastModifiedDate)
}
func (mqi MoveQueueItems) Len() int      { return len(mqi) }
func (mqi MoveQueueItems) Swap(i, j int) { mqi[i], mqi[j] = mqi[j], mqi[i] }

func sortQueueItemsByLastModifiedDate(moveQueueItems []models.MoveQueueItem) {
	sort.Sort(MoveQueueItems(moveQueueItems))
}

// Handle retrieves a list of all MoveQueueItems in the system in the moves queue
func (h ShowQueueHandler) Handle(params queueop.ShowQueueParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return queueop.NewShowQueueForbidden()
	}

	lifecycleState := params.QueueType

	MoveQueueItems, err := models.GetMoveQueueItems(h.DB(), lifecycleState)
	if err != nil {
		logger.Error("Loading Queue", zap.String("State", lifecycleState), zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	// Sorting the slice by LastModifiedDate so that the API results follow suit.
	sortQueueItemsByLastModifiedDate(MoveQueueItems)

	MoveQueueItemPayloads := make([]*internalmessages.MoveQueueItem, len(MoveQueueItems))
	for i, MoveQueueItem := range MoveQueueItems {
		MoveQueueItemPayload := payloadForMoveQueueItem(MoveQueueItem)
		MoveQueueItemPayloads[i] = MoveQueueItemPayload
	}
	return queueop.NewShowQueueOK().WithPayload(MoveQueueItemPayloads)
}
