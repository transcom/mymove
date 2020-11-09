package ghcapi

import (
	"fmt"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models/roles"
)

// GetMovesQueueHandler returns the moves for the TOO queue user via GET /queues/moves
type GetMovesQueueHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// FilterOption defines the type for the functional arguments used for private functions in MoveOrderFetcher
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("user is not authenticated with TOO office role")
		return queues.NewGetMovesQueueForbidden()
	}

	listMoveOrderParams := services.ListMoveOrderParams{
		Branch:                 params.Branch,
		MoveID:                 params.MoveID,
		DodID:                  params.DodID,
		LastName:               params.LastName,
		DestinationDutyStation: params.DestinationDutyStation,
		Status:                 params.Status,
		Page:                   params.Page,
		PerPage:                params.PerPage,
	}

	// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
	// for perPage.
	if params.Page == nil {
		listMoveOrderParams.Page = swag.Int64(1)
	}
	// Same for perPage
	if params.PerPage == nil {
		listMoveOrderParams.PerPage = swag.Int64(20)
	}

	orders, count, err := h.MoveOrderFetcher.ListMoveOrders(
		session.OfficeUserID,
		&listMoveOrderParams,
	)

	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(orders)

	result := &ghcmessages.QueueMovesResult{
		Page:       *listMoveOrderParams.Page,
		PerPage:    *listMoveOrderParams.PerPage,
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

// These are a bunch of private functions that are used to cobbble our list MoveOrders filters together.
// TODO: These are moving to the service objects, to TOO queue now and will happen for the TIO queue in: https://github.com/transcom/mymove/pull/5162
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
