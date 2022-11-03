package evaluationreport

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportDeleter struct {
}

func NewEvaluationReportDeleter() services.EvaluationReportDeleter {
	return &evaluationReportDeleter{}
}

func (o evaluationReportDeleter) DeleteEvaluationReport(appCtx appcontext.AppContext, reportID uuid.UUID) error {
	var report models.EvaluationReport
	err := appCtx.DB().Find(&report, reportID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(reportID, "attempting to delete evaluationReport")
		default:
			return apperror.NewQueryError("evaluationReport", err, "")
		}
	}

	sessionUserID := appCtx.Session().OfficeUserID

	if report.OfficeUserID != sessionUserID {
		appCtx.Logger().Info("Evaluation reports may only be edited by the user who created them.", zap.String("Evaluation Report ReportID", reportID.String()))

		return apperror.NewForbiddenError("Action not allowed")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Delete existing report_violations for this report
		existingReportViolations := models.ReportViolations{}
		err := appCtx.DB().Where("report_id in (?)", reportID).All(&existingReportViolations)
		if err != nil {
			return apperror.NewQueryError("EvaluationReport", err, "Unable to find existing report violations to remove")
		}
		err = appCtx.DB().Destroy(&existingReportViolations)
		if err != nil {
			return apperror.NewQueryError("EvaluationReport", err, "failed to delete existing report violations")
		}

		err = appCtx.DB().Destroy(&report)
		if err != nil {
			return apperror.NewQueryError("EvaluationReport", err, "failed to delete report")
		}
		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
