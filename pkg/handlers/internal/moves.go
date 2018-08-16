package internal

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/storage"
)

/*
 * --------------------------------------------
 * The code below is for the INTERNAL REST API.
 * --------------------------------------------
 */

func payloadForMoveModel(storer storage.FileStorer, order models.Order, move models.Move) (*internalmessages.MovePayload, error) {

	var ppmPayloads internalmessages.IndexPersonallyProcuredMovePayload
	for _, ppm := range move.PersonallyProcuredMoves {
		payload, err := payloadForPPMModel(storer, ppm)
		if err != nil {
			return nil, err
		}
		ppmPayloads = append(ppmPayloads, payload)
	}

	var SelectedMoveType internalmessages.SelectedMoveType
	if move.SelectedMoveType != nil {
		SelectedMoveType = internalmessages.SelectedMoveType(*move.SelectedMoveType)
	}

	var shipmentPayloads []*internalmessages.Shipment
	for _, shipment := range move.Shipments {
		payload := payloadForShipmentModel(shipment)
		shipmentPayloads = append(shipmentPayloads, payload)
	}

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               utils.FmtDateTime(move.CreatedAt),
		SelectedMoveType:        &SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      utils.FmtUUID(move.ID),
		UpdatedAt:               utils.FmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		OrdersID:                utils.FmtUUID(order.ID),
		Status:                  internalmessages.MoveStatus(move.Status),
		Shipments:               shipmentPayloads,
	}
	return movePayload, nil
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler utils.HandlerContext

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	ordersID, _ := uuid.FromString(params.OrdersID.String())

	orders, err := models.FetchOrderForUser(h.Db, session, ordersID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}

	move, verrs, err := orders.CreateNewMove(h.Db, params.CreateMovePayload.SelectedMoveType)
	if verrs.HasAny() || err != nil {
		if err == models.ErrCreateViolatesUniqueConstraint {
			h.Logger.Error("Failed to create Unique Record Locator")
		}
		return utils.ResponseForVErrors(h.Logger, verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.Storage, orders, *move)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	return moveop.NewCreateMoveCreated().WithPayload(movePayload)
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler utils.HandlerContext

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.Db, session, move.OrdersID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}

	movePayload, err := payloadForMoveModel(h.Storage, orders, *move)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	return moveop.NewShowMoveOK().WithPayload(movePayload)
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler utils.HandlerContext

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.Db, session, move.OrdersID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	payload := params.PatchMovePayload
	newSelectedMoveType := payload.SelectedMoveType

	if newSelectedMoveType != nil {
		stringSelectedMoveType := ""
		if newSelectedMoveType != nil {
			stringSelectedMoveType = string(*newSelectedMoveType)
			move.SelectedMoveType = &stringSelectedMoveType
		}
	}

	verrs, err := h.Db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.Logger, verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.Storage, orders, *move)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	return moveop.NewPatchMoveCreated().WithPayload(movePayload)
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler utils.HandlerContext

// Handle ... submit a move for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}

	err = move.Submit()
	if err != nil {
		h.Logger.Error("Failed to change move status to submit", zap.String("move_id", moveID.String()), zap.String("move_status", string(move.Status)))
		return utils.ResponseForError(h.Logger, err)
	}

	// Transaction to save move and dependencies
	verrs, err := models.SaveMoveDependencies(h.Db, move)
	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.Logger, verrs, err)
	}

	err = h.NotificationSender.SendNotification(
		notifications.NewMoveSubmitted(h.Db, h.Logger, session, moveID),
	)
	if err != nil {
		h.Logger.Error("problem sending email to user", zap.Error(err))
		return utils.ResponseForError(h.Logger, err)
	}

	movePayload, err := payloadForMoveModel(h.Storage, move.Orders, *move)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}
	return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload)
}
