package supportapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/support"
)

// ListMTOsHandler lists move task orders with the option to filter since a particular date
type ListMTOsHandler struct {
	handlers.HandlerConfig
	services.MoveTaskOrderFetcher
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h ListMTOsHandler) Handle(params movetaskorderops.ListMTOsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			searchParams := services.MoveTaskOrderFetcherParams{
				IncludeHidden: true,
			}
			if params.Since != nil {
				timeSince := time.Unix(*params.Since, 0)
				searchParams.Since = &timeSince
			}

			mtos, err := h.MoveTaskOrderFetcher.ListAllMoveTaskOrders(appCtx, &searchParams)

			if err != nil {
				appCtx.Logger().Error("Unable to fetch records:", zap.Error(err))
				return movetaskorderops.NewListMTOsInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.MoveTaskOrders(&mtos)

			return movetaskorderops.NewListMTOsOK().WithPayload(payload), nil
		})
}

// MakeMoveTaskOrderAvailableHandlerFunc updates the status of a Move Task Order
type MakeMoveTaskOrderAvailableHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderAvailabilityUpdater services.MoveTaskOrderUpdater
}

// Handle updates the prime availability of a MoveTaskOrder
func (h MakeMoveTaskOrderAvailableHandlerFunc) Handle(params movetaskorderops.MakeMoveTaskOrderAvailableParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			eTag := params.IfMatch

			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

			mto, err := h.moveTaskOrderAvailabilityUpdater.MakeAvailableToPrime(appCtx, moveTaskOrderID, eTag, false, false)

			if err != nil {
				appCtx.Logger().Error("supportapi.MakeMoveTaskOrderAvailableHandlerFunc error", zap.Error(err))
				switch typedErr := err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewMakeMoveTaskOrderAvailableNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return movetaskorderops.NewMakeMoveTaskOrderAvailableUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), typedErr.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewMakeMoveTaskOrderAvailablePreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return movetaskorderops.NewMakeMoveTaskOrderAvailableInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

			return movetaskorderops.NewMakeMoveTaskOrderAvailableOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// HideNonFakeMoveTaskOrdersHandlerFunc calls service to hide MTOs that are not using fake data
type HideNonFakeMoveTaskOrdersHandlerFunc struct {
	handlers.HandlerConfig
	services.MoveTaskOrderHider
}

// Handle hides any mto that doesnt have valid fake data
func (h HideNonFakeMoveTaskOrdersHandlerFunc) Handle(params movetaskorderops.HideNonFakeMoveTaskOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			hiddenMTOs, err := h.Hide(appCtx)
			if err != nil {
				appCtx.Logger().Error("supportapi.HideNonFakeMoveTaskOrdersHandlerFunc error", zap.Error(err))
				return movetaskorderops.NewHideNonFakeMoveTaskOrdersInternalServerError().WithPayload(
					payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}
			payload := payloads.MTOHideMovesResponse(hiddenMTOs)

			return movetaskorderops.NewHideNonFakeMoveTaskOrdersOK().WithPayload(payload), nil
		})
}

// GetMoveTaskOrderHandlerFunc returns the details for a particular Move Task Order
type GetMoveTaskOrderHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches an MTO from the database using its UUID
func (h GetMoveTaskOrderHandlerFunc) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
			searchParams := services.MoveTaskOrderFetcherParams{
				IncludeHidden:   true,
				MoveTaskOrderID: moveTaskOrderID,
			}
			mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(appCtx, &searchParams)
			if err != nil {
				appCtx.Logger().Error("primeapi.support.GetMoveTaskOrderHandler error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
			return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// CreateMoveTaskOrderHandler creates a move task order
type CreateMoveTaskOrderHandler struct {
	handlers.HandlerConfig
	moveTaskOrderCreator support.InternalMoveTaskOrderCreator
}

// Handle updates to move task order post-counseling
func (h CreateMoveTaskOrderHandler) Handle(params movetaskorderops.CreateMoveTaskOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveTaskOrder, err := h.moveTaskOrderCreator.InternalCreateMoveTaskOrder(appCtx, *params.Body)

			if err != nil {
				appCtx.Logger().Error("primeapi.support.CreateMoveTaskOrderHandler error", zap.Error(err))
				switch typedErr := err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewCreateMoveTaskOrderNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					errPayload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), typedErr.ValidationErrors)
					return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(errPayload), err
				case apperror.QueryError:
					// This error is generated when the validation passed but there was an error in creation
					// Usually this is due to a more complex dependency like a foreign key constraint
					return movetaskorderops.NewCreateMoveTaskOrderBadRequest().WithPayload(
						payloads.ClientError(handlers.SQLErrMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			moveTaskOrderPayload := payloads.MoveTaskOrder(moveTaskOrder)
			return movetaskorderops.NewCreateMoveTaskOrderCreated().WithPayload(moveTaskOrderPayload), nil
		})
}
