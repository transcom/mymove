package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

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

// FilterOption allows ListMoveOrders to pass in a number of functions
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("user is not authenticated with TOO office role")
		return queues.NewGetMovesQueueForbidden()
	}

	moveIDQuery := moveIDFilter(params)
	firstNameQuery := firstNameFilter(params)
	lastNameQuery := lastNameFilter(params)
	dutyStationQuery := destinationDutyStationFilter(params)

	orders, err := h.MoveOrderFetcher.ListMoveOrders(
		session.OfficeUserID,
		moveIDQuery,
		firstNameQuery,
		lastNameQuery,
		dutyStationQuery,
	)

	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(orders)

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

func lastNameFilter(params queues.GetMovesQueueParams) FilterOption {
	return func(query *pop.Query) {
		if params.LastName != nil {
			nameSearch := fmt.Sprintf("%s%%", *params.LastName)
			query = query.Where("service_members.last_name ILIKE ?", nameSearch)
		}
	}
}

func firstNameFilter(params queues.GetMovesQueueParams) FilterOption {
	return func(query *pop.Query) {
		if params.FirstName != nil {
			nameSearch := fmt.Sprintf("%s%%", *params.FirstName)
			query = query.Where("service_members.first_name ILIKE ?", nameSearch)
		}
	}
}

//func dodIDFilter(params queues.GetMovesQueueParams) FilterOption {
//	return func(query *pop.Query) {
//		if params.DodID != nil {
//			query = query.Where("service_members.dod_id ILIKE ?", params.DodID)
//		}
//	}
//}

func moveIDFilter(params queues.GetMovesQueueParams) FilterOption {
	return func(query *pop.Query) {
		if params.MoveID != nil {
			query = query.Where("moves.locator ILIKE ?", *params.MoveID)
		}
	}
}
func destinationDutyStationFilter(params queues.GetMovesQueueParams) FilterOption {
	return func(query *pop.Query) {
		if params.DestinationDutyStation != nil {
			nameSearch := fmt.Sprintf("%s%%", *params.DestinationDutyStation)
			query = query.Where("destination_duty_station.name ILIKE ?", nameSearch)
		}
	}
}
