package pptasapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	pptasop "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations/moves"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/pptasapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// PPTASReportsHandler lists reports with the option to filter since a particular date. Optimized ver.
type PPTASReportsHandler struct {
	handlers.HandlerConfig
	services.PPTASReportListFetcher
}

// Handle fetches all reports with the option to filter since a particular date. Optimized version.
func (h PPTASReportsHandler) Handle(params pptasop.PptasReportsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			var searchParams services.MoveTaskOrderFetcherParams
			if params.Since != nil {
				since := handlers.FmtDateTimePtrToPop(params.Since)
				searchParams.Since = &since
			}

			movesForReport, err := h.GetMovesForReportBuilder(appCtx, &searchParams)
			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching moves:", zap.Error(err))
				return pptasop.NewPptasReportsInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			pptasReports, err := h.BuildPPTASReportsFromMoves(appCtx, movesForReport)
			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching reports:", zap.Error(err))
				return pptasop.NewPptasReportsInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.PPTASReports(appCtx, &pptasReports)

			return pptasop.NewPptasReportsOK().WithPayload(payload), nil
		})
}
