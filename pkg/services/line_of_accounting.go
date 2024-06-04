package services

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// LinesOfAccountingFetcher is the exported interface for fetching lines of accounting
// based on a centralized location for business logic
//
//go:generate mockery --name LinesOfAccountingFetcher
type LinesOfAccountingFetcher interface {
	FetchLongLinesOfAccounting(serviceMemberAffiliation models.ServiceMemberAffiliation, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.LineOfAccounting, error)
}
