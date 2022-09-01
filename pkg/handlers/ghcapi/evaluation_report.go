package ghcapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/assets"
	evaluationReportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

// GetShipmentEvaluationReportsHandler gets a list of shipment evaluation reports for a given move
type GetShipmentEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportFetcher
}

// Handle handles GetShipmentEvaluationReports by move request
func (h GetShipmentEvaluationReportsHandler) Handle(params moveop.GetMoveShipmentEvaluationReportsListParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() {
				err := apperror.NewForbiddenError("not an office user")
				appCtx.Logger().Error(err.Error())
				return moveop.NewGetMoveShipmentEvaluationReportsListForbidden(), err
			}

			reports, err := h.FetchEvaluationReports(appCtx, models.EvaluationReportTypeShipment, handlers.FmtUUIDToPop(params.MoveID), appCtx.Session().OfficeUserID)
			if err != nil {
				return moveop.NewGetMoveShipmentEvaluationReportsListInternalServerError(), err
			}

			payload := payloads.EvaluationReportList(reports)
			return moveop.NewGetMoveShipmentEvaluationReportsListOK().WithPayload(payload), nil
		},
	)
}

// CreateEvaluationReport is the struct for creating an evaluation report
type CreateEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportCreator
}

//Handle is the handler for creating an evaluation report
func (h CreateEvaluationReportHandler) Handle(params evaluationReportop.CreateEvaluationReportParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			report := &models.EvaluationReport{
				ID: uuid.Must(uuid.NewV4()),
			}

			if params.Body != nil {
				payload := params.Body

				shipmentID := uuid.FromStringOrNil(payload.ShipmentID.String())
				report.Type = models.EvaluationReportTypeShipment
				report.ShipmentID = &shipmentID
			} else {
				report.Type = models.EvaluationReportTypeCounseling
			}

			if appCtx.Session() != nil {
				report.OfficeUserID = appCtx.Session().OfficeUserID
			}

			evaluationReport, err := h.CreateEvaluationReport(appCtx, report, params.Locator)
			if err != nil {
				appCtx.Logger().Error("Error creating evaluation report: ", zap.Error(err))
				return evaluationReportop.NewCreateEvaluationReportInternalServerError(), err
			}

			returnPayload := payloads.EvaluationReport(evaluationReport)

			return evaluationReportop.NewCreateEvaluationReportOK().WithPayload(returnPayload), nil
		})
}

// GetCounselingEvaluationReportsHandler gets a list of counseling evaluation reports for a given move
type GetCounselingEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportFetcher
}

// Handle handles GetCounselingEvaluationReports by move request
func (h GetCounselingEvaluationReportsHandler) Handle(params moveop.GetMoveCounselingEvaluationReportsListParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() {
				err := apperror.NewForbiddenError("not an office user")
				appCtx.Logger().Error(err.Error())
				return moveop.NewGetMoveCounselingEvaluationReportsListForbidden(), err
			}

			reports, err := h.FetchEvaluationReports(appCtx, models.EvaluationReportTypeCounseling, handlers.FmtUUIDToPop(params.MoveID), appCtx.Session().OfficeUserID)
			if err != nil {
				return moveop.NewGetMoveCounselingEvaluationReportsListInternalServerError(), err
			}

			payload := payloads.EvaluationReportList(reports)
			return moveop.NewGetMoveCounselingEvaluationReportsListOK().WithPayload(payload), nil
		},
	)
}

// DownloadEvaluationReportHandler is the struct for fetching an evaluation report by ID
type DownloadEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportFetcher
}

