package evaluationreport

import (
	"database/sql"
	"fmt"

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
		return apperror.NewForbiddenError("A report may only be saved by the office user that created it")
	}

	if originalReport.SubmittedAt != nil {
		return apperror.NewConflictError(evaluationReport.ID, "reports that have already been submitted cannot be updated")
	}

	evaluationReport.MoveID = originalReport.MoveID
	evaluationReport.ShipmentID = originalReport.ShipmentID
	evaluationReport.Type = originalReport.Type
	fmt.Println("REMARKS LOL")
	fmt.Println(evaluationReport.Remarks)
	verrs, err := appCtx.DB().ValidateAndSave(evaluationReport)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		return apperror.NewInvalidInputError(evaluationReport.ID, err, verrs, "")
	}
	return nil
}
