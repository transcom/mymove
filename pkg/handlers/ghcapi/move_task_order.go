package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// GetMoveTaskOrderHandler fetches a Move Task Order
type GetMoveTaskOrderHandler struct {
	handlers.HandlerConfig
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches a single MoveTaskOrder
func (h GetMoveTaskOrderHandler) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

			searchParams := services.MoveTaskOrderFetcherParams{
				IncludeHidden:   false,
				MoveTaskOrderID: moveTaskOrderID,
			}
			mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(appCtx, &searchParams)
			if err != nil {
				appCtx.Logger().Error("ghcapi.GetMoveTaskOrderHandler error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewGetMoveTaskOrderNotFound(), err
				case apperror.InvalidInputError:
					return movetaskorderops.NewGetMoveTaskOrderBadRequest(), err
				default:
					return movetaskorderops.NewGetMoveTaskOrderInternalServerError(), err
				}
			}
			moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
			return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			eTag := params.IfMatch

			// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
			// TODO to revisit?
			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

			serviceItemCodes := ghcmessages.MTOApprovalServiceItemCodes{}
			if params.ServiceItemCodes != nil {
				serviceItemCodes = *params.ServiceItemCodes
			}

			checker := movetaskorder.NewMoveTaskOrderChecker()
			availableBefore, err := checker.MTOAvailableToPrime(appCtx, moveTaskOrderID)
			if err != nil {
				return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError(), err
			}

			mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(appCtx, moveTaskOrderID, eTag,
				serviceItemCodes.ServiceCodeMS, serviceItemCodes.ServiceCodeCS)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Unable to complete request",
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						validate.NewErrors())
					return movetaskorderops.NewUpdateMoveTaskOrderStatusUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ConflictError:
					return movetaskorderops.NewUpdateMoveTaskOrderStatusConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError(), err
				}
			}

			if !availableBefore {
				availableAfter, checkErr := checker.MTOAvailableToPrime(appCtx, moveTaskOrderID)
				if checkErr != nil {
					return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError(), err
				}

				/* Do not send TOO approving and submitting service items email if BLUEBARK/SAFETY */
				if availableAfter && mto.Orders.CanSendEmailWithOrdersType() {
					emailErr := h.NotificationSender().SendNotification(appCtx,
						notifications.NewMoveIssuedToPrime(moveTaskOrderID),
					)
					if emailErr != nil {
						return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError(), err
					}
				}
			}

			moveTaskOrderPayload := payloads.Move(mto)

			// Audit attempt to make MTO available to prime
			_, err = audit.Capture(appCtx, mto, moveTaskOrderPayload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Auditing service error for making MTO available to Prime.", zap.Error(err))
				return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError(), err
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.MoveTaskOrderUpdateEventKey,
				MtoID:           mto.ID,
				UpdatedObjectID: mto.ID,
				EndpointKey:     event.GhcUpdateMoveTaskOrderStatusEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc could not generate the event")
			}

			return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// UpdateMTOStatusServiceCounselingCompletedHandlerFunc updates the status of a Move (MoveTaskOrder) to MoveStatusServiceCounselingCompleted
