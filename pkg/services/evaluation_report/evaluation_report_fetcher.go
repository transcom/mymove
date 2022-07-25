package evaluationreport

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportFetcher struct{}

func NewEvaluationReportFetcher() services.EvaluationReportFetcher {
	return &evaluationReportFetcher{}
}

func (f *evaluationReportFetcher) FetchEvaluationReports(appCtx appcontext.AppContext, reportType models.EvaluationReportType, moveID uuid.UUID, officeUserID uuid.UUID) (models.EvaluationReports, error) {
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

func (f *evaluationReportFetcher) FetchEvaluationReportByID(appCtx appcontext.AppContext, reportID uuid.UUID, officeUserID uuid.UUID) (*models.EvaluationReport, error) {
	var report models.EvaluationReport
	// Get the report by its ID, but don't return it if it's been soft-deleted.
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).EagerPreload("Move").Find(&report, reportID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reportID, "while looking for evaluation report")
		default:
			return nil, apperror.NewQueryError("EvaluationReport", err, "")
		}
	}
	// We shouldn't return the data if it's a draft (nil submitted_at) and the requester doesn't own it.
	if report.SubmittedAt == nil && report.OfficeUserID != officeUserID {
		return nil, apperror.NewForbiddenError("Draft evaluation reports are viewable only by their owner/creator.")
	}
	return &report, nil
}
