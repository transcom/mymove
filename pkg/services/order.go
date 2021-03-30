package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// OrderFetcher is the service object interface for FetchOrder
//go:generate mockery -name OrderFetcher
type OrderFetcher interface {
	FetchOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListOrders(officeUserID uuid.UUID, params *ListOrderParams) ([]models.Move, int, error)
}

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery -name OrderUpdater
type OrderUpdater interface {
	UpdateOrder(eTag string, order models.Order) (*models.Order, error)
}

// ListOrderParams is a public struct that's used to pass filter arguments to the ListOrders
type ListOrderParams struct {
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
