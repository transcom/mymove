package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// EvaluationReportFetcher is the service object interface for fetching all the evaluation reports for a move as a particular office user
//go:generate mockery --name EvaluationReportFetcher --disable-version-string
type EvaluationReportFetcher interface {
	FetchEvaluationReports(appCtx appcontext.AppContext, reportType models.EvaluationReportType, moveID uuid.UUID, officeUserID uuid.UUID) (models.EvaluationReports, error)
	FetchEvaluationReportByID(appCtx appcontext.AppContext, reportID uuid.UUID, officeUserID uuid.UUID) (*models.EvaluationReport, error)
}

//go:generate mockery --name EvaluationReportCreator --disable-version-string
type EvaluationReportCreator interface {
	CreateEvaluationReport(appCtx appcontext.AppContext, evaluationReport *models.EvaluationReport, locator string) (*models.EvaluationReport, error)
}

//go:generate mockery --name EvaluationReportUpdater --disable-version-string
type EvaluationReportUpdater interface {
	UpdateEvaluationReport(appCtx appcontext.AppContext, evaluationReport *models.EvaluationReport, officeUserID uuid.UUID, eTag string) error
}

//go:generate mockery --name EvaluationReportDeleter --disable-version-string
type EvaluationReportDeleter interface {
	DeleteEvaluationReport(appCtx appcontext.AppContext, reportID uuid.UUID) error
}
