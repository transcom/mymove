package ghcapi

import (
	"github.com/transcom/mymove/pkg/gen/ghcmessages"

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

// CreateEvaluationReportHandler is the struct for creating an evaluation report
type CreateEvaluationReportHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportCreator
}

//Handle is the handler for creating an evaluation report
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
