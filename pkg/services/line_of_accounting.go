package services

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// LineOfAccountingFetcher is the exported interface for fetching lines of accounting
// based on a centralized location for business logic
//
//go:generate mockery --name LineOfAccountingFetcher
type LineOfAccountingFetcher interface {
	FetchLongLinesOfAccounting(departmentIndicator models.DepartmentIndicator, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.LineOfAccounting, error)
	BuildFullLineOfAccountingString(loa models.LineOfAccounting) string
}
