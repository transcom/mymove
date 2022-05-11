package internalapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
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
		payload := payloads.MTOShipment(&copyOfHhg)
		hhgPayloads = append(hhgPayloads, payload)
	}

	var SelectedMoveType internalmessages.SelectedMoveType
	if move.SelectedMoveType != nil {
		SelectedMoveType = internalmessages.SelectedMoveType(*move.SelectedMoveType)
	}
	var SubmittedAt time.Time
	if move.SubmittedAt != nil {
		SubmittedAt = *move.SubmittedAt
	}

	movePayload := &internalmessages.MovePayload{
		CreatedAt:               handlers.FmtDateTime(move.CreatedAt),
		SubmittedAt:             handlers.FmtDateTime(SubmittedAt),
		SelectedMoveType:        &SelectedMoveType,
		Locator:                 swag.String(move.Locator),
		ID:                      handlers.FmtUUID(move.ID),
		UpdatedAt:               handlers.FmtDateTime(move.UpdatedAt),
		PersonallyProcuredMoves: ppmPayloads,
		MtoShipments:            hhgPayloads,
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
	handlers.HandlerContext
}

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID, _ := uuid.FromString(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			// Fetch orders for authorized user
			orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			payload := params.PatchMovePayload
			newSelectedMoveType := payload.SelectedMoveType

			if newSelectedMoveType != nil {
				stringSelectedMoveType := models.SelectedMoveType(*newSelectedMoveType)
				move.SelectedMoveType = &stringSelectedMoveType
			}

			verrs, err := appCtx.DB().ValidateAndUpdate(move)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(logger, verrs, err), err
			}
			movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			return moveop.NewPatchMoveCreated().WithPayload(movePayload), nil
		})
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler struct {
	handlers.HandlerContext
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
			err = h.MoveRouter.Submit(appCtx, move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			certificateParams := certop.NewCreateSignedCertificationParams()
			certificateParams.CreateSignedCertificationPayload = params.SubmitMoveForApprovalPayload.Certificate
			certificateParams.HTTPRequest = params.HTTPRequest
			certificateParams.MoveID = params.MoveID
			// Transaction to save move and dependencies
			verrs, err := h.saveMoveDependencies(appCtx, move, certificateParams, appCtx.Session().UserID)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(logger, verrs, err), err
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

// SaveMoveDependencies safely saves a Move status, ppmShipment status, mtoShipment status, orders statuses, signed certificate,
// and shipment GBLOCs.
func (h SubmitMoveHandler) saveMoveDependencies(appCtx appcontext.AppContext, move *models.Move, certificateParams certop.CreateSignedCertificationParams, userID uuid.UUID) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	date := time.Time(*certificateParams.CreateSignedCertificationPayload.Date)
	certType := models.SignedCertificationType(*certificateParams.CreateSignedCertificationPayload.CertificationType)
	newSignedCertification := models.SignedCertification{
		MoveID:                   uuid.FromStringOrNil(certificateParams.MoveID.String()),
		PersonallyProcuredMoveID: nil,
		CertificationType:        &certType,
		SubmittingUserID:         userID,
		CertificationText:        *certificateParams.CreateSignedCertificationPayload.CertificationText,
		Signature:                *certificateParams.CreateSignedCertificationPayload.Signature,
		Date:                     date,
	}

	if certificateParams.CreateSignedCertificationPayload.PersonallyProcuredMoveID != nil {
		ppmID := uuid.FromStringOrNil(certificateParams.CreateSignedCertificationPayload.PersonallyProcuredMoveID.String())
		newSignedCertification.PersonallyProcuredMoveID = &ppmID
	}

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		transactionError := errors.New("Rollback The transaction")
		// TODO: move creation of signed certification into a service
		verrs, err := txnAppCtx.DB().ValidateAndCreate(&newSignedCertification)
		if err != nil || verrs.HasAny() {
			responseError = fmt.Errorf("error saving signed certification: %w", err)
			responseVErrors.Append(verrs)
			return transactionError
		}

		// update ppmShipments and mtoShipments if needed
		for i := range move.MTOShipments {
			if move.MTOShipments[i].ShipmentType == models.MTOShipmentTypePPM {
				if verrs, err := txnAppCtx.DB().ValidateAndUpdate(move.MTOShipments[i].PPMShipment); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Updating PPMShipment")
					return transactionError
				}

				if verrs, err := txnAppCtx.DB().ValidateAndUpdate(&move.MTOShipments[i]); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Updating MTOShipment")
					return transactionError
				}
			}
		}

		if verrs, err := txnAppCtx.DB().ValidateAndSave(&move.Orders); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Orders")
			return transactionError
		}

		if verrs, err := txnAppCtx.DB().ValidateAndSave(move); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move")
			return transactionError
		}
		return nil
	})

	if transactionErr != nil {
		return responseVErrors, transactionErr
	}

	appCtx.Logger().Info("signedCertification created",
		zap.String("id", newSignedCertification.ID.String()),
		zap.String("moveId", newSignedCertification.MoveID.String()),
		zap.String("createdAt", newSignedCertification.CreatedAt.String()),
		zap.String("certification_type", string(*newSignedCertification.CertificationType)),
		zap.String("certification_text", newSignedCertification.CertificationText),
	)
	return responseVErrors, responseError
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
			ssfd.Obligations, err = ppmComputer.ComputeObligations(appCtx, ssfd, h.Planner())
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

			payload := ioutil.NopCloser(buf)
			filename := fmt.Sprintf("inline; filename=\"%s-%s-ssw-%s.pdf\"", *ssfd.ServiceMember.FirstName, *ssfd.ServiceMember.LastName, time.Now().Format("01-02-2006"))

			return moveop.NewShowShipmentSummaryWorksheetOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// ShowMoveDatesSummaryHandler returns a summary of the dates in the move process given a move date and move ID.
type ShowMoveDatesSummaryHandler struct {
	handlers.HandlerContext
}

// Handle returns a summary of the dates in the move process.
func (h ShowMoveDatesSummaryHandler) Handle(params moveop.ShowMoveDatesSummaryParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveDate := time.Time(params.MoveDate)
			moveID, _ := uuid.FromString(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// Attach move locator to logger
			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			summary, err := calculateMoveDatesFromMove(appCtx, h.Planner(), moveID, moveDate)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
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

			return moveop.NewShowMoveDatesSummaryOK().WithPayload(moveDatesSummary), nil
		})
}

// ShowShipmentSummaryWorksheetHandler returns a Shipment Summary Worksheet PDF
type ShowShipmentSummaryWorksheetHandler struct {
	handlers.HandlerContext
}

// SubmitAmendedOrdersHandler approves a move via POST /moves/{moveId}/submit
type SubmitAmendedOrdersHandler struct {
	handlers.HandlerContext
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

			err = h.MoveRouter.Submit(appCtx, move)
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
