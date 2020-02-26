package services

import (
	"github.com/gofrs/uuid"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"

	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//go:generate mockery -name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error)
	ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.MoveTaskOrder, error)
}

//MoveTaskOrderStatusUpdater is the service object interface for MakeAvailableToPrime
//go:generate mockery -name MoveTaskOrderStatusUpdater
type MoveTaskOrderStatusUpdater interface {
	MakeAvailableToPrime(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error)
	UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.MoveTaskOrder, error)
}
