package ghcapi

import (
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
	services.OrderFetcher
}

// FilterOption defines the type for the functional arguments used for private functions in OrderFetcher
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		appCtx.Logger().Error("user is not authenticated with TOO office role")
		return queues.NewGetMovesQueueForbidden()
	}

	ListOrderParams := services.ListOrderParams{
		Branch:                  params.Branch,
		Locator:                 params.Locator,
		DodID:                   params.DodID,
		LastName:                params.LastName,
		DestinationDutyLocation: params.DestinationDutyLocation,
		OriginDutyLocation:      params.OriginDutyLocation,
		Status:                  params.Status,
		Page:                    params.Page,
		PerPage:                 params.PerPage,
		Sort:                    params.Sort,
		Order:                   params.Order,
	}

	// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
	// for perPage.
	if params.Page == nil {
		ListOrderParams.Page = swag.Int64(1)
	}
	// Same for perPage
	if params.PerPage == nil {
		ListOrderParams.PerPage = swag.Int64(20)
	}

	moves, count, err := h.OrderFetcher.ListOrders(
		appCtx,
		appCtx.Session().OfficeUserID,
		&ListOrderParams,
	)

	if err != nil {
		appCtx.Logger().Error("error fetching list of moves for office user", zap.Error(err))
		return queues.NewGetMovesQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(moves)

	result := &ghcmessages.QueueMovesResult{
		Page:       *ListOrderParams.Page,
		PerPage:    *ListOrderParams.PerPage,
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

	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().Roles.HasRole(roles.RoleTypeTIO) {
		return queues.NewGetPaymentRequestsQueueForbidden()
	}

	listPaymentRequestParams := services.FetchPaymentRequestListParams{
		Branch:                 params.Branch,
		Locator:                params.Locator,
		DodID:                  params.DodID,
		LastName:               params.LastName,
		DestinationDutyStation: params.DestinationDutyLocation,
		Status:                 params.Status,
		Page:                   params.Page,
		PerPage:                params.PerPage,
		SubmittedAt:            handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
		Sort:                   params.Sort,
		Order:                  params.Order,
		OriginDutyLocation:     params.OriginDutyLocation,
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
		appCtx,
		appCtx.Session().OfficeUserID,
		&listPaymentRequestParams,
	)
	if err != nil {
		appCtx.Logger().Error("payment requests queue", zap.String("office_user_id", appCtx.Session().OfficeUserID.String()), zap.Error(err))
		return queues.NewGetPaymentRequestsQueueInternalServerError()
	}

	queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests)

	result := &ghcmessages.QueuePaymentRequestsResult{
		TotalCount:           int64(count),
		Page:                 int64(*listPaymentRequestParams.Page),
		PerPage:              int64(*listPaymentRequestParams.PerPage),
		QueuePaymentRequests: *queuePaymentRequests,
	}

	return queues.NewGetPaymentRequestsQueueOK().WithPayload(result)
}

// GetServicesCounselingQueueHandler returns the moves for the Service Counselor queue user via GET /queues/counselor
type GetServicesCounselingQueueHandler struct {
	handlers.HandlerContext
	services.OrderFetcher
}

// Handle returns the paginated list of moves for the services counselor
func (h GetServicesCounselingQueueHandler) Handle(params queues.GetServicesCounselingQueueParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
		appCtx.Logger().Error("user is not authenticated with an office role")
		return queues.NewGetServicesCounselingQueueForbidden()
	}

	ListOrderParams := services.ListOrderParams{
		Branch:             params.Branch,
		Locator:            params.Locator,
		DodID:              params.DodID,
		LastName:           params.LastName,
		OriginDutyLocation: params.OriginDutyLocation,
		OriginGBLOC:        params.OriginGBLOC,
		SubmittedAt:        handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
		RequestedMoveDate:  params.RequestedMoveDate,
		Page:               params.Page,
		PerPage:            params.PerPage,
		Sort:               params.Sort,
		Order:              params.Order,
	}

	if len(params.Status) == 0 {
		ListOrderParams.Status = []string{string(models.MoveStatusNeedsServiceCounseling)}
	} else {
		ListOrderParams.Status = params.Status
	}

	// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
	// for perPage.
	if params.Page == nil {
		ListOrderParams.Page = swag.Int64(1)
	}
	// Same for perPage
	if params.PerPage == nil {
		ListOrderParams.PerPage = swag.Int64(20)
	}

	moves, count, err := h.OrderFetcher.ListOrders(
		appCtx,
		appCtx.Session().OfficeUserID,
		&ListOrderParams,
	)

	if err != nil {
		appCtx.Logger().Error("error fetching list of moves for office user", zap.Error(err))
		return queues.NewGetServicesCounselingQueueInternalServerError()
	}

	queueMoves := payloads.QueueMoves(moves)

	result := &ghcmessages.QueueMovesResult{
		Page:       *ListOrderParams.Page,
		PerPage:    *ListOrderParams.PerPage,
		TotalCount: int64(count),
		QueueMoves: *queueMoves,
	}

	return queues.NewGetServicesCounselingQueueOK().WithPayload(result)
}
