package evaluationreport

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportUpdater struct {
}

func NewEvaluationReportUpdater() services.EvaluationReportUpdater {
	return &evaluationReportUpdater{}
}

// does this take an evaluation report ID separately or will we get it from the model?
// I think we'll find out which works better once we get to the handler
func (u evaluationReportUpdater) UpdateEvaluationReport(appCtx appcontext.AppContext, evaluationReport *models.EvaluationReport, officeUserID uuid.UUID, eTag string) error {
	var originalReport models.EvaluationReport
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Find(&originalReport, evaluationReport.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(evaluationReport.ID, "")
		default:
			return err
		}
	}

	if etag.GenerateEtag(originalReport.UpdatedAt) != eTag {
		return apperror.NewPreconditionFailedError(evaluationReport.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if officeUserID != originalReport.OfficeUserID {
		return apperror.NewForbiddenError("A report may only be saved by the user that created it")
	}

	if originalReport.SubmittedAt != nil {
		return apperror.NewConflictError(evaluationReport.ID, "reports that have already been submitted cannot be updated")
	}

	evaluationReport.MoveID = originalReport.MoveID
	evaluationReport.ShipmentID = originalReport.ShipmentID

	verrs, err := appCtx.DB().ValidateAndUpdate(evaluationReport)
	if verrs.HasAny() || err != nil {
		return err
	}
	return nil
}
