package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForMoveModel(storer storage.FileStorer, order models.Order, move models.Move) (*internalmessages.MovePayload, error) {

	var ppmPayloads internalmessages.IndexPersonallyProcuredMovePayload
	for _, ppm := range move.PersonallyProcuredMoves {
		payload, err := payloadForPPMModel(storer, ppm)
		if err != nil {
			return nil, err
		}
		ppmPayloads = append(ppmPayloads, payload)
	}

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               fmtDateTime(move.CreatedAt),
		SelectedMoveType:        move.SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      fmtUUID(move.ID),
		UpdatedAt:               fmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		OrdersID:                fmtUUID(order.ID),
		Status:                  internalmessages.MoveStatus(move.Status),
	}
	return movePayload, nil
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler HandlerContext

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	ordersID, _ := uuid.FromString(params.OrdersID.String())

	orders, err := models.FetchOrder(h.db, session, ordersID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	move, verrs, err := orders.CreateNewMove(h.db, params.CreateMovePayload.SelectedMoveType)
	if verrs.HasAny() || err != nil {
		if err == models.ErrCreateViolatesUniqueConstraint {
			h.logger.Error("Failed to create Unique Record Locator")
		}
		return responseForVErrors(h.logger, verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.storage, orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewCreateMoveCreated().WithPayload(movePayload)
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler HandlerContext

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrder(h.db, session, move.OrdersID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	movePayload, err := payloadForMoveModel(h.storage, orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewShowMoveOK().WithPayload(movePayload)
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler HandlerContext

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrder(h.db, session, move.OrdersID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	payload := params.PatchMovePayload
	newSelectedMoveType := payload.SelectedMoveType

	if newSelectedMoveType != nil {
		move.SelectedMoveType = newSelectedMoveType
	}

	verrs, err := h.db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.storage, orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewPatchMoveCreated().WithPayload(movePayload)
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler HandlerContext

// Handle ... submit a move for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	err = move.Submit()
	if err != nil {
		h.logger.Error("Failed to change move status to submit", zap.String("move_id", moveID.String()), zap.String("move_status", string(move.Status)))
		return responseForError(h.logger, err)
	}

	// Transaction to save move and dependencies
	verrs, err := models.SaveMoveStatuses(h.db, move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	movePayload, err := payloadForMoveModel(h.storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload)
}

// CancelMoveHandler cancels a move via $method $path
type CancelMoveHandler HandlerContext

// Handle cancels a move
func (h CancelMoveHandler) Handle(params moveop.SubmitMoveForCancellationParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	err = move.Cancel()
	if err != nil {
		h.logger.Error("Failed to change move status to cancel", zap.String("move_id", moveID.String()), zap.String("move_status", string(move.Status)))
		return responseForError(h.logger, err)
	}

	// Transaction to save move and dependencies
	verrs, err := models.SaveMoveStatuses(h.db, move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	movePayload, err := payloadForMoveModel(h.storage, move.Orders, *move)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewSubmitMoveForCancellationOK().WithPayload(movePayload)
}
