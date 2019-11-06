package movetaskorder

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"

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
	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, ErrInvalidInput{moveTaskOrderID, vErrors}
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
		return &models.MoveTaskOrder{}, ErrInvalidInput{moveTaskOrderID, vErrors}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

type updatePostCounselingInfo struct {
	db *pop.Connection
	fetchMoveTaskOrder
}

// NewMoveTaskOrderActualWeightUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderPostCounselingInfoUpdater(db *pop.Connection) services.MoveTaskOrderPostCounselingInfoUpdater {
	moveTaskOrderFetcher := fetchMoveTaskOrder{db}
	return &updatePostCounselingInfo{db, moveTaskOrderFetcher}
}

//UpdatePostCounselingInfo updates the actual weight of a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, scheduledMoveDate strfmt.Date, secondaryPickupAddress ghcmessages.Address, secondaryDeliveryAddress ghcmessages.Address, ppmIsIncluded bool) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}

	mto.ScheduledMoveDate = scheduledMoveDate
	mto.SecondaryPickupAddress.StreetAddress1 = *secondaryPickupAddress.StreetAddress1
	mto.SecondaryPickupAddress.City = *secondaryPickupAddress.City
	mto.SecondaryPickupAddress.State = *secondaryPickupAddress.State
	mto.SecondaryDeliveryAddress.StreetAddress1 = *secondaryDeliveryAddress.StreetAddress1
	mto.SecondaryDeliveryAddress.City = *secondaryDeliveryAddress.City
	mto.SecondaryDeliveryAddress.State = *secondaryDeliveryAddress.State
	mto.PpmIsIncluded = ppmIsIncluded

	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, ErrInvalidInput{moveTaskOrderID, vErrors}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

type updateDestinationAddress struct {
	db *pop.Connection
	fetchMoveTaskOrder
}

// NewMoveTaskOrderDestinationAddressUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderDestinationAddressUpdater(db *pop.Connection) services.MoveTaskOrderPostCounselingInfoUpdater {
	moveTaskOrderFetcher := fetchMoveTaskOrder{db}
	return &updateDestinationAddress{db, moveTaskOrderFetcher}
}

//UpdatePostCounselingInfo updates the actual weight of a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) UpdateDestinationAddress(moveTaskOrderID uuid.UUID, destinationAddress ghcmessages.Address) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}

	mto.DestinationAddress.StreetAddress1 = *destinationAddress.StreetAddress1
	mto.DestinationAddress.City = *destinationAddress.City
	mto.DestinationAddress.State = *destinationAddress.State

	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, ErrInvalidInput{moveTaskOrderID, vErrors}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}