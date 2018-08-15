package internal

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler utils.HandlerContext

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return officeop.NewApproveMoveForbidden()
	}
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	// Don't approve Move if orders are incomplete
	orders, ordersErr := models.FetchOrder(h.Db, move.OrdersID)
	if ordersErr != nil {
		return responseForError(h.Logger, ordersErr)
	}
	if orders.IsComplete() != true {
		return officeop.NewApprovePPMBadRequest()
	}

	move.Approve()

	verrs, err := h.Db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	// TODO: Save and/or update the move association status' (PPM, Reimbursement, Orders) a la Cancel handler

	movePayload, err := payloadForMoveModel(h.Storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return officeop.NewApproveMoveOK().WithPayload(movePayload)
}

// CancelMoveHandler cancels a move via POST /moves/{moveId}/cancel
type CancelMoveHandler utils.HandlerContext

// Handle ... cancels a Move from a request payload
func (h CancelMoveHandler) Handle(params officeop.CancelMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return officeop.NewCancelMoveForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	// Canceling move will result in canceled associated PPMs
	err = move.Cancel(*params.CancelMove.CancelReason)
	if err != nil {
		h.Logger.Error("Attempted to cancel move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
		return responseForError(h.Logger, err)
	}

	// Save move, orders, and PPMs statuses
	verrs, err := models.SaveMoveDependencies(h.Db, move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	err = h.NotificationSender.SendNotification(
		notifications.NewMoveCanceled(h.Db, h.Logger, session, moveID),
	)

	if err != nil {
		h.Logger.Error("problem sending email to user", zap.Error(err))
		return responseForError(h.Logger, err)
	}

	movePayload, err := payloadForMoveModel(h.Storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return officeop.NewCancelMoveOK().WithPayload(movePayload)
}

// ApprovePPMHandler approves a move via POST /personally_procured_moves/{personallyProcuredMoveId}/approve
type ApprovePPMHandler utils.HandlerContext

// Handle ... approves a Personally Procured Move from a request payload
func (h ApprovePPMHandler) Handle(params officeop.ApprovePPMParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return officeop.NewApprovePPMForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.Db, session, ppmID)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	moveID := ppm.MoveID
	ppm.Status = models.PPMStatusAPPROVED

	verrs, err := h.Db.ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	err = h.NotificationSender.SendNotification(
		notifications.NewMoveApproved(h.Db, h.Logger, session, moveID),
	)
	if err != nil {
		h.Logger.Error("problem sending email to user", zap.Error(err))
		return responseForError(h.Logger, err)
	}

	ppmPayload, err := payloadForPPMModel(h.Storage, *ppm)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return officeop.NewApprovePPMOK().WithPayload(ppmPayload)
}

// ApproveReimbursementHandler approves a move via POST /reimbursement/{reimbursementId}/approve
type ApproveReimbursementHandler utils.HandlerContext

// Handle ... approves a Reimbursement from a request payload
func (h ApproveReimbursementHandler) Handle(params officeop.ApproveReimbursementParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return officeop.NewApproveReimbursementForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	reimbursementID, _ := uuid.FromString(params.ReimbursementID.String())

	reimbursement, err := models.FetchReimbursement(h.Db, session, reimbursementID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	err = reimbursement.Approve()
	if err != nil {
		h.Logger.Error("Attempted to approve, got invalid transition", zap.Error(err), zap.String("reimbursement_status", string(reimbursement.Status)))
		return responseForError(h.Logger, err)
	}

	verrs, err := h.Db.ValidateAndUpdate(reimbursement)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	reimbursementPayload := payloadForReimbursementModel(reimbursement)
	return officeop.NewApproveReimbursementOK().WithPayload(reimbursementPayload)
}
