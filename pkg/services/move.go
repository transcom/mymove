package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// MoveListFetcher is the exported interface for fetching multiple moves
//go:generate mockery -name MoveListFetcher
type MoveListFetcher interface {
	FetchMoveList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Moves, error)
	FetchMoveCount(filters []QueryFilter) (int, error)
}

// MoveFetcher is the exported interface for fetching a move by locator
//go:generate mockery -name MoveFetcher
type MoveFetcher interface {
	FetchMove(locator string) (*models.Move, error)
}