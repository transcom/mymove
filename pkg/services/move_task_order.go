package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"

	"github.com/transcom/mymove/pkg/models"
)

// HiddenMove struct used to store the MTO ID and the reason that the move is being hidden.
type HiddenMove struct {
	MTOID  uuid.UUID
	Reason string
}

// HiddenMoves is the slice of HiddenMove to return in the handler call
type HiddenMoves []HiddenMove

// MoveTaskOrderHider is the service object interface for Hide
//go:generate mockery -name MoveTaskOrderHider
type MoveTaskOrderHider interface {
	Hide() (HiddenMoves, error)
}

// MoveTaskOrderCreator is the service object interface for CreateMoveTaskOrder
//go:generate mockery -name MoveTaskOrderCreator
type MoveTaskOrderCreator interface {
	CreateMoveTaskOrder(moveTaskOrder *models.Move) (*models.Move, *validate.Errors, error)
}

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//go:generate mockery -name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.Move, error)
	ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.Move, error)
	ListAllMoveTaskOrders(isAvailableToPrime bool, since *int64) (models.Moves, error)
}

//MoveTaskOrderUpdater is the service object interface for updating fields of a MoveTaskOrder
//go:generate mockery -name MoveTaskOrderUpdater
type MoveTaskOrderUpdater interface {
	MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string, includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error)
	UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.Move, error)
	ShowHide(moveTaskOrderID uuid.UUID, show *bool) (*models.Move, error)
}

//MoveTaskOrderChecker is the service object interface for checking if a MoveTaskOrder is in a certain state
//go:generate mockery -name MoveTaskOrderChecker
type MoveTaskOrderChecker interface {
	MTOAvailableToPrime(moveTaskOrderID uuid.UUID) (bool, error)
}
