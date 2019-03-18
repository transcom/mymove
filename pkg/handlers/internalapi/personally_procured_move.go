package internalapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForPPMModel(storer storage.FileStorer, personallyProcuredMove models.PersonallyProcuredMove) (*internalmessages.PersonallyProcuredMovePayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, personallyProcuredMove.AdvanceWorksheet)
	if err != nil {
		return nil, err
	}

	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                            handlers.FmtUUID(personallyProcuredMove.ID),
		MoveID:                        *handlers.FmtUUID(personallyProcuredMove.MoveID),
		CreatedAt:                     handlers.FmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:                     handlers.FmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:                          personallyProcuredMove.Size,
		WeightEstimate:                personallyProcuredMove.WeightEstimate,
		OriginalMoveDate:              handlers.FmtDatePtr(personallyProcuredMove.OriginalMoveDate),
		ActualMoveDate:                handlers.FmtDatePtr(personallyProcuredMove.ActualMoveDate),
		NetWeight:                     personallyProcuredMove.NetWeight,
		PickupPostalCode:              personallyProcuredMove.PickupPostalCode,
		HasAdditionalPostalCode:       personallyProcuredMove.HasAdditionalPostalCode,
		AdditionalPickupPostalCode:    personallyProcuredMove.AdditionalPickupPostalCode,
		DestinationPostalCode:         personallyProcuredMove.DestinationPostalCode,
		HasSit:                        personallyProcuredMove.HasSit,
		DaysInStorage:                 personallyProcuredMove.DaysInStorage,
		EstimatedStorageReimbursement: personallyProcuredMove.EstimatedStorageReimbursement,
		Status:                        internalmessages.PPMStatus(personallyProcuredMove.Status),
		HasRequestedAdvance:           &personallyProcuredMove.HasRequestedAdvance,
		Advance:                       payloadForReimbursementModel(personallyProcuredMove.Advance),
		AdvanceWorksheet:              documentPayload,
		Mileage:                       personallyProcuredMove.Mileage,
		TotalSitCost:                  handlers.FmtCost(personallyProcuredMove.TotalSITCost),
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
type CreatePersonallyProcuredMoveHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := params.CreatePersonallyProcuredMovePayload

	var advance *models.Reimbursement
	if payload.Advance != nil {
		a := models.BuildDraftReimbursement(unit.Cents(*payload.Advance.RequestedAmount), models.MethodOfReceipt(*payload.Advance.MethodOfReceipt))
		advance = &a
	}

	newPPM, verrs, err := move.CreatePPM(h.DB(),
		payload.Size,
		payload.WeightEstimate,
		(*time.Time)(payload.OriginalMoveDate),
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
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *newPPM)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(ppmPayload)
}

