package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForPPMModel(storer storage.FileStorer, personallyProcuredMove models.PersonallyProcuredMove) (*internalmessages.PersonallyProcuredMovePayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, personallyProcuredMove.AdvanceWorksheet)
	var hasProGear *string
	if personallyProcuredMove.HasProGear != nil {
		hpg := string(*personallyProcuredMove.HasProGear)
		hasProGear = &hpg
	}
	var hasProGearOverThousand *string
	if personallyProcuredMove.HasProGearOverThousand != nil {
		hpgot := string(*personallyProcuredMove.HasProGearOverThousand)
		hasProGearOverThousand = &hpgot
	}
	if err != nil {
		return nil, err
	}
	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:                            handlers.FmtUUID(personallyProcuredMove.ID),
		MoveID:                        *handlers.FmtUUID(personallyProcuredMove.MoveID),
		CreatedAt:                     handlers.FmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:                     handlers.FmtDateTime(personallyProcuredMove.UpdatedAt),
		WeightEstimate:                handlers.FmtPoundPtr(personallyProcuredMove.WeightEstimate),
		OriginalMoveDate:              handlers.FmtDatePtr(personallyProcuredMove.OriginalMoveDate),
		ActualMoveDate:                handlers.FmtDatePtr(personallyProcuredMove.ActualMoveDate),
		SubmitDate:                    handlers.FmtDateTimePtr(personallyProcuredMove.SubmitDate),
		ApproveDate:                   handlers.FmtDateTimePtr(personallyProcuredMove.ApproveDate),
		NetWeight:                     handlers.FmtPoundPtr(personallyProcuredMove.NetWeight),
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
		HasProGear:                    hasProGear,
		HasProGearOverThousand:        hasProGearOverThousand,
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
	if personallyProcuredMove.HasProGear != nil {
		hasProGear := string(*personallyProcuredMove.HasProGear)
		ppmPayload.HasProGear = &hasProGear
	}
	if personallyProcuredMove.HasProGearOverThousand != nil {
		hasProGearOverThousand := string(*personallyProcuredMove.HasProGearOverThousand)
		ppmPayload.HasProGearOverThousand = &hasProGearOverThousand
	}
	return &ppmPayload, nil
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler struct {
	handlers.HandlerConfig
}

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			payload := params.CreatePersonallyProcuredMovePayload

			var advance *models.Reimbursement
			if payload.Advance != nil {
				a := models.BuildDraftReimbursement(unit.Cents(*payload.Advance.RequestedAmount), models.MethodOfReceipt(*payload.Advance.MethodOfReceipt))
				advance = &a
			}

			destinationZip, err := GetDestinationDutyLocationPostalCode(appCtx, move.OrdersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newPPM, verrs, err := move.CreatePPM(appCtx.DB(),
				handlers.PoundPtrFromInt64Ptr(payload.WeightEstimate),
				(*time.Time)(payload.OriginalMoveDate),
				payload.PickupPostalCode,
				payload.HasAdditionalPostalCode,
				payload.AdditionalPickupPostalCode,
				&destinationZip,
				payload.HasSit,
				payload.DaysInStorage,
				payload.EstimatedStorageReimbursement,
				payload.HasRequestedAdvance,
				advance)

			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			ppmPayload, err := payloadForPPMModel(h.FileStorer(), *newPPM)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(ppmPayload), nil
		})
}

// IndexPersonallyProcuredMovesHandler returns a list of all the PPMs associated with this move.
type IndexPersonallyProcuredMovesHandler struct {
	handlers.HandlerConfig
}

// Handle handles the request
func (h IndexPersonallyProcuredMovesHandler) Handle(params ppmop.IndexPersonallyProcuredMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// The given move does belong to the current user.
			ppms := move.PersonallyProcuredMoves
			ppmsPayload := make(internalmessages.IndexPersonallyProcuredMovePayload, len(ppms))
			for i, ppm := range ppms {
				ppmPayload, err := payloadForPPMModel(h.FileStorer(), ppm)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				ppmsPayload[i] = ppmPayload
			}
			response := ppmop.NewIndexPersonallyProcuredMovesOK().WithPayload(ppmsPayload)
			return response, nil
		})
}

