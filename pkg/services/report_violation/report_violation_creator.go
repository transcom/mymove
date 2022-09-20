package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type reportViolationsCreator struct {
}

func NewReportViolationCreator() services.ReportViolationsCreator {
	return &reportViolationsCreator{}
}

// Remove all existing violations associations for a report and replace them with associations to the provided violations
func (u reportViolationsCreator) AssociateReportViolations(appCtx appcontext.AppContext, reportViolations *models.ReportViolations, reportID uuid.UUID) error {

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		// Delete all existing report_violations for this report
		existingReportViolations := models.ReportViolations{}
		err := appCtx.DB().Where("report_id in (?)", reportID).All(&existingReportViolations)
		if err != nil {
			return err
		}
		err = appCtx.DB().Destroy(existingReportViolations)
		if err != nil {
			return err
		}

		// Create new violations associations for the report
		if len(*reportViolations) > 0 {
			verrs, err := appCtx.DB().ValidateAndCreate(reportViolations)
			if verrs.Count() != 0 || err != nil {
				return err
			}
		}
		return nil
	})
	if txnErr != nil {
		return txnErr
	}
	return nil

}
