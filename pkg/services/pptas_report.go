package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PPTASReportListFetcher is the exported interface for fetching reports
//
//go:generate mockery --name PPTASReportListFetcher
type PPTASReportListFetcher interface {
	GetMovesForReportBuilder(appCtx appcontext.AppContext, params *MoveTaskOrderFetcherParams) (models.Moves, error)
	BuildPPTASReportsFromMoves(appCtx appcontext.AppContext, moves models.Moves) (models.PPTASReports, error)
}
