package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPPMModel(personallyProcuredMove models.PersonallyProcuredMove) internalmessages.PersonallyProcuredMovePayload {

	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                  fmtUUID(personallyProcuredMove.ID),
		CreatedAt:           fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:           fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:                personallyProcuredMove.Size,
		WeightEstimate:      personallyProcuredMove.WeightEstimate,
		EstimatedIncentive:  personallyProcuredMove.EstimatedIncentive,
		PlannedMoveDate:     fmtDatePtr(personallyProcuredMove.PlannedMoveDate),
		PickupZip:           personallyProcuredMove.PickupZip,
		AdditionalPickupZip: personallyProcuredMove.AdditionalPickupZip,
		DestinationZip:      personallyProcuredMove.DestinationZip,
		DaysInStorage:       personallyProcuredMove.DaysInStorage,
	}
	return ppmPayload
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreatePersonallyProcuredMovePayload
	newPPM, verrs, err := move.CreatePPM(h.db,
		payload.Size,
		payload.WeightEstimate,
		payload.EstimatedIncentive,
		(*time.Time)(payload.PlannedMoveDate),
		payload.PickupZip,
		payload.AdditionalPickupZip,
		payload.DestinationZip,
		payload.DaysInStorage)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload := payloadForPPMModel(*newPPM)
	return ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(&ppmPayload)
}

// IndexPersonallyProcuredMovesHandler returns a list of all the PPMs associated with this move.
type IndexPersonallyProcuredMovesHandler HandlerContext

// Handle handles the request
func (h IndexPersonallyProcuredMovesHandler) Handle(params ppmop.IndexPersonallyProcuredMovesParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// The given move does belong to the current user.
	ppms := move.PersonallyProcuredMoves
	ppmsPayload := make(internalmessages.IndexPersonallyProcuredMovePayload, len(ppms))
	for i, ppm := range ppms {
		ppmPayload := payloadForPPMModel(ppm)
		ppmsPayload[i] = &ppmPayload
	}
	response := ppmop.NewIndexPersonallyProcuredMovesOK().WithPayload(ppmsPayload)
	return response
}

func patchPPMWithPayload(ppm *models.PersonallyProcuredMove, payload *internalmessages.PatchPersonallyProcuredMovePayload) {

	if payload.Size != nil {
		ppm.Size = payload.Size
	}
	if payload.WeightEstimate != nil {
		ppm.WeightEstimate = payload.WeightEstimate
	}
	if payload.EstimatedIncentive != nil {
		ppm.EstimatedIncentive = payload.EstimatedIncentive
	}
	if payload.PlannedMoveDate != nil {
		ppm.PlannedMoveDate = (*time.Time)(payload.PlannedMoveDate)
	}
	if payload.PickupZip != nil {
		ppm.PickupZip = payload.PickupZip
	}
	if payload.AdditionalPickupZip != nil {
		ppm.AdditionalPickupZip = payload.AdditionalPickupZip
	}
	if payload.DestinationZip != nil {
		ppm.DestinationZip = payload.DestinationZip
	}
	if payload.DaysInStorage != nil {
		ppm.DaysInStorage = payload.DaysInStorage
	}
}

// PatchPersonallyProcuredMoveHandler Patchs a PPM
type PatchPersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.db, user, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	if ppm.MoveID != moveID {
		h.logger.Info("Move ID for PPM does not match requested PPM Move ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", ppm.MoveID.String()))
		// Response is empty here.
		fmt.Println("bad request response", ppmop.NewPatchPersonallyProcuredMoveBadRequest())
		return ppmop.NewPatchPersonallyProcuredMoveBadRequest()
	}

	patchPPMWithPayload(ppm, params.PatchPersonallyProcuredMovePayload)

	verrs, err := h.db.ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload := payloadForPPMModel(*ppm)
	return ppmop.NewPatchPersonallyProcuredMoveCreated().WithPayload(&ppmPayload)

}
