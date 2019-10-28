package movetaskorder

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

//ErrNotFound is returned when a given move task order is not found
type ErrNotFound struct {
	id uuid.UUID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("move task order id: %s not found", e.id.String())
}

//ErrInvalidInput is returned when an update to a move task order fails a validation rule
type ErrInvalidInput struct {
	id uuid.UUID
	error
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for move task order id: %s. %s", e.id.String(), e.error.Error())
}

type fetchMoveTaskOrder struct {
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &fetchMoveTaskOrder{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto := &models.MoveTaskOrder{}
	if err := f.db.Eager().Find(mto, moveTaskOrderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveTaskOrder{}, ErrNotFound{moveTaskOrderID}
		default:
			return &models.MoveTaskOrder{}, err
		}
	}
	return mto, nil
}

type updateMoveTaskOrderStatus struct {
	db *pop.Connection
	fetchMoveTaskOrder
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderStatusUpdater(db *pop.Connection) services.MoveTaskOrderStatusUpdater {
	moveTaskOrderFetcher := fetchMoveTaskOrder{db}
	return &updateMoveTaskOrderStatus{db, moveTaskOrderFetcher}
}

//UpdateMoveTaskOrderStatus updates the status of a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) UpdateMoveTaskOrderStatus(moveTaskOrderID uuid.UUID, status models.MoveTaskOrderStatus) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	mto.Status = status
	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, ErrInvalidInput{moveTaskOrderID, vErrors}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}
