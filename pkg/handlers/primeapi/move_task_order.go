package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/audit"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

// FetchMTOUpdatesHandler lists move task orders with the option to filter since a particular date
type FetchMTOUpdatesHandler struct {
	handlers.HandlerContext
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h FetchMTOUpdatesHandler) Handle(params movetaskorderops.FetchMTOUpdatesParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Where("is_available_to_prime = ?", true).Eager("PaymentRequests")
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError()
	}

	payload := payloads.MoveTaskOrders(&mtos)

	for _, mto := range mtos {
		// Audit attempt for prime to fetch move task orders
		_, err = audit.Capture(mto, nil, logger, session, params.HTTPRequest)
		if err != nil {
			logger.Error("Auditing service error for fetching MTO for Prime", zap.Error(err))
			return movetaskorderops.NewFetchMTOUpdatesInternalServerError()
		}
	}

	return movetaskorderops.NewFetchMTOUpdatesOK().WithPayload(payload)
}
