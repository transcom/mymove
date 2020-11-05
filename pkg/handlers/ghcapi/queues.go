package ghcapi

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// GetMovesQueueHandler returns the moves for the TOO queue user via GET /queues/moves
type GetMovesQueueHandler struct {
	handlers.HandlerContext
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
	// If we get a page argument let's pull it out and cast it to the expected int type.
	var page int
	if params.Page != nil {
		page = int(*params.Page)
	}

	var perPage int
	if params.PerPage != nil {
		perPage = int(*params.PerPage)
	}

	branchQuery := branchFilter(params.Branch)
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	moveStatusQuery := moveStatusFilter(params.Status)

	orders, count, err := h.MoveOrderFetcher.ListMoveOrders(
		session.OfficeUserID,
		&page,
		&perPage,
		branchQuery,
		moveIDQuery,
		lastNameQuery,
		dutyStationQuery,
		dodIDQuery,
		moveStatusQuery,
	)

	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(orders)
	// ToDo - May want to move this logic into the pop query later.
	// filter queueMoves by status
	queueMovesBeforeStatusFilterCount := len(*queueMoves)
	queueMoves = movesFilteredByStatus(params.Status, queueMoves)
	queueMovesAfterStatusFilterCount := len(*queueMoves)
	count = count - (queueMovesBeforeStatusFilterCount - queueMovesAfterStatusFilterCount)

	result := &ghcmessages.QueueMovesResult{
		Page:       int64(page),
		PerPage:    int64(perPage),
		TotalCount: int64(count),
		QueueMoves: *queueMoves,
	}

	return queues.NewGetMovesQueueOK().WithPayload(result)
}

// GetPaymentRequestsQueueHandler returns the payment requests for the TIO queue user via GET /queues/payment-requests
type GetPaymentRequestsQueueHandler struct {
	handlers.HandlerContext
	services.PaymentRequestListFetcher
}

// Handle returns the paginated list of payment requests for the TIO user
func (h GetPaymentRequestsQueueHandler) Handle(params queues.GetPaymentRequestsQueueParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.Roles.HasRole(roles.RoleTypeTIO) {
		return queues.NewGetPaymentRequestsQueueForbidden()
	}

	branchQuery := branchFilter(params.Branch)
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	statusQuery := paymentRequestsStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)

	paymentRequests, err := h.FetchPaymentRequestList(
		session.OfficeUserID,
		statusQuery,
		branchQuery,
		moveIDQuery,
		lastNameQuery,
		dutyStationQuery,
		dodIDQuery,
		submittedAtQuery,
	)
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

func branchFilter(branch *string) FilterOption {
	return func(query *pop.Query) {
		if branch != nil {
			query = query.Where("service_members.affiliation = ?", *branch)
		}
	}
}

func lastNameFilter(lastName *string) FilterOption {
	return func(query *pop.Query) {
		if lastName != nil {
			nameSearch := fmt.Sprintf("%s%%", *lastName)
			query = query.Where("service_members.last_name ILIKE ?", nameSearch)
		}
	}
}

func dodIDFilter(dodID *string) FilterOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query = query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func moveIDFilter(moveID *string) FilterOption {
	return func(query *pop.Query) {
		if moveID != nil {
			query = query.Where("moves.locator = ?", *moveID)
		}
	}
}
func destinationDutyStationFilter(destinationDutyStation *string) FilterOption {
	return func(query *pop.Query) {
		if destinationDutyStation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyStation)
			query = query.InnerJoin("duty_stations as destination_duty_station", "orders.new_duty_station_id = destination_duty_station.id").Where("destination_duty_station.name ILIKE ?", nameSearch)
		}
	}
}

func submittedAtFilter(submittedAt *string) FilterOption {
	return func(query *pop.Query) {
		if submittedAt != nil {
			query = query.Where("CAST(payment_requests.created_at AS DATE) = ?", *submittedAt)
		}
	}
}

func moveStatusFilter(statuses []string) FilterOption {
	return func(query *pop.Query) {
		if len(statuses) <= 0 {
			query = query.Where("moves.status NOT IN (?)", models.MoveStatusDRAFT, models.MoveStatusCANCELED)
		}
	}
}

// movesFilteredByStatus filters the status after the pop query call.
func movesFilteredByStatus(statuses []string, moves *ghcmessages.QueueMoves) *ghcmessages.QueueMoves {
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

// statusFilter filters the status after the pop query call.
func paymentRequestsStatusFilter(statuses []string) FilterOption {
	return func(query *pop.Query) {
		var translatedStatuses []string
		if len(statuses) > 0 {
			for _, status := range statuses {
				if strings.EqualFold(status, "Payment requested") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPending.String())

				}
				if strings.EqualFold(status, "reviewed") {
					translatedStatuses = append(translatedStatuses,
						models.PaymentRequestStatusReviewed.String(),
						models.PaymentRequestStatusSentToGex.String(),
						models.PaymentRequestStatusReceivedByGex.String())
				}
			}
			query = query.Where("payment_requests.status in (?)", translatedStatuses)
		}
	}

}