// Handle is the handler for fetching an evaluation report by ID
func (h DownloadEvaluationReportHandler) Handle(params evaluationReportop.DownloadEvaluationReportParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("DownloadEvaluationReport error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return evaluationReportop.NewDownloadEvaluationReportNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return evaluationReportop.NewDownloadEvaluationReportForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return evaluationReportop.NewDownloadEvaluationReportInternalServerError(), err
				default:
					return evaluationReportop.NewDownloadEvaluationReportInternalServerError(), err
				}
			}

			reportID := uuid.FromStringOrNil(params.ReportID.String())
			evaluationReport, err := h.FetchEvaluationReportByID(appCtx, reportID, appCtx.Session().OfficeUserID)
			if err != nil {
				return handleError(err)
			}
			page1Data := paperwork.FormatValuesEvaluationReportPage1(*evaluationReport)
			formFiller := paperwork.NewFormFiller()
			page1Layout := paperwork.EvaluationReportPage1Layout
			page1Template, err := assets.Asset(page1Layout.TemplateImagePath)
			if err != nil {
				appCtx.Logger().Error("Error reading page 1 template file", zap.String("asset", page1Layout.TemplateImagePath), zap.Error(err))
				return evaluationReportop.NewDownloadEvaluationReportInternalServerError(), err
			}
			page1Reader := bytes.NewReader(page1Template)
			err = formFiller.AppendPage(page1Reader, page1Layout.FieldsLayout, page1Data)
			if err != nil {
				appCtx.Logger().Error("Error appending page 1 to PDF", zap.Error(err))
				return evaluationReportop.NewDownloadEvaluationReportInternalServerError(), err
			}
			buf := new(bytes.Buffer)
			err = formFiller.Output(buf)
			if err != nil {
				appCtx.Logger().Error("Error writing out PDF", zap.Error(err))
				return evaluationReportop.NewDownloadEvaluationReportInternalServerError(), err
			}
			payload := ioutil.NopCloser(buf)
			filename := fmt.Sprintf("inline; filename=\"evalreport-%s-%s.pdf\"", evaluationReport.ID, time.Now().Format("01-02-2006"))
			return evaluationReportop.NewDownloadEvaluationReportOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// GetEvaluationReportHandler is the struct for fetching an evaluation report by ID
type GetEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportFetcher
}

// Handle is the handler for fetching an evaluation report by ID
func (h GetEvaluationReportHandler) Handle(params evaluationReportop.GetEvaluationReportParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetEvaluationReport error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return evaluationReportop.NewGetEvaluationReportNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return evaluationReportop.NewGetEvaluationReportForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return evaluationReportop.NewGetEvaluationReportInternalServerError(), err
				default:
					return evaluationReportop.NewGetEvaluationReportInternalServerError(), err
				}
			}

			reportID := uuid.FromStringOrNil(params.ReportID.String())
			evaluationReport, err := h.FetchEvaluationReportByID(appCtx, reportID, appCtx.Session().OfficeUserID)
			if err != nil {
				return handleError(err)
			}
			payload := payloads.EvaluationReport(evaluationReport)
			return evaluationReportop.NewGetEvaluationReportOK().WithPayload(payload), nil
		})
}

// DeleteEvaluationReportHandler is the struct for soft deleting evaluation reports
type DeleteEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportDeleter
}

// Handle is the handler function for soft deleting an evaluation report
func (h DeleteEvaluationReportHandler) Handle(params evaluationReportop.DeleteEvaluationReportParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			reportID := uuid.FromStringOrNil(string(params.ReportID.String()))
			err := h.DeleteEvaluationReport(appCtx, reportID)
			if err != nil {
				appCtx.Logger().Error("Error deleting evaluation report: ", zap.Error(err))
				return evaluationReportop.NewDeleteEvaluationReportInternalServerError(), err
			}

			return evaluationReportop.NewDeleteEvaluationReportNoContent(), nil
		})
}

type SaveEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportUpdater
}

func (h SaveEvaluationReportHandler) Handle(params evaluationReportop.SaveEvaluationReportParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			eTag := params.IfMatch
			payload := params.Body
			payload.ID = params.ReportID
			report := payloads.EvaluationReportFromUpdate(payload)

			if appCtx.Session() != nil {
				report.OfficeUserID = appCtx.Session().OfficeUserID
			}

			err := h.UpdateEvaluationReport(appCtx, report, appCtx.Session().OfficeUserID, eTag)
			if err != nil {
				appCtx.Logger().Error("Error saving evaluation report: ", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return evaluationReportop.NewSaveEvaluationReportNotFound(), err
				case apperror.PreconditionFailedError:
					return evaluationReportop.NewSaveEvaluationReportPreconditionFailed(), err
				case apperror.ForbiddenError:
					return evaluationReportop.NewSaveEvaluationReportForbidden(), err
				case apperror.ConflictError:
					return evaluationReportop.NewSaveEvaluationReportConflict(), err
				case apperror.InvalidInputError:
					return evaluationReportop.NewSaveEvaluationReportUnprocessableEntity(), err
				default:
					return evaluationReportop.NewSaveEvaluationReportInternalServerError(), err
				}
			}

			return evaluationReportop.NewSaveEvaluationReportNoContent(), nil
		})
}
