package evaluationreport

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type evaluationReportUpdater struct {
}

func NewEvaluationReportUpdater() services.EvaluationReportUpdater {
	return &evaluationReportUpdater{}
}

func (u evaluationReportUpdater) UpdateEvaluationReport(appCtx appcontext.AppContext, evaluationReport *models.EvaluationReport, officeUserID uuid.UUID, eTag string) error {
	var originalReport models.EvaluationReport
	err := appCtx.DB().Find(&originalReport, evaluationReport.ID)
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

	if evaluationReport.ViolationsObserved != nil && !*evaluationReport.ViolationsObserved {
		// Check to see if there are existing report_violations for this report
		existingReportViolations := models.ReportViolations{}
		err = appCtx.DB().Where("report_id = ?", evaluationReport.ID).All(&existingReportViolations)
		if err != nil {
			return apperror.NewQueryError("EvaluationReport", err, "Unable to find existing report violations")
		}
		// Delete the existing reportViolations
		if len(existingReportViolations) > 0 {
			err = appCtx.DB().Destroy(existingReportViolations)
			if err != nil {
				return apperror.NewQueryError("EvaluationReport", err, "failed to delete existing report violations")
			}
		}
	}

	verrs, err := appCtx.DB().ValidateAndSave(evaluationReport)
	if err != nil {
		return apperror.NewQueryError("EvaluationReport", err, "failed to save the evaluation report")
	}
	if verrs.HasAny() {
		return apperror.NewInvalidInputError(evaluationReport.ID, err, verrs, "")
	}
	return nil
}

// SubmitEvaluationReport sets a submittedAt value on an EvaluationReport. This 'finalizes' a report as ready for
// sharing with others.
func (u evaluationReportUpdater) SubmitEvaluationReport(appCtx appcontext.AppContext, evaluationReportID uuid.UUID, officeUserID uuid.UUID, eTag string) error {
	var originalReport models.EvaluationReport
	err := appCtx.DB().Find(&originalReport, evaluationReportID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(originalReport.ID, "")
		default:
			return err
		}
	}

	if etag.GenerateEtag(originalReport.UpdatedAt) != eTag {
		return apperror.NewPreconditionFailedError(originalReport.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if officeUserID != originalReport.OfficeUserID {
		return apperror.NewForbiddenError("A report may only be submitted by the office user that created it")
	}

	if originalReport.SubmittedAt != nil {
		return apperror.NewConflictError(originalReport.ID, "reports that have already been submitted cannot be updated")
	}
	// Make sure the report is ready for submission.
	err = isValidForSubmission(originalReport)
	if err != nil {
		return err
	}

	now := time.Now()
	originalReport.SubmittedAt = &now
	verrs, err := appCtx.DB().ValidateAndSave(&originalReport)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		return apperror.NewInvalidInputError(originalReport.ID, err, verrs, "")
	}

	return nil
}

// This function checks to ensure that required fields are filled out and that dependent field requirements are satisfied.
// It returns an error object if the report is out of alignment with a rule and returns nil if it's ready to go.
// TODO: Update this so it also applies rules for violations.
func isValidForSubmission(evaluationReport models.EvaluationReport) error {
	// Required field InspectionDate
	if evaluationReport.InspectionDate == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without an Inspection Date.",
				evaluationReport.ID))
	}
	// Required field InspectionType
	if evaluationReport.InspectionType == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without an Inspection Type.",
				evaluationReport.ID))
	}
	// Travel time required when inspection type is physical
	if *evaluationReport.InspectionType == models.EvaluationReportInspectionTypePhysical && evaluationReport.TravelTimeMinutes == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without travel time if the location is physical.",
				evaluationReport.ID))
	}
	// Required field location
	if evaluationReport.Location == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without a location.",
				evaluationReport.ID))
	}
	// LocationDescription is required when Location is Other
	if *evaluationReport.Location == models.EvaluationReportLocationTypeOther && evaluationReport.LocationDescription == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without location description if the location is other.",
				evaluationReport.ID))
	}
	// Required field EvaluationLengthMinutes
	if evaluationReport.EvaluationLengthMinutes == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without an evaluation length.",
				evaluationReport.ID))
	}
	// Required field ViolationsObserved
	if evaluationReport.ViolationsObserved == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without a value for violations observed.",
				evaluationReport.ID))
	}
	// Required field Remarks
	if evaluationReport.Remarks == nil {
		return errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Evaluation report with ID %s cannot be submitted without QAE remarks.",
				evaluationReport.ID))
	}

	return nil
}
