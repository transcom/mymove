package internalapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/transcom/mymove/pkg/cli"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
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

	var SelectedMoveType internalmessages.SelectedMoveType
	if move.SelectedMoveType != nil {
		SelectedMoveType = internalmessages.SelectedMoveType(*move.SelectedMoveType)
	}

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               handlers.FmtDateTime(move.CreatedAt),
		SelectedMoveType:        &SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      handlers.FmtUUID(move.ID),
		UpdatedAt:               handlers.FmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		OrdersID:                handlers.FmtUUID(order.ID),
		ServiceMemberID:         *handlers.FmtUUID(order.ServiceMemberID),
		Status:                  internalmessages.MoveStatus(move.Status),
	}
	return movePayload, nil
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.DB(), session, move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return moveop.NewShowMoveOK().WithPayload(movePayload)
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	// Fetch orders for authorized user
	orders, err := models.FetchOrderForUser(h.DB(), session, move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	payload := params.PatchMovePayload
	newSelectedMoveType := payload.SelectedMoveType

	if newSelectedMoveType != nil {
		stringSelectedMoveType := models.SelectedMoveType(*newSelectedMoveType)
		move.SelectedMoveType = &stringSelectedMoveType
	}

	verrs, err := h.DB().ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}
	movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return moveop.NewPatchMoveCreated().WithPayload(movePayload)
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler struct {
	handlers.HandlerContext
}

// Handle ... submit a move for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	submitDate := time.Time(*params.SubmitMoveForApprovalPayload.PpmSubmitDate)
	err = move.Submit(submitDate)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	// Transaction to save move and dependencies
	verrs, err := models.SaveMoveDependencies(h.DB(), move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	err = h.NotificationSender().SendNotification(
		ctx,
		notifications.NewMoveSubmitted(h.DB(), logger, session, moveID),
	)
	if err != nil {
		logger.Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	if h.HandlerContext.GetFeatureFlag(cli.FeatureFlagConvertPPMsToGHC) {
		if moID, convertErr := models.ConvertFromPPMToGHC(h.DB(), move.ID); convertErr != nil {
			logger.Error("PPM->GHC conversion error", zap.Error(err))
			// don't return an error here because we don't want this to break the endpoint
		} else {
			logger.Info("PPM->GHC conversion successful. New MoveOrder created.", zap.String("move_order_id", moID.String()))
		}
	} else {
		logger.Info("PPM->GHC conversion skipped (env var not set)")
	}

	movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload)
}

// Handle returns a generated PDF
func (h ShowShipmentSummaryWorksheetHandler) Handle(params moveop.ShowShipmentSummaryWorksheetParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ppmComputer := paperwork.NewSSWPPMComputer(rateengine.NewRateEngine(h.DB(), logger, *move))

	ssfd, err := models.FetchDataShipmentSummaryWorksheetFormData(h.DB(), session, moveID)
	if err != nil {
		logger.Error("Error fetching data for SSW", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	ssfd.PreparationDate = time.Time(params.PreparationDate)
	ssfd.Obligations, err = ppmComputer.ComputeObligations(ssfd, h.Planner())
	if err != nil {
		logger.Error("Error calculating obligations ", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	page1Data, page2Data, page3Data, err := models.FormatValuesShipmentSummaryWorksheet(ssfd)

	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	formFiller := paperwork.NewFormFiller()

	// page 1
	page1Layout := paperwork.ShipmentSummaryPage1Layout
	page1Template, err := assets.Asset(page1Layout.TemplateImagePath)

	if err != nil {
		logger.Error("Error reading page 1 template file", zap.String("asset", page1Layout.TemplateImagePath), zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	page1Reader := bytes.NewReader(page1Template)
	err = formFiller.AppendPage(page1Reader, page1Layout.FieldsLayout, page1Data)
	if err != nil {
		logger.Error("Error appending page 1 to PDF", zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	// page 2
	page2Layout := paperwork.ShipmentSummaryPage2Layout
	page2Template, err := assets.Asset(page2Layout.TemplateImagePath)

	if err != nil {
		logger.Error("Error reading page 2 template file", zap.String("asset", page2Layout.TemplateImagePath), zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	page2Reader := bytes.NewReader(page2Template)
	err = formFiller.AppendPage(page2Reader, page2Layout.FieldsLayout, page2Data)
	if err != nil {
		logger.Error("Error appending 2 page to PDF", zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	// page 3
	page3Layout := paperwork.ShipmentSummaryPage3Layout
	page3Template, err := assets.Asset(page3Layout.TemplateImagePath)

	if err != nil {
		logger.Error("Error reading page 3 template file", zap.String("asset", page3Layout.TemplateImagePath), zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	page3Reader := bytes.NewReader(page3Template)
	err = formFiller.AppendPage(page3Reader, page3Layout.FieldsLayout, page3Data)
	if err != nil {
		logger.Error("Error appending page 3 to PDF", zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	buf := new(bytes.Buffer)
	err = formFiller.Output(buf)
	if err != nil {
		logger.Error("Error writing out PDF", zap.Error(err))
		return moveop.NewShowShipmentSummaryWorksheetInternalServerError()
	}

	payload := ioutil.NopCloser(buf)
	filename := fmt.Sprintf("attachment; filename=\"%s-%s-ssw-%s.pdf\"", *ssfd.ServiceMember.FirstName, *ssfd.ServiceMember.LastName, time.Now().Format("01-02-2006"))

	return moveop.NewShowShipmentSummaryWorksheetOK().WithContentDisposition(filename).WithPayload(payload)
}

// ShowMoveDatesSummaryHandler returns a summary of the dates in the move process given a move date and move ID.
type ShowMoveDatesSummaryHandler struct {
	handlers.HandlerContext
}

// Handle returns a summary of the dates in the move process.
func (h ShowMoveDatesSummaryHandler) Handle(params moveop.ShowMoveDatesSummaryParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	moveDate := time.Time(params.MoveDate)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	_, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	summary, err := calculateMoveDatesFromMove(h.DB(), h.Planner(), moveID, moveDate)
	if err != nil {
		return handlers.ResponseForError(logger, err)
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

// ShowShipmentSummaryWorksheetHandler returns a Shipment Summary Worksheet PDF
type ShowShipmentSummaryWorksheetHandler struct {
	handlers.HandlerContext
}
