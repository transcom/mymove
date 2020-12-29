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

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery -name OrderUpdater
type OrderUpdater interface {
	UpdateOrder(eTag string, order models.Order) (*models.Order, error)
}

// ListMoveOrderParams is a public struct that's used to pass filter arguments to the ListMoveOrders
type ListMoveOrderParams struct {
	Branch                 *string
	Locator                *string
	DodID                  *string
	LastName               *string
	DestinationDutyStation *string
	Status                 []string
	Page                   *int64
	PerPage                *int64
	Sort                   *string
	Order                  *string
}
