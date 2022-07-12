package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetEvaluationReportsHandler gets a move by locator
type GetEvaluationReportsHandler struct {
	handlers.HandlerConfig
	services.EvaluationReportListFetcher
}

// Handle handles the getEvaluationReports by locator request
func (h GetEvaluationReportsHandler) Handle(params moveop.GetMoveEvaluationReportsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() {
				err := apperror.NewForbiddenError("not an office user")
				appCtx.Logger().Error(err.Error())
				return moveop.NewGetMoveEvaluationReportsForbidden(), err
			}

			reports, err := h.FetchEvaluationReports(appCtx, handlers.FmtUUIDToPop(params.MoveID), appCtx.Session().OfficeUserID)
			if err != nil {
				return moveop.NewGetMoveEvaluationReportsInternalServerError(), err
			}

			payload := payloads.EvaluationReports(reports)
			return moveop.NewGetMoveEvaluationReportsOK().WithPayload(payload), nil
		},
	)
}
