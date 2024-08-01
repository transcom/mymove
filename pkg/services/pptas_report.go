package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PPTASReportListFetcher is the exported interface for fetching reports
//
//go:generate mockery --name ReportListFetcher
type PPTASReportListFetcher interface {
	BuildPPTASReportsFromMoves(appCtx appcontext.AppContext, params *MoveTaskOrderFetcherParams) (models.PPTASReports, error)
}
