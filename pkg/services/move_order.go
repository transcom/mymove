package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveOrderFetcher is the service object interface for FetchMoveOrder
//go:generate mockery -name MoveOrderFetcher
type MoveOrderFetcher interface {
	FetchMoveOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListMoveOrders() ([]models.Order, error)
}
