package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"

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

//MoveTaskOrderPrimeEstimatedWeightUpdater is the service object interface for UpdatePrimeEstimatedWeight
//go:generate mockery -name MoveTaskOrderPrimeEstimatedWeightUpdater
type MoveTaskOrderPrimeEstimatedWeightUpdater interface {
	UpdatePrimeEstimatedWeight(moveTaskOrderID uuid.UUID, primeEstimatedWeight unit.Pound, updateTime time.Time) (*models.MoveTaskOrder, error)
}

type PostCounselingInformation struct {
	PPMIsIncluded                                    bool
	ScheduledMoveDate                                time.Time
	SecondaryDeliveryAddress, SecondaryPickupAddress string
}

//MoveTaskOrderPrimePostCounselingUpdater is the service object interface for UpdateMoveTaskOrderPostCounselingInformation
//go:generate mockery -name MoveTaskOrderPrimePostCounselingUpdater
type MoveTaskOrderPrimePostCounselingUpdater interface {
	UpdateMoveTaskOrderPostCounselingInformation(moveTaskOrderID uuid.UUID, postCounselingInformation PostCounselingInformation) (*models.MoveTaskOrder, error)
}

type MoveTaskOrderDestinationAddressUpdater interface {
	UpdateMoveTaskOrderDestinationAddress(moveTaskOrderID uuid.UUID, destinationAddress *models.Address) (*models.MoveTaskOrder, error)
}
