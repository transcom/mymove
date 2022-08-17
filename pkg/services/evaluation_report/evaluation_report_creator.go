package evaluationreport

import (
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type evaluationReportCreator struct {
}

func NewEvaluationReportCreator() services.EvaluationReportCreator {
	return &evaluationReportCreator{}
}

func (o evaluationReportCreator) CreateEvaluationReport(appCtx appcontext.AppContext, report *models.EvaluationReport, locator string) (*models.EvaluationReport, error) {

	// check if it is a shipment or counseling report
	reportType := report.Type

	// counseling
	if reportType == models.EvaluationReportTypeCounseling {
		// get moveID via locator & make sure it exists
		var move models.Move
		err := appCtx.DB().Q().Where("locator = ?", locator).First(&move)
		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				return nil, models.ErrFetchNotFound
			}
			return nil, err
		}

		report.MoveID = move.ID
		report.Move = move
	}

	// shipment
	if reportType == models.EvaluationReportTypeShipment {
		// Need to get the Shipment for some report fields
		var shipment models.MTOShipment
		err := appCtx.DB().Q().Where("id = ?", report.ShipmentID).First(&shipment)
		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				return nil, models.ErrFetchNotFound
			}
			return nil, err
		}

		report.Shipment = &shipment
		report.MoveID = shipment.MoveTaskOrderID
	}

	verrs, err := appCtx.DB().ValidateAndCreate(report)
	if verrs.Count() != 0 || err != nil {
		return nil, err
	}

	return report, err
}
