package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrderFetcher is the service object interface for MoveTaskOrderFetch
//go:generate mockery -name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error)
}

//MoveTaskOrderFetcher is the service object interface for MoveTaskOrderFetch
//go:generate mockery -name MoveTaskOrderStatusUpdater
type MoveTaskOrderStatusUpdater interface {
	UpdateMoveTaskOrderStatus(moveTaskOrderID uuid.UUID, status models.MoveTaskOrderStatus) (*models.MoveTaskOrder, error)
}
