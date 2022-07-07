package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// EvaluationReportListFetcher is the service object interface for fetching all the evaluation reports for a move as a particular office user
//go:generate mockery --name EvaluationReportListFetcher --disable-version-string
type EvaluationReportListFetcher interface {
	FetchEvaluationReports(appCtx appcontext.AppContext, moveID uuid.UUID, officeUserID uuid.UUID) (models.EvaluationReports, error)
}
