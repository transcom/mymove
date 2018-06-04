package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForPPMModel(personallyProcuredMove models.PersonallyProcuredMove) internalmessages.PersonallyProcuredMovePayload {

	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                            fmtUUID(personallyProcuredMove.ID),
		CreatedAt:                     fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:                     fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:                          personallyProcuredMove.Size,
		WeightEstimate:                personallyProcuredMove.WeightEstimate,
		EstimatedIncentive:            personallyProcuredMove.EstimatedIncentive,
		PlannedMoveDate:               fmtDatePtr(personallyProcuredMove.PlannedMoveDate),
		PickupPostalCode:              personallyProcuredMove.PickupPostalCode,
		HasAdditionalPostalCode:       personallyProcuredMove.HasAdditionalPostalCode,
		AdditionalPickupPostalCode:    personallyProcuredMove.AdditionalPickupPostalCode,
		DestinationPostalCode:         personallyProcuredMove.DestinationPostalCode,
		HasSit:                        personallyProcuredMove.HasSit,
		DaysInStorage:                 personallyProcuredMove.DaysInStorage,
		EstimatedStorageReimbursement: personallyProcuredMove.EstimatedStorageReimbursement,
		Status:              internalmessages.PPMStatus(personallyProcuredMove.Status),
		HasRequestedAdvance: &personallyProcuredMove.HasRequestedAdvance,
		Advance:             payloadForReimbursementModel(personallyProcuredMove.Advance),
	}
	return ppmPayload
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreatePersonallyProcuredMovePayload

	var advance *models.Reimbursement
	if payload.Advance != nil {
		a := models.BuildDraftReimbursement(unit.Cents(*payload.Advance.RequestedAmount), models.MethodOfReceipt(*payload.Advance.MethodOfReceipt))
		advance = &a
	}

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
		payload.EstimatedStorageReimbursement,
		payload.HasRequestedAdvance,
		advance)

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
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
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
			ppm.EstimatedStorageReimbursement = nil
		} else if *payload.HasSit == true {
			ppm.DaysInStorage = payload.DaysInStorage
			ppm.EstimatedStorageReimbursement = payload.EstimatedStorageReimbursement
		}
		ppm.HasSit = payload.HasSit
	}

	if payload.HasRequestedAdvance != nil {
		ppm.HasRequestedAdvance = *payload.HasRequestedAdvance
	} else if payload.Advance != nil {
		ppm.HasRequestedAdvance = true
	}
	if ppm.HasRequestedAdvance {
		if payload.Advance != nil {
			methodOfReceipt := models.MethodOfReceipt(*payload.Advance.MethodOfReceipt)
			requestedAmount := unit.Cents(*payload.Advance.RequestedAmount)

			if ppm.Advance != nil {
				ppm.Advance.MethodOfReceipt = methodOfReceipt
				ppm.Advance.RequestedAmount = requestedAmount
			} else {
				var advance models.Reimbursement
				if ppm.Status == models.PPMStatusDRAFT {
					advance = models.BuildDraftReimbursement(requestedAmount, methodOfReceipt)
				} else {
					advance = models.BuildRequestedReimbursement(requestedAmount, methodOfReceipt)
				}
				ppm.Advance = &advance
			}
		}
	}
}

// PatchPersonallyProcuredMoveHandler Patchs a PPM
type PatchPersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.db, session, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	if ppm.MoveID != moveID {
		h.logger.Info("Move ID for PPM does not match requested PPM Move ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", ppm.MoveID.String()))
		return ppmop.NewPatchPersonallyProcuredMoveBadRequest()
	}

	patchPPMWithPayload(ppm, params.PatchPersonallyProcuredMovePayload)

	verrs, err := models.SavePersonallyProcuredMove(h.db, ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload := payloadForPPMModel(*ppm)
	return ppmop.NewPatchPersonallyProcuredMoveCreated().WithPayload(&ppmPayload)

}
