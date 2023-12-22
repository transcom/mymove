package internalapi

import (
	"bytes"
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
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/etag"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/services"
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

			moveID, _ := uuid.FromString(params.MoveID.String())

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			ppmComputer := paperwork.NewSSWPPMComputer(rateengine.NewRateEngine(*move))

			ssfd, err := models.FetchDataShipmentSummaryWorksheetFormData(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				logger.Error("Error fetching data for SSW", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			ssfd.PreparationDate = time.Time(params.PreparationDate)
			ssfd.Obligations, err = ppmComputer.ComputeObligations(appCtx, ssfd, h.DTODPlanner())
			if err != nil {
				logger.Error("Error calculating obligations ", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			page1Data, page2Data, page3Data, err := models.FormatValuesShipmentSummaryWorksheet(ssfd)

			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			formFiller := paperwork.NewFormFiller()

			// page 1
			page1Layout := paperwork.ShipmentSummaryPage1Layout
			page1Template, err := assets.Asset(page1Layout.TemplateImagePath)

			if err != nil {
				appCtx.Logger().Error("Error reading page 1 template file", zap.String("asset", page1Layout.TemplateImagePath), zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			page1Reader := bytes.NewReader(page1Template)
			err = formFiller.AppendPage(page1Reader, page1Layout.FieldsLayout, page1Data)
			if err != nil {
				appCtx.Logger().Error("Error appending page 1 to PDF", zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			// page 2
			page2Layout := paperwork.ShipmentSummaryPage2Layout
			page2Template, err := assets.Asset(page2Layout.TemplateImagePath)

			if err != nil {
				appCtx.Logger().Error("Error reading page 2 template file", zap.String("asset", page2Layout.TemplateImagePath), zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			page2Reader := bytes.NewReader(page2Template)
			err = formFiller.AppendPage(page2Reader, page2Layout.FieldsLayout, page2Data)
			if err != nil {
				appCtx.Logger().Error("Error appending 2 page to PDF", zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			// page 3
			page3Layout := paperwork.ShipmentSummaryPage3Layout
			page3Template, err := assets.Asset(page3Layout.TemplateImagePath)

			if err != nil {
				appCtx.Logger().Error("Error reading page 3 template file", zap.String("asset", page3Layout.TemplateImagePath), zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			page3Reader := bytes.NewReader(page3Template)
			err = formFiller.AppendPage(page3Reader, page3Layout.FieldsLayout, page3Data)
			if err != nil {
				appCtx.Logger().Error("Error appending page 3 to PDF", zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			buf := new(bytes.Buffer)
			err = formFiller.Output(buf)
			if err != nil {
				appCtx.Logger().Error("Error writing out PDF", zap.Error(err))
				return moveop.NewShowShipmentSummaryWorksheetInternalServerError(), err
			}

			payload := io.NopCloser(buf)
			filename := fmt.Sprintf("inline; filename=\"%s-%s-ssw-%s.pdf\"", *ssfd.ServiceMember.FirstName, *ssfd.ServiceMember.LastName, time.Now().Format("01-02-2006"))

			return moveop.NewShowShipmentSummaryWorksheetOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// ShowShipmentSummaryWorksheetHandler returns a Shipment Summary Worksheet PDF
type ShowShipmentSummaryWorksheetHandler struct {
	handlers.HandlerConfig
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
