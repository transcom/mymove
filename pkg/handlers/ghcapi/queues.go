package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// GetMovesQueueHandler returns the moves for the TOO queue user via GET /queues/moves
type GetMovesQueueHandler struct {
	handlers.HandlerContext
	services.OfficeUserFetcher
	services.MoveOrderFetcher
}

// FilterOption defines the type for the functional arguments passed to ListMoveOrders
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("user is not authenticated with TOO office role")
		return queues.NewGetMovesQueueForbidden()
	}

	branchQuery := branchFilter(params)

	orders, err := h.MoveOrderFetcher.ListMoveOrders(session.OfficeUserID, branchQuery)
	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(orders)
	// ToDo - May want to move this logic into the pop query later.
	// filter queueMoves by status
	queueMoves = statusFilter(params.Status, queueMoves)

	result := &ghcmessages.QueueMovesResult{
		Page:       0,
		PerPage:    0,
		TotalCount: int64(len(*queueMoves)),
		QueueMoves: *queueMoves,
	}

	return queues.NewGetMovesQueueOK().WithPayload(result)
}

// GetPaymentRequestsQueueHandler returns the payment requests for the TIO queue user via GET /queues/payment-requests
type GetPaymentRequestsQueueHandler struct {
	handlers.HandlerContext
	services.OfficeUserFetcher
	services.PaymentRequestListFetcher
}

// Handle returns the paginated list of payment requests for the TIO user
func (h GetPaymentRequestsQueueHandler) Handle(params queues.GetPaymentRequestsQueueParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.Roles.HasRole(roles.RoleTypeTIO) {
		return queues.NewGetPaymentRequestsQueueForbidden()
	}

	paymentRequests, err := h.FetchPaymentRequestList(session.OfficeUserID)
	if err != nil {
		logger.Error("payment requests queue", zap.String("office_user_id", session.OfficeUserID.String()), zap.Error(err))
		return queues.NewGetPaymentRequestsQueueInternalServerError()
	}

	queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests)

	result := &ghcmessages.QueuePaymentRequestsResult{
		TotalCount:           int64(len(*queuePaymentRequests)),
		QueuePaymentRequests: *queuePaymentRequests,
	}

	return queues.NewGetPaymentRequestsQueueOK().WithPayload(result)
}

func branchFilter(params queues.GetMovesQueueParams) FilterOption {
	return func(query *pop.Query) {
		if params.Branch != nil {
			query = query.Where("orders.department_indicator = ?", *params.Branch)
		}
	}
}

// statusFilter filters the status after the pop query call.
func statusFilter(statuses []string, moves *ghcmessages.QueueMoves) *ghcmessages.QueueMoves {
	if len(statuses) <= 0 || moves == nil {
		return moves
	}

	ret := make(ghcmessages.QueueMoves, 0)
	// New move, Approvals requested, and Move approved statuses
	// convert into a map to make it easier to lookup
	statusMap := make(map[string]string, 0)
	for _, status := range statuses {
		statusMap[status] = status
	}

	// then include only the moves based on status filter
	// and exclude DRAFT and CANCELLED
	for _, move := range *moves {
		if _, ok := statusMap[string(move.Status)]; ok && string(move.Status) != string(models.MoveStatusCANCELED) &&
			string(move.Status) != string(models.MoveStatusDRAFT) {
			ret = append(ret, move)
		}
	}

	return &ret
}
