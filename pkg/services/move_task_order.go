package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"

	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrderCreator is the service object interface for CreateMoveTaskOrder
//go:generate mockery -name MoveTaskOrderCreator
type MoveTaskOrderCreator interface {
	CreateMoveTaskOrder(moveTaskOrder *models.MoveTaskOrder) (*models.MoveTaskOrder, *validate.Errors, error)
}

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//go:generate mockery -name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error)
	ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.MoveTaskOrder, error)
	ListAllMoveTaskOrders(isAvailableToPrime bool, since *int64) (models.MoveTaskOrders, error)
}

//MoveTaskOrderUpdater is the service object interface for MakeAvailableToPrime
//go:generate mockery -name MoveTaskOrderUpdater
type MoveTaskOrderUpdater interface {
	MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string) (*models.MoveTaskOrder, error)
	UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.MoveTaskOrder, error)
}
