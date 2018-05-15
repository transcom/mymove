package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForMoveModel(order models.Order, move models.Move) internalmessages.MovePayload {

	var ppmPayloads internalmessages.IndexPersonallyProcuredMovePayload
	for _, ppm := range move.PersonallyProcuredMoves {
		payload := payloadForPPMModel(ppm)
		ppmPayloads = append(ppmPayloads, &payload)
	}

	movePayload := internalmessages.MovePayload{
		CreatedAt:               fmtDateTime(move.CreatedAt),
		SelectedMoveType:        move.SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      fmtUUID(move.ID),
		UpdatedAt:               fmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		OrdersID:                fmtUUID(order.ID),
	}
	return movePayload
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler HandlerContext

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	// Get orders for authorized user
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	ordersID, _ := uuid.FromString(params.OrdersID.String())
	orders, err := models.FetchOrder(h.db, user, reqApp, ordersID)
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
	movePayload := payloadForMoveModel(orders, *move)
	return moveop.NewCreateMoveCreated().WithPayload(&movePayload)
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler HandlerContext

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, reqApp, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrder(h.db, user, reqApp, move.OrdersID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	movePayload := payloadForMoveModel(orders, *move)
	return moveop.NewShowMoveOK().WithPayload(&movePayload)
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler HandlerContext

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, reqApp, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrder(h.db, user, reqApp, move.OrdersID)
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
	movePayload := payloadForMoveModel(orders, *move)
	return moveop.NewPatchMoveCreated().WithPayload(&movePayload)
}
