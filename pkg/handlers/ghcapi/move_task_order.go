package ghcapi

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services/event"

	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
)

// GetMoveTaskOrderHandler fetches a Move Task Order
type GetMoveTaskOrderHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches a single MoveTaskOrder
func (h GetMoveTaskOrderHandler) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		logger.Error("ghcapi.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound()
		case apperror.InvalidInputError:
			return movetaskorderops.NewGetMoveTaskOrderBadRequest()
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	eTag := params.IfMatch

	// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
	// TODO to revisit?
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	serviceItemCodes := ghcmessages.MTOApprovalServiceItemCodes{}
	if params.ServiceItemCodes != nil {
		serviceItemCodes = *params.ServiceItemCodes
	}

	mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(appCtx, moveTaskOrderID, eTag,
		serviceItemCodes.ServiceCodeMS, serviceItemCodes.ServiceCodeCS)

	if err != nil {
		logger.Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc error", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return movetaskorderops.NewUpdateMoveTaskOrderStatusUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.ConflictError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}

	moveTaskOrderPayload := payloads.Move(mto)

	// Audit attempt to make MTO available to prime
	_, err = audit.Capture(mto, moveTaskOrderPayload, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for making MTO available to Prime.", zap.Error(err))
		return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.MoveTaskOrderUpdateEventKey,
		MtoID:           mto.ID,
		UpdatedObjectID: mto.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdateMoveTaskOrderStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc could not generate the event")
	}

	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}

// UpdateMTOStatusServiceCounselingCompletedHandlerFunc updates the status of a Move (MoveTaskOrder) to MoveStatusServiceCounselingCompleted
type UpdateMTOStatusServiceCounselingCompletedHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a Move (MoveTaskOrder). Slightly different from UpdateMoveTaskOrderStatusHandlerFunc,
// this handler will update the Move status without making it available to the Prime and without creating basic service items.
func (h UpdateMTOStatusServiceCounselingCompletedHandlerFunc) Handle(params movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	eTag := params.IfMatch

	// TODO - Revisit authorization for Service Counselor role
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	mto, err := h.moveTaskOrderStatusUpdater.UpdateStatusServiceCounselingCompleted(appCtx, moveTaskOrderID, eTag)

	if err != nil {
		logger.Error("ghcapi.UpdateMTOStatusServiceCounselingCompletedHandlerFunc error", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.ConflictError:
			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedInternalServerError()
		}
	}

	moveTaskOrderPayload := payloads.Move(mto)

	// Audit
	_, err = audit.Capture(mto, moveTaskOrderPayload, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for transitioning Move status to Service Counseling Completed.", zap.Error(err))
		return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedInternalServerError()
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.MoveTaskOrderUpdateEventKey,
		MtoID:           mto.ID,
		UpdatedObjectID: mto.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdateMTOStatusServiceCounselingCompletedEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdateMTOStatusServiceCounselingCompletedHandlerFunc could not generate the event")
	}

	return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedOK().WithPayload(moveTaskOrderPayload)
}

// UpdateMTOReviewedBillableWeightsAtHandlerFunc provides timestamp for a Move's (MoveTaskOrder's) ReviewedBillableWeightsAt field
type UpdateMTOReviewedBillableWeightsAtHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the timestamp for a Move's (MoveTaskOrder's) ReviewedBillableWeightsAt field
func (h UpdateMTOReviewedBillableWeightsAtHandlerFunc) Handle(params movetaskorderops.UpdateMTOReviewedBillableWeightsAtParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	eTag := params.IfMatch

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	mto, err := h.moveTaskOrderStatusUpdater.UpdateReviewedBillableWeightsAt(appCtx, moveTaskOrderID, eTag)

	if err != nil {
		logger.Error("ghcapi.UpdateMTOReviewedBillableWeightsAtHandlerFunc error", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.ConflictError:
			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtInternalServerError()
		}
	}

	moveTaskOrderPayload := payloads.Move(mto)

	// Audit
	_, err = audit.Capture(mto, moveTaskOrderPayload, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error updating the move's billableWeightsReviewedAt field.", zap.Error(err))
		return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtInternalServerError()
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.MoveTaskOrderUpdateEventKey,
		MtoID:           mto.ID,
		UpdatedObjectID: mto.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdateMTOReviewedBillableWeightsEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdateMTOReviewedBillableWeightsAtHandlerFunc could not generate the event")
	}

	return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtOK().WithPayload(moveTaskOrderPayload)
}
