package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//go:generate mockery -name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error)
}

//MoveTaskOrderStatusUpdater is the service object interface for UpdateMoveTaskOrderStatus
//go:generate mockery -name MoveTaskOrderStatusUpdater
type MoveTaskOrderStatusUpdater interface {
	UpdateMoveTaskOrderStatus(moveTaskOrderID uuid.UUID, status models.MoveTaskOrderStatus) (*models.MoveTaskOrder, error)
}

//MoveTaskOrderActualWeightUpdater is the service object interface for UpdateMoveTaskOrderActualWeight
//go:generate mockery -name MoveTaskOrderActualWeightUpdater
type MoveTaskOrderActualWeightUpdater interface {
	UpdateMoveTaskOrderActualWeight(moveTaskOrderID uuid.UUID, actualWeight int64) (*models.MoveTaskOrder, error)
}
