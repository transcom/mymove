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

// ListReportsHandler lists reports with the option to filter since a particular date. Optimized ver.
type ListReportsHandler struct {
	handlers.HandlerConfig
	services.ReportListFetcher
}

// Handle fetches all reports with the option to filter since a particular date. Optimized version.
func (h ListReportsHandler) Handle(params pptasop.ListReportsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			var searchParams services.FetchPaymentRequestListParams
			// if params.Since != nil {
			// 	since := handlers.FmtDateTimePtrToPop(params.Since)
			// 	searchParams.Since = &since
			// }

			reports, err := h.FetchReportList(appCtx, &searchParams)
			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching reports:", zap.Error(err))
				return pptasop.NewListReportsInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.ListReports(&reports)

			return pptasop.NewListReportsOK().WithPayload(payload), nil
		})
}
