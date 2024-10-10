package services

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// TransportationAccountingCodeFetcher is the exported interface for fetching transportation accounting codes
// based on a centralized location for business logic
//
//go:generate mockery --name TransportationAccountingCodeFetcher
type TransportationAccountingCodeFetcher interface {
	FetchOrderTransportationAccountingCodes(departmentIndicator models.DepartmentIndicator, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.TransportationAccountingCode, error)
}
