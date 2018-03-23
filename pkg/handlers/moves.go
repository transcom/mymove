package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForMoveModel(user models.User, move models.Move) internalmessages.MovePayload {
	movePayload := internalmessages.MovePayload{
		CreatedAt:        fmtDateTime(move.CreatedAt),
		SelectedMoveType: move.SelectedMoveType,
		ID:               fmtUUID(move.ID),
		UpdatedAt:        fmtDateTime(move.UpdatedAt),
		UserID:           fmtUUID(user.ID),
	}
	return movePayload
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler HandlerContext

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	var response middleware.Responder
	// Get user id from context
	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response = moveop.NewCreateMoveUnauthorized()
		return response
	}

	// Create a new move for an authenticated user
	newMove := models.Move{
		UserID:           user.ID,
		SelectedMoveType: params.CreateMovePayload.SelectedMoveType,
	}
	if verrs, err := h.db.ValidateAndCreate(&newMove); verrs.HasAny() || err != nil {
		if verrs.HasAny() {
			h.logger.Error("DB Validation", zap.Error(verrs))
		} else {
			h.logger.Error("DB Insertion", zap.Error(err))
		}
		response = moveop.NewCreateMoveBadRequest()
	} else {
		movePayload := payloadForMoveModel(user, newMove)
		response = moveop.NewCreateMoveCreated().WithPayload(&movePayload)
	}
	return response
}

// IndexMovesHandler returns a list of all moves
type IndexMovesHandler HandlerContext

// Handle retrieves a list of all moves in the system belonging to the logged in user
func (h IndexMovesHandler) Handle(params moveop.IndexMovesParams) middleware.Responder {
	var response middleware.Responder

	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response = moveop.NewIndexMovesUnauthorized()
		return response
	}

	moves, err := models.GetMovesForUserID(h.db, user.ID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = moveop.NewIndexMovesBadRequest()
	} else {
		movePayloads := make(internalmessages.IndexMovesPayload, len(moves))
		for i, move := range moves {
			movePayload := payloadForMoveModel(user, move)
			movePayloads[i] = &movePayload
		}
		response = moveop.NewIndexMovesOK().WithPayload(movePayloads)
	}
	return response
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler HandlerContext

// Handle ... patches a new Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	var response middleware.Responder
	// Get user id from context
	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response = moveop.NewPatchMoveUnauthorized()
		return response
	}
	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	// Validate that this move belongs to the current user
	moveResult, err := models.GetMoveForUser(h.db, user.ID, moveID)
	if err != nil {
		h.logger.Error("DB Error checking on move validity", zap.Error(err))
		response = moveop.NewPatchMoveInternalServerError()
	} else if !moveResult.IsValid() {
		switch errCode := moveResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound:
			response = moveop.NewPatchMoveNotFound()
		case models.FetchErrorForbidden:
			response = moveop.NewPatchMoveForbidden()
		default:
			h.logger.Fatal("This case statement is no longer exhaustive!")
		}
	} else { // The given move does belong to the current user.
		move := moveResult.Move()
		payload := params.PatchMovePayload
		newSelectedMoveType := payload.SelectedMoveType

		move.SelectedMoveType = newSelectedMoveType

		if verrs, err := h.db.ValidateAndUpdate(&move); verrs.HasAny() || err != nil {
			if verrs.HasAny() {
				h.logger.Error("DB Validation", zap.Error(verrs))
			} else {
				h.logger.Error("DB Update", zap.Error(err))
			}
			response = moveop.NewPatchMoveBadRequest()
		} else {
			movePayload := payloadForMoveModel(user, move)
			response = moveop.NewPatchMoveCreated().WithPayload(&movePayload)
		}
	}
	return response
}
