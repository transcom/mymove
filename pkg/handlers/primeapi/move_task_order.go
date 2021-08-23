package primeapi

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

// FetchMTOUpdatesHandler lists move task orders with the option to filter since a particular date
type FetchMTOUpdatesHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h FetchMTOUpdatesHandler) Handle(params movetaskorderops.FetchMTOUpdatesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	searchParams := services.MoveTaskOrderFetcherParams{
		IsAvailableToPrime: true,
	}
	if params.Since != nil {
		timeSince := time.Unix(*params.Since, 0)
		searchParams.Since = &timeSince
	}

	mtos, err := h.MoveTaskOrderFetcher.ListAllMoveTaskOrders(appCtx, &searchParams)

	if err != nil {
		logger.Error("Unexpected error while fetching records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewFetchMTOUpdatesOK().WithPayload(payload)
}

// ListMovesHandler lists move task orders with the option to filter since a particular date. Optimized ver.
type ListMovesHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle fetches all move task orders with the option to filter since a particular date. Optimized version.
func (h ListMovesHandler) Handle(params movetaskorderops.ListMovesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	var searchParams services.MoveTaskOrderFetcherParams
	if params.Since != nil {
		since := handlers.FmtDateTimePtrToPop(params.Since)
		searchParams.Since = &since
	}

	mtos, err := h.MoveTaskOrderFetcher.ListPrimeMoveTaskOrders(appCtx, &searchParams)

	if err != nil {
		logger.Error("Unexpected error while fetching moves:", zap.Error(err))
		return movetaskorderops.NewListMovesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}

	payload := payloads.ListMoves(&mtos)

	return movetaskorderops.NewListMovesOK().WithPayload(payload)
}

// UpdateMTOPostCounselingInformationHandler updates the move task order with post-counseling information
type UpdateMTOPostCounselingInformationHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MoveTaskOrderUpdater
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// GetMoveTaskOrderHandlerFunc returns the details for a particular Move Task Order
type GetMoveTaskOrderHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches an MTO from the database using its UUID or move code
func (h GetMoveTaskOrderHandlerFunc) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	searchParams := services.MoveTaskOrderFetcherParams{
		IsAvailableToPrime: true,
	}

	// Add either ID or Locator to search params
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveID)
	if moveTaskOrderID != uuid.Nil {
		searchParams.MoveTaskOrderID = moveTaskOrderID
	} else {
		searchParams.Locator = params.MoveID
	}

	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		logger.Error("primeapi.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

	return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}

// Handle updates to move task order post-counseling
func (h UpdateMTOPostCounselingInformationHandler) Handle(params movetaskorderops.UpdateMTOPostCounselingInformationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOPostCounselingInformationHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, mtoID)

	if err != nil {
		logger.Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
		return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity().WithPayload(
			payloads.ValidationError(err.Error(), h.GetTraceID(), nil))
	}

	if !mtoAvailableToPrime {
		logger.Error("primeapi.UpdateMTOPostCounselingInformationHandler error - MTO is not available to Prime")
		return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound().WithPayload(payloads.ClientError(
			handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", mtoID), h.GetTraceID()))
	}

	mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(appCtx, mtoID, params.Body, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), e.ValidationErrors))
		default:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}
	mtoPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload)
}
