package internalapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return officeop.NewApproveMoveForbidden()
	}
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching move", zap.String("move_id", moveID.String()))
	}
	// Don't approve Move if orders are incomplete
	orders, ordersErr := models.FetchOrder(h.DB(), move.OrdersID)
	if ordersErr != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching order", zap.String("move_id", move.OrdersID.String()))
	}
	if orders.IsComplete() != true {
		return officeop.NewApprovePPMBadRequest()
	}

	err = move.Approve()
	if err != nil {
		h.Logger().Info("Attempted to approve move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
		return h.RespondAndTraceError(ctx, err, "error approving move")
	}

	verrs, err := h.DB().ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error validating and updating move")
	}

	// TODO: Save and/or update the move association status' (PPM, Reimbursement, Orders) a la Cancel handler

	movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching payload for move")
	}
	return officeop.NewApproveMoveOK().WithPayload(movePayload)
}

// CancelMoveHandler cancels a move via POST /moves/{moveId}/cancel
type CancelMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... cancels a Move from a request payload
func (h CancelMoveHandler) Handle(params officeop.CancelMoveParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return officeop.NewCancelMoveForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching move", zap.String("move_id", moveID.String()))
	}

	// Canceling move will result in canceled associated PPMs
	err = move.Cancel(*params.CancelMove.CancelReason)
	if err != nil {
		h.Logger().Error("Attempted to cancel move, got invalid transition", zap.Error(err), zap.String("move_status", string(move.Status)))
		return h.RespondAndTraceError(ctx, err, "error cancelling move")
	}

	// Save move, orders, and PPMs statuses
	verrs, err := models.SaveMoveDependencies(h.DB(), move)
	if err != nil || verrs.HasAny() {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error saving move dependencies")
	}

	err = h.NotificationSender().SendNotification(
		ctx,
		notifications.NewMoveCanceled(h.DB(), h.Logger(), session, moveID),
	)

	if err != nil {
		h.Logger().Error("problem sending email to user", zap.Error(err))
		return h.RespondAndTraceError(ctx, err, "error sending email to user")

	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching payload for move")
	}
	return officeop.NewCancelMoveOK().WithPayload(movePayload)
}

// ApprovePPMHandler approves a move via POST /personally_procured_moves/{personallyProcuredMoveId}/approve
type ApprovePPMHandler struct {
	handlers.HandlerContext
}

// Handle ... approves a Personally Procured Move from a request payload
func (h ApprovePPMHandler) Handle(params officeop.ApprovePPMParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return officeop.NewApprovePPMForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching personally procured move", zap.String("personally_procured_move_id", ppmID.String()))
	}
	moveID := ppm.MoveID
	err = ppm.Approve()
	if err != nil {
		h.Logger().Error("Attempted to approve PPM, got invalid transition", zap.Error(err), zap.String("move_status", string(ppm.Status)))
		return h.RespondAndTraceError(ctx, err, "error approving personally procured move")

	}

	verrs, err := h.DB().ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error validating and updating personally procured move")
	}

	err = h.NotificationSender().SendNotification(
		ctx,
		notifications.NewMoveApproved(h.DB(), h.Logger(), session, moveID),
	)
	if err != nil {
		h.Logger().Error("problem sending email to user", zap.Error(err))
		return h.RespondAndTraceError(ctx, err, "error sending email to user")
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching payload for personally procured move", zap.String("personally_procured_move_id", ppm.ID.String()))
	}
	return officeop.NewApprovePPMOK().WithPayload(ppmPayload)
}

// ApproveReimbursementHandler approves a move via POST /reimbursement/{reimbursementId}/approve
type ApproveReimbursementHandler struct {
	handlers.HandlerContext
}

// Handle ... approves a Reimbursement from a request payload
func (h ApproveReimbursementHandler) Handle(params officeop.ApproveReimbursementParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return officeop.NewApproveReimbursementForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	reimbursementID, _ := uuid.FromString(params.ReimbursementID.String())

	reimbursement, err := models.FetchReimbursement(h.DB(), session, reimbursementID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching reimbursement", zap.String("reimbursement_id", reimbursementID.String()))
	}

	err = reimbursement.Approve()
	if err != nil {
		h.Logger().Error("Attempted to approve, got invalid transition", zap.Error(err), zap.String("reimbursement_status", string(reimbursement.Status)))
		return h.RespondAndTraceError(ctx, err, "error approving reimbursement")
	}

	verrs, err := h.DB().ValidateAndUpdate(reimbursement)
	if err != nil || verrs.HasAny() {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error validating and updating reimbursement")
	}

	reimbursementPayload := payloadForReimbursementModel(reimbursement)
	return officeop.NewApproveReimbursementOK().WithPayload(reimbursementPayload)
}
