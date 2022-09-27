package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type reportViolationFetcher struct{}

func NewReportViolationFetcher() services.ReportViolationFetcher {
	return &reportViolationFetcher{}
}

func (f *reportViolationFetcher) FetchReportViolationsByReportID(appCtx appcontext.AppContext, reportID uuid.UUID) (models.ReportViolations, error) {
	reportViolations := models.ReportViolations{}
	if reportID == uuid.Nil {
		return nil, apperror.NewBadDataError("reportID must be provided")
	}

	err := appCtx.DB().
		EagerPreload("Violation").
		Where("report_id = ?", reportID).
		All(&reportViolations)

	if err != nil {
		return nil, apperror.NewQueryError("ReportViolation", err, "")
	}

	return reportViolations, nil
}
