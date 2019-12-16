package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

// ListMoveTaskOrdersHandler lists move task orders with the option to filter since a particular date
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h ListMoveTaskOrdersHandler) Handle(params movetaskorderops.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Q()
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewListMoveTaskOrdersOK().WithPayload(payload)
}
