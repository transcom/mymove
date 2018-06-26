package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler HandlerContext

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	move.Status = models.MoveStatusAPPROVED

	verrs, err := h.db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	movePayload, err := payloadForMoveModel(h.storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return officeop.NewApproveMoveOK().WithPayload(movePayload)
}

// CancelMoveHandler cancels a move via POST /moves/{moveId}/cancel
type CancelMoveHandler HandlerContext

// Handle ... cancels a Move from a request payload
func (h CancelMoveHandler) Handle(params officeop.CancelMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Canceling move will result in canceled associated PPMs
	err = move.Cancel(*params.Reason)
	if err != nil {
		h.logger.Error("Attempted to cancel move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
		return responseForError(h.logger, err)
	}

	// Save move, orders, and PPMs statuses
	verrs, err := models.SaveMoveStatuses(h.db, move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	err = notifications.SendNotification(
		notifications.NewMoveCanceled(h.db, h.logger, session, moveID),
		h.sesService,
	)

	if err != nil {
		h.logger.Error("problem sending email to user", zap.Error(err))
		return responseForError(h.logger, err)
	}

	movePayload, err := payloadForMoveModel(h.storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return officeop.NewCancelMoveOK().WithPayload(movePayload)
}

// ApprovePPMHandler approves a move via POST /personally_procured_moves/{personallyProcuredMoveId}/approve
type ApprovePPMHandler HandlerContext

// Handle ... approves a Personally Procured Move from a request payload
func (h ApprovePPMHandler) Handle(params officeop.ApprovePPMParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.db, session, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	moveID := ppm.MoveID
	ppm.Status = models.PPMStatusAPPROVED

	verrs, err := h.db.ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	err = notifications.SendNotification(
		notifications.NewMoveApproved(h.db, h.logger, session, moveID),
		h.sesService,
	)
	if err != nil {
		h.logger.Error("problem sending email to user", zap.Error(err))
		return responseForError(h.logger, err)
	}

	ppmPayload, err := payloadForPPMModel(h.storage, *ppm)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return officeop.NewApprovePPMOK().WithPayload(ppmPayload)
}

// ApproveReimbursementHandler approves a move via POST /reimbursement/{reimbursementId}/approve
type ApproveReimbursementHandler HandlerContext

// Handle ... approves a Reimbursement from a request payload
func (h ApproveReimbursementHandler) Handle(params officeop.ApproveReimbursementParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return officeop.NewApproveReimbursementUnauthorized()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	reimbursementID, _ := uuid.FromString(params.ReimbursementID.String())

	reimbursement, err := models.FetchReimbursement(h.db, session, reimbursementID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	err = reimbursement.Approve()
	if err != nil {
		h.logger.Error("Attempted to approve, got invalid transition", zap.Error(err), zap.String("reimbursement_status", string(reimbursement.Status)))
		return responseForError(h.logger, err)
	}

	verrs, err := h.db.ValidateAndUpdate(reimbursement)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	reimbursementPayload := payloadForReimbursementModel(reimbursement)
	return officeop.NewApproveReimbursementOK().WithPayload(reimbursementPayload)
}
