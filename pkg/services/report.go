package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ReportListFetcher is the exported interface for fetching reports
//
//go:generate mockery --name ReportListFetcher
type ReportListFetcher interface {
	FetchReportList(appCtx appcontext.AppContext, params *FetchPaymentRequestListParams) (models.Reports, error)
}
