package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	evaluationReportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetShipmentEvaluationReportsHandler gets a list of shipment evaluation reports for a given move
type GetShipmentEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportListFetcher
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

// GetCounselingEvaluationReportsHandler gets a list of counseling evaluation reports for a given move
type GetCounselingEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportListFetcher
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

type CreateEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportCreator
}

type DeleteEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportDeleter
}

func (h CreateEvaluationReportHandler) Handle(params evaluationReportop.CreateEvaluationReportForShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			shipmentID := uuid.FromStringOrNil(payload.ShipmentID.String())
			report := &models.EvaluationReport{
				ShipmentID: &shipmentID,
				Type:       models.EvaluationReportTypeShipment,
				ID:         uuid.Must(uuid.NewV4()),
			}

			if appCtx.Session() != nil {
				report.OfficeUserID = appCtx.Session().OfficeUserID
			}

			evaluationReport, err := h.CreateEvaluationReport(appCtx, report)
			if err != nil {
				appCtx.Logger().Error("Error creating evaluation report: ", zap.Error(err))
				return evaluationReportop.NewCreateEvaluationReportForShipmentInternalServerError(), err
			}

			returnPayload := payloads.EvaluationReport(evaluationReport)

			return evaluationReportop.NewCreateEvaluationReportForShipmentOK().WithPayload(returnPayload), nil
		})
}

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
