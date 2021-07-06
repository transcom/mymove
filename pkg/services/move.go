package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveListFetcher is the exported interface for fetching multiple moves
//go:generate mockery --name MoveListFetcher --disable-version-string
type MoveListFetcher interface {
	FetchMoveList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Moves, error)
	FetchMoveCount(filters []QueryFilter) (int, error)
}

// MoveFetcher is the exported interface for fetching a move by locator
//go:generate mockery --name MoveFetcher --disable-version-string
type MoveFetcher interface {
	FetchMove(locator string, searchParams *MoveFetcherParams) (*models.Move, error)
}

// MoveFetcherParams is  public struct that's used to pass filter arguments to
// MoveFetcher queries
type MoveFetcherParams struct {
	IncludeHidden bool // indicates if a hidden/disabled move can be returned
}

// MoveRouter is the exported interface for routing moves at different stages
//go:generate mockery --name MoveRouter --disable-version-string
type MoveRouter interface {
	Approve(move *models.Move) error
	ApproveAmendedOrders(moveID uuid.UUID, orderID uuid.UUID) (models.Move, error)
	Cancel(reason string, move *models.Move) error
	CompleteServiceCounseling(move *models.Move) error
	SendToOfficeUser(move *models.Move) error
	Submit(move *models.Move) error
	SetLogger(logger Logger)
}
