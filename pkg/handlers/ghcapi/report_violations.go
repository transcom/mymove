package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	reportViolationop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/report_violations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// AssociateReportViolationsHandler is the struct report violations
type AssociateReportViolationsHandler struct {
	handlers.HandlerConfig
	services.ReportViolationsCreator
}

//Handle is the handler for associating violations with reports
func (h AssociateReportViolationsHandler) Handle(params reportViolationop.AssociateReportViolationsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			reportID := handlers.FmtUUIDToPop(params.ReportID)

			var reportViolations models.ReportViolations
			for _, violation := range params.Body.Violations {
				violatinoID, err := uuid.FromString(violation.String())
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Error parsing violation id: %s", violatinoID.String()), zap.Error(err))
					return reportViolationop.NewAssociateReportViolationsInternalServerError(), err
				}
				reportViolation := models.ReportViolation{
					ReportID:    reportID,
					ViolationID: violatinoID,
				}
				reportViolations = append(reportViolations, reportViolation)
			}

			err := h.AssociateReportViolations(appCtx, &reportViolations, reportID)
			if err != nil {
				appCtx.Logger().Error("Error associating report violations: ", zap.Error(err))
				return reportViolationop.NewAssociateReportViolationsInternalServerError(), err
			}

			return reportViolationop.NewAssociateReportViolationsNoContent(), nil
		})
}

// Get gets a list of PWS violations for a report
type GetReportViolationsHandler struct {
	handlers.HandlerConfig
	services.ReportViolationFetcher
}

// Handle gets a list of PWS violations for a report
func (h GetReportViolationsHandler) Handle(params reportViolationop.GetReportViolationsByReportIDParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			reports, err := h.FetchReportViolationsByReportID(appCtx, handlers.FmtUUIDToPop(params.ReportID))
			if err != nil {
				appCtx.Logger().Error("Error fetching report violations: ", zap.Error(err))
				return reportViolationop.NewGetReportViolationsByReportIDInternalServerError(), err
			}

			payload := payloads.ReportViolations(reports)
			return reportViolationop.NewGetReportViolationsByReportIDOK().WithPayload(payload), nil
		},
	)
}
