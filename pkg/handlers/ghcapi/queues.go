package ghcapi

import (
	"fmt"

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

	branchQuery := branchFilter(params.Branch)
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)

	orders, err := h.MoveOrderFetcher.ListMoveOrders(
		session.OfficeUserID,
		branchQuery,
		moveIDQuery,
		lastNameQuery,
		dutyStationQuery,
		dodIDQuery,
	)

	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(orders)
	// ToDo - May want to move this logic into the pop query later.
	// filter queueMoves by status
	queueMoves = moveStatusFilter(params.Status, queueMoves)

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

	branchQuery := branchFilter(params.Branch)
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)

	paymentRequests, err := h.FetchPaymentRequestList(
		session.OfficeUserID,
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

	queuePaymentRequests = paymentRequestsStatusFilter(params.Status, queuePaymentRequests)

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
			// some datetime conversion to compare YYYY-MM-DD to DateTime which may involve translating DateTime to Date
			query = query.Where("CAST(payment_requests.created_at AS DATE) = ?", *submittedAt)
		}
	}
}

// statusFilter filters the status after the pop query call.
func moveStatusFilter(statuses []string, moves *ghcmessages.QueueMoves) *ghcmessages.QueueMoves {
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
func paymentRequestsStatusFilter(statuses []string, paymentRequests *ghcmessages.QueuePaymentRequests) *ghcmessages.QueuePaymentRequests {
	if len(statuses) <= 0 || paymentRequests == nil {
		return paymentRequests
	}

	ret := make(ghcmessages.QueuePaymentRequests, 0)
	// New move, Approvals requested, and Move approved statuses
	// convert into a map to make it easier to lookup
	statusMap := make(map[string]string, 0)
	for _, status := range statuses {
		statusMap[status] = status
	}

	// then include only the moves based on status filter
	// and exclude DRAFT and CANCELLED
	for _, paymentRequest := range *paymentRequests {
		if _, ok := statusMap[string(paymentRequest.Status)]; ok && string(paymentRequest.Status) != string(models.MoveStatusCANCELED) &&
			string(paymentRequest.Status) != string(models.MoveStatusDRAFT) {
			ret = append(ret, paymentRequest)
		}
	}

	return &ret
}
