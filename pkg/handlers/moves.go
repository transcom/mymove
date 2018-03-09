package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gorilla/context"
	"github.com/markbates/pop"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForMoveModel(move models.Move) internalmessages.MovePayload {
	movePayload := internalmessages.MovePayload{
		CreatedAt:        fmtDateTime(move.CreatedAt),
		SelectedMoveType: swag.String(move.SelectedMoveType),
		ID:               fmtUUID(move.ID),
		UpdatedAt:        fmtDateTime(move.UpdatedAt),
	}
	return movePayload
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewCreateMoveHandler returns a new CreateMoveHandler
func NewCreateMoveHandler(db *pop.Connection, logger *zap.Logger) CreateMoveHandler {
	return CreateMoveHandler{
		db:     db,
		logger: logger,
	}
}

// CreateMoveHandler creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	// Get user id from context
	fmt.Println("HIT IT")
	fmt.Println("!!!", params.HTTPRequest)
	//fmt.Println("CONTEXT", context.GetAll(params.Ht))
	fmt.Println(context.GetAll(params.HTTPRequest))
	//newMove := models.Move{
	//	SelectedMoveType:  *params.CreateMovePayload.SelectedMoveType,
	//}
	var response middleware.Responder
	//if _, err := h.db.ValidateAndCreate(&newMove); err != nil {
	//	h.logger.Error("DB Insertion", zap.Error(err))
	//	response = moveop.NewCreateMoveBadRequest()
	//} else {
	//	movePayload := payloadForMoveModel(newMove)
	//	response = moveop.NewCreateMoveCreated().WithPayload(&movePayload)
	//
	//}
	return response
}
