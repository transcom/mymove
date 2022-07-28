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

func (o evaluationReportCreator) CreateEvaluationReport(appCtx appcontext.AppContext, report *models.EvaluationReport) (*models.EvaluationReport, error) {

	// Note: this assumes we are creating a shipment report. When adding counceling eval reports this will need tweaked.
	// Need to get the Shipment for some report fields
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Where("id = ?", report.ShipmentID).First(&shipment)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, models.ErrFetchNotFound
		}
		return nil, err
	}
	report.MoveID = shipment.MoveTaskOrderID

	verrs, err := appCtx.DB().ValidateAndCreate(report)
	if verrs.Count() != 0 || err != nil {
		return nil, err
	}

	return report, err
}
