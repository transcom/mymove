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
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForPPMModel(storage FileStorer, personallyProcuredMove models.PersonallyProcuredMove) (*internalmessages.PersonallyProcuredMovePayload, error) {

	documentPayload, err := payloadForDocumentModel(storage, personallyProcuredMove.AdvanceWorksheet)
	if err != nil {
		return nil, err
	}

	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                            fmtUUID(personallyProcuredMove.ID),
		CreatedAt:                     fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:                     fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:                          personallyProcuredMove.Size,
		WeightEstimate:                personallyProcuredMove.WeightEstimate,
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
		AdvanceWorksheet:    documentPayload,
		Mileage:             personallyProcuredMove.Mileage,
	}
	if personallyProcuredMove.IncentiveEstimateMin != nil {
		min := (*personallyProcuredMove.IncentiveEstimateMin).Int64()
		ppmPayload.IncentiveEstimateMin = &min
	}
	if personallyProcuredMove.IncentiveEstimateMax != nil {
		max := (*personallyProcuredMove.IncentiveEstimateMax).Int64()
		ppmPayload.IncentiveEstimateMax = &max
	}
	if personallyProcuredMove.PlannedSITMax != nil {
		max := (*personallyProcuredMove.PlannedSITMax).Int64()
		ppmPayload.PlannedSitMax = &max
	}
	if personallyProcuredMove.SITMax != nil {
		max := (*personallyProcuredMove.SITMax).Int64()
		ppmPayload.SitMax = &max
	}
	return &ppmPayload, nil
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

	ppmPayload, err := payloadForPPMModel(h.storage, *newPPM)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(ppmPayload)
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
		ppmPayload, err := payloadForPPMModel(h.storage, ppm)
		if err != nil {
			return responseForError(h.logger, err)
		}
		ppmsPayload[i] = ppmPayload
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

	originPtr := params.PatchPersonallyProcuredMovePayload.PickupPostalCode
	destinationPtr := params.PatchPersonallyProcuredMovePayload.DestinationPostalCode
	weightPtr := params.PatchPersonallyProcuredMovePayload.WeightEstimate

	// Figure out if we have values to compare and, if so, whether the new or old value
	// should be used in the calculation
	origin, originChanged, originOK := stringForComparison(ppm.PickupPostalCode, originPtr)
	destination, destinationChanged, destinationOK := stringForComparison(ppm.DestinationPostalCode, destinationPtr)
	weight, weightChanged, weightOK := int64ForComparison(ppm.WeightEstimate, weightPtr)

	patchPPMWithPayload(ppm, params.PatchPersonallyProcuredMovePayload)

	if originOK && destinationOK && weightOK && (originChanged || destinationChanged || weightChanged) {
		h.logger.Info("updating PPM calculated fields",
			zap.String("originZip", origin),
			zap.String("destinationZip", destination),
			zap.Int64("weight", weight),
		)
		err = h.updateCalculatedFields(ppm, origin, destination)
		if err != nil {
			h.logger.Error("Unable to set calculated fields on PPM", zap.Error(err))
		}
	} else {
		h.logger.Info("not recalculating cached PPM fields",
			zap.String("originZip", origin),
			zap.String("destinationZip", destination),
		)
	}

	verrs, err := models.SavePersonallyProcuredMove(h.db, ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload, err := payloadForPPMModel(h.storage, *ppm)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return ppmop.NewPatchPersonallyProcuredMoveCreated().WithPayload(ppmPayload)

}

func stringForComparison(previousValue, newValue *string) (value string, valueChanged bool, canCompare bool) {
	if newValue != nil {
		if previousValue != nil {
			return *newValue, *previousValue != *newValue, true
		}
		return *newValue, true, true
	}
	if previousValue != nil {
		return *previousValue, false, true
	}

	return "", false, false
}

func int64ForComparison(previousValue, newValue *int64) (value int64, valueChanged bool, canCompare bool) {
	if newValue != nil {
		if previousValue != nil {
			return *newValue, *previousValue != *newValue, true
		}
		return *newValue, true, true
	}
	if previousValue != nil {
		return *previousValue, false, true
	}

	return 0, false, false
}

func (h PatchPersonallyProcuredMoveHandler) updateCalculatedFields(ppm *models.PersonallyProcuredMove, newOrigin string, newDestination string) error {
	re := rateengine.NewRateEngine(h.db, h.logger, h.planner)
	daysInSIT := 0
	if ppm.HasSit != nil && *ppm.HasSit && ppm.DaysInStorage != nil {
		daysInSIT = int(*ppm.DaysInStorage)
	}

	lhDiscount, sitDiscount, err := PPMDiscountFetch(h.db, h.logger, newOrigin, newDestination, *ppm.PlannedMoveDate)
	if err != nil {
		return err
	}

	cost, err := re.ComputePPM(unit.Pound(*ppm.WeightEstimate), newOrigin, newDestination, *ppm.PlannedMoveDate, daysInSIT, lhDiscount, sitDiscount)
	if err != nil {
		return err
	}

	mileage := int64(cost.LinehaulCostComputation.Mileage)
	ppm.Mileage = &mileage
	ppm.PlannedSITMax = &cost.SITFee
	ppm.SITMax = &cost.SITMax
	min := cost.GCC.MultiplyFloat64(0.95)
	max := cost.GCC.MultiplyFloat64(1.05)
	ppm.IncentiveEstimateMin = &min
	ppm.IncentiveEstimateMax = &max

	return nil
}
