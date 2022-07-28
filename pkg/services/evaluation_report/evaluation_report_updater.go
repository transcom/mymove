package evaluationreport

import (
	"github.com/gofrs/uuid"

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
	return nil
}
