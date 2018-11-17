package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/awardqueue"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/server"
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

	var SelectedMoveType internalmessages.SelectedMoveType
	if move.SelectedMoveType != nil {
		SelectedMoveType = internalmessages.SelectedMoveType(*move.SelectedMoveType)
	}

	var shipmentPayloads []*internalmessages.Shipment
	for _, shipment := range move.Shipments {
		payload, err := payloadForShipmentModel(shipment)
		if err != nil {
			return nil, err
		}
		shipmentPayloads = append(shipmentPayloads, payload)
	}

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               handlers.FmtDateTime(move.CreatedAt),
		SelectedMoveType:        &SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      handlers.FmtUUID(move.ID),
		UpdatedAt:               handlers.FmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		OrdersID:                handlers.FmtUUID(order.ID),
		Status:                  internalmessages.MoveStatus(move.Status),
		Shipments:               shipmentPayloads,
	}
	return movePayload, nil
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	ordersID, _ := uuid.FromString(params.OrdersID.String())

	orders, err := models.FetchOrderForUser(h.DB(), session, ordersID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	selectedMoveType := models.SelectedMoveType(*params.CreateMovePayload.SelectedMoveType)
	move, verrs, err := orders.CreateNewMove(h.DB(), &selectedMoveType)
	if verrs.HasAny() || err != nil {
		if err == models.ErrCreateViolatesUniqueConstraint {
			h.Logger().Error("Failed to create Unique Record Locator")
		}
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return moveop.NewCreateMoveCreated().WithPayload(movePayload)
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.DB(), session, move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return moveop.NewShowMoveOK().WithPayload(movePayload)
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.DB(), session, move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	payload := params.PatchMovePayload
	newSelectedMoveType := payload.SelectedMoveType

	if newSelectedMoveType != nil {
		if newSelectedMoveType != nil {
			stringSelectedMoveType := models.SelectedMoveType(*newSelectedMoveType)
			move.SelectedMoveType = &stringSelectedMoveType
		}
	}

	verrs, err := h.DB().ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return moveop.NewPatchMoveCreated().WithPayload(movePayload)
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... submit a move for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)
	ctx := params.HTTPRequest.Context()
	ctx, span := beeline.StartSpan(ctx, "SubmitMoveHandler")
	defer span.Send()

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())
	span.AddField("move_id", moveID)

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = move.Submit()
	span.AddField("move-status", string(move.Status))
	if err != nil {
		h.HoneyZapLogger().TraceError(ctx, "Failed to change move status to submit",
			zap.String("move_id", moveID.String()),
			zap.String("move_status", string(move.Status)))
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Transaction to save move and dependencies
	verrs, err := models.SaveMoveDependencies(h.DB(), move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	err = h.NotificationSender().SendNotification(
		notifications.NewMoveSubmitted(h.DB(), h.Logger(), session, h.FetchServiceMember(), moveID),
	)
	if err != nil {
		h.Logger().Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if len(move.Shipments) > 0 {
		go awardqueue.NewAwardQueue(h.DB(), h.HoneyZapLogger()).Run(ctx)
	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload)
}

// ShowMoveDatesSummaryHandler returns a summary of the dates in the move process given a move date and move ID.
type ShowMoveDatesSummaryHandler struct {
	handlers.HandlerContext
}

// Handle returns a summary of the dates in the move process.
func (h ShowMoveDatesSummaryHandler) Handle(params moveop.ShowMoveDatesSummaryParams) middleware.Responder {
	moveDate := time.Time(params.MoveDate)
	moveID, _ := uuid.FromString(params.MoveID.String())

	summary, err := calculateMoveDatesFromMove(h.DB(), h.Planner(), moveID, moveDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	moveDatesSummary := &internalmessages.MoveDatesSummary{
		ID:       swag.String(params.MoveID.String() + ":" + params.MoveDate.String()),
		MoveID:   &params.MoveID,
		MoveDate: &params.MoveDate,
		Pack:     handlers.FmtDateSlice(summary.PackDays),
		Pickup:   handlers.FmtDateSlice(summary.PickupDays),
		Transit:  handlers.FmtDateSlice(summary.TransitDays),
		Delivery: handlers.FmtDateSlice(summary.DeliveryDays),
		Report:   handlers.FmtDateSlice(summary.ReportDays),
	}

	return moveop.NewShowMoveDatesSummaryOK().WithPayload(moveDatesSummary)
}
