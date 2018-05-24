package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPPMModel(personallyProcuredMove models.PersonallyProcuredMove) internalmessages.PersonallyProcuredMovePayload {

	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                         fmtUUID(personallyProcuredMove.ID),
		CreatedAt:                  fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:                  fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:                       personallyProcuredMove.Size,
		WeightEstimate:             personallyProcuredMove.WeightEstimate,
		EstimatedIncentive:         personallyProcuredMove.EstimatedIncentive,
		PlannedMoveDate:            fmtDatePtr(personallyProcuredMove.PlannedMoveDate),
		PickupPostalCode:           personallyProcuredMove.PickupPostalCode,
		AdditionalPickupPostalCode: personallyProcuredMove.AdditionalPickupPostalCode,
		HasAdditionalPostalCode:    personallyProcuredMove.HasAdditionalPostalCode,
		DestinationPostalCode:      personallyProcuredMove.DestinationPostalCode,
		HasSit:                     personallyProcuredMove.HasSit,
		DaysInStorage:              personallyProcuredMove.DaysInStorage,
		Status:                     internalmessages.PPMStatus(personallyProcuredMove.Status),
	}
	return ppmPayload
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, reqApp, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreatePersonallyProcuredMovePayload
	newPPM, verrs, err := move.CreatePPM(h.db,
		payload.Size,
		payload.WeightEstimate,
		payload.EstimatedIncentive,
		(*time.Time)(payload.PlannedMoveDate),
		payload.PickupPostalCode,
		payload.HasAdditionalPostalCode,
		payload.AdditionalPickupPostalCode,
		payload.DestinationPostalCode,
		payload.HasSit,
		payload.DaysInStorage,
		false,
		nil)

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
	// #nosec User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, reqApp, moveID)
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
	if payload.PickupPostalCode != nil {
		ppm.PickupPostalCode = payload.PickupPostalCode
	}
	if payload.HasAdditionalPostalCode != nil {
		if *payload.HasAdditionalPostalCode == false {
			ppm.AdditionalPickupPostalCode = nil
		} else if *payload.HasAdditionalPostalCode == true {
			ppm.AdditionalPickupPostalCode = payload.AdditionalPickupPostalCode
		}
		ppm.HasAdditionalPostalCode = payload.HasAdditionalPostalCode
	}
	if payload.DestinationPostalCode != nil {
		ppm.DestinationPostalCode = payload.DestinationPostalCode
	}
	if payload.HasSit != nil {
		if *payload.HasSit == false {
			ppm.DaysInStorage = nil
		} else if *payload.HasSit == true {
			ppm.DaysInStorage = payload.DaysInStorage
		}
		ppm.HasSit = payload.HasSit
	}

}

// PatchPersonallyProcuredMoveHandler Patchs a PPM
type PatchPersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.db, user, reqApp, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	if ppm.MoveID != moveID {
		h.logger.Info("Move ID for PPM does not match requested PPM Move ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", ppm.MoveID.String()))
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
