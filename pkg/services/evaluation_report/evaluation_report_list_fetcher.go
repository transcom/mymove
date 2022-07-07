package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type evaluationReportListFetcher struct{}

func NewEvaluationReportListFetcher() services.EvaluationReportListFetcher {
	return &evaluationReportListFetcher{}
}

func (f *evaluationReportListFetcher) FetchEvaluationReports(appCtx appcontext.AppContext, moveID uuid.UUID, officeUserID uuid.UUID) (models.EvaluationReports, error) {
	return nil, nil
}
