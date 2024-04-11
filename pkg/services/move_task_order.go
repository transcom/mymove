package services

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type MoveOrderAmendmentAvailableSinceCount struct {
	MoveID              uuid.UUID
	Total               int
	AvailableSinceTotal int
}

type MoveOrderAmendmentAvailableSinceCounts []MoveOrderAmendmentAvailableSinceCount

// HiddenMove struct used to store the MTO ID and the reason that the move is being hidden.
type HiddenMove struct {
	MTOID  uuid.UUID
	Reason string
}

// HiddenMoves is the slice of HiddenMove to return in the handler call
type HiddenMoves []HiddenMove

// MoveTaskOrderHider is the service object interface for Hide
//
//go:generate mockery --name MoveTaskOrderHider
type MoveTaskOrderHider interface {
	Hide(appCtx appcontext.AppContext) (HiddenMoves, error)
}

// MoveTaskOrderCreator is the service object interface for CreateMoveTaskOrder
//
//go:generate mockery --name MoveTaskOrderCreator
type MoveTaskOrderCreator interface {
	CreateMoveTaskOrder(appCtx appcontext.AppContext, moveTaskOrder *models.Move) (*models.Move, *validate.Errors, error)
}

// MoveTaskOrderFetcher is the service object interface for FetchMoveTaskOrder
//
//go:generate mockery --name MoveTaskOrderFetcher
type MoveTaskOrderFetcher interface {
	FetchMoveTaskOrder(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams) (*models.Move, error)
	ListAllMoveTaskOrders(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams) (models.Moves, error)
	ListPrimeMoveTaskOrders(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams) (models.Moves, error)
	ListNewPrimeMoveTaskOrders(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams) (models.Moves, int, error)
	GetMove(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams, eagerAssociations ...string) (*models.Move, error)
	ListPrimeMoveTaskOrdersV2(appCtx appcontext.AppContext, searchParams *MoveTaskOrderFetcherParams) (models.Moves, MoveOrderAmendmentAvailableSinceCounts, error)
}

// MoveTaskOrderUpdater is the service object interface for updating fields of a MoveTaskOrder
//
//go:generate mockery --name MoveTaskOrderUpdater
type MoveTaskOrderUpdater interface {
	MakeAvailableToPrime(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string, includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error)
	UpdatePostCounselingInfo(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error)
	UpdateStatusServiceCounselingCompleted(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error)
	UpdateReviewedBillableWeightsAt(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error)
	UpdateTIORemarks(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string, remarks string) (*models.Move, error)
	ShowHide(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, show *bool) (*models.Move, error)
	UpdatePPMType(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID) (*models.Move, error)
}

// MoveTaskOrderChecker is the service object interface for checking if a MoveTaskOrder is in a certain state
//
//go:generate mockery --name MoveTaskOrderChecker
type MoveTaskOrderChecker interface {
	MTOAvailableToPrime(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID) (bool, error)
}

// MoveTaskOrderFetcherParams is a public struct that's used to pass filter arguments to
// ListAllMoveTaskOrders, and FetchMoveTaskOrder queries
type MoveTaskOrderFetcherParams struct {
	IsAvailableToPrime       bool       // indicates if all MTOs returned must be Prime-available
	IncludeHidden            bool       // indicates if hidden/disabled MTOs should be included in the output
	Since                    *time.Time // if filled, only MTOs that have been updated after this timestamp will be returned
	MoveTaskOrderID          uuid.UUID  // ID of the move task order
	Locator                  string     // the locator is a unique string that identifies the move
	MoveCode                 *string    // moveCode that is passed in when searching through prime moves
	ID                       *string    // id of the move that is sent in when searching through prime moves
	ExcludeExternalShipments bool       // indicates if external vendor shipments should be returned
	Page                     *int64
	PerPage                  *int64
}
