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
	CreateEvaluationReport(appCtx appcontext.AppContext, evaluationReport *models.EvaluationReport) (*models.EvaluationReport, error)
}
