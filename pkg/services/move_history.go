package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MoveHistoryFetcher is the exported interface for fetching a move by locator
//go:generate mockery --name MoveHistoryFetcher --disable-version-string
type MoveHistoryFetcher interface {
	FetchMoveHistory(appCtx appcontext.AppContext, locator string) (*models.MoveHistory, error)
}
