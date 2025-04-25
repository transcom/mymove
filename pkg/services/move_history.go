package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MoveHistoryFetcher is the exported interface for fetching a move by locator
//
//go:generate mockery --name MoveHistoryFetcher
type MoveHistoryFetcher interface {
	FetchMoveHistory(appCtx appcontext.AppContext, params *FetchMoveHistoryParams, useDatabaseProcInstead bool) (*models.MoveHistory, int64, error)
}

type FetchMoveHistoryParams struct {
	Locator string
	Page    *int64
	PerPage *int64
}
