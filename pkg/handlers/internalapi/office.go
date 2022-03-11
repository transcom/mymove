package internalapi

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler struct {
	handlers.HandlerContext
	services.MoveRouter
}

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() {
				return officeop.NewApproveMoveForbidden()
			}

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}
			// Don't approve Move if orders are incomplete
			orders, ordersErr := models.FetchOrder(appCtx.DB(), move.OrdersID)
			if ordersErr != nil {
				return handlers.ResponseForError(appCtx.Logger(), ordersErr)
			}
			if !orders.IsComplete() {
				return officeop.NewApprovePPMBadRequest()
			}

			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))
			err = h.MoveRouter.Approve(appCtx, move)
			if err != nil {
				logger.Info("Attempted to approve move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
				return handlers.ResponseForError(logger, err)
			}

			verrs, err := appCtx.DB().ValidateAndUpdate(move)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(logger, verrs, err)
			}

			// TODO: Save and/or update the move association status' (PPM, Reimbursement, Orders) a la Cancel handler

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err)
			}
			return officeop.NewApproveMoveOK().WithPayload(movePayload)
		})
}

// CancelMoveHandler cancels a move via POST /moves/{moveId}/cancel
type CancelMoveHandler struct {
	handlers.HandlerContext
	services.MoveRouter
}

// Handle ... cancels a Move from a request payload
func (h CancelMoveHandler) Handle(params officeop.CancelMoveParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() {
		return officeop.NewCancelMoveForbidden()
	}

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))
	// Canceling move will result in canceled associated PPMs
	err = h.MoveRouter.Cancel(appCtx, *params.CancelMove.CancelReason, move)
	if err != nil {
		logger.Error("Attempted to cancel move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
		return handlers.ResponseForError(logger, err)
	}

	// Save move, orders, and PPMs statuses
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	err = h.NotificationSender().SendNotification(appCtx,
		notifications.NewMoveCanceled(moveID),
	)

	if err != nil {
		logger.Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return officeop.NewCancelMoveOK().WithPayload(movePayload)
}

// ApprovePPMHandler approves a move via POST /personally_procured_moves/{personallyProcuredMoveId}/approve
type ApprovePPMHandler struct {
	handlers.HandlerContext
}

// Handle ... approves a Personally Procured Move from a request payload
func (h ApprovePPMHandler) Handle(params officeop.ApprovePPMParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() {
		return officeop.NewApprovePPMForbidden()
	}

	ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	moveID := ppm.MoveID
	var approveDate time.Time
	if params.ApprovePersonallyProcuredMovePayload.ApproveDate != nil {
		approveDate = time.Time(*params.ApprovePersonallyProcuredMovePayload.ApproveDate)
	}
	err = ppm.Approve(approveDate)
	if err != nil {
		appCtx.Logger().Error("Attempted to approve PPM, got invalid transition", zap.Error(err), zap.String("move_status", string(ppm.Status)))
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}

	err = h.NotificationSender().SendNotification(appCtx,
		notifications.NewMoveApproved(h.HandlerContext.AppNames().MilServername, moveID),
	)
	if err != nil {
		appCtx.Logger().Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return officeop.NewApprovePPMOK().WithPayload(ppmPayload)
}

// ApproveReimbursementHandler approves a move via POST /reimbursement/{reimbursementId}/approve
type ApproveReimbursementHandler struct {
	handlers.HandlerContext
}

// Handle ... approves a Reimbursement from a request payload
func (h ApproveReimbursementHandler) Handle(params officeop.ApproveReimbursementParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() {
		return officeop.NewApproveReimbursementForbidden()
	}

	reimbursementID, err := uuid.FromString(params.ReimbursementID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	reimbursement, err := models.FetchReimbursement(appCtx.DB(), appCtx.Session(), reimbursementID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	err = reimbursement.Approve()
	if err != nil {
		appCtx.Logger().Error("Attempted to approve, got invalid transition", zap.Error(err), zap.String("reimbursement_status", string(reimbursement.Status)))
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(reimbursement)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}

	reimbursementPayload := payloadForReimbursementModel(reimbursement)
	return officeop.NewApproveReimbursementOK().WithPayload(reimbursementPayload)
}
