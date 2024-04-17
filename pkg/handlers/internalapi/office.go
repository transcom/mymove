package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() {
				return officeop.NewApproveMoveForbidden(), apperror.NewForbiddenError("user must be office user")
			}

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// Don't approve Move if orders are incomplete
			orders, ordersErr := models.FetchOrder(appCtx.DB(), move.OrdersID)
			if ordersErr != nil {
				return handlers.ResponseForError(appCtx.Logger(), ordersErr), ordersErr
			}
			if !orders.IsComplete() {
				return officeop.NewApproveMoveBadRequest(), apperror.NewBadDataError("order must be complete")
			}

			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))
			err = h.MoveRouter.Approve(appCtx, move)
			if err != nil {
				logger.Info("Attempted to approve move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
				return handlers.ResponseForError(logger, err), err
			}

			verrs, err := appCtx.DB().ValidateAndUpdate(move)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(logger, verrs, err), err
			}

			// TODO: Save and/or update the move association status' (PPM, Reimbursement, Orders) a la Cancel handler

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return officeop.NewApproveMoveOK().WithPayload(movePayload), nil
		})
}

// CancelMoveHandler cancels a move via POST /moves/{moveId}/cancel
type CancelMoveHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... cancels a Move from a request payload
func (h CancelMoveHandler) Handle(params officeop.CancelMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() {
				sessionErr := apperror.NewSessionError(
					"user is not authorized NewCancelMoveForbidden",
				)
				appCtx.Logger().Error(sessionErr.Error())
				return officeop.NewCancelMoveForbidden(), sessionErr
			}

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))
			// Canceling move will result in canceled associated PPMs
			err = h.MoveRouter.Cancel(appCtx, *params.CancelMove.CancelReason, move)
			if err != nil {
				logger.Error("Attempted to cancel move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
				return handlers.ResponseForError(logger, err), err
			}

			// Save move, orders, and PPMs statuses
			verrs, err := models.SaveMoveDependencies(appCtx.DB(), move)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(logger, verrs, err), err
			}

			/* Don't send emails to BLUEBARK moves */
			if move.Orders.OrdersType != "BLUEBARK" {
				err = h.NotificationSender().SendNotification(appCtx,
					notifications.NewMoveCanceled(moveID),
				)
			}

			if err != nil {
				logger.Error("problem sending email to user", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return officeop.NewCancelMoveOK().WithPayload(movePayload), nil
		})
}

// ApproveReimbursementHandler approves a move via POST /reimbursement/{reimbursementId}/approve
type ApproveReimbursementHandler struct {
	handlers.HandlerConfig
}

// Handle ... approves a Reimbursement from a request payload
func (h ApproveReimbursementHandler) Handle(params officeop.ApproveReimbursementParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() {
				return officeop.NewApproveReimbursementForbidden(), apperror.NewForbiddenError("user must be office user")
			}

			reimbursementID, err := uuid.FromString(params.ReimbursementID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			reimbursement, err := models.FetchReimbursement(appCtx.DB(), appCtx.Session(), reimbursementID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			err = reimbursement.Approve()
			if err != nil {
				appCtx.Logger().Error("Attempted to approve, got invalid transition", zap.Error(err), zap.String("reimbursement_status", string(reimbursement.Status)))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			verrs, err := appCtx.DB().ValidateAndUpdate(reimbursement)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			reimbursementPayload := payloadForReimbursementModel(reimbursement)
			return officeop.NewApproveReimbursementOK().WithPayload(reimbursementPayload), nil
		})
}
