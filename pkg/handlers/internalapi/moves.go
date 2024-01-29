package internalapi

import (
	"fmt"
	"io"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	"github.com/transcom/mymove/pkg/storage"
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

func payloadForMoveModel(storer storage.FileStorer, order models.Order, move models.Move) (*internalmessages.MovePayload, error) {

	var ppmPayloads internalmessages.IndexPersonallyProcuredMovePayload
	for _, ppm := range move.PersonallyProcuredMoves {
		payload, err := payloadForPPMModel(storer, ppm)
		if err != nil {
			return nil, err
		}
		ppmPayloads = append(ppmPayloads, payload)
	}

	var hhgPayloads internalmessages.MTOShipments
	for _, hhg := range move.MTOShipments {
		copyOfHhg := hhg // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload := payloads.MTOShipment(storer, &copyOfHhg)
		hhgPayloads = append(hhgPayloads, payload)
	}

	var SubmittedAt time.Time
	if move.SubmittedAt != nil {
		SubmittedAt = *move.SubmittedAt
	}

	eTag := etag.GenerateEtag(move.UpdatedAt)

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               handlers.FmtDateTime(move.CreatedAt),
		SubmittedAt:             handlers.FmtDateTime(SubmittedAt),
		Locator:                 models.StringPointer(move.Locator),
		ID:                      handlers.FmtUUID(move.ID),
		UpdatedAt:               handlers.FmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		MtoShipments:            hhgPayloads,
		OrdersID:                handlers.FmtUUID(order.ID),
		ServiceMemberID:         *handlers.FmtUUID(order.ServiceMemberID),
		Status:                  internalmessages.MoveStatus(move.Status),
		ETag:                    &eTag,
	}

	if move.CloseoutOffice != nil {
		movePayload.CloseoutOffice = payloads.TransportationOffice(*move.CloseoutOffice)
	}
	return movePayload, nil
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, _ := uuid.FromString(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// Fetch orders for authorized user
			orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return moveop.NewShowMoveOK().WithPayload(movePayload), nil
		})
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler struct {
	handlers.HandlerConfig
	services.MoveCloseoutOfficeUpdater
}

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("PatchMoveHandler error", zap.Error(err))
				switch errors.Cause(err) {
				case models.ErrFetchForbidden:
					return moveop.NewPatchMoveForbidden(), err
				case models.ErrFetchNotFound:
					return moveop.NewPatchMoveNotFound(), err
				default:
					switch err.(type) {
					case apperror.NotFoundError:
						return moveop.NewPatchMoveNotFound(), err
					case apperror.PreconditionFailedError:
						return moveop.NewPatchMovePreconditionFailed(), err
					default:
						return moveop.NewPatchMoveInternalServerError(), err
					}
				}
			}

			if !appCtx.Session().IsMilApp() || !appCtx.Session().IsServiceMember() {
				return moveop.NewPatchMoveUnauthorized(), nil
			}

			moveID := uuid.FromStringOrNil(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handleError(err)
			}

			// Fetch orders for authorized user
			orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
			if err != nil {
				return handleError(err)
			}

			closeoutOfficeID := uuid.FromStringOrNil(params.PatchMovePayload.CloseoutOfficeID.String())
			move, err = h.MoveCloseoutOfficeUpdater.UpdateCloseoutOffice(appCtx, move.Locator, closeoutOfficeID, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
			if err != nil {
				return handleError(err)
			}

			return moveop.NewPatchMoveOK().WithPayload(movePayload), nil
		})
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... submit a move to TOO for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID, _ := uuid.FromString(params.MoveID.String())

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			newSignedCertification := payloads.SignedCertificationFromSubmit(params.SubmitMoveForApprovalPayload, appCtx.Session().UserID, params.MoveID)
			err = h.MoveRouter.Submit(appCtx, move, newSignedCertification)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			err = h.NotificationSender().SendNotification(appCtx,
				notifications.NewMoveSubmitted(moveID),
			)
			if err != nil {
				logger.Error("problem sending email to user", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload), nil
		})
}

// Handle returns a generated PDF
func (h ShowShipmentSummaryWorksheetHandler) Handle(params moveop.ShowShipmentSummaryWorksheetParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			logger := appCtx.Logger()

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				logger.Error("Error fetching PPMShipment", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			ssfd, err := h.SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(appCtx, appCtx.Session(), ppmShipmentID)
			if err != nil {
				logger.Error("Error fetching data for SSW", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			ssfd.Obligations, err = h.SSWPPMComputer.ComputeObligations(appCtx, *ssfd, h.DTODPlanner())
			if err != nil {
				logger.Error("Error calculating obligations ", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			page1Data, page2Data := h.SSWPPMComputer.FormatValuesShipmentSummaryWorksheet(*ssfd)
			if err != nil {
				logger.Error("Error formatting data for SSW", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			ppmGenerator := shipmentsummaryworksheet.NewSSWPPMGenerator()
			SSWPPMWorksheet, SSWPDFInfo, err := ppmGenerator.FillSSWPDFForm(page1Data, page2Data)
			if err != nil {
				return nil, err
			}
			if SSWPDFInfo.PageCount != 2 {
				return nil, errors.Wrap(err, "SSWGenerator output a corrupted or incorretly altered PDF")
			}
			payload := io.NopCloser(SSWPPMWorksheet)
			filename := fmt.Sprintf("inline; filename=\"%s-%s-ssw-%s.pdf\"", *ssfd.ServiceMember.FirstName, *ssfd.ServiceMember.LastName, time.Now().Format("01-02-2006"))

			return moveop.NewShowShipmentSummaryWorksheetOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// ShowShipmentSummaryWorksheetHandler returns a Shipment Summary Worksheet PDF
type ShowShipmentSummaryWorksheetHandler struct {
	handlers.HandlerConfig
	services.SSWPPMComputer
}

// SubmitAmendedOrdersHandler approves a move via POST /moves/{moveId}/submit
type SubmitAmendedOrdersHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... submit a move to TOO for approval
func (h SubmitAmendedOrdersHandler) Handle(params moveop.SubmitAmendedOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID, _ := uuid.FromString(params.MoveID.String())

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			err = h.MoveRouter.RouteAfterAmendingOrders(appCtx, move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			responseVErrors := validate.NewErrors()
			var responseError error

			if verrs, saveErr := appCtx.DB().ValidateAndSave(move); verrs.HasAny() || saveErr != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(saveErr, "Error Saving Move")
			}

			if responseVErrors.HasAny() {
				return handlers.ResponseForVErrors(logger, responseVErrors, responseError), responseError
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return moveop.NewSubmitAmendedOrdersOK().WithPayload(movePayload), nil
		})
}
