package services

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveOrderFetcher is the service object interface for FetchMoveOrder
//go:generate mockery -name MoveOrderFetcher
type MoveOrderFetcher interface {
	FetchMoveOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListMoveOrders(officeUserID uuid.UUID, options ...func(query *pop.Query)) ([]models.Order, error)
}

//MoveOrderUpdater is the service object interface for updating fields of a MoveOrder
//go:generate mockery -name MoveOrderUpdater
type MoveOrderUpdater interface {
	UpdateMoveOrder(moveOrderID uuid.UUID, eTag string, moveOrder models.Order) (*models.Order, error)
}
