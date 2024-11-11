package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name ReportViolationFetcher
type ReportViolationFetcher interface {
	FetchReportViolationsByReportID(appCtx appcontext.AppContext, reportID uuid.UUID) (models.ReportViolations, error)
}

//go:generate mockery --name ReportViolationsCreator
type ReportViolationsCreator interface {
	AssociateReportViolations(appCtx appcontext.AppContext, reportViolations *models.ReportViolations, reportID uuid.UUID) error
}

//go:generate mockery --name ReportViolationsAddAppeal
type ReportViolationsAddAppeal interface {
	AddAppealToViolation(appCtx appcontext.AppContext, reportID uuid.UUID, reportViolationID uuid.UUID, officeUserID uuid.UUID, remarks string, appealStatus string) (models.GsrAppeal, error)
}