func patchPPMWithPayload(ppm *models.PersonallyProcuredMove, payload *internalmessages.PatchPersonallyProcuredMovePayload) {

	if payload.WeightEstimate != nil {
		ppm.WeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.WeightEstimate)
	}
	if payload.IncentiveEstimateMax != nil {
		incentiveEstimateMax := unit.Cents(int(*payload.IncentiveEstimateMax))
		ppm.IncentiveEstimateMax = &incentiveEstimateMax
	}
	if payload.IncentiveEstimateMin != nil {
		incentiveEstimateMin := unit.Cents(int(*payload.IncentiveEstimateMin))
		ppm.IncentiveEstimateMin = &incentiveEstimateMin
	}
	if payload.NetWeight != nil {
		ppm.NetWeight = handlers.PoundPtrFromInt64Ptr(payload.NetWeight)
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
		if !*payload.HasAdditionalPostalCode {
			ppm.AdditionalPickupPostalCode = nil
		} else if *payload.HasAdditionalPostalCode {
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
	if payload.HasProGear != nil {
		hasProGear := models.ProGearStatus(*payload.HasProGear)
		ppm.HasProGear = &hasProGear
	}
	if payload.HasProGearOverThousand != nil {
		hasProGearOverThousand := models.ProGearStatus(*payload.HasProGearOverThousand)
		ppm.HasProGearOverThousand = &hasProGearOverThousand
	}
}

// PatchPersonallyProcuredMoveHandler Patches a PPM
type PatchPersonallyProcuredMoveHandler struct {
	handlers.HandlerConfig
}

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			if ppm.MoveID != moveID {
				errMsg := "Move ID for PPM does not match requested PPM Move ID"
				appCtx.Logger().Info(errMsg, zap.String("requested move_id", moveID.String()), zap.String("actual move_id", ppm.MoveID.String()))
				return ppmop.NewPatchPersonallyProcuredMoveBadRequest(), apperror.NewBadDataError(errMsg)
			}

			patchPPMWithPayload(ppm, params.PatchPersonallyProcuredMovePayload)

			verrs, err := models.SavePersonallyProcuredMove(appCtx.DB(), ppm)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ppmop.NewPatchPersonallyProcuredMoveOK().WithPayload(ppmPayload), nil
		})
}

// SubmitPersonallyProcuredMoveHandler Submits a PPM
type SubmitPersonallyProcuredMoveHandler struct {
	handlers.HandlerConfig
}

// Handle Submits a PPM to change its status to SUBMITTED
func (h SubmitPersonallyProcuredMoveHandler) Handle(params ppmop.SubmitPersonallyProcuredMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var submitDate time.Time
			if params.SubmitPersonallyProcuredMovePayload.SubmitDate != nil {
				submitDate = time.Time(*params.SubmitPersonallyProcuredMovePayload.SubmitDate)
			}
			err = ppm.Submit(submitDate)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			verrs, err := models.SavePersonallyProcuredMove(appCtx.DB(), ppm)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			return ppmop.NewSubmitPersonallyProcuredMoveOK().WithPayload(ppmPayload), nil
		})
}

// RequestPPMPaymentHandler requests a payment for a PPM
type RequestPPMPaymentHandler struct {
	handlers.HandlerConfig
}

// Handle is the handler
func (h RequestPPMPaymentHandler) Handle(params ppmop.RequestPPMPaymentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			err = ppm.RequestPayment()
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			verrs, err := models.SavePersonallyProcuredMove(appCtx.DB(), ppm)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			ppmPayload, err := payloadForPPMModel(h.FileStorer(), *ppm)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ppmop.NewRequestPPMPaymentOK().WithPayload(ppmPayload), nil
		})
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
	handlers.HandlerConfig
}

// Handle is the handler
func (h RequestPPMExpenseSummaryHandler) Handle(params ppmop.RequestPPMExpenseSummaryParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

			// Fetch all approved expense documents for a PPM
			status := models.MoveDocumentStatusOK
			moveDocsExpense, err := models.FetchMoveDocuments(appCtx.DB(), appCtx.Session(), ppmID, &status, models.MoveDocumentTypeEXPENSE, false)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			expenseSummaryPayload := buildExpenseSummaryPayload(moveDocsExpense)

			return ppmop.NewRequestPPMExpenseSummaryOK().WithPayload(&expenseSummaryPayload), nil
		})
}
