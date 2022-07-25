package evaluationreport

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportDeleter struct {
}

func NewEvaluationReportDeleter() services.EvaluationReportDeleter {
	return &evaluationReportDeleter{}
}

func (o evaluationReportDeleter) DeleteEvaluationReport(appCtx appcontext.AppContext, reportID uuid.UUID) error {
	var report models.EvaluationReport
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Find(&report, reportID)
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
		appCtx.Logger().Warn("Evaluation reports may only be edited by the user who created them.", zap.String("Evaluation Report ReportID", reportID.String()))

		return apperror.NewForbiddenError("Action not allowed")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err := utilities.SoftDestroy(appCtx.DB(), &report)
		if err != nil {
			switch err.Error() {
			case "error updating model":
				return apperror.NewUnprocessableEntityError("while updating model")
			default:
				return apperror.NewInternalServerError("failed attempt to soft delete model")
			}
		}
		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
