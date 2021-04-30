package primeapi

import (
	"fmt"

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

	searchParams := services.MoveTaskOrderFetcherParams{
		IsAvailableToPrime: true,
		Since:              params.Since,
	}
	mtos, err := h.MoveTaskOrderFetcher.ListAllMoveTaskOrders(&searchParams)

	if err != nil {
		logger.Error("Unexpected error while fetching records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewFetchMTOUpdatesOK().WithPayload(payload)
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

// Handle fetches an MTO from the database using its UUID
func (h GetMoveTaskOrderHandlerFunc) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	searchParams := services.MoveTaskOrderFetcherParams{
		IsAvailableToPrime: true,
	}
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID, &searchParams)
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
	mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOPostCounselingInformationHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(mtoID)

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

	mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(mtoID, params.Body, eTag)
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
