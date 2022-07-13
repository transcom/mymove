package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportListFetcher struct{}

func NewEvaluationReportListFetcher() services.EvaluationReportListFetcher {
	return &evaluationReportListFetcher{}
}

func (f *evaluationReportListFetcher) FetchEvaluationReports(appCtx appcontext.AppContext, reportType models.EvaluationReportType, moveID uuid.UUID, officeUserID uuid.UUID) (models.EvaluationReports, error) {
	reports := models.EvaluationReports{}
	if moveID == uuid.Nil {
		return nil, apperror.NewBadDataError("moveID must be provided")
	}
	if officeUserID == uuid.Nil {
		return nil, apperror.NewBadDataError("officeUserID must be provided")
	}

	err := appCtx.DB().
		Scope(utilities.ExcludeDeletedScope()).
		Where("move_id = ?", moveID).
		Where("type = ?", reportType).
		Where("(submitted_at IS NOT NULL OR office_user_id = ?)", officeUserID).
		Order("submitted_at ASC, created_at ASC").
		All(&reports)

	if err != nil {
		return nil, apperror.NewQueryError("EvaluationReport", err, "")
	}
	return reports, nil
}
