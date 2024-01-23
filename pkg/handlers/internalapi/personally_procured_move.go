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
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForPPMModel(storer storage.FileStorer, personallyProcuredMove models.PersonallyProcuredMove) (*internalmessages.PersonallyProcuredMovePayload, error) {

	documentPayload, err := payloads.PayloadForDocumentModel(storer, personallyProcuredMove.AdvanceWorksheet)
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