// IndexPersonallyProcuredMovesHandler returns a list of all the PPMs associated with this move.
type IndexPersonallyProcuredMovesHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h IndexPersonallyProcuredMovesHandler) Handle(params ppmop.IndexPersonallyProcuredMovesParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	// The given move does belong to the current user.
	ppms := move.PersonallyProcuredMoves
	ppmsPayload := make(internalmessages.IndexPersonallyProcuredMovePayload, len(ppms))
	for i, ppm := range ppms {
		ppmPayload, err := payloadForPPMModel(h.FileStorer(), ppm)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
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
	if payload.NetWeight != nil {
		ppm.NetWeight = payload.NetWeight
	}
	if payload.OriginalMoveDate != nil {
		ppm.OriginalMoveDate = (*time.Time)(payload.OriginalMoveDate)
	}
	if payload.ActualMoveDate != nil {
		ppm.ActualMoveDate = (*time.Time)(payload.ActualMoveDate)
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
		ppm.HasSit = payload.HasSit
	}

	if payload.TotalSitCost != nil {
		cost := unit.Cents(*payload.TotalSitCost)
		ppm.TotalSITCost = &cost
	}

	if payload.DaysInStorage != nil {
		ppm.DaysInStorage = payload.DaysInStorage
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

// PatchPersonallyProcuredMoveHandler Patches a PPM
type PatchPersonallyProcuredMoveHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	if ppm.MoveID != moveID {
		h.Logger().Info("Move ID for PPM does not match requested PPM Move ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", ppm.MoveID.String()))
		return ppmop.NewPatchPersonallyProcuredMoveBadRequest()
	}

	needsEstimatesRecalculated := h.ppmNeedsEstimatesRecalculated(ppm, params.PatchPersonallyProcuredMovePayload)

	patchPPMWithPayload(ppm, params.PatchPersonallyProcuredMovePayload)

	if needsEstimatesRecalculated {
		err = h.updateEstimates(ppm)
		if err != nil {
			h.Logger().Error("Unable to set calculated fields on PPM", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	}

	verrs, err := models.SavePersonallyProcuredMove(h.DB(), ppm)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ppmop.NewPatchPersonallyProcuredMoveOK().WithPayload(ppmPayload)
}

// ppmNeedsEstimatesRecalculated determines whether the fields that comprise
// the PPM incentive and SIT estimate calculations have changed, necessitating a recalculation
func (h PatchPersonallyProcuredMoveHandler) ppmNeedsEstimatesRecalculated(ppm *models.PersonallyProcuredMove, patch *internalmessages.PatchPersonallyProcuredMovePayload) bool {
	originPtr := patch.PickupPostalCode
	destinationPtr := patch.DestinationPostalCode
	weightPtr := patch.WeightEstimate
	datePtr := patch.OriginalMoveDate
	daysPtr := patch.DaysInStorage

	// Figure out if we have values to compare and, if so, whether the new or old value
	// should be used in the calculation
	origin, originChanged, originOK := stringForComparison(ppm.PickupPostalCode, originPtr)
	destination, destinationChanged, destinationOK := stringForComparison(ppm.DestinationPostalCode, destinationPtr)
	weight, weightChanged, weightOK := int64ForComparison(ppm.WeightEstimate, weightPtr)
	date, dateChanged, dateOK := dateForComparison(ppm.OriginalMoveDate, (*time.Time)(datePtr))
	daysInStorage, daysChanged, _ := int64ForComparison(ppm.DaysInStorage, daysPtr)

	// We don't care if daysInStorage is OK, since we just want to meet the minimum bar to recalculate
	valuesOK := originOK && destinationOK && weightOK && dateOK
	valuesChanged := originChanged || destinationChanged || weightChanged || dateChanged || daysChanged

	needsUpdate := valuesOK && valuesChanged

	if needsUpdate {
		h.Logger().Info("updating PPM calculated fields",
			zap.String("originZip", origin),
			zap.String("destinationZip", destination),
			zap.Int64("weight", weight),
			zap.Time("date", date),
			zap.Int64("daysInStorage", daysInStorage),
		)
	}

	return needsUpdate
}

// SubmitPersonallyProcuredMoveHandler Submits a PPM
type SubmitPersonallyProcuredMoveHandler struct {
	handlers.HandlerContext
}

// Handle Submits a PPM to change its status to SUBMITTED
func (h SubmitPersonallyProcuredMoveHandler) Handle(params ppmop.SubmitPersonallyProcuredMoveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)

	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = ppm.Submit()

	verrs, err := models.SavePersonallyProcuredMove(h.DB(), ppm)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)

	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ppmop.NewSubmitPersonallyProcuredMoveOK().WithPayload(ppmPayload)
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

func dateForComparison(previousValue, newValue *time.Time) (value time.Time, valueChanged bool, canCompare bool) {
	if newValue != nil {
		if previousValue != nil {
			return *newValue, previousValue.Equal(*newValue), true
		}
		return *newValue, true, true
	}
	if previousValue != nil {
		return *previousValue, false, true
	}

	return value, false, false
}

func (h PatchPersonallyProcuredMoveHandler) updateEstimates(ppm *models.PersonallyProcuredMove) error {
	re := rateengine.NewRateEngine(h.DB(), h.Logger())
	daysInSIT := 0
	if ppm.HasSit != nil && *ppm.HasSit && ppm.DaysInStorage != nil {
		daysInSIT = int(*ppm.DaysInStorage)
	}

	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(h.DB(), h.Logger(), *ppm.PickupPostalCode, *ppm.DestinationPostalCode, *ppm.OriginalMoveDate)
	if err != nil {
		return err
	}

	// Update SIT estimate
	if ppm.HasSit != nil && *ppm.HasSit == true {
		cwtWeight := unit.Pound(*ppm.WeightEstimate).ToCWT()
		sitZip3 := rateengine.Zip5ToZip3(*ppm.DestinationPostalCode)
		sitComputation, err := re.SitCharge(cwtWeight, daysInSIT, sitZip3, *ppm.OriginalMoveDate, true)
		if err != nil {
			return err
		}
		sitCharge := float64(sitComputation.ApplyDiscount(lhDiscount, sitDiscount))
		reimbursementString := fmt.Sprintf("$%.2f", sitCharge/100)
		ppm.EstimatedStorageReimbursement = &reimbursementString
	}

	distanceMiles, err := h.Planner().Zip5TransitDistance(*ppm.PickupPostalCode, *ppm.DestinationPostalCode)
	if err != nil {
		return err
	}

	cost, err := re.ComputePPM(unit.Pound(*ppm.WeightEstimate), *ppm.PickupPostalCode, *ppm.DestinationPostalCode, distanceMiles, *ppm.OriginalMoveDate, daysInSIT, lhDiscount, sitDiscount)
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

// RequestPPMPaymentHandler requests a payment for a PPM
type RequestPPMPaymentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h RequestPPMPaymentHandler) Handle(params ppmop.RequestPPMPaymentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = ppm.RequestPayment()
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	verrs, err := models.SavePersonallyProcuredMove(h.DB(), ppm)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ppmop.NewRequestPPMPaymentOK().WithPayload(ppmPayload)

}

func buildExpenseSummaryPayload(moveDocsExpense []models.MoveDocument) internalmessages.ExpenseSummaryPayload {

	expenseSummaryPayload := internalmessages.ExpenseSummaryPayload{
		GrandTotal: &internalmessages.ExpenseSummaryPayloadGrandTotal{
			PaymentMethodTotals: &internalmessages.PaymentMethodsTotals{},
		},
		Categories: []*internalmessages.CategoryExpenseSummary{},
	}

	if len(moveDocsExpense) < 1 {
		return expenseSummaryPayload
	}

	catMap := map[internalmessages.MovingExpenseType]*internalmessages.CategoryExpenseSummary{}

	for _, moveDoc := range moveDocsExpense {
		// First add up grand totals by payment type and grand total
		expenseDoc := moveDoc.MovingExpenseDocument
		amount := expenseDoc.RequestedAmountCents.Int64()
		methodTotals := expenseSummaryPayload.GrandTotal.PaymentMethodTotals
		switch expenseDoc.PaymentMethod {
		case "OTHER":
			methodTotals.OTHER += amount
		case "GTCC":
			methodTotals.GTCC += amount
		}
		expenseSummaryPayload.GrandTotal.Total += amount

		// Build categories by expense type
		expenseType := internalmessages.MovingExpenseType(string(expenseDoc.MovingExpenseType))
		// Check if expense type exists in catMap - increment values if so
		if CategoryExpenseSummary, ok := catMap[expenseType]; ok {
			switch expenseDoc.PaymentMethod {
			case "OTHER":
				CategoryExpenseSummary.PaymentMethods.OTHER += amount
			case "GTCC":
				CategoryExpenseSummary.PaymentMethods.GTCC += amount
			}
			CategoryExpenseSummary.Total += amount
		} else { // initialize CategoryExpenseSummary
			var otherAmt, gtccAmt int64
			switch expenseDoc.PaymentMethod {
			case "OTHER":
				otherAmt = amount
			case "GTCC":
				gtccAmt = amount
			}
			catMap[expenseType] = &internalmessages.CategoryExpenseSummary{
				Category: expenseType,
				Total:    amount,
				PaymentMethods: &internalmessages.PaymentMethodsTotals{
					OTHER: otherAmt,
					GTCC:  gtccAmt,
				},
			}
		}
	}
	for _, catExpenseSummary := range catMap {
		expenseSummaryPayload.Categories = append(
			expenseSummaryPayload.Categories, catExpenseSummary)
	}
	return expenseSummaryPayload
}

// RequestPPMExpenseSummaryHandler requests
type RequestPPMExpenseSummaryHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h RequestPPMExpenseSummaryHandler) Handle(params ppmop.RequestPPMExpenseSummaryParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	// Fetch all approved expense documents for a PPM
	moveDocsExpense, err := models.FetchApprovedMovingExpenseDocuments(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	expenseSummaryPayload := buildExpenseSummaryPayload(moveDocsExpense)

	return ppmop.NewRequestPPMExpenseSummaryOK().WithPayload(&expenseSummaryPayload)
}
