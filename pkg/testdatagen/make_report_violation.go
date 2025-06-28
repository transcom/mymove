package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeReportViolation(db *pop.Connection, assertions Assertions) (models.ReportViolation, error) {

	report := assertions.Report
	if isZeroUUID(assertions.Report.ID) {
		var err error
		report, err = makeEvaluationReport(db, assertions)
		if err != nil {
			return models.ReportViolation{}, err
		}
	}

	violation := assertions.Violation
	if isZeroUUID(assertions.Violation.ID) {
		violation = MakePWSViolation(db, assertions)
	}

	reportViolation := models.ReportViolation{
		ID:          uuid.Must(uuid.NewV4()),
		ReportID:    report.ID,
		Violation:   violation,
		ViolationID: violation.ID,
	}

	mergeModels(&reportViolation, assertions.ReportViolation)
	mustCreate(db, &reportViolation, assertions.Stub)

	return reportViolation, nil
}
