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
//go:generate mockery --name MoveTaskOrderHider --disable-version-string
type MoveTaskOrderHider interface {
	Hide() (HiddenMoves, error)
}

// MoveTaskOrderCreator is the service object interface for CreateMoveTaskOrder
//go:generate mockery --name MoveTaskOrderCreator --disable-version-string
type MoveTaskOrderCreator interface {
	CreateMoveTaskOrder(moveTaskOrder *models.Move) (*models.Move, *validate.Errors, error)
}

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//go:generate mockery --name MoveTaskOrderFetcher --disable-version-string
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(moveTaskOrderID uuid.UUID, searchParams *MoveTaskOrderFetcherParams) (*models.Move, error)
	ListMoveTaskOrders(orderID uuid.UUID, searchParams *MoveTaskOrderFetcherParams) ([]models.Move, error)
	ListAllMoveTaskOrders(searchParams *MoveTaskOrderFetcherParams) (models.Moves, error)
}

//MoveTaskOrderUpdater is the service object interface for updating fields of a MoveTaskOrder
//go:generate mockery --name MoveTaskOrderUpdater --disable-version-string
type MoveTaskOrderUpdater interface {
	MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string, includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error)
	UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.Move, error)
	UpdateStatusServiceCounselingCompleted(moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error)
	ShowHide(moveTaskOrderID uuid.UUID, show *bool) (*models.Move, error)
}

//MoveTaskOrderChecker is the service object interface for checking if a MoveTaskOrder is in a certain state
//go:generate mockery --name MoveTaskOrderChecker --disable-version-string
type MoveTaskOrderChecker interface {
	MTOAvailableToPrime(moveTaskOrderID uuid.UUID) (bool, error)
}

// MoveTaskOrderFetcherParams is a public struct that's used to pass filter arguments to
// ListMoveTaskOrders, ListAllMoveTaskOrders, and FetchMoveTaskOrder queries
type MoveTaskOrderFetcherParams struct {
	IsAvailableToPrime bool   // indicates if all MTOs returned must be Prime-available
	IncludeHidden      bool   // indicates if hidden/disabled MTOs should be included in the output
	Since              *int64 // if filled, only MTOs that have been updated after this timestamp will be returned
}
