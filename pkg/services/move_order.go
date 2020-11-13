package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveOrderFetcher is the service object interface for FetchMoveOrder
//go:generate mockery -name MoveOrderFetcher
type MoveOrderFetcher interface {
	FetchMoveOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListMoveOrders(officeUserID uuid.UUID, params *ListMoveOrderParams) ([]models.Move, int, error)
}

//MoveOrderUpdater is the service object interface for updating fields of a MoveOrder
//go:generate mockery -name MoveOrderUpdater
type MoveOrderUpdater interface {
	UpdateMoveOrder(moveOrderID uuid.UUID, eTag string, moveOrder models.Order) (*models.Order, error)
}

// ListMoveOrderParams is a public struct that's used to pass filter arguments to the ListMoveOrders
type ListMoveOrderParams struct {
	Branch                 *string
	MoveID                 *string
	DodID                  *string
	LastName               *string
	DestinationDutyStation *string
	Status                 []string
	Page                   *int64
	PerPage                *int64
}