type UpdateMTOStatusServiceCounselingCompletedHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a Move (MoveTaskOrder). Slightly different from UpdateMoveTaskOrderStatusHandlerFunc,
// this handler will update the Move status without making it available to the Prime and without creating basic service items.
func (h UpdateMTOStatusServiceCounselingCompletedHandlerFunc) Handle(params movetaskorderops.UpdateMTOStatusServiceCounselingCompletedParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateMTOStatusServiceCounselingCompletedHandlerFunc error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Unable to complete request",
						err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest),
						validate.NewErrors())
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedUnprocessableEntity().
						WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ConflictError:
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				return handleError(apperror.NewForbiddenError("is not a Services Counselor"))
			}

			eTag := params.IfMatch
			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
			mto, err := h.moveTaskOrderStatusUpdater.UpdateStatusServiceCounselingCompleted(appCtx, moveTaskOrderID, eTag)

			if err != nil {
				return handleError(err)
			}

			moveTaskOrderPayload := payloads.Move(mto)

			// Audit
			_, err = audit.Capture(appCtx, mto, moveTaskOrderPayload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Auditing service error for transitioning Move status to Service Counseling Completed.", zap.Error(err))
				return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedInternalServerError(), err
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.MoveTaskOrderUpdateEventKey,
				MtoID:           mto.ID,
				UpdatedObjectID: mto.ID,
				EndpointKey:     event.GhcUpdateMTOStatusServiceCounselingCompletedEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOStatusServiceCounselingCompletedHandlerFunc could not generate the event")
			}

			/* Do not send SC Move Details Submitted email if orders type is BLUEBARK/SAFETY */
			if mto.Orders.CanSendEmailWithOrdersType() {
				err = h.NotificationSender().SendNotification(appCtx, notifications.NewMoveCounseled(moveTaskOrderID))
				if err != nil {
					appCtx.Logger().Error("problem sending email to user", zap.Error(err))
				}
			}

			return movetaskorderops.NewUpdateMTOStatusServiceCounselingCompletedOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// UpdateMTOReviewedBillableWeightsAtHandlerFunc provides timestamp for a Move's (MoveTaskOrder's) ReviewedBillableWeightsAt field
type UpdateMTOReviewedBillableWeightsAtHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the timestamp for a Move's (MoveTaskOrder's) ReviewedBillableWeightsAt field
func (h UpdateMTOReviewedBillableWeightsAtHandlerFunc) Handle(params movetaskorderops.UpdateMTOReviewedBillableWeightsAtParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			eTag := params.IfMatch

			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

			mto, err := h.moveTaskOrderStatusUpdater.UpdateReviewedBillableWeightsAt(appCtx, moveTaskOrderID, eTag)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOReviewedBillableWeightsAtHandlerFunc error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtPreconditionFailed().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				case apperror.ConflictError:
					return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtConflict().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				default:
					return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtInternalServerError(), err
				}
			}

			moveTaskOrderPayload := payloads.Move(mto)

			// Audit
			_, err = audit.Capture(appCtx, mto, moveTaskOrderPayload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Auditing service error updating the move's billableWeightsReviewedAt field.", zap.Error(err))
				return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtInternalServerError(), err
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.MoveTaskOrderUpdateEventKey,
				MtoID:           mto.ID,
				UpdatedObjectID: mto.ID,
				EndpointKey:     event.GhcUpdateMTOReviewedBillableWeightsEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOReviewedBillableWeightsAtHandlerFunc could not generate the event")
			}

			return movetaskorderops.NewUpdateMTOReviewedBillableWeightsAtOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// UpdateMoveTIORemarksHandlerFunc updates a Move's (MoveTaskOrder's) TIORemarks field
type UpdateMoveTIORemarksHandlerFunc struct {
	handlers.HandlerConfig
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates a Move's (MoveTaskOrder's) TIORemarks field
func (h UpdateMoveTIORemarksHandlerFunc) Handle(params movetaskorderops.UpdateMoveTIORemarksParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			eTag := params.IfMatch

			moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
			remarks := params.Body.TioRemarks

			mto, err := h.moveTaskOrderStatusUpdater.UpdateTIORemarks(appCtx, moveTaskOrderID, eTag, *remarks)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMoveTIORemarksHandlerFunc error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewUpdateMoveTIORemarksNotFound(), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewUpdateMoveTIORemarksPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return movetaskorderops.NewUpdateMoveTIORemarksInternalServerError(), err
				}
			}

			moveTaskOrderPayload := payloads.Move(mto)

			// Audit
			_, err = audit.Capture(appCtx, mto, moveTaskOrderPayload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Auditing service error updating the move's TioRemarks field.", zap.Error(err))
				return movetaskorderops.NewUpdateMoveTIORemarksInternalServerError(), err
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.MoveTaskOrderUpdateEventKey,
				MtoID:           mto.ID,
				UpdatedObjectID: mto.ID,
				EndpointKey:     event.GhcUpdateMoveTIORemarksEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMoveTIORemarksHandlerFunc could not generate the event")
			}

			return movetaskorderops.NewUpdateMoveTIORemarksOK().WithPayload(moveTaskOrderPayload), nil
		})
}
