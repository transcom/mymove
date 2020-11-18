package ghcapi

import (
	"github.com/go-openapi/swag"

	"github.com/gobuffalo/pop/v5"

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

	moves, count, err := h.MoveOrderFetcher.ListMoveOrders(
		session.OfficeUserID,
		&listMoveOrderParams,
	)

	if err != nil {
		logger.Error("error fetching list of move orders for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(moves)

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

	listPaymentRequestParams := services.FetchPaymentRequestListParams{
		Branch:                 params.Branch,
		MoveID:                 params.MoveID,
		DodID:                  params.DodID,
		LastName:               params.LastName,
		DestinationDutyStation: params.DestinationDutyStation,
		Status:                 params.Status,
		Page:                   params.Page,
		PerPage:                params.PerPage,
		SubmittedAt:            params.SubmittedAt,
		Sort:                   params.Sort,
		Order:                  params.Order,
	}

	if params.Sort == nil {
		sortBy := "created_at"
		listPaymentRequestParams.Sort = &sortBy
	}

	if params.Order == nil {
		orderBy := true
		listPaymentRequestParams.Order = &orderBy
	}

	// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
	// for perPage.
	if params.Page == nil {
		listPaymentRequestParams.Page = swag.Int64(1)
	}
	// Same for perPage
	if params.PerPage == nil {
		listPaymentRequestParams.PerPage = swag.Int64(20)
	}

	paymentRequests, count, err := h.FetchPaymentRequestList(
		session.OfficeUserID,
		&listPaymentRequestParams,
	)
	if err != nil {
		logger.Error("payment requests queue", zap.String("office_user_id", session.OfficeUserID.String()), zap.Error(err))
		return queues.NewGetPaymentRequestsQueueInternalServerError()
	}

	queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests)

	result := &ghcmessages.QueuePaymentRequestsResult{
		TotalCount:           int64(count),
		Page:                 int64(*listPaymentRequestParams.Page),
		PerPage:              int64(*listPaymentRequestParams.PerPage),
		QueuePaymentRequests: *queuePaymentRequests,
		Sort:                 *listPaymentRequestParams.Sort,
		Order:                *listPaymentRequestParams.Order,
	}

	return queues.NewGetPaymentRequestsQueueOK().WithPayload(result)
}
