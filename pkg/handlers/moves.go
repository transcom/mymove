package handlers

import (
	"github.com/go-openapi/runtime/middleware"
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
type IndexMovesHandler struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewIndexMovesHandler returns a new IndexMovesHandler
func NewIndexMovesHandler(db *pop.Connection, logger *zap.Logger) IndexMovesHandler {
	return IndexMovesHandler{
		db:     db,
		logger: logger,
	}
}

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
