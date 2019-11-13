package movetaskorder

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

//ErrNotFound is returned when a given move task order is not found
type ErrNotFound struct {
	id uuid.UUID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("move task order id: %s not found", e.id.String())
}

type errInvalidInput struct {
	id uuid.UUID
	error
	validationErrors map[string][]string
}

//ErrInvalidInput is returned when an update to a move task order fails a validation rule
type ErrInvalidInput struct {
	errInvalidInput
}

func NewErrInvalidInput(id uuid.UUID, err error, validationErrors map[string][]string) ErrInvalidInput {
	return ErrInvalidInput{
		errInvalidInput{
			id:               id,
			error:            err,
			validationErrors: validationErrors,
		},
	}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for move task order id: %s. %s", e.id.String(), e.InvalidFields())
}

func (e ErrInvalidInput) InvalidFields() map[string]string {
	es := make(map[string]string)
	if e.validationErrors == nil {
		return es
	}
	for k, v := range e.validationErrors {
		es[k] = strings.Join(v, " ")
	}
	return es
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

// NewMoveTaskOrderStatusUpdater creates a new struct with the service dependencies
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

	// not sure if DRAFT is the right status but going with this for now
	if status == models.MoveTaskOrderStatusDraft {
		mto.AvailableToPrimeDate = time.Now()
	}

	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, NewErrInvalidInput(moveTaskOrderID, err, vErrors.Errors)
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

type updateMoveTaskOrderEstimatedWeight struct {
	db *pop.Connection
	fetchMoveTaskOrder
}

func NewMoveTaskOrderEstimatedWeightUpdater(db *pop.Connection) services.MoveTaskOrderPrimeEstimatedWeightUpdater {
	moveTaskOrderFetcher := fetchMoveTaskOrder{db}
	return &updateMoveTaskOrderEstimatedWeight{db, moveTaskOrderFetcher}
}

//UpdateEstimatedWeight updates the estimated weight of a MoveTaskOrder for a given UUID
func (f updateMoveTaskOrderEstimatedWeight) UpdatePrimeEstimatedWeight(moveTaskOrderID uuid.UUID, primeEstimatedWeight unit.Pound, updateTime time.Time) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	mto.PrimeEstimatedWeight = &primeEstimatedWeight
	mto.PrimeEstimatedWeightRecordedDate = &updateTime
	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, NewErrInvalidInput(moveTaskOrderID, err, vErrors.Errors)
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

type updateMoveTaskOrderActualWeight struct {
	db *pop.Connection
	fetchMoveTaskOrder
}

// NewMoveTaskOrderActualWeightUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderActualWeightUpdater(db *pop.Connection) services.MoveTaskOrderActualWeightUpdater {
	moveTaskOrderFetcher := fetchMoveTaskOrder{db}
	return &updateMoveTaskOrderActualWeight{db, moveTaskOrderFetcher}
}

//UpdateMoveTaskOrderActualWeight updates the actual weight of a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) UpdateMoveTaskOrderActualWeight(moveTaskOrderID uuid.UUID, actualWeight int64) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	weight := unit.Pound(actualWeight)
	mto.ActualWeight = &weight

	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, NewErrInvalidInput(moveTaskOrderID, err, vErrors.Errors)
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}
