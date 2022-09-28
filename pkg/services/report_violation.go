package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name ReportViolationFetcher --disable-version-string
type ReportViolationFetcher interface {
	FetchReportViolationsByReportID(appCtx appcontext.AppContext, reportID uuid.UUID) (models.ReportViolations, error)
}

//go:generate mockery --name ReportViolationsCreator --disable-version-string
type ReportViolationsCreator interface {
	AssociateReportViolations(appCtx appcontext.AppContext, reportViolations *models.ReportViolations, reportID uuid.UUID) error
}
