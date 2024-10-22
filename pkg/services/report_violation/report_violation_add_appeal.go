package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type reportViolationAddAppeal struct{}

func NewReportViolationsAddAppeal() services.ReportViolationsAddAppeal {
	return &reportViolationAddAppeal{}
}

func (f *reportViolationAddAppeal) AddAppealToViolation(appCtx appcontext.AppContext, reportID uuid.UUID, reportViolationID uuid.UUID, officeUserID uuid.UUID, remarks string, appealStatus string) (models.EvaluationReport, error) {
	reportViolations := models.EvaluationReport{}
	if reportID == uuid.Nil {
		return models.EvaluationReport{}, apperror.NewBadDataError("reportID must be provided")
	}
	if reportViolationID == uuid.Nil {
		return models.EvaluationReport{}, apperror.NewBadDataError("reportViolationID must be provided")
	}
	if officeUserID == uuid.Nil {
		return models.EvaluationReport{}, apperror.NewBadDataError("officeUserID must be provided")
	}

	var appealDecision models.AppealStatus
	if appealStatus == "sustained" {
		appealDecision = models.AppealStatusSustained
	} else {
		appealDecision = models.AppealStatusRejected
	}

	violationAppeal := models.GsrAppeal{
		EvaluationReportID:      reportID,
		ReportViolationID:       &reportViolationID,
		OfficeUserID:            officeUserID,
		IsSeriousIncidentAppeal: models.BoolPointer(false),
		AppealStatus:            appealDecision,
		Remarks:                 remarks,
	}

	verrs, err := appCtx.DB().ValidateAndCreate(&violationAppeal)
	if verrs != nil && verrs.HasAny() {
		return models.EvaluationReport{}, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while associating report violations.")
	} else if err != nil {
		return models.EvaluationReport{}, apperror.NewQueryError("reportViolations", err, "")
	}

	return reportViolations, nil
}
