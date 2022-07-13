package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
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

// Handle handles getShipmentEvaluationReports by move request
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

			payload := payloads.EvaluationReports(reports)
			return moveop.NewGetMoveShipmentEvaluationReportsListOK().WithPayload(payload), nil
		},
	)
}

// GetCounselingEvaluationReportsHandler gets a list of shipment evaluation reports for a given move
type GetCounselingEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportListFetcher
}

// Handle handles getCounselingEvaluationReports by move request
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

			payload := payloads.EvaluationReports(reports)
			return moveop.NewGetMoveCounselingEvaluationReportsListOK().WithPayload(payload), nil
		},
	)
}
